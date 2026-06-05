package auth

import (
	"net/http"
	"strings"

	"delicias-da-lu-service.com/mod/internal/entity/user"
	authUC "delicias-da-lu-service.com/mod/internal/usecase/auth"
	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog/log"
)

type AuthHandler interface {
	Login(c *echo.Context) error
	Refresh(c *echo.Context) error
}

type authHandlerImpl struct {
	authUseCase authUC.AuthUseCase
}

func NewAuthHandler(authUseCase authUC.AuthUseCase) AuthHandler {
	return authHandlerImpl{
		authUseCase: authUseCase,
	}
}

func (h authHandlerImpl) Login(c *echo.Context) error {
	var req user.LoginRequest
	if err := c.Bind(&req); err != nil {
		log.Error().Err(err).Msg("failed to parse login request")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	response, err := h.authUseCase.Login(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}

func (h authHandlerImpl) Refresh(c *echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "missing authorization header",
		})
	}

	// Extract token (assuming "Bearer <token>" format)
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid authorization header",
		})
	}

	token := parts[1]
	newToken, err := h.authUseCase.RefreshToken(c.Request().Context(), token)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user.RefreshTokenResponse{
		Token: newToken,
	})
}
