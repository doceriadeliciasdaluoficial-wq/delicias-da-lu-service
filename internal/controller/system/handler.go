package system

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type Handler interface {
	Get(e *echo.Context) error
}

type handlerImpl struct {
}

func NewHandler() Handler {
	return handlerImpl{}
}

func (ref handlerImpl) Get(e *echo.Context) error {
	return e.JSON(http.StatusOK, map[string]string{
		"admin":     "Docerias da Lu",
		"email":     "doceriadeliciasdaluoficial@gmail.com",
		"instagram": "delicias.lu.oficial",
	})
}
