package errorList

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"time"

	"delicias-da-lu-service.com/mod/internal/entity/issue"
	"delicias-da-lu-service.com/mod/internal/platform/problemdetails"
	"delicias-da-lu-service.com/mod/internal/repository/errorFirestore"
)

type errorRecorderImpl struct {
	errorRepository errorFirestore.ErrorRepository
}

func NewErrorRecorder(repository errorFirestore.ErrorRepository) problemdetails.ErrorRecorder {
	return errorRecorderImpl{
		errorRepository: repository,
	}
}

func (ref errorRecorderImpl) Record(ctx context.Context, record problemdetails.ErrorRecord) (string, error) {
	typeIdentifier := extractTypeIdentifier(record.Type, record.Title)
	baseURL := strings.TrimRight(record.BaseURL, "/")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	if record.OccurredAt.IsZero() {
		record.OccurredAt = time.Now().UTC()
	}

	traceID := record.TraceID
	if traceID == "" {
		traceID = generateTraceID()
	}

	typeHTML := buildTypeHTML(typeHTMLData{
		Identifier:   typeIdentifier,
		Title:        record.Title,
		Detail:       record.Detail,
		Resolution:   "Please contact support with your trace ID if you need assistance.",
		Status:       record.Status,
		BaseURL:      baseURL,
		SupportEmail: "doceriadeliciasdaluoficial@gmail.com",
		UpdatedAt:    time.Now().UTC(),
	})

	if err := ref.errorRepository.UpsertErrorType(ctx, typeIdentifier, issue.ErrorType{
		Html:      typeHTML,
		UpdatedAt: time.Now().UTC(),
	}); err != nil {
		return "", err
	}

	instanceHTML := buildInstanceHTML(instanceHTMLData{
		Title:      record.Title,
		Detail:     record.Detail,
		Type:       record.Type,
		Status:     record.Status,
		RequestURL: record.RequestURL,
		Method:     record.Method,
		UserAgent:  record.UserAgent,
		TraceID:    traceID,
		OccurredAt: record.OccurredAt,
		TypeLink:   fmt.Sprintf("%s/v1/error?filter=type&identifier=%s", baseURL, typeIdentifier),
	})

	instance := issue.ErrorInstance{
		RequestBody: record.RequestBody,
		RequestDate: record.OccurredAt,
		RequestURL:  record.RequestURL,
		Status:      record.Status,
		Title:       record.Title,
		TraceID:     traceID,
		Type:        record.Type,
		UserAgent:   record.UserAgent,
		Html:        instanceHTML,
	}

	if err := ref.errorRepository.CreateErrorInstance(ctx, traceID, instance); err != nil {
		return "", err
	}

	return traceID, nil
}

func extractTypeIdentifier(typeValue, title string) string {
	if typeValue != "" {
		if parsed, err := url.Parse(typeValue); err == nil {
			if identifier := parsed.Query().Get("identifier"); identifier != "" {
				return identifier
			}
			if strings.Trim(parsed.Path, "/") != "" {
				parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
				return parts[len(parts)-1]
			}
		}
	}

	if title != "" {
		return strings.ToLower(strings.ReplaceAll(title, " ", "-"))
	}

	return "unexpectedUnhandledError"
}

func generateTraceID() string {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(buf)
}
