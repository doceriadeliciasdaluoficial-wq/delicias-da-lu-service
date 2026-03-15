package controller

import (
	"delicias-da-lu-service.com/mod/internal/controller/system"
	"delicias-da-lu-service.com/mod/internal/controller/user"
	"delicias-da-lu-service.com/mod/internal/platform/problemdetails"

	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog/log"
)

const PORT = ":8080"

type APIServer interface {
	Start() error
	AddRoutes(testeHandler system.Handler, userHandler user.UserHandler) error
}

type apiServerImpl struct {
	server *echo.Echo
}

func NewAPIServer() APIServer {
	return apiServerImpl{
		server: echo.New(),
	}
}

func (ref apiServerImpl) Start() error {
	ref.server.HTTPErrorHandler = problemdetails.ErrorHandler

	if err := ref.server.Start(PORT); err != nil {
		log.Error().Err(err).Msg("server exiting with error")
		return err
	}
	return nil
}

func (ref apiServerImpl) AddRoutes(systemHandler system.Handler, userHandler user.UserHandler) error {

	group := ref.server.Group("/v1")

	group.GET("", systemHandler.Get)
	group.GET("/error", systemHandler.GetError)

	group.POST("/users", userHandler.Create)
	group.GET("/users", userHandler.Get)
	group.PUT("/users/:id", userHandler.Update)
	group.DELETE("/users/:id", userHandler.Delete)

	return nil
}
