package home

import (
	"net/http"

	"delicias-da-lu-service.com/mod/internal/entity/config"
	homeUC "delicias-da-lu-service.com/mod/internal/usecase/home"
	"github.com/labstack/echo/v5"
)

type HomeHandler interface {
	GetFeaturedCakes(c *echo.Context) error
	GetFeaturedCakeByID(c *echo.Context) error
	CreateFeaturedCake(c *echo.Context) error
	UpdateFeaturedCake(c *echo.Context) error
	DeleteFeaturedCake(c *echo.Context) error
}

type homeHandlerImpl struct {
	homeUseCase homeUC.HomeUseCase
}

func NewHomeHandler(homeUseCase homeUC.HomeUseCase) HomeHandler {
	return homeHandlerImpl{
		homeUseCase: homeUseCase,
	}
}

func (h homeHandlerImpl) GetFeaturedCakes(c *echo.Context) error {
	cakes, err := h.homeUseCase.GetFeaturedCakes(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	if cakes == nil {
		cakes = []config.FeaturedCake{}
	}
	return c.JSON(http.StatusOK, cakes)
}

func (h homeHandlerImpl) GetFeaturedCakeByID(c *echo.Context) error {
	id := c.Param("id")

	cake, err := h.homeUseCase.GetFeaturedCakeByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, cake)
}

func (h homeHandlerImpl) CreateFeaturedCake(c *echo.Context) error {
	var cake config.FeaturedCake
	if err := c.Bind(&cake); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	createdCake, err := h.homeUseCase.CreateFeaturedCake(c.Request().Context(), &cake)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, createdCake)
}

func (h homeHandlerImpl) UpdateFeaturedCake(c *echo.Context) error {
	id := c.Param("id")

	var cake config.FeaturedCake
	if err := c.Bind(&cake); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	updatedCake, err := h.homeUseCase.UpdateFeaturedCake(c.Request().Context(), id, &cake)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, updatedCake)
}

func (h homeHandlerImpl) DeleteFeaturedCake(c *echo.Context) error {
	id := c.Param("id")

	err := h.homeUseCase.DeleteFeaturedCake(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}
