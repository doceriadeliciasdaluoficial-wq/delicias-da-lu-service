package documentvalidator

import (
	"net/http"

	"delicias-da-lu-service.com/mod/internal/platform/problemdetails"
)

var (
	ErrInvalidDocumentLenght = problemdetails.Error{
		Type:       "https://delicias-da-lu-514609008596.southamerica-east1.run.app/v1/error/filter=type&identifier=invalidDocumentLenght",
		Title:      "Invalid Document Lenght",
		Detail:     "The provided document does not have any known lenght (CPF nor CNPJ)",
		HTTPStatus: http.StatusBadRequest,

		Instance: "https://delicias-da-lu-514609008596.southamerica-east1.run.app/v1/error/invalidDocumentLenght/1",

		Severity: problemdetails.Err,
	}
)
