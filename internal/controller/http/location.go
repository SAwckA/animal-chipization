package http

import (
	"animal-chipization/internal/domain"
	"animal-chipization/internal/errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type locationUsecase interface {
	CreateLocation(lat, lon float64) (*domain.Location, error)
	GetLocation(locationID int) *domain.Location
	UpdateLocation(location *domain.Location) error
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
		locations.POST("/",
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

	location := h.usecase.GetLocation(pointID)

	if location == nil {
		newErrorResponse(c, http.StatusNotFound, "Location not found", nil)
		return
	}

	c.JSON(http.StatusOK, location)
}

func (h *LocationHandler) createLocation(c *gin.Context) {
	var newLocation *domain.Location

	if err := c.BindJSON(&newLocation); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid request data", nil)
		return
	}

	if newLocation.Latitude < -90 || newLocation.Latitude > 90 || newLocation.Longitude < -180 || newLocation.Longitude > 180 {
		newErrorResponse(c, http.StatusBadRequest, "Invalid request data", nil)
		return
	}

	location, err := h.usecase.CreateLocation(newLocation.Latitude, newLocation.Longitude)

	if err != nil {
		newErrorResponse(c, http.StatusConflict, "Location with this point already exists", nil)
		return
	}

	c.JSON(http.StatusOK, location)
}

func (h *LocationHandler) updateLocation(c *gin.Context) {

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

	var newLocation *domain.Location
	if err := c.BindJSON(&newLocation); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid request data", nil)
		return
	}

	newLocation.ID = pointID

	if newLocation.Latitude < -90 || newLocation.Latitude > 90 || newLocation.Longitude < -180 || newLocation.Longitude > 180 || newLocation.Latitude == 0 || newLocation.Longitude == 0 {
		newErrorResponse(c, http.StatusBadRequest, "Invalid request data", nil)
		return
	}

	err = h.usecase.UpdateLocation(newLocation)

	if err == errors.ErrNotFound {
		newErrorResponse(c, http.StatusNotFound, "Location at this id not found", nil)
		return
	}

	if err == errors.ErrAlreadyExist {
		newErrorResponse(c, http.StatusConflict, "Location with this latitude and longitude already exist", nil)
		return
	}

	c.JSON(http.StatusOK, newLocation)
}

func (h *LocationHandler) deleteLocation(c *gin.Context) {

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

	err = h.usecase.DeleteLocation(pointID)

	if err == errors.ErrNotFound {
		newErrorResponse(c, http.StatusNotFound, "Locaton not found", nil)
		return
	}

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Location linked with animal", nil)
		return
	}
}
