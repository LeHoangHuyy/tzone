package handle_uploads

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func SaveImage(c *gin.Context, file *multipart.FileHeader) (string, error) {
	
	uploadDir := "uploads/devices"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
	filePath := filepath.Join(uploadDir, fileName)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return "/" + filePath, nil
}
