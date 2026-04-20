package handler

import (
	"net/http"
	"strings"

	"github.com/LuuDinhTheTai/tzone/internal/dto"
	"github.com/LuuDinhTheTai/tzone/internal/service"
	"github.com/LuuDinhTheTai/tzone/util/response"
	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	favoriteService *service.FavoriteService
}

func NewFavoriteHandler(favoriteService *service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{favoriteService: favoriteService}
}

func getAuthUserID(c *gin.Context) (string, bool) {
	userIDValue, ok := c.Get("user_id")
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized", nil)
		return "", false
	}

	userID, ok := userIDValue.(string)
	if !ok || strings.TrimSpace(userID) == "" {
		response.Error(c, http.StatusUnauthorized, "unauthorized", nil)
		return "", false
	}

	return userID, true
}

func (h *FavoriteHandler) GetFavorites(c *gin.Context) {
	userID, ok := getAuthUserID(c)
	if !ok {
		return
	}

	favorites, err := h.favoriteService.GetFavorites(userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "favorites retrieved successfully", favorites)
}

func (h *FavoriteHandler) AddFavorite(c *gin.Context) {
	userID, ok := getAuthUserID(c)
	if !ok {
		return
	}

	var req dto.AddFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	favorites, err := h.favoriteService.AddFavorite(c.Request.Context(), userID, req.DeviceID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "favorite added successfully", favorites)
}

func (h *FavoriteHandler) RemoveFavorite(c *gin.Context) {
	userID, ok := getAuthUserID(c)
	if !ok {
		return
	}

	deviceID := c.Param("deviceId")
	if strings.TrimSpace(deviceID) == "" {
		response.Error(c, http.StatusBadRequest, "deviceId is required", nil)
		return
	}

	favorites, err := h.favoriteService.RemoveFavorite(userID, deviceID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "favorite removed successfully", favorites)
}

func (h *FavoriteHandler) SyncFavorites(c *gin.Context) {
	userID, ok := getAuthUserID(c)
	if !ok {
		return
	}

	var req dto.SyncFavoritesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	favorites, err := h.favoriteService.SyncFavorites(c.Request.Context(), userID, req.DeviceIDs)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "favorites synced successfully", favorites)
}
