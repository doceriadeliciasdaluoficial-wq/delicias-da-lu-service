package middleware

import (
	"net/http"
	"strings"

	"delicias-da-lu-service.com/mod/internal/platform/logging"
	"delicias-da-lu-service.com/mod/internal/platform/problemdetails"
	"delicias-da-lu-service.com/mod/internal/usecase/auth"
	"github.com/labstack/echo/v5"
)

func JWTMiddleware(authUseCase auth.AuthUseCase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				logging.WarnEventFromEcho(c).
					Str("auth_header", "missing").
					Msg("authorization header missing")
				return problemdetails.NewErrorWithStackTrace(problemdetails.Error{
					Type:       "https://delicias-da-lu-service.com/docs/errors/unauthorized",
					Title:      "Unauthorized",
					Detail:     "Authorization header is missing",
					HTTPStatus: http.StatusUnauthorized,
					Instance:   c.Request().RequestURI,
					Severity:   problemdetails.Err,
				})
			}

			// Extract token from "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				logging.WarnEventFromEcho(c).
					Str("auth_header", "invalid_format").
					Msg("authorization header invalid format")
				return problemdetails.NewErrorWithStackTrace(problemdetails.Error{
					Type:       "https://delicias-da-lu-service.com/docs/errors/unauthorized",
					Title:      "Unauthorized",
					Detail:     "Invalid authorization header format",
					HTTPStatus: http.StatusUnauthorized,
					Instance:   c.Request().RequestURI,
					Severity:   problemdetails.Err,
				})
			}

			token := parts[1]

			// Validate token
			userID, err := authUseCase.ValidateToken(c.Request().Context(), token)
			if err != nil {
				logging.WarnEventFromEcho(c).
					Err(err).
					Msg("token validation failed")
				return err
			}

			// Store user ID in context for later use
			c.Set("userID", userID)

			return next(c)
		}
	}
}
