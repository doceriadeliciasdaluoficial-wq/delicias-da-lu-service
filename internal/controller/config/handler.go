package config

import (
	"net/http"

	"delicias-da-lu-service.com/mod/internal/entity/config"
	"delicias-da-lu-service.com/mod/internal/platform/logging"
	configUC "delicias-da-lu-service.com/mod/internal/usecase/config"
	"github.com/labstack/echo/v5"
)

type ConfigHandler interface {
	GetPublic(c *echo.Context) error
	GetAdmin(c *echo.Context) error
	Update(c *echo.Context) error
}

type configHandlerImpl struct {
	configUseCase configUC.ConfigUseCase
}

func NewConfigHandler(configUseCase configUC.ConfigUseCase) ConfigHandler {
	return configHandlerImpl{
		configUseCase: configUseCase,
	}
}

func (h configHandlerImpl) GetPublic(c *echo.Context) error {
	cfg, err := h.configUseCase.Get(c.Request().Context())
	if err != nil {
		logging.ErrorEventFromEcho(c, err).
			Msg("failed to get config")
		return err
	}

	return c.JSON(http.StatusOK, cfg)
}

func (h configHandlerImpl) GetAdmin(c *echo.Context) error {
	cfg, err := h.configUseCase.Get(c.Request().Context())
	if err != nil {
		logging.ErrorEventFromEcho(c, err).
			Msg("failed to get config")
		return err
	}

	return c.JSON(http.StatusOK, cfg)
}

func (h configHandlerImpl) Update(c *echo.Context) error {
	var cfg config.SiteConfig
	if err := c.Bind(&cfg); err != nil {
		logging.WarnEventFromEcho(c).
			Err(err).
			Msg("invalid config payload")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	logging.DebugEventFromEcho(c).
		Msg("config update requested")

	updated, err := h.configUseCase.Update(c.Request().Context(), &cfg)
	if err != nil {
		logging.ErrorEventFromEcho(c, err).
			Msg("failed to update config")
		return err
	}

	return c.JSON(http.StatusOK, updated)
}
