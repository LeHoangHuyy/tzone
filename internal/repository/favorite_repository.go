package repository

import (
	"github.com/LuuDinhTheTai/tzone/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FavoriteRepository struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) *FavoriteRepository {
	return &FavoriteRepository{db: db}
}

func (r *FavoriteRepository) ListDeviceIDsByUserID(userID uuid.UUID) ([]string, error) {
	var favorites []model.Favorite
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&favorites).Error; err != nil {
		return nil, err
	}

	deviceIDs := make([]string, 0, len(favorites))
	for _, favorite := range favorites {
		deviceIDs = append(deviceIDs, favorite.DeviceID)
	}

	return deviceIDs, nil
}

func (r *FavoriteRepository) ListExistingDeviceIDs(userID uuid.UUID, deviceIDs []string) (map[string]struct{}, error) {
	existing := make([]string, 0)
	if len(deviceIDs) == 0 {
		return map[string]struct{}{}, nil
	}

	if err := r.db.Model(&model.Favorite{}).
		Where("user_id = ? AND device_id IN ?", userID, deviceIDs).
		Pluck("device_id", &existing).Error; err != nil {
		return nil, err
	}

	set := make(map[string]struct{}, len(existing))
	for _, id := range existing {
		set[id] = struct{}{}
	}

	return set, nil
}

func (r *FavoriteRepository) Add(userID uuid.UUID, deviceID string) error {
	favorite := model.Favorite{
		UserID:   userID,
		DeviceID: deviceID,
	}

	return r.db.Where("user_id = ? AND device_id = ?", userID, deviceID).
		FirstOrCreate(&favorite).Error
}

func (r *FavoriteRepository) AddBulk(userID uuid.UUID, deviceIDs []string) error {
	if len(deviceIDs) == 0 {
		return nil
	}

	favorites := make([]model.Favorite, 0, len(deviceIDs))
	for _, deviceID := range deviceIDs {
		favorites = append(favorites, model.Favorite{UserID: userID, DeviceID: deviceID})
	}

	return r.db.Create(&favorites).Error
}

func (r *FavoriteRepository) Remove(userID uuid.UUID, deviceID string) error {
	return r.db.Where("user_id = ? AND device_id = ?", userID, deviceID).Delete(&model.Favorite{}).Error
}
