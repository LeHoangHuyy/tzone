package dto

import (
	"mime/multipart"

	"github.com/LuuDinhTheTai/tzone/internal/model"
)

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
)

type PaginationQuery struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

func (q *PaginationQuery) Normalize() {
	if q.Page < 1 {
		q.Page = DefaultPage
	}

	if q.Limit < 1 {
		q.Limit = DefaultLimit
	}

	if q.Limit > MaxLimit {
		q.Limit = MaxLimit
	}
}

type SearchQuery struct {
	Name string `form:"name" binding:"required,min=1,max=100"`
	PaginationQuery
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	OTP      string `json:"otp" binding:"required,len=6"`
}

type SendOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	OTP         string `json:"otp" binding:"required,len=6"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
	OTP         string `json:"otp" binding:"required,len=6"`
}

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
	ImageUrl       string               `json:"imageUrl" binding:"required,min=1,max=512"`
	Specifications model.Specifications `json:"specifications"`
}

// UpdateDeviceRequest represents the request body for updating a device
type UpdateDeviceRequest struct {
	BrandID        string               `json:"brand_id" binding:"required"`
	ModelName      string               `json:"model_name" binding:"required,min=1,max=100"`
	ImageUrl       string               `json:"imageUrl" binding:"required,min=1,max=512"`
	Specifications model.Specifications `json:"specifications"`
}

// CreateDeviceFormRequest represents the form data for creating a new device
type CreateDeviceFormRequest struct {
	BrandID        string                `form:"brand_id" binding:"required"`
	ModelName      string                `form:"model_name" binding:"required,min=1,max=100"`
	Image          *multipart.FileHeader `form:"image" binding:"required"`
	Specifications string                `form:"specifications"`
}

// UpdateDeviceFormRequest represents the form data for updating a device
type UpdateDeviceFormRequest struct {
	BrandID        string                `form:"brand_id" binding:"required"`
	ModelName      string                `form:"model_name" binding:"required,min=1,max=100"`
	Image          *multipart.FileHeader `form:"image"`
	Specifications string                `form:"specifications"`
}
