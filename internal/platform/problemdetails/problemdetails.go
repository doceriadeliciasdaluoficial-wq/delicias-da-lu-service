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

	"delicias-da-lu-service.com/mod/internal/platform/logging"
	"github.com/labstack/echo/v5"
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

	OriginFile string `json:"-"`
	OriginLine int    `json:"-"`
	OriginFunc string `json:"-"`

	CallerFile string `json:"-"`
	CallerLine int    `json:"-"`
	CallerFunc string `json:"-"`

	StackTrace []uintptr `json:"-"`
	Severity   int       `json:"-"`

	Err error `json:"-"`
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

type frameLocation struct {
	File string
	Line int
	Func string
}

func resolveFrame(pc uintptr) frameLocation {
	if pc == 0 {
		return frameLocation{}
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return frameLocation{}
	}

	file, line := fn.FileLine(pc)
	return frameLocation{
		File: file,
		Line: line,
		Func: fn.Name(),
	}
}

func resolveCaller(skip int) frameLocation {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return frameLocation{}
	}
	name := ""
	if fn := runtime.FuncForPC(pc); fn != nil {
		name = fn.Name()
	}
	return frameLocation{
		File: file,
		Line: line,
		Func: name,
	}
}

func originFromStack(pcs []uintptr) frameLocation {
	if len(pcs) == 0 {
		return frameLocation{}
	}
	return resolveFrame(pcs[0])
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
	caller := resolveCaller(1)

	if !errors.As(err, &problemdetailsError) {
		origin := resolveCaller(2)
		logging.ErrorEventFromEcho(e, err).
			Str("origin_file", origin.File).
			Int("origin_line", origin.Line).
			Str("origin_func", origin.Func).
			Str("caller_file", caller.File).
			Int("caller_line", caller.Line).
			Str("caller_func", caller.Func).
			Str("request_method", e.Request().Method).
			Str("request_uri", e.Request().RequestURI).
			Msg("error response handled")
		problemdetailsError = Error{
			Type:       "unexpectedUnhandledError",
			Title:      "UnexpectedError",
			Detail:     "An untreatable an unrecognized error was found, please contact support. Specific error can be found on '#/internal'",
			HTTPStatus: http.StatusInternalServerError,
			OriginFile: origin.File,
			OriginLine: origin.Line,
			OriginFunc: origin.Func,
			CallerFile: caller.File,
			CallerLine: caller.Line,
			CallerFunc: caller.Func,
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
		origin := originFromStack(problemdetailsError.StackTrace)
		if origin.File == "" {
			origin = resolveCaller(2)
		}
		logging.ErrorEventFromEcho(e, problemdetailsError).
			Str("error_type", problemdetailsError.Type).
			Str("error_title", problemdetailsError.Title).
			Str("error_instance", problemdetailsError.Instance).
			Int("status", problemdetailsError.HTTPStatus).
			Str("origin_file", origin.File).
			Int("origin_line", origin.Line).
			Str("origin_func", origin.Func).
			Str("caller_file", caller.File).
			Int("caller_line", caller.Line).
			Str("caller_func", caller.Func).
			Str("request_method", e.Request().Method).
			Str("request_uri", e.Request().RequestURI).
			Msg("error response handled")
		problemdetailsError.OriginFile = origin.File
		problemdetailsError.OriginLine = origin.Line
		problemdetailsError.OriginFunc = origin.Func
		problemdetailsError.CallerFile = caller.File
		problemdetailsError.CallerLine = caller.Line
		problemdetailsError.CallerFunc = caller.Func
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
			TraceID:     logging.TraceIDFromContext(e.Request().Context()),
		}

		instanceID, recordErr := errorRecorder.Record(e.Request().Context(), record)
		if recordErr != nil {
			logging.ErrorEventFromEcho(e, recordErr).
				Msg("failed to record error instance")
		} else if instanceID != "" && (problemdetailsError.Instance == "" || !strings.Contains(problemdetailsError.Instance, "filter=instance")) {
			problemdetailsError.Instance = fmt.Sprintf("%s/v1/error?filter=instance&identifier=%s", baseURL, instanceID)
		}
	}

	e.JSON(problemdetailsError.HTTPStatus, problemdetailsError)
}
