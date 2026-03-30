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

func NewDeviceService(mongoDbRepo *repository.BrandRepository) *DeviceService {
	return &DeviceService{
		mongoDbRepo: mongoDbRepo,
	}
}

// CreateDevice creates a new device
func (s *DeviceService) CreateDevice(ctx context.Context, req dto.CreateDeviceRequest) (*dto.DeviceResponse, error) {
	log.Printf("🔄 Creating device: %s", req.ModelName)

	brandID, err := bson.ObjectIDFromHex(req.BrandID)
	if err != nil {
		return nil, fmt.Errorf("invalid brand ID format %w", err)
	}

	device := &model.Device{
		ID:             bson.NewObjectID(),
		ModelName:      req.ModelName,
		ImageUrl:       req.ImageUrl,
		Specifications: req.Specifications,
	}

	err = s.mongoDbRepo.AddDeviceToBrand(ctx, brandID, device)
	if err != nil {
		return nil, fmt.Errorf("failed to add device to brand: %w", err)
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

// GetAllDevices retrieves all devices
func (s *DeviceService) GetAllDevices(ctx context.Context) (*dto.DeviceListResponse, error) {
	log.Printf("🔄 Fetching all devices")

	devices, _, err := s.mongoDbRepo.GetAllDevices(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all devices: %w", err)
	}

	var deviceResponses []dto.DeviceResponse
	for _, device := range devices {
		deviceResponses = append(deviceResponses, dto.DeviceResponse{
			ID:             device.ID.Hex(),
			ModelName:      device.ModelName,
			ImageUrl:       device.ImageUrl,
			Specifications: device.Specifications,
		})
	}

	response := &dto.DeviceListResponse{
		Devices: deviceResponses,
		Total:   len(deviceResponses),
	}

	log.Printf("✅ Retrieved %d devices", response.Total)
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

	newBrandID, err := bson.ObjectIDFromHex(req.BrandID)
	if err != nil {
		return nil, fmt.Errorf("invalid brand ID format %w", err)
	}

	deviceToUpdate := &model.Device{
		ID:             oldDevice.ID,
		ModelName:      req.ModelName,
		Specifications: req.Specifications,
	}

	// Process array data inside brand
	if oldBrandIDHex != req.BrandID {
		oldBrandObjID, _ := bson.ObjectIDFromHex(oldBrandIDHex)

		// Remove device from old brand
		if err := s.mongoDbRepo.RemoveDeviceFromBrand(ctx, oldBrandObjID, oldDevice.ID); err != nil {
			return nil, fmt.Errorf("failed to remove device from old brand: %w", err)
		}
		// Add a device to new brand
		if err := s.mongoDbRepo.AddDeviceToBrand(ctx, newBrandID, deviceToUpdate); err != nil {
			return nil, fmt.Errorf("failed to add device to new brand: %w", err)
		}
	} else {
		// Only update information in the current brand
		if err := s.mongoDbRepo.UpdateDeviceInBrand(ctx, newBrandID, deviceToUpdate); err != nil {
			return nil, fmt.Errorf("failed to update device in brand doc: %w", err)
		}
	}

	response := &dto.DeviceResponse{
		ID:             deviceToUpdate.ID.Hex(),
		BrandID:        req.BrandID,
		ModelName:      deviceToUpdate.ModelName,
		ImageUrl:       deviceToUpdate.ImageUrl,
		Specifications: deviceToUpdate.Specifications,
	}

	log.Printf("✅ Brand updated successfully: %s", response.ModelName)
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

	brandObjID, _ := bson.ObjectIDFromHex(brandIDHex)

	// Remove device from brand's array
	err = s.mongoDbRepo.RemoveDeviceFromBrand(ctx, brandObjID, device.ID)
	if err != nil {
		return fmt.Errorf("failed to remove deleted device from brand: %w", err)
	}

	log.Printf("✅ Device deleted successfully")
	return nil
}
