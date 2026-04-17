package service

import (
	"context"
	"fmt"
	"log"

	"github.com/LuuDinhTheTai/tzone/internal/dto"
	"github.com/LuuDinhTheTai/tzone/internal/model"
	"github.com/LuuDinhTheTai/tzone/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type DeviceService struct {
	mongoDbRepo *repository.BrandRepository
}

func (s *DeviceService) addDeviceToBrand(ctx context.Context, brandID string, device *model.Device) error {
	objID, err := bson.ObjectIDFromHex(brandID)
	if err != nil {
		return fmt.Errorf("invalid brand ID format: %w", err)
	}

	if err := s.mongoDbRepo.AddDeviceToBrand(ctx, objID, device); err != nil {
		return fmt.Errorf("failed to add device to brand: %w", err)
	}

	return nil
}

func (s *DeviceService) updateDeviceInBrand(ctx context.Context, brandID string, device *model.Device) error {
	objID, err := bson.ObjectIDFromHex(brandID)
	if err != nil {
		return fmt.Errorf("invalid brand ID format: %w", err)
	}

	if err := s.mongoDbRepo.UpdateDeviceInBrand(ctx, objID, device); err != nil {
		return fmt.Errorf("failed to update device in brand doc: %w", err)
	}

	return nil
}

func (s *DeviceService) removeDeviceFromBrand(ctx context.Context, brandID string, deviceID string) error {
	objBrandID, err := bson.ObjectIDFromHex(brandID)
	if err != nil {
		return fmt.Errorf("invalid brand ID format: %w", err)
	}

	objDeviceID, err := bson.ObjectIDFromHex(deviceID)
	if err != nil {
		return fmt.Errorf("invalid device ID format: %w", err)
	}

	if err := s.mongoDbRepo.RemoveDeviceFromBrand(ctx, objBrandID, objDeviceID); err != nil {
		return fmt.Errorf("failed to remove device from brand: %w", err)
	}

	return nil
}

func NewDeviceService(mongoDbRepo *repository.BrandRepository) *DeviceService {
	return &DeviceService{
		mongoDbRepo: mongoDbRepo,
	}
}

// CreateDevice creates a new device
func (s *DeviceService) CreateDevice(ctx context.Context, req dto.CreateDeviceRequest) (*dto.DeviceResponse, error) {
	log.Printf("🔄 Creating device: %s", req.ModelName)

	device := &model.Device{
		ID:             bson.NewObjectID(),
		ModelName:      req.ModelName,
		ImageUrl:       req.ImageUrl,
		Specifications: req.Specifications,
	}

	if err := s.addDeviceToBrand(ctx, req.BrandID, device); err != nil {
		return nil, err
	}

	response := &dto.DeviceResponse{
		ID:             device.ID.Hex(),
		BrandID:        req.BrandID,
		ModelName:      device.ModelName,
		ImageUrl:       device.ImageUrl,
		Specifications: device.Specifications,
	}

	log.Printf("✅ Device created successfully: %s", response.ModelName)
	return response, nil
}

// GetDeviceById retrieves a device by ID
func (s *DeviceService) GetDeviceById(ctx context.Context, id string) (*dto.DeviceResponse, error) {
	log.Printf("🔄 Fetching device with ID: %s", id)

	device, brandIDHex, err := s.mongoDbRepo.GetDeviceById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	response := &dto.DeviceResponse{
		ID:             device.ID.Hex(),
		BrandID:        brandIDHex,
		ModelName:      device.ModelName,
		ImageUrl:       device.ImageUrl,
		Specifications: device.Specifications,
	}

	return response, nil
}

// GetAllDevices retrieves paginated devices
func (s *DeviceService) GetAllDevices(ctx context.Context, page int, limit int) (*dto.DeviceListResponse, error) {
	log.Printf("🔄 Fetching devices (page=%d, limit=%d)", page, limit)

	devices, total, err := s.mongoDbRepo.GetAllDevices(ctx, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get all devices: %w", err)
	}

	var deviceResponses []dto.DeviceResponse
	for _, device := range devices {
		deviceResponses = append(deviceResponses, dto.DeviceResponse{
			ID:             device.Device.ID.Hex(),
			BrandID:        device.BrandID.Hex(),
			ModelName:      device.Device.ModelName,
			ImageUrl:       device.Device.ImageUrl,
			Specifications: device.Device.Specifications,
		})
	}

	response := &dto.DeviceListResponse{
		Devices:    deviceResponses,
		Total:      int(total),
		Pagination: buildPaginationMeta(total, page, limit),
	}

	log.Printf("✅ Retrieved %d devices", response.Total)
	return response, nil
}

// GetDevicesByBrandId retrieves paginated devices for a specific brand
func (s *DeviceService) GetDevicesByBrandId(ctx context.Context, brandID string, page int, limit int) (*dto.DeviceListResponse, error) {
	log.Printf("🔄 Fetching devices by brand ID: %s (page=%d, limit=%d)", brandID, page, limit)

	objBrandID, err := bson.ObjectIDFromHex(brandID)
	if err != nil {
		return nil, fmt.Errorf("invalid brand ID format: %w", err)
	}

	devices, total, err := s.mongoDbRepo.GetDevicesByBrandID(ctx, objBrandID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get devices by brand ID: %w", err)
	}

	var deviceResponses []dto.DeviceResponse
	for _, device := range devices {
		deviceResponses = append(deviceResponses, dto.DeviceResponse{
			ID:             device.Device.ID.Hex(),
			BrandID:        device.BrandID.Hex(),
			ModelName:      device.Device.ModelName,
			ImageUrl:       device.Device.ImageUrl,
			Specifications: device.Device.Specifications,
		})
	}

	response := &dto.DeviceListResponse{
		Devices:    deviceResponses,
		Total:      int(total),
		Pagination: buildPaginationMeta(total, page, limit),
	}

	log.Printf("✅ Retrieved %d devices for brand %s", response.Total, brandID)
	return response, nil
}

// UpdateDevice updates existing device and handles brand changing
func (s *DeviceService) UpdateDevice(ctx context.Context, id string, req dto.UpdateDeviceRequest) (*dto.DeviceResponse, error) {
	log.Printf("🔄 Updating device with ID: %s", id)

	// Get old information to check brandID
	oldDevice, oldBrandIDHex, err := s.mongoDbRepo.GetDeviceById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("device not found: %w", err)
	}

	imageUrl := oldDevice.ImageUrl
	if req.ImageUrl != "" {
		imageUrl = req.ImageUrl
	}

	deviceToUpdate := &model.Device{
		ID:             oldDevice.ID,
		ModelName:      req.ModelName,
		ImageUrl:       imageUrl,
		Specifications: req.Specifications,
	}

	// Process array data inside brand
	if oldBrandIDHex != req.BrandID {
		// Remove device from old brand
		if err := s.removeDeviceFromBrand(ctx, oldBrandIDHex, oldDevice.ID.Hex()); err != nil {
			return nil, fmt.Errorf("failed to remove device from old brand: %w", err)
		}
		// Add a device to new brand
		if err := s.addDeviceToBrand(ctx, req.BrandID, deviceToUpdate); err != nil {
			return nil, fmt.Errorf("failed to add device to new brand: %w", err)
		}
	} else {
		// Only update information in the current brand
		if err := s.updateDeviceInBrand(ctx, req.BrandID, deviceToUpdate); err != nil {
			return nil, err
		}
	}

	response := &dto.DeviceResponse{
		ID:             deviceToUpdate.ID.Hex(),
		BrandID:        req.BrandID,
		ModelName:      deviceToUpdate.ModelName,
		ImageUrl:       deviceToUpdate.ImageUrl,
		Specifications: deviceToUpdate.Specifications,
	}

	log.Printf("✅ Device updated successfully: %s", response.ModelName)
	return response, nil
}

// DeleteDevice deletes a device by ID
func (s *DeviceService) DeleteDevice(ctx context.Context, id string) error {
	log.Printf("🔄 Deleting device with ID: %s", id)

	// Get the device before erasing to get the brandID
	device, brandIDHex, err := s.mongoDbRepo.GetDeviceById(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find device before deleting: %w", err)
	}

	// Remove device from brand's array
	if err := s.removeDeviceFromBrand(ctx, brandIDHex, device.ID.Hex()); err != nil {
		return fmt.Errorf("failed to remove deleted device from brand: %w", err)
	}

	log.Printf("✅ Device deleted successfully")
	return nil
}
