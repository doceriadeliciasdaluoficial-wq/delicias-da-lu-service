package menu

import (
	"net/http"

	"delicias-da-lu-service.com/mod/internal/entity/menu"
	menuUC "delicias-da-lu-service.com/mod/internal/usecase/menu"
	"github.com/labstack/echo/v5"
)

type MenuHandler interface {
	GetAll(c *echo.Context) error
	GetByID(c *echo.Context) error
	Create(c *echo.Context) error
	Update(c *echo.Context) error
	Delete(c *echo.Context) error
	UpdateOrder(c *echo.Context) error
}

type menuHandlerImpl struct {
	menuUseCase menuUC.MenuUseCase
}

func NewMenuHandler(menuUseCase menuUC.MenuUseCase) MenuHandler {
	return menuHandlerImpl{
		menuUseCase: menuUseCase,
	}
}

func (h menuHandlerImpl) GetAll(c *echo.Context) error {
	var active *bool
	if c.QueryParam("active") != "" {
		activeVal := c.QueryParam("active") == "true"
		active = &activeVal
	}

	category := c.QueryParam("category")

	items, err := h.menuUseCase.GetAll(c.Request().Context(), active, category)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	if items == nil {
		items = []menu.MenuItem{}
	}

	return c.JSON(http.StatusOK, items)
}

func (h menuHandlerImpl) GetByID(c *echo.Context) error {
	categoryID := c.Param("category")
	itemID := c.Param("id")

	item, err := h.menuUseCase.GetByID(c.Request().Context(), categoryID, itemID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, item)
}

func (h menuHandlerImpl) Create(c *echo.Context) error {
	categoryID := c.Param("category")

	var item menu.MenuItem
	if err := c.Bind(&item); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	item.Category = categoryID

	createdItem, err := h.menuUseCase.Create(c.Request().Context(), categoryID, &item)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, createdItem)
}

func (h menuHandlerImpl) Update(c *echo.Context) error {
	categoryID := c.Param("category")
	itemID := c.Param("id")

	var item menu.MenuItem
	if err := c.Bind(&item); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	updatedItem, err := h.menuUseCase.Update(c.Request().Context(), categoryID, itemID, &item)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, updatedItem)
}

func (h menuHandlerImpl) Delete(c *echo.Context) error {
	categoryID := c.Param("category")
	itemID := c.Param("id")

	if err := h.menuUseCase.Delete(c.Request().Context(), categoryID, itemID); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h menuHandlerImpl) UpdateOrder(c *echo.Context) error {
	categoryID := c.Param("category")
	itemID := c.Param("id")

	var req struct {
		Order int `json:"order"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	item, err := h.menuUseCase.UpdateOrder(c.Request().Context(), categoryID, itemID, req.Order)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, item)
}
