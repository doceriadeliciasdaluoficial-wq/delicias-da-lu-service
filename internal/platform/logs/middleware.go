package logs

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog/log"
)

func NewLoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		var requestBody []byte
		if _, err := c.Request().Body.Read(requestBody); err != nil {
			log.Warn().Err(err).Msg("error reading request body for logging")
			return next(c)
		}

		ctx := c.Request().Context()
		log.Ctx(ctx).Info().
			Str("method", c.Request().Method).
			Str("URL", c.Request().Host).
			Bytes("request_body", requestBody).
			Msg("request received")

		err := next(c)
		if err != nil {

		}

		_, code := echo.ResolveResponseStatus(c.Response(), err)

		log.Ctx(ctx).Info().
			Str("response_code", http.StatusText(code)).
			Msg("response sent")

		return nil
	}
}
