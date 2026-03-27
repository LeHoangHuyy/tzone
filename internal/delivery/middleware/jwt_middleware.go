package middleware

import (
	"net/http"
	"strings"

	"github.com/LuuDinhTheTai/tzone/internal/service"
	"github.com/LuuDinhTheTai/tzone/util/jwt"
	"github.com/LuuDinhTheTai/tzone/util/response"

	"github.com/gin-gonic/gin"
)

// JWTAuth middleware verifies the JWT token present in the Authorization header
func JWTAuth(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "authorization header is required", nil)
			c.Abort()
			return
		}

		// Check basic Bearer format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "invalid authorization format, expected Bearer <token>", nil)
			c.Abort()
			return
		}

		// Extract physical token
		tokenString := parts[1]

		// Validate token
		userID, err := jwt.ValidateToken(tokenString)
		if err != nil {
			if err == jwt.ErrExpiredToken {
				// Try auto-refresh if client provided Refresh-Token header
				rtString := c.GetHeader("Refresh-Token")
				if rtString != "" && authService != nil {
					newAccessToken, newRefreshToken, newUserID, errRefresh := authService.RefreshToken(rtString)
					if errRefresh == nil {
						// Auto refresh succeeded
						c.Header("X-New-Access-Token", newAccessToken)
						c.Header("X-New-Refresh-Token", newRefreshToken)
						
						// Attach user_id and continue successfully
						c.Set("user_id", newUserID.String())
						c.Next()
						return
					}
				}
				// Refresh failed or no Refresh-Token provided
				response.Error(c, http.StatusUnauthorized, "token has expired", nil)
				c.Abort()
				return
			}
			
			response.Error(c, http.StatusUnauthorized, "invalid token", nil)
			c.Abort()
			return
		}

		// Attach user_id to context if token is valid
		c.Set("user_id", userID.String())

		// Continue processing normally
		c.Next()
	}
}
