package dto

import "github.com/LuuDinhTheTai/tzone/internal/model"

// CreateBrandRequest represents the request body for creating a new brand
type CreateBrandRequest struct {
	Name string `json:"brand_name" binding:"required,min=1,max=100"`
}

// UpdateBrandRequest represents the request body for updating a brand
type UpdateBrandRequest struct {
	Name string `json:"brand_name" binding:"required,min=1,max=100"`
}

// CreateDeviceRequest represents the request body for creating a new device
type CreateDeviceRequest struct {
	BrandID        string               `json:"brand_id" binding:"required"`
	ModelName      string               `json:"model_name" binding:"required,min=1,max=100"`
	ImageUrl       string               `json:"imageUrl" binding:"required,min=1,max=100"`
	Specifications model.Specifications `json:"specifications"`
}

// UpdateDeviceRequest represents the request body for updating a device
type UpdateDeviceRequest struct {
	BrandID        string               `json:"brand_id" binding:"required"`
	ModelName      string               `json:"model_name" binding:"required,min=1,max=100"`
	Specifications model.Specifications `json:"specifications"`
}
