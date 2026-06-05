package problemdetails

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

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

type ErrorRecord struct {
	Type        string
	Title       string
	Detail      string
	Status      int
	RequestURL  string
	Method      string
	UserAgent   string
	RequestBody string
	TraceID     string
	BaseURL     string
	OccurredAt  time.Time
}

type ErrorRecorder interface {
	Record(ctx context.Context, record ErrorRecord) (string, error)
}

var errorRecorder ErrorRecorder

func SetErrorRecorder(recorder ErrorRecorder) {
	errorRecorder = recorder
}

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
	problemdetailsError := Error{}
	if !errors.As(err, &problemdetailsError) {
		log.Error().Err(err).Msg("error response handled")
		problemdetailsError = Error{
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
			Instance: "https://delicias-da-lu-514609008596.southamerica-east1.run.app/v1/error/unexpectedUnhandledError/",
			Err:      err,
		}
	} else {
		log.Error().Err(problemdetailsError).Msg("error response handled")
	}

	if problemdetailsError.HTTPStatus == 0 {
		problemdetailsError.HTTPStatus = http.StatusInternalServerError
	}

	scheme := e.Scheme()
	host := e.Request().Host
	baseURL := strings.TrimRight(fmt.Sprintf("%s://%s", scheme, host), "/")
	if host == "" || baseURL == "://" {
		baseURL = "http://localhost:8080"
	}

	if errorRecorder != nil {
		record := ErrorRecord{
			Type:        problemdetailsError.Type,
			Title:       problemdetailsError.Title,
			Detail:      problemdetailsError.Detail,
			Status:      problemdetailsError.HTTPStatus,
			RequestURL:  e.Request().RequestURI,
			Method:      e.Request().Method,
			UserAgent:   e.Request().UserAgent(),
			BaseURL:     baseURL,
			OccurredAt:  time.Now().UTC(),
			RequestBody: "",
			TraceID:     "",
		}

		instanceID, recordErr := errorRecorder.Record(e.Request().Context(), record)
		if recordErr != nil {
			log.Error().Err(recordErr).Msg("failed to record error instance")
		} else if instanceID != "" && (problemdetailsError.Instance == "" || !strings.Contains(problemdetailsError.Instance, "filter=instance")) {
			problemdetailsError.Instance = fmt.Sprintf("%s/v1/error?filter=instance&identifier=%s", baseURL, instanceID)
		}
	}

	e.JSON(problemdetailsError.HTTPStatus, problemdetailsError)
}
