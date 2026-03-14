package system

import (
	"net/http"

	"delicias-da-lu-service.com/mod/internal/platform/problemdetails"
	"delicias-da-lu-service.com/mod/internal/repository/errorFirestore"

	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog/log"
)

type Handler interface {
	Get(e *echo.Context) error
	GetError(e *echo.Context) error
}

type handlerImpl struct {
	errorRepository errorFirestore.ErrorRepository
}

func NewHandler(repository errorFirestore.ErrorRepository) Handler {
	return handlerImpl{
		errorRepository: repository,
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
		Msg("handling error request")

	var content string
	var err error
	switch filterQueryParameter {
	case "type":
		content, err = ref.errorRepository.GetTypeOfErrorByIdentifier(e.Request().Context(), identifierQueryParameter)
		if err != nil {
			log.Error().Err(err).Msg("error fetching error type")
			return err
		}
	case "instance":
		content, err = ref.errorRepository.GetInstanceOfErrorByIdentifier(e.Request().Context(), identifierQueryParameter)
		if err != nil {
			log.Error().Err(err).Msg("error fetching error instance")
			return err
		}
	default:
		log.Warn().Str("filter", filterQueryParameter).Msg("invalid filter query parameter")
		return problemdetails.NewErrorWithStackTrace(problemdetails.Error{
			Type:       "localhost:8080/v1/error?filter=type&identifier=invalidFilter",
			Title:      "Invalid Filter",
			Detail:     "The provided filter query parameter is invalid. Valid values are 'type' and 'instance'",
			HTTPStatus: http.StatusBadRequest,
			Instance:   "localhost:8080/v1/error/invalidFilter/1",
			Severity:   problemdetails.Err,
		})
	}
	return e.HTMLBlob(http.StatusOK, []byte(content))
}
