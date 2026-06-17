package middleware

import (
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func CORSMiddleware() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "https://deliciasdaluoficial.com.br", "https://delicias-da-lu.com", "https://www.delicias-da-lu.com", "https://delicias-da-lu-frontend-514609008596.us-east1.run.app"},
		AllowMethods:     []string{"GET", "HEAD", "PUT", "PATCH", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Requested-With"},
		ExposeHeaders:    []string{"Link"},
		MaxAge:           300,
		AllowCredentials: true,
	})
}
