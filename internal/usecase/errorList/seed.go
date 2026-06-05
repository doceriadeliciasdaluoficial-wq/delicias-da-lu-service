package errorList

import (
	"context"
	"net/http"
	"time"

	"delicias-da-lu-service.com/mod/internal/entity/issue"
)

type errorTypeDefinition struct {
	Identifier   string
	Title        string
	Detail       string
	Resolution   string
	Status       int
	SupportEmail string
}

func (ref errorListUseCaseImpl) SeedErrorTypes(ctx context.Context) error {
	for _, def := range defaultErrorTypes() {
		html := buildTypeHTML(typeHTMLData{
			Identifier:   def.Identifier,
			Title:        def.Title,
			Detail:       def.Detail,
			Resolution:   def.Resolution,
			Status:       def.Status,
			SupportEmail: def.SupportEmail,
			UpdatedAt:    time.Now().UTC(),
		})

		err := ref.errorRepository.UpsertErrorType(ctx, def.Identifier, issue.ErrorType{
			Html:      html,
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func defaultErrorTypes() []errorTypeDefinition {
	supportEmail := "doceriadeliciasdaluoficial@gmail.com"
	return []errorTypeDefinition{
		{
			Identifier:   "not-found",
			Title:        "Not Found",
			Detail:       "The requested resource could not be found.",
			Resolution:   "Verify that the resource identifier is correct and exists in the system.",
			Status:       http.StatusNotFound,
			SupportEmail: supportEmail,
		},
		{
			Identifier:   "invalid-type",
			Title:        "Invalid Type",
			Detail:       "The provided type is not supported.",
			Resolution:   "Review the documentation to see the list of supported types and correct your request.",
			Status:       http.StatusBadRequest,
			SupportEmail: supportEmail,
		},
		{
			Identifier:   "invalid-credentials",
			Title:        "Invalid Credentials",
			Detail:       "The username or password is incorrect.",
			Resolution:   "Verify your username and password are correct. If you forgot your password, use the password recovery option.",
			Status:       http.StatusUnauthorized,
			SupportEmail: supportEmail,
		},
		{
			Identifier:   "invalid-token",
			Title:        "Invalid Token",
			Detail:       "The provided token is invalid or expired.",
			Resolution:   "Obtain a new authentication token by logging in with your credentials.",
			Status:       http.StatusUnauthorized,
			SupportEmail: supportEmail,
		},
		{
			Identifier:   "unauthorized",
			Title:        "Unauthorized",
			Detail:       "You are not authorized to access this resource.",
			Resolution:   "Check your permissions or contact support if you believe you should have access to this resource.",
			Status:       http.StatusUnauthorized,
			SupportEmail: supportEmail,
		},
		{
			Identifier:   "type-not-found",
			Title:        "Error Type Not Found",
			Detail:       "No error type found for the provided identifier.",
			Resolution:   "Verify the error type identifier and try again. If the error persists, contact support.",
			Status:       http.StatusNotFound,
			SupportEmail: supportEmail,
		},
		{
			Identifier:   "instance-not-found",
			Title:        "Error Instance Not Found",
			Detail:       "No error instance found for the provided identifier.",
			Resolution:   "Verify the trace ID is correct. Old error instances may be deleted from the system.",
			Status:       http.StatusNotFound,
			SupportEmail: supportEmail,
		},
		{
			Identifier:   "invalidDocumentLenght",
			Title:        "Invalid Document Length",
			Detail:       "The provided document does not have any known length (CPF nor CNPJ).",
			Resolution:   "Verify the sent document and try correcting it. Make sure the parameters are being named correctly and that the document is clean and has the correct length (11 digits for CPF, 14 for CNPJ).",
			Status:       http.StatusBadRequest,
			SupportEmail: supportEmail,
		},
		{
			Identifier:   "invalidFilter",
			Title:        "Invalid Filter",
			Detail:       "The filter query parameter is invalid. Valid values are 'type' and 'instance'.",
			Resolution:   "Use only 'type' or 'instance' as valid filter values, and provide the corresponding identifier parameter.",
			Status:       http.StatusBadRequest,
			SupportEmail: supportEmail,
		},
		{
			Identifier:   "unexpectedUnhandledError",
			Title:        "Unexpected Error",
			Detail:       "An unexpected error occurred. Please contact support.",
			Resolution:   "This is an internal server error. Please share your trace ID with support so we can investigate the issue.",
			Status:       http.StatusInternalServerError,
			SupportEmail: supportEmail,
		},
	}
}
