package system

import (
	"net/http"

	"delicias-da-lu-service.com/mod/internal/platform/logging"
	"delicias-da-lu-service.com/mod/internal/platform/problemdetails"
	"delicias-da-lu-service.com/mod/internal/usecase/errorList"

	"github.com/labstack/echo/v5"
)

type Handler interface {
	Get(e *echo.Context) error
	GetError(e *echo.Context) error
	CreateErrorTypes(e *echo.Context) error
	DeleteErrorType(e *echo.Context) error
}

type handlerImpl struct {
	errorUsecase errorList.ErrorListUseCase
}

func NewHandler(usecase errorList.ErrorListUseCase) Handler {
	return handlerImpl{
		errorUsecase: usecase,
	}
}

func (ref handlerImpl) Get(e *echo.Context) error {
	return e.JSON(http.StatusOK, map[string]string{
		"admin":     "Docerias da Lu",
		"email":     "doceriadeliciasdaluoficial@gmail.com",
		"instagram": "delicias.lu.oficial",
	})
}

func (ref handlerImpl) GetError(e *echo.Context) error {
	filterQueryParameter := e.QueryParam("filter")
	identifierQueryParameter := e.QueryParam("identifier")

	logging.DebugEventFromEcho(e).
		Str("filter", filterQueryParameter).
		Str("identifier", identifierQueryParameter).
		Msg("error documentation requested")

	switch filterQueryParameter {
	case "type":
		content, err := ref.errorUsecase.GetTypeOfErrorByIdentifier(e.Request().Context(), identifierQueryParameter)
		if err != nil {
			logging.ErrorEventFromEcho(e, err).
				Str("identifier", identifierQueryParameter).
				Msg("failed to fetch error type")
			return err
		}
		return e.HTML(http.StatusOK, content)
	case "instance":
		htmlContent, err := ref.errorUsecase.GetInstanceHTMLByIdentifier(e.Request().Context(), identifierQueryParameter)
		if err != nil {
			logging.ErrorEventFromEcho(e, err).
				Str("identifier", identifierQueryParameter).
				Msg("failed to fetch error instance")
			return err
		}
		return e.HTML(http.StatusOK, htmlContent)
	default:
		logging.WarnEventFromEcho(e).
			Str("filter", filterQueryParameter).
			Msg("invalid filter query parameter")
		return problemdetails.NewErrorWithStackTrace(problemdetails.Error{
			Type:       "https://delicias-da-lu-514609008596.southamerica-east1.run.app/v1/error?filter=type&identifier=invalidFilter",
			Title:      "Invalid Filter",
			Detail:     "The provided filter query parameter is invalid. Valid values are 'type' and 'instance'",
			HTTPStatus: http.StatusBadRequest,
			Instance:   "https://delicias-da-lu-514609008596.southamerica-east1.run.app/v1/error/invalidFilter/",
			Severity:   problemdetails.Err,
		})
	}
}

func (ref handlerImpl) CreateErrorTypes(e *echo.Context) error {
	logging.DebugEventFromEcho(e).
		Msg("create error types requested")

	var requests []errorList.CreateErrorTypeRequest
	if err := e.Bind(&requests); err != nil {
		logging.WarnEventFromEcho(e).
			Err(err).
			Msg("invalid request payload for create error types")
		return problemdetails.NewErrorWithStackTrace(problemdetails.Error{
			Type:       "https://delicias-da-lu-514609008596.southamerica-east1.run.app/v1/error?filter=type&identifier=invalid-type",
			Title:      "Invalid Request",
			Detail:     "The request payload is invalid",
			HTTPStatus: http.StatusBadRequest,
			Severity:   problemdetails.Err,
		})
	}

	if err := ref.errorUsecase.CreateErrorTypesFromList(e.Request().Context(), requests); err != nil {
		logging.ErrorEventFromEcho(e, err).
			Msg("failed to create error types")
		return err
	}

	logging.InfoEventFromEcho(e).
		Int("count", len(requests)).
		Msg("error types created successfully")

	return e.JSON(http.StatusCreated, map[string]interface{}{
		"message": "error types created successfully",
		"count":   len(requests),
	})
}

func (ref handlerImpl) DeleteErrorType(e *echo.Context) error {
	identifier := e.Param("identifier")

	logging.DebugEventFromEcho(e).
		Str("identifier", identifier).
		Msg("delete error type requested")

	if err := ref.errorUsecase.DeleteErrorType(e.Request().Context(), identifier); err != nil {
		logging.ErrorEventFromEcho(e, err).
			Str("identifier", identifier).
			Msg("failed to delete error type")
		return err
	}

	logging.InfoEventFromEcho(e).
		Str("identifier", identifier).
		Msg("error type deleted successfully")

	return e.JSON(http.StatusOK, map[string]interface{}{
		"message":    "error type deleted successfully",
		"identifier": identifier,
	})
}
