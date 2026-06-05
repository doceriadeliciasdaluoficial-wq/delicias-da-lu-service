package order

import (
	"net/http"
	"strconv"

	"delicias-da-lu-service.com/mod/internal/entity/order"
	orderUC "delicias-da-lu-service.com/mod/internal/usecase/order"
	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog/log"
)

type OrderHandler interface {
	GetAll(c *echo.Context) error
	GetByID(c *echo.Context) error
	Create(c *echo.Context) error
	UpdateStatus(c *echo.Context) error
}

type orderHandlerImpl struct {
	orderUseCase orderUC.OrderUseCase
}

func NewOrderHandler(orderUseCase orderUC.OrderUseCase) OrderHandler {
	return orderHandlerImpl{
		orderUseCase: orderUseCase,
	}
}

func (h orderHandlerImpl) GetAll(c *echo.Context) error {
	status := c.QueryParam("status")
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	offset := 0
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	result, err := h.orderUseCase.GetAll(c.Request().Context(), status, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("failed to get orders")
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func (h orderHandlerImpl) GetByID(c *echo.Context) error {
	id := c.Param("id")

	ord, err := h.orderUseCase.GetByID(c.Request().Context(), id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("failed to get order")
		return err
	}

	return c.JSON(http.StatusOK, ord)
}

func (h orderHandlerImpl) Create(c *echo.Context) error {
	var ord order.Order
	if err := c.Bind(&ord); err != nil {
		log.Error().Err(err).Msg("failed to parse order")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	created, err := h.orderUseCase.Create(c.Request().Context(), &ord)
	if err != nil {
		log.Error().Err(err).Msg("failed to create order")
		return err
	}

	return c.JSON(http.StatusCreated, created)
}

func (h orderHandlerImpl) UpdateStatus(c *echo.Context) error {
	id := c.Param("id")

	var req struct {
		Status string `json:"status"`
	}

	if err := c.Bind(&req); err != nil {
		log.Error().Err(err).Msg("failed to parse status update request")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	updated, err := h.orderUseCase.UpdateStatus(c.Request().Context(), id, req.Status)
	if err != nil {
		log.Error().Err(err).Str("id", id).Str("status", req.Status).Msg("failed to update order status")
		return err
	}

	return c.JSON(http.StatusOK, updated)
}
