package handler

import (
	"log"
	"net/http"

	"github.com/LuuDinhTheTai/tzone/internal/dto"
	"github.com/LuuDinhTheTai/tzone/internal/service"
	"github.com/LuuDinhTheTai/tzone/util/response"
	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	deviceService *service.DeviceService
}

func NewDeviceHandler(deviceService *service.DeviceService) *DeviceHandler {
	return &DeviceHandler{
		deviceService: deviceService,
	}
}

// CreateDevice  handles POST request to create a new device
// @Summary Create a new device
// @Description Create a new device with the provided name
// @Tags devices
// @Accept json
// @Produce json
// @Param device body dto.CreateDeviceRequest true "Device information"
// @Success 201 {object} response.ApiResponse{data=dto.DeviceResponse}
// @Failure 400 {object} response.ApiResponse
// @Failure 500 {object} response.ApiResponse
// @Router /api/v1/devices [post]
func (h *DeviceHandler) CreateDevice(ctx *gin.Context) {
	var req dto.CreateDeviceRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ Invalid request body: %v", err)
		response.Error(ctx, http.StatusBadRequest, "Invalid request body", []response.ErrorResponse{
			{Field: "", Error: err.Error()},
		})
		return
	}

	device, err := h.deviceService.CreateDevice(ctx.Request.Context(), req)
	if err != nil {
		log.Printf("❌ Failed to create device: %v", err)
		response.Error(ctx, http.StatusInternalServerError, "Failed to create device", []response.ErrorResponse{
			{Field: "server", Error: err.Error()},
		})
		return
	}

	response.Success(ctx, http.StatusCreated, "Device created successfully", device)
}

// GetDeviceById handles GET request to retrieve a device by ID
// @Summary Get a device by ID
// @Description Get device details by device ID
// @Tags devices
// @Accept json
// @Produce json
// @Param id path string true "Device ID"
// @Success 200 {object} response.ApiResponse{data=dto.DeviceResponse}
// @Failure 400 {object} response.ApiResponse
// @Failure 404 {object} response.ApiResponse
// @Failure 500 {object} response.ApiResponse
// @Router /api/v1/devices/{id} [get]
func (h *DeviceHandler) GetDeviceById(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		log.Printf("❌ Device ID is required")
		response.Error(ctx, http.StatusBadRequest, "Device ID is required", []response.ErrorResponse{
			{Field: "id", Error: "id parameter is missing"},
		})
		return
	}

	device, err := h.deviceService.GetDeviceById(ctx.Request.Context(), id)
	if err != nil {
		log.Printf("❌ Failed to get device: %v", err)
		response.Error(ctx, http.StatusNotFound, "Device not found", []response.ErrorResponse{
			{Field: "id", Error: err.Error()},
		})
		return
	}

	response.Success(ctx, http.StatusOK, "Device retrieved successfully", device)
}

// GetAllDevices handles GET request to retrieve all devices
// @Summary Get all devices
// @Description Get a list of all devices
// @Tags devices
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page (max 100)" default(10)
// @Success 200 {object} response.ApiResponse{data=dto.DeviceListResponse}
// @Failure 500 {object} response.ApiResponse
// @Router /api/v1/devices [get]
func (h *DeviceHandler) GetAllDevices(ctx *gin.Context) {
	var pagination dto.PaginationQuery
	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		log.Printf("❌ Invalid pagination query: %v", err)
		response.Error(ctx, http.StatusBadRequest, "Invalid pagination query", []response.ErrorResponse{
			{Field: "query", Error: err.Error()},
		})
		return
	}

	pagination.Normalize()

	devices, err := h.deviceService.GetAllDevices(ctx.Request.Context(), pagination.Page, pagination.Limit)
	if err != nil {
		log.Printf("❌ Failed to get all devices: %v", err)
		response.Error(ctx, http.StatusInternalServerError, "Failed to retrieve devices", []response.ErrorResponse{
			{Field: "server", Error: err.Error()},
		})
		return
	}

	response.Success(ctx, http.StatusOK, "Devices retrieved successfully", devices)
}

// GetDevicesByBrandId handles GET request to retrieve devices by brand ID
// @Summary Get devices by brand ID
// @Description Get a paginated list of devices for a specific brand
// @Tags devices
// @Accept json
// @Produce json
// @Param brandId path string true "Brand ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page (max 100)" default(10)
// @Success 200 {object} response.ApiResponse{data=dto.DeviceListResponse}
// @Failure 400 {object} response.ApiResponse
// @Failure 500 {object} response.ApiResponse
// @Router /api/v1/devices/brand/{brandId} [get]
func (h *DeviceHandler) GetDevicesByBrandId(ctx *gin.Context) {
	brandID := ctx.Param("brandId")
	if brandID == "" {
		log.Printf("❌ Brand ID is required")
		response.Error(ctx, http.StatusBadRequest, "Brand ID is required", []response.ErrorResponse{
			{Field: "brandId", Error: "brandId parameter is missing"},
		})
		return
	}

	var pagination dto.PaginationQuery
	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		log.Printf("❌ Invalid pagination query: %v", err)
		response.Error(ctx, http.StatusBadRequest, "Invalid pagination query", []response.ErrorResponse{
			{Field: "query", Error: err.Error()},
		})
		return
	}

	pagination.Normalize()

	devices, err := h.deviceService.GetDevicesByBrandId(ctx.Request.Context(), brandID, pagination.Page, pagination.Limit)
	if err != nil {
		log.Printf("❌ Failed to get devices by brand ID: %v", err)
		response.Error(ctx, http.StatusInternalServerError, "Failed to retrieve devices by brand", []response.ErrorResponse{
			{Field: "server", Error: err.Error()},
		})
		return
	}

	response.Success(ctx, http.StatusOK, "Devices retrieved successfully", devices)
}

// UpdateDevice handles PUT request to update a device
// @Summary Update a device
// @Description Update device information by ID
// @Tags devices
// @Accept json
// @Produce json
// @Param id path string true "Device ID"
// @Param device body dto.UpdateDeviceRequest true "Updated device information"
// @Success 200 {object} response.ApiResponse{data=dto.DeviceResponse}
// @Failure 400 {object} response.ApiResponse
// @Failure 404 {object} response.ApiResponse
// @Failure 500 {object} response.ApiResponse
// @Router /api/v1/devices/{id} [put]
func (h *DeviceHandler) UpdateDevice(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		log.Printf("❌ Device ID is required")
		response.Error(ctx, http.StatusBadRequest, "Device ID is required", []response.ErrorResponse{
			{Field: "id", Error: "id parameter is missing"},
		})
		return
	}

	var req dto.UpdateDeviceRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ Invalid request body: %v", err)
		response.Error(ctx, http.StatusBadRequest, "Invalid request body", []response.ErrorResponse{
			{Field: "", Error: err.Error()},
		})
		return
	}

	device, err := h.deviceService.UpdateDevice(ctx.Request.Context(), id, req)
	if err != nil {
		log.Printf("❌ Failed to update device: %v", err)
		response.Error(ctx, http.StatusInternalServerError, "Failed to update device", []response.ErrorResponse{
			{Field: "server", Error: err.Error()},
		})
		return
	}

	response.Success(ctx, http.StatusOK, "Device updated successfully", device)
}

// DeleteDevice handles DELETE request to delete a device
// @Summary Delete a device
// @Description Delete a device by ID
// @Tags devices
// @Accept json
// @Produce json
// @Param id path string true "Device ID"
// @Success 200 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Failure 404 {object} response.ApiResponse
// @Failure 500 {object} response.ApiResponse
// @Router /api/v1/devices/{id} [delete]
func (h *DeviceHandler) DeleteDevice(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		log.Printf("❌ Device ID is required")
		response.Error(ctx, http.StatusBadRequest, "Device ID is required", []response.ErrorResponse{
			{Field: "id", Error: "id parameter is missing"},
		})
		return
	}

	err := h.deviceService.DeleteDevice(ctx.Request.Context(), id)
	if err != nil {
		log.Printf("❌ Failed to delete device: %v", err)
		response.Error(ctx, http.StatusInternalServerError, "Failed to delete device", []response.ErrorResponse{
			{Field: "server", Error: err.Error()},
		})
		return
	}

	response.Success(ctx, http.StatusOK, "Device deleted successfully", nil)
}
