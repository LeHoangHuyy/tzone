package handler

import (
	"log"
	"net/http"

	"github.com/LuuDinhTheTai/tzone/internal/dto"
	"github.com/LuuDinhTheTai/tzone/internal/service"
	"github.com/LuuDinhTheTai/tzone/util/response"
	"github.com/gin-gonic/gin"
)

type BrandHandler struct {
	brandService *service.BrandService
}

func NewBrandHandler(brandService *service.BrandService) *BrandHandler {
	return &BrandHandler{
		brandService: brandService,
	}
}

// CreateBrand handles POST request to create a new brand
// @Summary Create a new brand
// @Description Create a new brand with the provided name
// @Tags brands
// @Accept json
// @Produce json
// @Param brand body dto.CreateBrandRequest true "Brand information"
// @Success 201 {object} response.ApiResponse{data=dto.BrandResponse}
// @Failure 400 {object} response.ApiResponse
// @Failure 500 {object} response.ApiResponse
// @Router /api/v1/brands [post]
func (h *BrandHandler) CreateBrand(ctx *gin.Context) {
	var req dto.CreateBrandRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ Invalid request body: %v", err)
		response.Error(ctx, http.StatusBadRequest, "Invalid request body", []response.ErrorResponse{
			{Field: "brand_name", Error: err.Error()},
		})
		return
	}

	brand, err := h.brandService.CreateBrand(ctx.Request.Context(), req)
	if err != nil {
		log.Printf("❌ Failed to create brand: %v", err)
		response.Error(ctx, http.StatusInternalServerError, "Failed to create brand", []response.ErrorResponse{
			{Field: "server", Error: err.Error()},
		})
		return
	}

	response.Success(ctx, http.StatusCreated, "Brand created successfully", brand)
}

// GetBrandById handles GET request to retrieve a brand by ID
// @Summary Get a brand by ID
// @Description Get brand details by brand ID
// @Tags brands
// @Accept json
// @Produce json
// @Param id path string true "Brand ID"
// @Success 200 {object} response.ApiResponse{data=dto.BrandResponse}
// @Failure 400 {object} response.ApiResponse
// @Failure 404 {object} response.ApiResponse
// @Failure 500 {object} response.ApiResponse
// @Router /api/v1/brands/{id} [get]
func (h *BrandHandler) GetBrandById(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		log.Printf("❌ Brand ID is required")
		response.Error(ctx, http.StatusBadRequest, "Brand ID is required", []response.ErrorResponse{
			{Field: "id", Error: "id parameter is missing"},
		})
		return
	}

	brand, err := h.brandService.GetBrandById(ctx.Request.Context(), id)
	if err != nil {
		log.Printf("❌ Failed to get brand: %v", err)
		response.Error(ctx, http.StatusNotFound, "Brand not found", []response.ErrorResponse{
			{Field: "id", Error: err.Error()},
		})
		return
	}

	response.Success(ctx, http.StatusOK, "Brand retrieved successfully", brand)
}

// GetAllBrands handles GET request to retrieve all brands
// @Summary Get all brands
// @Description Get a list of all brands
// @Tags brands
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponse{data=dto.BrandListResponse}
// @Failure 500 {object} response.ApiResponse
// @Router /api/v1/brands [get]
func (h *BrandHandler) GetAllBrands(ctx *gin.Context) {
	brands, err := h.brandService.GetAllBrands(ctx.Request.Context())
	if err != nil {
		log.Printf("❌ Failed to get all brands: %v", err)
		response.Error(ctx, http.StatusInternalServerError, "Failed to retrieve brands", []response.ErrorResponse{
			{Field: "server", Error: err.Error()},
		})
		return
	}

	response.Success(ctx, http.StatusOK, "Brands retrieved successfully", brands)
}

// UpdateBrand handles PUT request to update a brand
// @Summary Update a brand
// @Description Update brand information by ID
// @Tags brands
// @Accept json
// @Produce json
// @Param id path string true "Brand ID"
// @Param brand body dto.UpdateBrandRequest true "Updated brand information"
// @Success 200 {object} response.ApiResponse{data=dto.BrandResponse}
// @Failure 400 {object} response.ApiResponse
// @Failure 404 {object} response.ApiResponse
// @Failure 500 {object} response.ApiResponse
// @Router /api/v1/brands/{id} [put]
func (h *BrandHandler) UpdateBrand(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		log.Printf("❌ Brand ID is required")
		response.Error(ctx, http.StatusBadRequest, "Brand ID is required", []response.ErrorResponse{
			{Field: "id", Error: "id parameter is missing"},
		})
		return
	}

	var req dto.UpdateBrandRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ Invalid request body: %v", err)
		response.Error(ctx, http.StatusBadRequest, "Invalid request body", []response.ErrorResponse{
			{Field: "brand_name", Error: err.Error()},
		})
		return
	}

	brand, err := h.brandService.UpdateBrand(ctx.Request.Context(), id, req)
	if err != nil {
		log.Printf("❌ Failed to update brand: %v", err)
		response.Error(ctx, http.StatusInternalServerError, "Failed to update brand", []response.ErrorResponse{
			{Field: "server", Error: err.Error()},
		})
		return
	}

	response.Success(ctx, http.StatusOK, "Brand updated successfully", brand)
}

// DeleteBrand handles DELETE request to delete a brand
// @Summary Delete a brand
// @Description Delete a brand by ID
// @Tags brands
// @Accept json
// @Produce json
// @Param id path string true "Brand ID"
// @Success 200 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Failure 404 {object} response.ApiResponse
// @Failure 500 {object} response.ApiResponse
// @Router /api/v1/brands/{id} [delete]
func (h *BrandHandler) DeleteBrand(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		log.Printf("❌ Brand ID is required")
		response.Error(ctx, http.StatusBadRequest, "Brand ID is required", []response.ErrorResponse{
			{Field: "id", Error: "id parameter is missing"},
		})
		return
	}

	err := h.brandService.DeleteBrand(ctx.Request.Context(), id)
	if err != nil {
		log.Printf("❌ Failed to delete brand: %v", err)
		response.Error(ctx, http.StatusInternalServerError, "Failed to delete brand", []response.ErrorResponse{
			{Field: "server", Error: err.Error()},
		})
		return
	}

	response.Success(ctx, http.StatusOK, "Brand deleted successfully", nil)
}
