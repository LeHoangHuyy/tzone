package dto

import (
	"mime/multipart"
	"strings"

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

type DeviceFinderQuery struct {
	Name        string `form:"name"`
	BrandID     string `form:"brand_id"`
	OS          string `form:"os"`
	Chipset     string `form:"chipset"`
	CPU         string `form:"cpu"`
	GPU         string `form:"gpu"`
	Memory      string `form:"memory"`
	DisplaySize string `form:"display_size"`
	Battery     string `form:"battery"`
	NFC         string `form:"nfc"`
	PaginationQuery
}

func (q *DeviceFinderQuery) Normalize() {
	q.PaginationQuery.Normalize()
	q.Name = strings.TrimSpace(q.Name)
	q.BrandID = strings.TrimSpace(q.BrandID)
	q.OS = strings.TrimSpace(q.OS)
	q.Chipset = strings.TrimSpace(q.Chipset)
	q.CPU = strings.TrimSpace(q.CPU)
	q.GPU = strings.TrimSpace(q.GPU)
	q.Memory = strings.TrimSpace(q.Memory)
	q.DisplaySize = strings.TrimSpace(q.DisplaySize)
	q.Battery = strings.TrimSpace(q.Battery)
	q.NFC = strings.TrimSpace(q.NFC)
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

type AddFavoriteRequest struct {
	DeviceID string `json:"device_id" binding:"required"`
}

type SyncFavoritesRequest struct {
	DeviceIDs []string `json:"device_ids" binding:"required"`
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
