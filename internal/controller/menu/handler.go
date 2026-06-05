package menu

import (
	"net/http"

	"delicias-da-lu-service.com/mod/internal/entity/menu"
	menuUC "delicias-da-lu-service.com/mod/internal/usecase/menu"
	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog/log"
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
	category := c.QueryParam("category")
	activeStr := c.QueryParam("active")

	var active *bool
	if activeStr != "" {
		val := activeStr == "true"
		active = &val
	}

	items, err := h.menuUseCase.GetAll(c.Request().Context(), active, category)
	if err != nil {
		log.Error().Err(err).Msg("failed to get menu items")
		return err
	}

	if items == nil {
		items = []menu.MenuItem{}
	}

	return c.JSON(http.StatusOK, items)

}
func (h menuHandlerImpl) GetByID(c *echo.Context) error {
	id := c.Param("id")

	item, err := h.menuUseCase.GetByID(c.Request().Context(), id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("failed to get menu item")
		return err
	}

	return c.JSON(http.StatusOK, item)

}
func (h menuHandlerImpl) Create(c *echo.Context) error {
	var item menu.MenuItem
	if err := c.Bind(&item); err != nil {
		log.Error().Err(err).Msg("failed to parse menu item")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	created, err := h.menuUseCase.Create(c.Request().Context(), &item)
	if err != nil {
		log.Error().Err(err).Msg("failed to create menu item")
		return err
	}

	return c.JSON(http.StatusCreated, created)

}
func (h menuHandlerImpl) Update(c *echo.Context) error {
	id := c.Param("id")

	var item menu.MenuItem
	if err := c.Bind(&item); err != nil {
		log.Error().Err(err).Msg("failed to parse menu item")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	updated, err := h.menuUseCase.Update(c.Request().Context(), id, &item)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("failed to update menu item")
		return err
	}

	return c.JSON(http.StatusOK, updated)

}
func (h menuHandlerImpl) Delete(c *echo.Context) error {
	id := c.Param("id")

	if err := h.menuUseCase.Delete(c.Request().Context(), id); err != nil {
		log.Error().Err(err).Str("id", id).Msg("failed to delete menu item")
		return err
	}

	return c.NoContent(http.StatusNoContent)

}
func (h menuHandlerImpl) UpdateOrder(c *echo.Context) error {
	id := c.Param("id")

	var req struct {
		Order int `json:"order"`
	}

	if err := c.Bind(&req); err != nil {
		log.Error().Err(err).Msg("failed to parse order update request")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	updated, err := h.menuUseCase.UpdateOrder(c.Request().Context(), id, req.Order)
	if err != nil {
		log.Error().Err(err).Str("id", id).Int("order", req.Order).Msg("failed to update menu item order")
		return err
	}

	return c.JSON(http.StatusOK, updated)

}
