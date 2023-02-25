package http

import (
	"animal-chipization/internal/domain"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const pointIDParam = "pointId"

type locationUsecase interface {
	CreateLocation(lat, lon float64) (*domain.Location, error)
	GetLocation(locationID int) (*domain.Location, error)
	UpdateLocation(locationID int, location *domain.Location) (*domain.Location, error)
	DeleteLocation(locationID int) error
}

type authMiddleware interface {
	ckeckAuthHeaderMiddleware(ctx *gin.Context)
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
		locations.GET("/:pointId",
			h.middleware.ckeckAuthHeaderMiddleware,
			h.getLocationPoint,
		)
		locations.POST("",
			h.middleware.authMiddleware,
			h.createLocation,
		)
		locations.PUT("/:pointId",
			h.middleware.authMiddleware,
			h.updateLocation,
		)
		locations.DELETE("/:pointId",
			h.middleware.authMiddleware,
			h.deleteLocation,
		)
	}

	return router
}

func (h *LocationHandler) getLocationPoint(c *gin.Context) {

	pointIDString := c.Param("pointId")

	if pointIDString == "null" || pointIDString == "" {
		newErrorResponse(c, http.StatusBadRequest, "Invalid pointId", nil)
		return
	}

	pointID, err := strconv.Atoi(pointIDString)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid pointId", nil)
		return
	}

	if pointID <= 0 {
		newErrorResponse(c, http.StatusBadRequest, "Invalid pointId", nil)
		return
	}

	location, err := h.usecase.GetLocation(pointID)

	switch errors.Unwrap(err) {
	case domain.ErrNotFound:
		notFoundResponse(c, err.Error())

	case domain.ErrAlreadyExist:
		conflictResponse(c, err.Error())

	case domain.ErrLinked:
		badRequest(c, err.Error())

	default:
		c.JSON(http.StatusOK, location)
	}
}

func (h *LocationHandler) createLocation(c *gin.Context) {
	var newLocation *domain.Location

	if err := c.BindJSON(&newLocation); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	location, err := h.usecase.CreateLocation(*newLocation.Latitude, *newLocation.Longitude)

	switch errors.Unwrap(err) {
	case domain.ErrNotFound:
		notFoundResponse(c, err.Error())

	case domain.ErrAlreadyExist:
		conflictResponse(c, err.Error())

	case domain.ErrLinked:
		badRequest(c, err.Error())

	default:
		c.JSON(http.StatusCreated, location)
	}
}

func (h *LocationHandler) updateLocation(c *gin.Context) {

	pointID, err := validateID(c.Copy(), pointIDParam)
	if err != nil {
		badRequest(c, err.Error())
		return
	}

	var newLocation *domain.Location
	if err := c.BindJSON(&newLocation); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid request data", nil)
		return
	}

	newLocation, err = h.usecase.UpdateLocation(pointID, newLocation)

	switch errors.Unwrap(err) {
	case domain.ErrNotFound:
		notFoundResponse(c, err.Error())

	case domain.ErrAlreadyExist:
		conflictResponse(c, err.Error())

	case domain.ErrLinked:
		badRequest(c, err.Error())

	default:
		c.JSON(http.StatusOK, newLocation)
	}
}

func (h *LocationHandler) deleteLocation(c *gin.Context) {

	pointID, err := validateID(c.Copy(), pointIDParam)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid pointId", nil)
		return
	}

	err = h.usecase.DeleteLocation(pointID)

	switch errors.Unwrap(err) {
	case domain.ErrNotFound:
		notFoundResponse(c, err.Error())

	case domain.ErrAlreadyExist:
		conflictResponse(c, err.Error())

	case domain.ErrLinked:
		badRequest(c, err.Error())

	default:
		c.JSON(http.StatusOK, nil)
	}

}
