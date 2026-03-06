package controller

import (
	"delicias-da-lu-service/internal/controller/teste"

	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog/log"
)

const PORT = ":8080"

type APIServer interface {
	Start() error
	AddRoutes(testeHandler teste.Handler) error
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
	if err := ref.server.Start(PORT); err != nil {
		log.Error().Err(err).Msg("server exiting with error")
		return err
	}
	return nil
}

func (ref apiServerImpl) AddRoutes(testeHandler teste.Handler) error {

	group := ref.server.Group("/v1")

	group.POST("", testeHandler.Create)

	return nil
}
