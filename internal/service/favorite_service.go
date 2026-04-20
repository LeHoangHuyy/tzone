package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/LuuDinhTheTai/tzone/internal/dto"
	"github.com/LuuDinhTheTai/tzone/internal/repository"
	"github.com/google/uuid"
)

type FavoriteService struct {
	favoriteRepo *repository.FavoriteRepository
	deviceRepo   *repository.BrandRepository
}

func NewFavoriteService(favoriteRepo *repository.FavoriteRepository, deviceRepo *repository.BrandRepository) *FavoriteService {
	return &FavoriteService{
		favoriteRepo: favoriteRepo,
		deviceRepo:   deviceRepo,
	}
}

func parseUserID(userID string) (uuid.UUID, error) {
	parsed, err := uuid.Parse(strings.TrimSpace(userID))
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user id")
	}
	return parsed, nil
}

func normalizeDeviceIDs(deviceIDs []string) []string {
	seen := make(map[string]struct{}, len(deviceIDs))
	normalized := make([]string, 0, len(deviceIDs))
	for _, id := range deviceIDs {
		trimmed := strings.TrimSpace(id)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		normalized = append(normalized, trimmed)
	}
	return normalized
}

func (s *FavoriteService) GetFavorites(userID string) (*dto.FavoriteListResponse, error) {
	parsedUserID, err := parseUserID(userID)
	if err != nil {
		return nil, err
	}

	deviceIDs, err := s.favoriteRepo.ListDeviceIDsByUserID(parsedUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get favorites: %w", err)
	}

	return &dto.FavoriteListResponse{DeviceIDs: deviceIDs}, nil
}

func (s *FavoriteService) AddFavorite(ctx context.Context, userID string, deviceID string) (*dto.FavoriteListResponse, error) {
	parsedUserID, err := parseUserID(userID)
	if err != nil {
		return nil, err
	}

	deviceID = strings.TrimSpace(deviceID)
	if deviceID == "" {
		return nil, fmt.Errorf("device_id is required")
	}

	if _, _, err := s.deviceRepo.GetDeviceById(ctx, deviceID); err != nil {
		return nil, fmt.Errorf("device not found")
	}

	if err := s.favoriteRepo.Add(parsedUserID, deviceID); err != nil {
		return nil, fmt.Errorf("failed to add favorite: %w", err)
	}

	return s.GetFavorites(userID)
}

func (s *FavoriteService) RemoveFavorite(userID string, deviceID string) (*dto.FavoriteListResponse, error) {
	parsedUserID, err := parseUserID(userID)
	if err != nil {
		return nil, err
	}

	deviceID = strings.TrimSpace(deviceID)
	if deviceID == "" {
		return nil, fmt.Errorf("device_id is required")
	}

	if err := s.favoriteRepo.Remove(parsedUserID, deviceID); err != nil {
		return nil, fmt.Errorf("failed to remove favorite: %w", err)
	}

	return s.GetFavorites(userID)
}

func (s *FavoriteService) SyncFavorites(ctx context.Context, userID string, deviceIDs []string) (*dto.FavoriteListResponse, error) {
	parsedUserID, err := parseUserID(userID)
	if err != nil {
		return nil, err
	}

	normalizedIDs := normalizeDeviceIDs(deviceIDs)
	if len(normalizedIDs) == 0 {
		return s.GetFavorites(userID)
	}

	existingMap, err := s.favoriteRepo.ListExistingDeviceIDs(parsedUserID, normalizedIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing favorites: %w", err)
	}

	toInsert := make([]string, 0)
	for _, deviceID := range normalizedIDs {
		if _, exists := existingMap[deviceID]; exists {
			continue
		}

		if _, _, err := s.deviceRepo.GetDeviceById(ctx, deviceID); err != nil {
			continue
		}

		toInsert = append(toInsert, deviceID)
	}

	if err := s.favoriteRepo.AddBulk(parsedUserID, toInsert); err != nil {
		return nil, fmt.Errorf("failed to sync favorites: %w", err)
	}

	return s.GetFavorites(userID)
}
