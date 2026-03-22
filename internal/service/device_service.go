package service

import (
	"context"
	"fmt"
	"log"

	"github.com/LuuDinhTheTai/tzone/internal/dto"
	"github.com/LuuDinhTheTai/tzone/internal/model"
	"github.com/LuuDinhTheTai/tzone/internal/repository"
)

type DeviceService struct {
	mongoDbRepo *repository.DeviceRepository
}

func NewDeviceService(mongoDbRepo *repository.DeviceRepository) *DeviceService {
	return &DeviceService{
		mongoDbRepo: mongoDbRepo,
	}
}

// CreateDevice creates a new device
func (s *DeviceService) CreateDevice(ctx context.Context, req dto.CreateDeviceRequest) (*dto.DeviceResponse, error) {
	log.Printf("🔄 Creating device: %s", req.ModelName)

	device := &model.Device{
		ModelName:      req.ModelName,
		ImageUrl:       req.ImageUrl,
		Specifications: req.Specifications,
	}

	createdDevice, err := s.mongoDbRepo.CreateDevice(ctx, device)
	if err != nil {
		return nil, fmt.Errorf("failed to create device: %w", err)
	}

	response := &dto.DeviceResponse{
		ID:             createdDevice.ID.Hex(),
		ModelName:      createdDevice.ModelName,
		ImageUrl:       createdDevice.ImageUrl,
		Specifications: createdDevice.Specifications,
	}

	log.Printf("✅ Device created successfully: %s", response.ModelName)
	return response, nil
}

// GetDeviceById retrieves a device by ID
func (s *DeviceService) GetDeviceById(ctx context.Context, id string) (*dto.DeviceResponse, error) {
	log.Printf("🔄 Fetching device with ID: %s", id)

	device, err := s.mongoDbRepo.GetDeviceById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	response := &dto.DeviceResponse{
		ID:             device.ID.Hex(),
		ModelName:      device.ModelName,
		ImageUrl:       device.ImageUrl,
		Specifications: device.Specifications,
	}

	return response, nil
}

// GetAllDevices retrieves all devices
func (s *DeviceService) GetAllDevices(ctx context.Context) (*dto.DeviceListResponse, error) {
	log.Printf("🔄 Fetching all devices")

	devices, err := s.mongoDbRepo.GetAllDevices(ctx)
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

// UpdateDevice updates existing device
func (s *DeviceService) UpdateDevice(ctx context.Context, id string, req dto.UpdateDeviceRequest) (*dto.DeviceResponse, error) {
	log.Printf("🔄 Updating device with ID: %s", id)

	device := &model.Device{
		ModelName:      req.ModelName,
		Specifications: req.Specifications,
	}

	updatedBrand, err := s.mongoDbRepo.UpdateDevice(ctx, id, device)
	if err != nil {
		return nil, fmt.Errorf("failed to update device: %w", err)
	}

	response := &dto.DeviceResponse{
		ID:             updatedBrand.ID.Hex(),
		ModelName:      updatedBrand.ModelName,
		ImageUrl:       updatedBrand.ImageUrl,
		Specifications: updatedBrand.Specifications,
	}

	log.Printf("✅ Brand updated successfully: %s", response.ModelName)
	return response, nil
}

// DeleteDevice deletes a device by ID
func (s *DeviceService) DeleteDevice(ctx context.Context, id string) error {
	log.Printf("🔄 Deleting device with ID: %s", id)

	err := s.mongoDbRepo.DeleteDevice(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete device: %w", err)
	}

	log.Printf("✅ Device deleted successfully")
	return nil
}
