package service

import (
	"errors"

	"time"

	"github.com/LuuDinhTheTai/tzone/internal/model"
	"github.com/LuuDinhTheTai/tzone/internal/repository"

	"github.com/LuuDinhTheTai/tzone/util/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// auth
type AuthService struct {
	userRepo  repository.UserRepository
	tokenRepo repository.RefreshTokenRepository
}

func NewAuthService(userRepo repository.UserRepository, tokenRepo repository.RefreshTokenRepository) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

// register
func (s *AuthService) Register(email string, password string) error {

	// check email tồn tại
	existing, _ := s.userRepo.FindByEmail(email)
	if existing != nil {
		return errors.New("email already exists")
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return err
	}

	user := model.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hash),
	}

	return s.userRepo.Create(&user)
}

// login
func (s *AuthService) Login(email string, password string) (string, string, *model.User, error) {

	user, err := s.userRepo.FindByEmail(email)

	if err != nil {
		return "", "", nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)

	if err != nil {
		return "", "", nil, errors.New("invalid email or password")
	}

	jti := uuid.New()
	accessToken, refreshToken, err := jwt.GenerateTokenPair(user.ID, jti)
	if err != nil {
		return "", "", nil, errors.New("failed to generate tokens")
	}

	// Save Refresh Token in DB
	rtRecord := &model.RefreshToken{
		ID:        jti,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := s.tokenRepo.Create(rtRecord); err != nil {
		return "", "", nil, errors.New("failed to save session")
	}

	return accessToken, refreshToken, user, nil
}

// RefreshToken handles generating a new token pair from a valid refresh token
func (s *AuthService) RefreshToken(tokenString string) (string, string, uuid.UUID, error) {
	userID, jti, err := jwt.ValidateRefreshToken(tokenString)
	if err != nil {
		return "", "", uuid.Nil, errors.New("invalid or expired refresh token")
	}

	// Check if this JTI exists in the database
	_, err = s.tokenRepo.FindByID(jti)
	if err != nil {
		// ALARM: The token is structurally valid but NOT in DB!
		// This likely means it was already used (Token Reuse) or forged.
		// Security action: Revoke ALL active sessions for this user.
		_ = s.tokenRepo.DeleteAllByUserID(userID)
		return "", "", uuid.Nil, errors.New("security breach detected: token reuse. All sessions revoked")
	}

	// Consume the old Refresh Token (Rotation)
	_ = s.tokenRepo.DeleteByID(jti)

	// Issue a new token pair
	newJTI := uuid.New()
	newAccessToken, newRefreshToken, err := jwt.GenerateTokenPair(userID, newJTI)
	if err != nil {
		return "", "", uuid.Nil, errors.New("failed to generate new tokens")
	}

	// Save new Refresh Token in DB
	rtRecord := &model.RefreshToken{
		ID:        newJTI,
		UserID:    userID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := s.tokenRepo.Create(rtRecord); err != nil {
		return "", "", uuid.Nil, errors.New("failed to save new session")
	}

	return newAccessToken, newRefreshToken, userID, nil
}

// Logout consumes a refresh token to end the session
func (s *AuthService) Logout(tokenString string) error {
	_, jti, err := jwt.ValidateRefreshToken(tokenString)
	if err != nil {
		return errors.New("invalid refresh token")
	}

	// Delete from DB regardless, preventing further use
	return s.tokenRepo.DeleteByID(jti)
}
