package problemdetails

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"

	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog/log"
)

const (
	Debug = iota
	Info
	Warn
	Err
	Fatal
)

type Error struct {
	Type       string `json:"type"`
	Title      string `json:"title"`
	Detail     string `json:"detail,omitempty"`
	HTTPStatus int    `json:"status"`

	Errors []ErrorDetails `json:"errors,omitempty"`

	Instance string `json:"instance"`

	StackTrace []uintptr `json:"-"`
	Severity   int       `json:"-"`

	Err error `json:"internal,omitempty"`
}

func (ref Error) Error() string {
	if marshaled, err := json.Marshal(ref); err != nil {
		return "{\"detail\":\"" + ref.Detail + "\", \"status\":" + fmt.Sprint(ref.HTTPStatus) + "}"
	} else {
		return string(marshaled)
	}
}

type ErrorDetails struct {
	Detail  string `json:"detail"`
	Pointer string `json:"pointer"`
}

func GetStackTrace() []uintptr {
	pcs := make([]uintptr, 0)

	runtime.Callers(3, pcs)

	return pcs
}

func NewErrorWithStackTrace(err Error) Error {
	pcs := make([]uintptr, 0)

	runtime.Callers(2, pcs)

	err.StackTrace = pcs

	return err
}

func ErrorHandler(e *echo.Context, err error) {
	var problemdetailsError Error = Error{}
	if !errors.As(err, &problemdetailsError) {
		log.Error().Err(err).Msg("error response handled")
		e.JSON(http.StatusInternalServerError, Error{
			Type:       "unexpectedUnhandledError",
			Title:      "UnexpectedError",
			Detail:     "An untreatable an unrecognized error was found, please contact support. Specific error can be found on '#/internal'",
			HTTPStatus: http.StatusInternalServerError,

			Errors: []ErrorDetails{
				{
					Detail:  "doceriadeliciasdaluoficial@gmail.com",
					Pointer: "email",
				},
			},

			Instance: "localhost:8080/v1/error/unexpectedUnhandledError/",

			Err: err,
		})
		return
	}
	log.Error().Err(problemdetailsError).Msg("error response handled")
	e.JSON(problemdetailsError.HTTPStatus, problemdetailsError)

	return
}
