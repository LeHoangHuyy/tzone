package route

import (
	"github.com/LuuDinhTheTai/tzone/internal/delivery/handler"
	"github.com/LuuDinhTheTai/tzone/internal/delivery/middleware"
	"github.com/LuuDinhTheTai/tzone/internal/service"
	"github.com/gin-gonic/gin"
)

func MapDeviceRoutes(r *gin.Engine, deviceHandler *handler.DeviceHandler, permissionService *service.PermissionService) {
	deviceGroup := r.Group("/api/v1/devices")
	deviceGroup.Use(middleware.APIRateLimit())
	{
		deviceGroup.GET("", deviceHandler.GetAllDevices)
		deviceGroup.GET("/:id", deviceHandler.GetDeviceById)

		// Protected endpoints
		deviceGroup.POST("", middleware.JWTAuth(), middleware.RBACAuth(permissionService), deviceHandler.CreateDevice)
		deviceGroup.PUT("/:id", middleware.JWTAuth(), middleware.RBACAuth(permissionService), deviceHandler.UpdateDevice)
		deviceGroup.DELETE("/:id", middleware.JWTAuth(), middleware.RBACAuth(permissionService), deviceHandler.DeleteDevice)
	}
}
