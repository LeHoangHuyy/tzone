package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// GenerateTokenPair generates an access token (30m) and a refresh token (7 days) with a unique JTI
func GenerateTokenPair(userID uuid.UUID, jti uuid.UUID) (string, string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "default_secret_key"
	}

	// 1. Generate Access Token (30 mins)
	atClaims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(30 * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
	}
	atToken := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessToken, err := atToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}

	// 2. Generate Refresh Token (7 days)
	rtClaims := jwt.MapClaims{
		"user_id": userID.String(),
		"jti":     jti.String(),
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	rtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err := rtToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ValidateToken parses and validates an Access JWT token
func ValidateToken(tokenString string) (uuid.UUID, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "default_secret_key"
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return uuid.Nil, ErrExpiredToken
		}
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return uuid.Nil, ErrInvalidToken
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return uuid.Nil, ErrInvalidToken
		}

		return userID, nil
	}

	return uuid.Nil, ErrInvalidToken
}

// ValidateRefreshToken parses and validates a Refresh Token and returns UserID and JTI
func ValidateRefreshToken(tokenString string) (uuid.UUID, uuid.UUID, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "default_secret_key"
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return uuid.Nil, uuid.Nil, ErrExpiredToken
		}
		return uuid.Nil, uuid.Nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, ok1 := claims["user_id"].(string)
		jtiStr, ok2 := claims["jti"].(string)

		if !ok1 || !ok2 {
			return uuid.Nil, uuid.Nil, ErrInvalidToken
		}

		userID, err1 := uuid.Parse(userIDStr)
		jti, err2 := uuid.Parse(jtiStr)
		if err1 != nil || err2 != nil {
			return uuid.Nil, uuid.Nil, ErrInvalidToken
		}

		return userID, jti, nil
	}

	return uuid.Nil, uuid.Nil, ErrInvalidToken
}
