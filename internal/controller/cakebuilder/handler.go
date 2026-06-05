package cakebuilder

import (
	"net/http"

	"delicias-da-lu-service.com/mod/internal/entity/cakebuilder"
	"delicias-da-lu-service.com/mod/internal/entity/config"
	"delicias-da-lu-service.com/mod/internal/platform/logging"
	cakebuildUC "delicias-da-lu-service.com/mod/internal/usecase/cakebuilder"
	"github.com/labstack/echo/v5"
)

type CakeBuilderHandler interface {
	GetAll(c *echo.Context) error
	GetByType(c *echo.Context) error
	GetByID(c *echo.Context) error
	Create(c *echo.Context) error
	Update(c *echo.Context) error
	Delete(c *echo.Context) error
}

type cakeBuilderHandlerImpl struct {
	cakeBuilderUseCase cakebuildUC.CakeBuilderUseCase
}

func NewCakeBuilderHandler(cakeBuilderUseCase cakebuildUC.CakeBuilderUseCase) CakeBuilderHandler {
	return cakeBuilderHandlerImpl{
		cakeBuilderUseCase: cakeBuilderUseCase,
	}
}

func (h cakeBuilderHandlerImpl) GetAll(c *echo.Context) error {
	activeStr := c.QueryParam("active")

	var active *bool
	if activeStr != "" {
		val := activeStr == "true"
		active = &val
	}

	logging.DebugEventFromEcho(c).
		Str("active", activeStr).
		Msg("cake builder get all requested")

	components, err := h.cakeBuilderUseCase.GetAll(c.Request().Context(), active)
	if err != nil {
		logging.ErrorEventFromEcho(c, err).
			Str("active", activeStr).
			Msg("failed to get cake builder components")
		return err
	}

	result := &config.CakeBuilderConfig{
		Massas:     components["massas"],
		Recheios:   components["recheios"],
		Coberturas: components["coberturas"],
		Decoracoes: components["decoracoes"],
	}

	return c.JSON(http.StatusOK, result)
}

func (h cakeBuilderHandlerImpl) GetByType(c *echo.Context) error {
	componentType := c.Param("type")
	activeStr := c.QueryParam("active")

	var active *bool
	if activeStr != "" {
		val := activeStr == "true"
		active = &val
	}

	logging.DebugEventFromEcho(c).
		Str("component_type", componentType).
		Str("active", activeStr).
		Msg("cake builder get by type requested")

	components, err := h.cakeBuilderUseCase.GetByType(c.Request().Context(), componentType, active)
	if err != nil {
		logging.ErrorEventFromEcho(c, err).
			Str("component_type", componentType).
			Msg("failed to get cake builder components by type")
		return err
	}

	if components == nil {
		components = []cakebuilder.CakeBuilderComponent{}
	}

	return c.JSON(http.StatusOK, components)
}

func (h cakeBuilderHandlerImpl) GetByID(c *echo.Context) error {
	componentType := c.Param("type")
	id := c.Param("id")
	logging.DebugEventFromEcho(c).
		Str("component_type", componentType).
		Str("component_id", id).
		Msg("cake builder get by id requested")

	component, err := h.cakeBuilderUseCase.GetByID(c.Request().Context(), componentType, id)
	if err != nil {
		logging.WarnEventFromEcho(c).
			Err(err).
			Str("component_type", componentType).
			Str("component_id", id).
			Msg("failed to get cake builder component")
		return err
	}

	return c.JSON(http.StatusOK, component)
}

func (h cakeBuilderHandlerImpl) Create(c *echo.Context) error {
	componentType := c.Param("type")

	var component cakebuilder.CakeBuilderComponent
	if err := c.Bind(&component); err != nil {
		logging.WarnEventFromEcho(c).
			Err(err).
			Str("component_type", componentType).
			Msg("invalid cake builder payload")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	logging.DebugEventFromEcho(c).
		Str("component_type", componentType).
		Str("component_name", component.Name).
		Msg("cake builder create requested")

	component.Type = componentType

	created, err := h.cakeBuilderUseCase.Create(c.Request().Context(), &component)
	if err != nil {
		logging.ErrorEventFromEcho(c, err).
			Str("component_type", componentType).
			Msg("failed to create cake builder component")
		return err
	}

	return c.JSON(http.StatusCreated, created)
}

func (h cakeBuilderHandlerImpl) Update(c *echo.Context) error {
	componentType := c.Param("type")
	id := c.Param("id")

	var component cakebuilder.CakeBuilderComponent
	if err := c.Bind(&component); err != nil {
		logging.WarnEventFromEcho(c).
			Err(err).
			Str("component_type", componentType).
			Str("component_id", id).
			Msg("invalid cake builder payload")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	logging.DebugEventFromEcho(c).
		Str("component_type", componentType).
		Str("component_id", id).
		Str("component_name", component.Name).
		Msg("cake builder update requested")

	component.Type = componentType

	updated, err := h.cakeBuilderUseCase.Update(c.Request().Context(), id, &component)
	if err != nil {
		logging.ErrorEventFromEcho(c, err).
			Str("component_type", componentType).
			Str("component_id", id).
			Msg("failed to update cake builder component")
		return err
	}

	return c.JSON(http.StatusOK, updated)
}

func (h cakeBuilderHandlerImpl) Delete(c *echo.Context) error {
	componentType := c.Param("type")
	id := c.Param("id")
	logging.DebugEventFromEcho(c).
		Str("component_type", componentType).
		Str("component_id", id).
		Msg("cake builder delete requested")

	if err := h.cakeBuilderUseCase.Delete(c.Request().Context(), componentType, id); err != nil {
		logging.ErrorEventFromEcho(c, err).
			Str("component_type", componentType).
			Str("component_id", id).
			Msg("failed to delete cake builder component")
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
