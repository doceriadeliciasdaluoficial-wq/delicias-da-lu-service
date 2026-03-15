package user

import "github.com/labstack/echo/v5"

type UserHandler interface {
	Create(*echo.Context) error
	Get(*echo.Context) error
	Update(*echo.Context) error
	Delete(*echo.Context) error
}

type userHandlerImpl struct {
}

func NewUserHandler() UserHandler {
	return userHandlerImpl{}
}

func (ref userHandlerImpl) Create(ctx *echo.Context) error {
	return nil
}

func (ref userHandlerImpl) Get(ctx *echo.Context) error {
	return nil
}

func (ref userHandlerImpl) Update(ctx *echo.Context) error {
	return nil
}

func (ref userHandlerImpl) Delete(ctx *echo.Context) error {
	return nil
}
