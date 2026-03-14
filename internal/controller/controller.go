package controller

import (
	"delicias-da-lu-service.com/mod/internal/controller/system"
	"delicias-da-lu-service.com/mod/internal/platform/problemdetails"

	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog/log"
)

const PORT = ":8080"

type APIServer interface {
	Start() error
	AddRoutes(testeHandler system.Handler) error
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

func (ref apiServerImpl) AddRoutes(testeHandler system.Handler) error {

	group := ref.server.Group("/v1")

	group.GET("", testeHandler.Get)
	group.GET("/error", testeHandler.GetError)

	return nil
}
