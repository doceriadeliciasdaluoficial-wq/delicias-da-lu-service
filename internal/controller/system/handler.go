package system

import (
	"net/http"

	"delicias-da-lu-service.com/mod/internal/platform/problemdetails"
	"delicias-da-lu-service.com/mod/internal/usecase/errorList"

	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog/log"
)

type Handler interface {
	Get(e *echo.Context) error
	GetError(e *echo.Context) error
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

	log.Debug().
		Str("filter", filterQueryParameter).
		Str("identifier", identifierQueryParameter).
		Msg("handling GET error request")

	switch filterQueryParameter {
	case "type":
		content, err := ref.errorUsecase.GetTypeOfErrorByIdentifier(e.Request().Context(), identifierQueryParameter)
		if err != nil {
			log.Error().Err(err).Msg("error fetching error type")
			return err
		}
		return e.HTML(http.StatusOK, content)
	case "instance":
		htmlContent, err := ref.errorUsecase.GetInstanceHTMLByIdentifier(e.Request().Context(), identifierQueryParameter)
		if err != nil {
			log.Error().Err(err).Msg("error fetching error instance")
			return err
		}
		return e.HTML(http.StatusOK, htmlContent)
	default:
		log.Warn().Str("filter", filterQueryParameter).Msg("invalid filter query parameter")
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
