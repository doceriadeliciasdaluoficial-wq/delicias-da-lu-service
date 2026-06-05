package middleware

import (
	"net/http"
	"time"

	"delicias-da-lu-service.com/mod/internal/platform/logging"
	"github.com/labstack/echo/v5"
)

func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			start := time.Now()

			err := next(c)

			status := http.StatusOK
			size := int64(0)
			if resp, unwrapErr := echo.UnwrapResponse(c.Response()); unwrapErr == nil && resp != nil {
				if resp.Status != 0 {
					status = resp.Status
				}
				size = resp.Size
			}

			traceID := TraceIDFromEcho(c)
			latency := time.Since(start)
			event := logging.InfoEventFromEcho(c)
			if err != nil {
				event = logging.ErrorEventFromEcho(c, err)
			}

			event.
				Str("request_method", c.Request().Method).
				Str("request_uri", c.Request().RequestURI).
				Str("path", c.Path()).
				Int("status", status).
				Int64("duration_ms", latency.Milliseconds()).
				Int64("bytes", size).
				Str("remote_ip", c.RealIP()).
				Str("user_agent", c.Request().UserAgent()).
				Str("trace_id", traceID).
				Msg("request completed")

			return err
		}
	}
}
