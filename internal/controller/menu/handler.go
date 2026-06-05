package menu

import (
	"net/http"

	"delicias-da-lu-service.com/mod/internal/entity/menu"
	"delicias-da-lu-service.com/mod/internal/platform/logging"
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
	category := c.QueryParam("category")
	activeStr := c.QueryParam("active")

	var active *bool
	if activeStr != "" {
		val := activeStr == "true"
		active = &val
	}

	logging.DebugEventFromEcho(c).
		Str("category", category).
		Str("active", activeStr).
		Msg("menu get all requested")

	items, err := h.menuUseCase.GetAll(c.Request().Context(), active, category)
	if err != nil {
		logging.ErrorEventFromEcho(c, err).
			Str("category", category).
			Str("active", activeStr).
			Msg("failed to get menu items")
		return err
	}

	if items == nil {
		items = []menu.MenuItem{}
	}

	return c.JSON(http.StatusOK, items)

}
func (h menuHandlerImpl) GetByID(c *echo.Context) error {
	id := c.Param("id")
	logging.DebugEventFromEcho(c).
		Str("menu_item_id", id).
		Msg("menu get by id requested")

	item, err := h.menuUseCase.GetByID(c.Request().Context(), id)
	if err != nil {
		logging.WarnEventFromEcho(c).
			Err(err).
			Str("menu_item_id", id).
			Msg("failed to get menu item")
		return err
	}

	return c.JSON(http.StatusOK, item)

}
func (h menuHandlerImpl) Create(c *echo.Context) error {
	var item menu.MenuItem
	if err := c.Bind(&item); err != nil {
		logging.WarnEventFromEcho(c).
			Err(err).
			Msg("invalid menu item payload")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	logging.DebugEventFromEcho(c).
		Str("menu_item_name", item.Name).
		Str("category", item.Category).
		Msg("menu create requested")

	created, err := h.menuUseCase.Create(c.Request().Context(), &item)
	if err != nil {
		logging.ErrorEventFromEcho(c, err).
			Str("menu_item_name", item.Name).
			Str("category", item.Category).
			Msg("failed to create menu item")
		return err
	}

	return c.JSON(http.StatusCreated, created)

}
func (h menuHandlerImpl) Update(c *echo.Context) error {
	id := c.Param("id")

	var item menu.MenuItem
	if err := c.Bind(&item); err != nil {
		logging.WarnEventFromEcho(c).
			Err(err).
			Str("menu_item_id", id).
			Msg("invalid menu item payload")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	logging.DebugEventFromEcho(c).
		Str("menu_item_id", id).
		Str("menu_item_name", item.Name).
		Msg("menu update requested")

	updated, err := h.menuUseCase.Update(c.Request().Context(), id, &item)
	if err != nil {
		logging.ErrorEventFromEcho(c, err).
			Str("menu_item_id", id).
			Msg("failed to update menu item")
		return err
	}

	return c.JSON(http.StatusOK, updated)

}
func (h menuHandlerImpl) Delete(c *echo.Context) error {
	id := c.Param("id")
	logging.DebugEventFromEcho(c).
		Str("menu_item_id", id).
		Msg("menu delete requested")

	if err := h.menuUseCase.Delete(c.Request().Context(), id); err != nil {
		logging.ErrorEventFromEcho(c, err).
			Str("menu_item_id", id).
			Msg("failed to delete menu item")
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
		logging.WarnEventFromEcho(c).
			Err(err).
			Str("menu_item_id", id).
			Msg("invalid menu order payload")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	logging.DebugEventFromEcho(c).
		Str("menu_item_id", id).
		Int("order", req.Order).
		Msg("menu order update requested")

	updated, err := h.menuUseCase.UpdateOrder(c.Request().Context(), id, req.Order)
	if err != nil {
		logging.ErrorEventFromEcho(c, err).
			Str("menu_item_id", id).
			Int("order", req.Order).
			Msg("failed to update menu item order")
		return err
	}

	return c.JSON(http.StatusOK, updated)

}
