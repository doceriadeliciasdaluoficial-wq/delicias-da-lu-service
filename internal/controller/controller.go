package controller

import (
	"delicias-da-lu-service.com/mod/internal/controller/auth"
	"delicias-da-lu-service.com/mod/internal/controller/cakebuilder"
	"delicias-da-lu-service.com/mod/internal/controller/config"
	"delicias-da-lu-service.com/mod/internal/controller/contact"
	"delicias-da-lu-service.com/mod/internal/controller/health"
	"delicias-da-lu-service.com/mod/internal/controller/menu"
	"delicias-da-lu-service.com/mod/internal/controller/order"
	"delicias-da-lu-service.com/mod/internal/controller/system"
	"delicias-da-lu-service.com/mod/internal/platform/middleware"
	"delicias-da-lu-service.com/mod/internal/platform/problemdetails"
	authUsecase "delicias-da-lu-service.com/mod/internal/usecase/auth"

	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog/log"
)

const PORT = ":8080"

type APIServer interface {
	Start() error
	RegisterRoutes(
		healthHandler health.HealthHandler,
		authHandler auth.AuthHandler,
		authUseCase authUsecase.AuthUseCase,
		menuHandler menu.MenuHandler,
		cakeBuilderHandler cakebuilder.CakeBuilderHandler,
		contactHandler contact.ContactHandler,
		configHandler config.ConfigHandler,
		orderHandler order.OrderHandler,
		systemHandler system.Handler,
	) error
}

type apiServerImpl struct {
	server *echo.Echo
}

func NewAPIServer() APIServer {
	server := echo.New()
	server.Use(middleware.TraceIDMiddleware())
	server.Use(middleware.RequestLogger())
	return apiServerImpl{
		server: server,
	}
}

func (ref apiServerImpl) Start() error {
	ref.server.HTTPErrorHandler = problemdetails.ErrorHandler

	if err := ref.server.Start(PORT); err != nil {
		log.Error().Err(err).Str("component", "api_server").Msg("server exiting with error")
		return err
	}
	return nil
}

func (ref apiServerImpl) RegisterRoutes(
	healthHandler health.HealthHandler,
	authHandler auth.AuthHandler,
	authUseCase authUsecase.AuthUseCase,
	menuHandler menu.MenuHandler,
	cakeBuilderHandler cakebuilder.CakeBuilderHandler,
	contactHandler contact.ContactHandler,
	configHandler config.ConfigHandler,
	orderHandler order.OrderHandler,
	systemHandler system.Handler,
) error {
	v1 := ref.server.Group("/v1")

	// Health check
	v1.GET("/health", healthHandler.Check)

	// Error documentation
	v1.GET("/error", systemHandler.GetError)

	// Auth endpoints
	v1.POST("/auth/login", authHandler.Login)
	v1.POST("/auth/refresh", authHandler.Refresh, middleware.JWTMiddleware(authUseCase))

	// Public config/menu endpoints
	v1.GET("/config/public", configHandler.GetPublic)
	v1.GET("/menu/items", menuHandler.GetAll)
	v1.GET("/menu/items/:id", menuHandler.GetByID)
	v1.GET("/cake-builder", cakeBuilderHandler.GetAll)
	v1.GET("/cake-builder/:type", cakeBuilderHandler.GetByType)
	v1.GET("/cake-builder/:type/:id", cakeBuilderHandler.GetByID)
	v1.GET("/contacts", contactHandler.Get)

	// Order endpoints (public for creation, admin for listing)
	v1.POST("/orders", orderHandler.Create)
	v1.GET("/orders/:id", orderHandler.GetByID)

	// Admin protected endpoints
	adminGroup := v1.Group("", middleware.JWTMiddleware(authUseCase))

	// Config admin endpoints
	adminGroup.GET("/config/admin", configHandler.GetAdmin)
	adminGroup.PUT("/config/admin", configHandler.Update)

	// Menu admin endpoints
	adminGroup.POST("/menu/items", menuHandler.Create)
	adminGroup.PUT("/menu/items/:id", menuHandler.Update)
	adminGroup.DELETE("/menu/items/:id", menuHandler.Delete)
	adminGroup.PATCH("/menu/items/:id/order", menuHandler.UpdateOrder)

	// CakeBuilder admin endpoints
	adminGroup.POST("/cake-builder/:type", cakeBuilderHandler.Create)
	adminGroup.PUT("/cake-builder/:type/:id", cakeBuilderHandler.Update)
	adminGroup.DELETE("/cake-builder/:type/:id", cakeBuilderHandler.Delete)

	// Contacts admin endpoints
	adminGroup.PUT("/contacts", contactHandler.Update)

	// Orders admin endpoints
	adminGroup.GET("/orders", orderHandler.GetAll)
	adminGroup.PUT("/orders/:id", orderHandler.UpdateStatus)

	return nil
}
