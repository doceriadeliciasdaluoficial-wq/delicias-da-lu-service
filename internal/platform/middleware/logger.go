package middleware

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog/log"
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

			userID, _ := c.Get("userID").(string)
			latency := time.Since(start)
			event := log.Info()
			if err != nil {
				event = log.Error().Err(err)
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
				Str("user_id", userID).
				Msg("request completed")

			return err
		}
	}
}
