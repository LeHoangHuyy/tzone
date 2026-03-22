package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"maragu.dev/gomponents"
)

type ApiResponse struct {
	Success bool            `json:"success"`
	Code    int             `json:"code"`
	Message string          `json:"message,omitempty"`
	Data    interface{}     `json:"data,omitempty"`
	Errors  []ErrorResponse `json:"errors,omitempty"`
	Meta    MetaResponse    `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Field string `json:"field,omitempty"`
	Error string `json:"error"`
}

type MetaResponse struct {
	Timestamp string `json:"timestamp"`
	RequestId string `json:"request_id"`
}

func Success(ctx *gin.Context, code int, msg string, data interface{}) {
	ctx.JSON(code, ApiResponse{
		Success: true,
		Code:    code,
		Message: msg,
		Data:    data,
	})
}

func Error(ctx *gin.Context, code int, msg string, errorResponse []ErrorResponse) {
	ctx.JSON(code, ApiResponse{
		Success: false,
		Code:    code,
		Message: msg,
		Errors:  errorResponse,
	})
}

func HTML(ctx *gin.Context, node gomponents.Node) {
	ctx.Header("Content-Type", "text/html; charset=utf-8")

	ctx.Status(http.StatusOK)

	err := node.Render(ctx.Writer)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}
