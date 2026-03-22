package dto

import "github.com/LuuDinhTheTai/tzone/internal/model"

// CreateDeviceRequest represents the request body for creating a new device
type CreateDeviceRequest struct {
	ModelName      string               `json:"model_name" binding:"required,min=1,max=100"`
	ImageUrl       string               `json:"imageUrl" binding:"required,min=1,max=100"`
	Specifications model.Specifications `json:"specifications"`
}

// UpdateDeviceRequest represents the request body for updating a device
type UpdateDeviceRequest struct {
	ModelName      string               `json:"model_name" binding:"required,min=1,max=100"`
	Specifications model.Specifications `json:"specifications"`
}
