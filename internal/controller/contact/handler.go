package contact

import (
	"net/http"

	"delicias-da-lu-service.com/mod/internal/entity/contact"
	"delicias-da-lu-service.com/mod/internal/platform/logging"
	contactUC "delicias-da-lu-service.com/mod/internal/usecase/contact"
	"github.com/labstack/echo/v5"
)

type ContactHandler interface {
	Get(c *echo.Context) error
	Update(c *echo.Context) error
}

type contactHandlerImpl struct {
	contactUseCase contactUC.ContactUseCase
}

func NewContactHandler(contactUseCase contactUC.ContactUseCase) ContactHandler {
	return contactHandlerImpl{
		contactUseCase: contactUseCase,
	}
}

func (h contactHandlerImpl) Get(c *echo.Context) error {
	cnt, err := h.contactUseCase.Get(c.Request().Context())
	if err != nil {
		logging.ErrorEventFromEcho(c, err).
			Msg("failed to get contacts")
		return err
	}

	return c.JSON(http.StatusOK, cnt)
}

func (h contactHandlerImpl) Update(c *echo.Context) error {
	var cnt contact.Contact
	if err := c.Bind(&cnt); err != nil {
		logging.WarnEventFromEcho(c).
			Err(err).
			Msg("invalid contact payload")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	logging.DebugEventFromEcho(c).
		Str("email", cnt.Email).
		Str("instagram", cnt.Instagram).
		Msg("contacts update requested")

	updated, err := h.contactUseCase.Update(c.Request().Context(), &cnt)
	if err != nil {
		logging.ErrorEventFromEcho(c, err).
			Msg("failed to update contacts")
		return err
	}

	return c.JSON(http.StatusOK, updated)
}
