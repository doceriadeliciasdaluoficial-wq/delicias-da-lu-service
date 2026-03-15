package user

import (
	"strings"
	"time"

	"delicias-da-lu-service.com/mod/internal/entity/user"
	"delicias-da-lu-service.com/mod/internal/platform/documentvalidator"
	"delicias-da-lu-service.com/mod/internal/platform/problemdetails"
	"delicias-da-lu-service.com/mod/internal/platform/validator"
)

type userValidator struct {
	Validations []func(user.User) problemdetails.ErrorDetails
}

func NewUserValidator() validator.Validator[user.User] {
	userValidator := userValidator{}
	userValidator.Validations = append(userValidator.Validations, validateName, validateBirthday, validateDocument, validateZipCode)
	return userValidator
}

func (ref userValidator) Validate(i user.User) error {
	validationErrors := problemdetails.Error{
		Type:   "localhost:8080/v1/error/filter=type&identifier=validationErrors",
		Title:  "Validation Errors",
		Detail: "One or more validation errors occurred",

		Instance: "localhost:8080/v1/error/validationErrors/1",

		Severity: problemdetails.Err,
	}
	for _, validation := range ref.Validations {
		if err := validation(i); err != (problemdetails.ErrorDetails{}) {
			validationErrors.Errors = append(validationErrors.Errors, err)
		}
	}
	if len(validationErrors.Errors) > 0 {
		return problemdetails.NewErrorWithStackTrace(validationErrors)
	}
	return nil
}

func validateName(u user.User) problemdetails.ErrorDetails {
	if u.Name == "" {
		return problemdetails.ErrorDetails{
			Pointer: "#/name",
			Detail:  "Name is required",
		}
	}
	if len(u.Name) < 2 {
		return problemdetails.ErrorDetails{
			Pointer: "#/name",
			Detail:  "Name must be at least 2 characters long",
		}
	}
	if len(u.Name) > 100 {
		return problemdetails.ErrorDetails{
			Pointer: "#/name",
			Detail:  "Name must be less than 100 characters long",
		}
	}
	return problemdetails.ErrorDetails{}
}

func validateBirthday(u user.User) problemdetails.ErrorDetails {
	if u.Birthday == "" {
		return problemdetails.ErrorDetails{
			Pointer: "#/birthday",
			Detail:  "Birthday is required",
		}
	}
	parsedTime, err := time.Parse("2006-12-24", u.Birthday)
	if err != nil {
		return problemdetails.ErrorDetails{
			Pointer: "#/birthday",
			Detail:  "Birthday must be in the format YYYY-MM-DD",
		}
	}
	if parsedTime.After(time.Now()) {
		return problemdetails.ErrorDetails{
			Pointer: "#/birthday",
			Detail:  "Birthday cannot be in the future",
		}
	}
	return problemdetails.ErrorDetails{}
}

func validateDocument(u user.User) problemdetails.ErrorDetails {
	if u.Document == "" {
		return problemdetails.ErrorDetails{
			Pointer: "#/document",
			Detail:  "Document is required",
		}
	}
	documentValidator := documentvalidator.NewDocumentValidator(documentvalidator.Document(u.Document))
	if ok, _ := documentValidator.IsValid(documentvalidator.Document(u.Document), 0); !ok {
		return problemdetails.ErrorDetails{
			Pointer: "#/document",
			Detail:  "Document must be a valid CPF or CNPJ",
		}
	}
	return problemdetails.ErrorDetails{}
}

func validateZipCode(u user.User) problemdetails.ErrorDetails {
	if u.ZipCode == "" {
		return problemdetails.ErrorDetails{
			Pointer: "#/zip_code",
			Detail:  "Zip code is required",
		}
	}
	u.ZipCode = strings.ReplaceAll(u.ZipCode, "-", "")
	if len(u.ZipCode) != 8 {
		return problemdetails.ErrorDetails{
			Pointer: "#/zip_code",
			Detail:  "Zip code must be 8 characters long",
		}
	}
	return problemdetails.ErrorDetails{}
}
