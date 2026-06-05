package errorList

import (
	"context"
	"net/http"
	"time"

	"delicias-da-lu-service.com/mod/internal/entity/issue"
)

type errorTypeDefinition struct {
	Identifier string
	Title      string
	Detail     string
	Status     int
}

func (ref errorListUseCaseImpl) SeedErrorTypes(ctx context.Context) error {
	for _, def := range defaultErrorTypes() {
		html := buildTypeHTML(typeHTMLData{
			Identifier: def.Identifier,
			Title:      def.Title,
			Detail:     def.Detail,
			Status:     def.Status,
			UpdatedAt:  time.Now().UTC(),
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
	return []errorTypeDefinition{
		{
			Identifier: "not-found",
			Title:      "Not Found",
			Detail:     "The requested resource could not be found.",
			Status:     http.StatusNotFound,
		},
		{
			Identifier: "invalid-type",
			Title:      "Invalid Type",
			Detail:     "The provided type is not supported.",
			Status:     http.StatusBadRequest,
		},
		{
			Identifier: "invalid-credentials",
			Title:      "Invalid Credentials",
			Detail:     "The username or password is incorrect.",
			Status:     http.StatusUnauthorized,
		},
		{
			Identifier: "invalid-token",
			Title:      "Invalid Token",
			Detail:     "The provided token is invalid or expired.",
			Status:     http.StatusUnauthorized,
		},
		{
			Identifier: "unauthorized",
			Title:      "Unauthorized",
			Detail:     "You are not authorized to access this resource.",
			Status:     http.StatusUnauthorized,
		},
		{
			Identifier: "type-not-found",
			Title:      "Error Type Not Found",
			Detail:     "No error type found for the provided identifier.",
			Status:     http.StatusNotFound,
		},
		{
			Identifier: "instance-not-found",
			Title:      "Error Instance Not Found",
			Detail:     "No error instance found for the provided identifier.",
			Status:     http.StatusNotFound,
		},
		{
			Identifier: "invalidDocumentLenght",
			Title:      "Invalid Document Lenght",
			Detail:     "The provided document does not have a known length (CPF or CNPJ).",
			Status:     http.StatusBadRequest,
		},
		{
			Identifier: "invalidFilter",
			Title:      "Invalid Filter",
			Detail:     "The filter query parameter is invalid. Valid values are 'type' and 'instance'.",
			Status:     http.StatusBadRequest,
		},
		{
			Identifier: "unexpectedUnhandledError",
			Title:      "Unexpected Error",
			Detail:     "An unexpected error occurred. Please contact support.",
			Status:     http.StatusInternalServerError,
		},
	}
}
