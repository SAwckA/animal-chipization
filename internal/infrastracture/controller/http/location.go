package http

import (
	"animal-chipization/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

const pointIDParam = "pointId"

type locationUsecase interface {
	CreateLocation(lat, lon float64) (*domain.Location, error)
	GetLocation(locationID int) (*domain.Location, error)
	UpdateLocation(locationID int, location *domain.Location) (*domain.Location, error)
	DeleteLocation(locationID int) error
}

type authMiddleware interface {
	checkAuthHeaderMiddleware(ctx *gin.Context)
	authMiddleware(ctx *gin.Context)
}

type LocationHandler struct {
	usecase    locationUsecase
	middleware authMiddleware
}

func NewLocationHandler(usecase locationUsecase, middleware authMiddleware) *LocationHandler {
	return &LocationHandler{usecase: usecase, middleware: middleware}
}

func (h *LocationHandler) InitRoutes(router *gin.Engine) *gin.Engine {

	locations := router.Group("/locations")
	{
		locations.Use(h.middleware.checkAuthHeaderMiddleware)
		locations.GET("/:pointId",
			errorHandlerWrap(h.getLocationPoint),
		)
		locations.POST("",
			h.middleware.authMiddleware,
			errorHandlerWrap(h.createLocation),
		)
		locations.PUT("/:pointId",
			h.middleware.authMiddleware,
			errorHandlerWrap(h.updateLocation),
		)
		locations.DELETE("/:pointId",
			h.middleware.authMiddleware,
			errorHandlerWrap(h.deleteLocation),
		)
	}

	return router
}

func (h *LocationHandler) getLocationPoint(c *gin.Context) error {
	pointID, err := validateID(c.Copy(), pointIDParam)
	if err != nil {
		return err
	}

	location, err := h.usecase.GetLocation(pointID)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, location)
	return nil
}

func (h *LocationHandler) createLocation(c *gin.Context) error {
	var newLocation *domain.Location

	if err := c.BindJSON(&newLocation); err != nil {
		return NewErrBind(err)
	}

	location, err := h.usecase.CreateLocation(*newLocation.Latitude, *newLocation.Longitude)
	if err != nil {
		return err
	}

	c.JSON(http.StatusCreated, location)
	return nil
}

func (h *LocationHandler) updateLocation(c *gin.Context) error {

	pointID, err := validateID(c.Copy(), pointIDParam)
	if err != nil {
		return err
	}

	var newLocation *domain.Location
	if err := c.BindJSON(&newLocation); err != nil {
		return NewErrBind(err)
	}

	newLocation, err = h.usecase.UpdateLocation(pointID, newLocation)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, newLocation)
	return nil
}

func (h *LocationHandler) deleteLocation(c *gin.Context) error {

	pointID, err := validateID(c.Copy(), pointIDParam)

	if err != nil {
		return err
	}

	err = h.usecase.DeleteLocation(pointID)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, nil)
	return nil
}
