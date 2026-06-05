package health

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type HealthHandler interface {
	Check(c *echo.Context) error
}

type healthHandlerImpl struct{}

func NewHealthHandler() HealthHandler {
	return healthHandlerImpl{}
}

func (h healthHandlerImpl) Check(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}
