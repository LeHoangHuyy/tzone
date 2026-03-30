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

type BrandService struct {
	mongoDbRepo *repository.BrandRepository
}

func NewBrandService(mongoDbRepo *repository.BrandRepository) *BrandService {
	return &BrandService{
		mongoDbRepo: mongoDbRepo,
	}
}

// CreateBrand creates a new brand
func (s *BrandService) CreateBrand(ctx context.Context, req dto.CreateBrandRequest) (*dto.BrandResponse, error) {
	log.Printf("🔄 Creating brand: %s", req.Name)

	brand := &model.Brand{
		Name: req.Name,
	}

	createdBrand, err := s.mongoDbRepo.CreateBrand(ctx, brand)
	if err != nil {
		return nil, fmt.Errorf("failed to create brand: %w", err)
	}

	response := &dto.BrandResponse{
		Id:   createdBrand.Id.Hex(),
		Name: createdBrand.Name,
	}

	log.Printf("✅ Brand created successfully: %s", response.Name)
	return response, nil
}

// GetBrandById retrieves a brand by ID
func (s *BrandService) GetBrandById(ctx context.Context, id string) (*dto.BrandResponse, error) {
	log.Printf("🔄 Fetching brand with ID: %s", id)

	brand, err := s.mongoDbRepo.GetBrandById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get brand: %w", err)
	}

	response := &dto.BrandResponse{
		Id:   brand.Id.Hex(),
		Name: brand.Name,
	}

	return response, nil
}

// GetAllBrands retrieves all brands
func (s *BrandService) GetAllBrands(ctx context.Context) (*dto.BrandListResponse, error) {
	log.Printf("🔄 Fetching all brands")

	brands, err := s.mongoDbRepo.GetAllBrands(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all brands: %w", err)
	}

	var brandResponses []dto.BrandResponse
	for _, brand := range brands {
		brandResponses = append(brandResponses, dto.BrandResponse{
			Id:   brand.Id.Hex(),
			Name: brand.Name,
		})
	}

	response := &dto.BrandListResponse{
		Brands: brandResponses,
		Total:  len(brandResponses),
	}

	log.Printf("✅ Retrieved %d brands", response.Total)
	return response, nil
}

// UpdateBrand updates an existing brand
func (s *BrandService) UpdateBrand(ctx context.Context, id string, req dto.UpdateBrandRequest) (*dto.BrandResponse, error) {
	log.Printf("🔄 Updating brand with ID: %s", id)

	brand := &model.Brand{
		Name: req.Name,
	}

	updatedBrand, err := s.mongoDbRepo.UpdateBrand(ctx, id, brand)
	if err != nil {
		return nil, fmt.Errorf("failed to update brand: %w", err)
	}

	response := &dto.BrandResponse{
		Id:   updatedBrand.Id.Hex(),
		Name: updatedBrand.Name,
	}

	log.Printf("✅ Brand updated successfully: %s", response.Name)
	return response, nil
}

// DeleteBrand deletes a brand by ID
func (s *BrandService) DeleteBrand(ctx context.Context, id string) error {
	log.Printf("🔄 Deleting brand with ID: %s", id)

	//  Get brand information before to check
	brand, err := s.mongoDbRepo.GetBrandById(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find brand to delete: %w", err)
	}

	// Don't delete brand if still contains any device
	if len(brand.Devices) > 0 {
		return fmt.Errorf("cannot delete brand because it still contains %d device(s)", len(brand.Devices))
	}

	err = s.mongoDbRepo.DeleteBrand(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete brand: %w", err)
	}

	log.Printf("✅ Brand deleted successfully")
	return nil
}

// AddDeviceToBrand implements BrandUpdater interface
func (s *BrandService) AddDeviceToBrand(ctx context.Context, brandID string, device *model.Device) error {
	objID, err := bson.ObjectIDFromHex(brandID)
	if err != nil {
		return fmt.Errorf("invalid brand ID format: %w", err)
	}
	return s.mongoDbRepo.AddDeviceToBrand(ctx, objID, device)
}

// UpdateDeviceInBrand implements BrandUpdater interface
func (s *BrandService) UpdateDeviceInBrand(ctx context.Context, brandID string, device *model.Device) error {
	objID, err := bson.ObjectIDFromHex(brandID)
	if err != nil {
		return fmt.Errorf("invalid brand ID format: %w", err)
	}
	return s.mongoDbRepo.UpdateDeviceInBrand(ctx, objID, device)
}

// RemoveDeviceFromBrand implements BrandUpdater interface
func (s *BrandService) RemoveDeviceFromBrand(ctx context.Context, brandID string, deviceID string) error {
	objBrandID, err := bson.ObjectIDFromHex(brandID)
	if err != nil {
		return fmt.Errorf("invalid brand ID format: %w", err)
	}
	objDeviceID, err := bson.ObjectIDFromHex(deviceID)
	if err != nil {
		return fmt.Errorf("invalid device ID format: %w", err)
	}
	return s.mongoDbRepo.RemoveDeviceFromBrand(ctx, objBrandID, objDeviceID)
}
