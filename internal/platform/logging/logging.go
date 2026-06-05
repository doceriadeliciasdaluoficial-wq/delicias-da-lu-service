package logging

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"runtime"

	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type contextKey string

const traceIDKey contextKey = "trace_id"

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

func TraceIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if value := ctx.Value(traceIDKey); value != nil {
		if traceID, ok := value.(string); ok {
			return traceID
		}
	}
	return ""
}

func GenerateTraceID() string {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return ""
	}
	return hex.EncodeToString(buf)
}

func LoggerFromContext(ctx context.Context) zerolog.Logger {
	logger := log.With()
	if traceID := TraceIDFromContext(ctx); traceID != "" {
		logger = logger.Str("trace_id", traceID)
	}
	return logger.Logger()
}

func LoggerFromEcho(c *echo.Context) zerolog.Logger {
	if c == nil {
		return log.Logger
	}
	logger := LoggerFromContext(c.Request().Context())
	if userID, ok := c.Get("userID").(string); ok && userID != "" {
		logger = logger.With().Str("user_id", userID).Logger()
	}
	return logger
}

func ErrorEvent(ctx context.Context, err error) *zerolog.Event {
	logger := LoggerFromContext(ctx)
	logger = withCaller(logger, 2)
	return logger.Error().Err(err)
}

func WarnEvent(ctx context.Context) *zerolog.Event {
	logger := LoggerFromContext(ctx)
	logger = withCaller(logger, 2)
	return logger.Warn()
}

func InfoEvent(ctx context.Context) *zerolog.Event {
	logger := LoggerFromContext(ctx)
	logger = withCaller(logger, 2)
	return logger.Info()
}

func DebugEvent(ctx context.Context) *zerolog.Event {
	logger := LoggerFromContext(ctx)
	logger = withCaller(logger, 2)
	return logger.Debug()
}

func ErrorEventFromEcho(c *echo.Context, err error) *zerolog.Event {
	logger := LoggerFromEcho(c)
	logger = withCaller(logger, 2)
	return logger.Error().Err(err)
}

func WarnEventFromEcho(c *echo.Context) *zerolog.Event {
	logger := LoggerFromEcho(c)
	logger = withCaller(logger, 2)
	return logger.Warn()
}

func InfoEventFromEcho(c *echo.Context) *zerolog.Event {
	logger := LoggerFromEcho(c)
	logger = withCaller(logger, 2)
	return logger.Info()
}

func DebugEventFromEcho(c *echo.Context) *zerolog.Event {
	logger := LoggerFromEcho(c)
	logger = withCaller(logger, 2)
	return logger.Debug()
}

func withCaller(logger zerolog.Logger, skip int) zerolog.Logger {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return logger
	}
	funcName := ""
	if fn := runtime.FuncForPC(pc); fn != nil {
		funcName = fn.Name()
	}
	return logger.With().
		Str("caller_file", file).
		Int("caller_line", line).
		Str("caller_func", funcName).
		Logger()
}

// ErrorWithContext logs an error with full context: trace_id, user_id, and caller location
func ErrorWithContext(ctx context.Context, err error) *zerolog.Event {
	logger := LoggerFromContext(ctx)
	logger = withCaller(logger, 2)
	return logger.Error().Err(err)
}

// ErrorWithEchoContext logs an error with Echo context: trace_id, user_id, and caller location
func ErrorWithEchoContext(c *echo.Context, err error) *zerolog.Event {
	logger := LoggerFromEcho(c)
	logger = withCaller(logger, 2)
	return logger.Error().Err(err)
}
