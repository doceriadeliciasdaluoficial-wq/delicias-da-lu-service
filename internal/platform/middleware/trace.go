package middleware

import (
	"delicias-da-lu-service.com/mod/internal/platform/logging"
	"github.com/labstack/echo/v5"
)

const traceHeader = "X-Trace-Id"

func TraceIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			traceID := c.Request().Header.Get(traceHeader)
			if traceID == "" {
				traceID = logging.GenerateTraceID()
			}
			if traceID != "" {
				c.Set("trace_id", traceID)
				c.Response().Header().Set(traceHeader, traceID)
				ctx := logging.WithTraceID(c.Request().Context(), traceID)
				c.SetRequest(c.Request().WithContext(ctx))
			}

			return next(c)
		}
	}
}

func TraceIDFromEcho(c *echo.Context) string {
	if c == nil {
		return ""
	}
	if value, ok := c.Get("trace_id").(string); ok {
		return value
	}
	return c.Request().Header.Get(traceHeader)
}
