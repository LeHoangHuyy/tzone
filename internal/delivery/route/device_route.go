package route

import (
	"github.com/LuuDinhTheTai/tzone/internal/delivery/handler"
	"github.com/gin-gonic/gin"
)

func MapDeviceRoutes(r *gin.Engine, deviceHandler *handler.DeviceHandler) {
	deviceGroup := r.Group("/api/v1/devices")
	{
		deviceGroup.POST("", deviceHandler.CreateDevice)       // Create a new brand
		deviceGroup.GET("", deviceHandler.GetAllDevices)       // Get all brands
		deviceGroup.GET("/:id", deviceHandler.GetDeviceById)   // Get a brand by ID
		deviceGroup.PUT("/:id", deviceHandler.UpdateDevice)    // Update a brand by ID
		deviceGroup.DELETE("/:id", deviceHandler.DeleteDevice) // Delete a brand by ID
	}

}
