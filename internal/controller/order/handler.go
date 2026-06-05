package order

import (
	"net/http"
	"strconv"

	"delicias-da-lu-service.com/mod/internal/entity/order"
	"delicias-da-lu-service.com/mod/internal/platform/logging"
	orderUC "delicias-da-lu-service.com/mod/internal/usecase/order"
	"github.com/labstack/echo/v5"
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

	logging.DebugEventFromEcho(c).
		Str("status", status).
		Int("limit", limit).
		Int("offset", offset).
		Msg("order list requested")

	result, err := h.orderUseCase.GetAll(c.Request().Context(), status, limit, offset)
	if err != nil {
		logging.ErrorEventFromEcho(c, err).
			Str("status", status).
			Int("limit", limit).
			Int("offset", offset).
			Msg("failed to get orders")
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func (h orderHandlerImpl) GetByID(c *echo.Context) error {
	id := c.Param("id")
	logging.DebugEventFromEcho(c).
		Str("order_id", id).
		Msg("order get by id requested")

	ord, err := h.orderUseCase.GetByID(c.Request().Context(), id)
	if err != nil {
		logging.WarnEventFromEcho(c).
			Err(err).
			Str("order_id", id).
			Msg("failed to get order")
		return err
	}

	return c.JSON(http.StatusOK, ord)
}

func (h orderHandlerImpl) Create(c *echo.Context) error {
	var ord order.Order
	if err := c.Bind(&ord); err != nil {
		logging.WarnEventFromEcho(c).
			Err(err).
			Msg("invalid order payload")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	logging.DebugEventFromEcho(c).
		Str("customer_name", ord.CustomerInfo.Name).
		Str("customer_phone", ord.CustomerInfo.Phone).
		Msg("order create requested")

	created, err := h.orderUseCase.Create(c.Request().Context(), &ord)
	if err != nil {
		logging.ErrorEventFromEcho(c, err).
			Str("customer_name", ord.CustomerInfo.Name).
			Msg("failed to create order")
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
		logging.WarnEventFromEcho(c).
			Err(err).
			Str("order_id", id).
			Msg("invalid order status payload")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	logging.DebugEventFromEcho(c).
		Str("order_id", id).
		Str("status", req.Status).
		Msg("order status update requested")

	updated, err := h.orderUseCase.UpdateStatus(c.Request().Context(), id, req.Status)
	if err != nil {
		logging.ErrorEventFromEcho(c, err).
			Str("order_id", id).
			Str("status", req.Status).
			Msg("failed to update order status")
		return err
	}

	return c.JSON(http.StatusOK, updated)
}
