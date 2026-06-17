package upload

import (
	"fmt"
	"net/http"
	"path/filepath"

	"delicias-da-lu-service.com/mod/internal/platform/logging"
	"delicias-da-lu-service.com/mod/internal/repository/storage"
	"github.com/labstack/echo/v5"
)

type UploadHandler interface {
	UploadImage(c *echo.Context) error
	DeleteImage(c *echo.Context) error
}

type uploadHandlerImpl struct {
	storageRepo storage.StorageRepository
}

func NewUploadHandler(storageRepo storage.StorageRepository) UploadHandler {
	return uploadHandlerImpl{
		storageRepo: storageRepo,
	}
}

func (h uploadHandlerImpl) UploadImage(c *echo.Context) error {
	uploadType := c.QueryParam("type")
	category := c.QueryParam("category")
	id := c.QueryParam("id")

	if uploadType == "" || id == "" {
		logging.WarnEventFromEcho(c).Msg("missing upload parameters")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing required parameters: type, id",
		})
	}

	file, err := c.FormFile("file")
	if err != nil {
		logging.WarnEventFromEcho(c).Err(err).Msg("no file provided")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "No file provided",
		})
	}

	if file.Size > 5*1024*1024 {
		logging.WarnEventFromEcho(c).Msg("file too large")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "File too large (max 5MB)",
		})
	}

	contentType := file.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/webp" {
		logging.WarnEventFromEcho(c).Str("contentType", contentType).Msg("invalid file type")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid file type. Only JPEG, PNG, and WebP allowed",
		})
	}

	var filePath string
	switch uploadType {
	case "menu":
		if category == "" {
			logging.WarnEventFromEcho(c).Msg("category required for menu uploads")
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Category required for menu uploads",
			})
		}
		ext := filepath.Ext(file.Filename)
		filePath = fmt.Sprintf("menu/%s/%s%s", category, id, ext)

	case "cakeBuilder":
		if category == "" {
			logging.WarnEventFromEcho(c).Msg("type required for cakeBuilder uploads")
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Type required for cakeBuilder uploads",
			})
		}
		ext := filepath.Ext(file.Filename)
		filePath = fmt.Sprintf("cakeBuilder/%s/%s%s", category, id, ext)

	case "home":
		ext := filepath.Ext(file.Filename)
		filePath = fmt.Sprintf("home/featured/%s%s", id, ext)

	default:
		logging.WarnEventFromEcho(c).Str("type", uploadType).Msg("invalid upload type")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid type. Must be: menu, cakeBuilder, or home",
		})
	}

	ctx := c.Request().Context()
	url, err := h.storageRepo.UploadFile(ctx, filePath, file)
	if err != nil {
		logging.ErrorEventFromEcho(c, err).Str("path", filePath).Msg("failed to upload file")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to upload file",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"url":  url,
		"path": filePath,
	})
}

func (h uploadHandlerImpl) DeleteImage(c *echo.Context) error {
	path := c.QueryParam("path")
	if path == "" {
		logging.WarnEventFromEcho(c).Msg("missing path parameter")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing required parameter: path",
		})
	}

	ctx := c.Request().Context()
	if err := h.storageRepo.DeleteFile(ctx, path); err != nil {
		logging.ErrorEventFromEcho(c, err).Str("path", path).Msg("failed to delete file")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete file",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "File deleted successfully",
	})
}
