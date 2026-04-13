package route

import (
	"github.com/LuuDinhTheTai/tzone/internal/delivery/handler"
	"github.com/LuuDinhTheTai/tzone/internal/delivery/middleware"
	"github.com/LuuDinhTheTai/tzone/internal/service"
	"github.com/gin-gonic/gin"
)

func MapBrandRoutes(r *gin.Engine, brandHandler *handler.BrandHandler, permissionService *service.PermissionService) {
	brandGroup := r.Group("/api/v1/brands")
	brandGroup.Use(middleware.APIRateLimit())
	{
		brandGroup.GET("", brandHandler.GetAllBrands)
		brandGroup.GET("/:id", brandHandler.GetBrandById)

		// Protected endpoints
		brandGroup.POST("", middleware.JWTAuth(), middleware.RBACAuth(permissionService), brandHandler.CreateBrand)
		brandGroup.PUT("/:id", middleware.JWTAuth(), middleware.RBACAuth(permissionService), brandHandler.UpdateBrand)
		brandGroup.DELETE("/:id", middleware.JWTAuth(), middleware.RBACAuth(permissionService), brandHandler.DeleteBrand)
	}
}
