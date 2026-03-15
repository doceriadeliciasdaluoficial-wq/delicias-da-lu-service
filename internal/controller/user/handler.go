package user

import (
	"net/http"

	"delicias-da-lu-service.com/mod/internal/entity/user"
	userUsecase "delicias-da-lu-service.com/mod/internal/usecase/user"

	"github.com/labstack/echo/v5"
)

type UserHandler interface {
	Create(*echo.Context) error
	Get(*echo.Context) error
	Update(*echo.Context) error
	Delete(*echo.Context) error
}

type userHandlerImpl struct {
	userUsecase userUsecase.UserUseCase
}

func NewUserHandler(userUsecase userUsecase.UserUseCase) UserHandler {
	return userHandlerImpl{
		userUsecase: userUsecase,
	}
}

func (ref userHandlerImpl) Create(ctx *echo.Context) error {
	var user user.User
	if err := ctx.Bind(&user); err != nil {
		return err
	}
	createdUser, err := ref.userUsecase.Create(ctx.Request().Context(), &user)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, createdUser)
}

func (ref userHandlerImpl) Get(ctx *echo.Context) error {
	queryParameterField := ctx.QueryParam("field")
	queryParameterValue := ctx.QueryParam("value")
	users, err := ref.userUsecase.Get(ctx.Request().Context(), queryParameterField, queryParameterValue)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, users)
}

func (ref userHandlerImpl) Update(ctx *echo.Context) error {
	id := ctx.Param("id")
	var user user.User
	if err := ctx.Bind(&user); err != nil {
		return err
	}
	updatedUser, err := ref.userUsecase.Update(ctx.Request().Context(), id, &user)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, updatedUser)
}

func (ref userHandlerImpl) Delete(ctx *echo.Context) error {
	id := ctx.Param("id")
	err := ref.userUsecase.Delete(ctx.Request().Context(), id)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusNoContent, nil)
}
