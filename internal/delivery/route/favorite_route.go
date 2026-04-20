package route

import (
	"github.com/LuuDinhTheTai/tzone/internal/delivery/handler"
	"github.com/LuuDinhTheTai/tzone/internal/delivery/middleware"
	"github.com/gin-gonic/gin"
)

func MapFavoriteRoutes(r *gin.Engine, favoriteHandler *handler.FavoriteHandler) {
	favoriteGroup := r.Group("/api/v1/favorites")
	favoriteGroup.Use(middleware.APIRateLimit(), middleware.JWTAuth())
	{
		favoriteGroup.GET("", favoriteHandler.GetFavorites)
		favoriteGroup.POST("", favoriteHandler.AddFavorite)
		favoriteGroup.DELETE("/:deviceId", favoriteHandler.RemoveFavorite)
		favoriteGroup.POST("/sync", favoriteHandler.SyncFavorites)
	}
}
