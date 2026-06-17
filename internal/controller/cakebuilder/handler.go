package cakebuilder

import (
	"net/http"

	"delicias-da-lu-service.com/mod/internal/entity/cakebuilder"
	cakeUC "delicias-da-lu-service.com/mod/internal/usecase/cakebuilder"
	"github.com/labstack/echo/v5"
)

type CakeBuilderHandler interface {
	GetAll(c *echo.Context) error
	GetByType(c *echo.Context) error
	GetByID(c *echo.Context) error
	Create(c *echo.Context) error
	Update(c *echo.Context) error
	Delete(c *echo.Context) error
	UpdateOrder(c *echo.Context) error
}

type cakeBuilderHandlerImpl struct {
	cakeBuilderUseCase cakeUC.CakeBuilderUseCase
}

func NewCakeBuilderHandler(cakeBuilderUseCase cakeUC.CakeBuilderUseCase) CakeBuilderHandler {
	return cakeBuilderHandlerImpl{
		cakeBuilderUseCase: cakeBuilderUseCase,
	}
}

func (h cakeBuilderHandlerImpl) GetAll(c *echo.Context) error {
	var active *bool
	if c.QueryParam("active") != "" {
		activeVal := c.QueryParam("active") == "true"
		active = &activeVal
	}

	items, err := h.cakeBuilderUseCase.GetAll(c.Request().Context(), active)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, items)
}

func (h cakeBuilderHandlerImpl) GetByType(c *echo.Context) error {
	cType := c.Param("type")

	var active *bool
	if c.QueryParam("active") != "" {
		activeVal := c.QueryParam("active") == "true"
		active = &activeVal
	}

	components, err := h.cakeBuilderUseCase.GetByType(c.Request().Context(), cType, active)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	if components == nil {
		components = []cakebuilder.CakeBuilderComponent{}
	}

	return c.JSON(http.StatusOK, components)
}

func (h cakeBuilderHandlerImpl) GetByID(c *echo.Context) error {
	cType := c.Param("type")
	id := c.Param("id")

	component, err := h.cakeBuilderUseCase.GetByID(c.Request().Context(), cType, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, component)
}

func (h cakeBuilderHandlerImpl) Create(c *echo.Context) error {
	cType := c.Param("type")

	var component cakebuilder.CakeBuilderComponent
	if err := c.Bind(&component); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	component.Type = cType

	createdComponent, err := h.cakeBuilderUseCase.Create(c.Request().Context(), cType, &component)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, createdComponent)
}

func (h cakeBuilderHandlerImpl) Update(c *echo.Context) error {
	cType := c.Param("type")
	id := c.Param("id")

	var component cakebuilder.CakeBuilderComponent
	if err := c.Bind(&component); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	updatedComponent, err := h.cakeBuilderUseCase.Update(c.Request().Context(), cType, id, &component)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, updatedComponent)
}

func (h cakeBuilderHandlerImpl) Delete(c *echo.Context) error {
	cType := c.Param("type")
	id := c.Param("id")

	if err := h.cakeBuilderUseCase.Delete(c.Request().Context(), cType, id); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h cakeBuilderHandlerImpl) UpdateOrder(c *echo.Context) error {
	cType := c.Param("type")
	id := c.Param("id")

	var req struct {
		Order int `json:"order"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	component, err := h.cakeBuilderUseCase.UpdateOrder(c.Request().Context(), cType, id, req.Order)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, component)
}
