package teste

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type Handler interface {
	Create(e *echo.Context) error
}

type handlerImpl struct {
}

func NewHandler() Handler {
	return handlerImpl{}
}

func (ref handlerImpl) Create(e *echo.Context) error {
	return e.JSON(http.StatusCreated, map[string]string{
		"message": "success",
	})
}
