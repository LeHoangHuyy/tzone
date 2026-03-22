package route

import (
	"github.com/LuuDinhTheTai/tzone/internal/delivery/handler"
	"github.com/gin-gonic/gin"
)

func MapBrandRoutes(r *gin.Engine, brandHandler *handler.BrandHandler) {
	brandGroup := r.Group("/api/v1/brands")
	{
		brandGroup.POST("", brandHandler.CreateBrand)       // Create a new brand
		brandGroup.GET("", brandHandler.GetAllBrands)       // Get all brands
		brandGroup.GET("/:id", brandHandler.GetBrandById)   // Get a brand by ID
		brandGroup.PUT("/:id", brandHandler.UpdateBrand)    // Update a brand by ID
		brandGroup.DELETE("/:id", brandHandler.DeleteBrand) // Delete a brand by ID
	}
}
