package route

import (
	"github.com/LuuDinhTheTai/tzone/internal/delivery/handler"
	"github.com/LuuDinhTheTai/tzone/internal/delivery/middleware"
	"github.com/gin-gonic/gin"
)

func MapAuthRoutes(r *gin.Engine, h *handler.AuthHandler) {

	auth := r.Group("/auth")
	auth.Use(middleware.AuthRateLimit())

	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
	auth.POST("/refresh", h.RefreshToken)
	auth.POST("/logout", h.Logout)
}
