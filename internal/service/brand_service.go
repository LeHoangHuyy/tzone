package service

import (
	"context"
	"fmt"
	"log"

	"github.com/LuuDinhTheTai/tzone/internal/dto"
	"github.com/LuuDinhTheTai/tzone/internal/model"
	"github.com/LuuDinhTheTai/tzone/internal/repository"
)

type BrandService struct {
	mongoDbRepo *repository.BrandRepository
}

func NewBrandService(mongoDbRepo *repository.BrandRepository) *BrandService {
	return &BrandService{
		mongoDbRepo: mongoDbRepo,
	}
}

func buildPaginationMeta(total int64, page int, limit int) dto.PaginationMeta {
	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	return dto.PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    int64(page*limit) < total,
		HasPrev:    page > 1,
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

// GetAllBrands retrieves paginated brands
func (s *BrandService) GetAllBrands(ctx context.Context, page int, limit int) (*dto.BrandListResponse, error) {
	log.Printf("🔄 Fetching brands (page=%d, limit=%d)", page, limit)

	brands, total, err := s.mongoDbRepo.GetAllBrands(ctx, page, limit)
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
		Brands:     brandResponses,
		Total:      int(total),
		Pagination: buildPaginationMeta(total, page, limit),
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
