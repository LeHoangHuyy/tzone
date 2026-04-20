package dto

import "github.com/LuuDinhTheTai/tzone/internal/model"

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// BrandResponse represents the response structure for a brand
type BrandResponse struct {
	Id   string `json:"id"`
	Name string `json:"brand_name"`
}

// BrandListResponse represents the response structure for a list of brands
type BrandListResponse struct {
	Brands     []BrandResponse `json:"brands"`
	Total      int             `json:"total"`
	Pagination PaginationMeta  `json:"pagination"`
}

// DeviceResponse represents the response structure for a device
type DeviceResponse struct {
	ID             string               `json:"id,omitempty"`
	BrandID        string               `json:"brand_id,omitempty"`
	ModelName      string               `json:"model_name,omitempty"`
	ImageUrl       string               `json:"imageUrl,omitempty"`
	Specifications model.Specifications `json:"specifications,omitempty"`
}

// DeviceListResponse represents the response structure for a list of devices
type DeviceListResponse struct {
	Devices    []DeviceResponse `json:"devices"`
	Total      int              `json:"total"`
	Pagination PaginationMeta   `json:"pagination"`
}

type FavoriteListResponse struct {
	DeviceIDs []string `json:"device_ids"`
}
