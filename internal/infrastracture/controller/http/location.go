package http

import (
	"animal-chipization/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

const pointIDParam = "pointId"

type locationUsecase interface {
	Location(id int) (*domain.Location, error)
	Create(lat, lon float64) (*domain.Location, error)
	Update(id int, location *domain.Location) (*domain.Location, error)
	Delete(id int) error
}

type LocationHandler struct {
	usecase locationUsecase
	auth    authMiddleware
}

func NewLocationHandler(usecase locationUsecase, auth authMiddleware) *LocationHandler {
	return &LocationHandler{
		usecase: usecase,
		auth:    auth,
	}
}

func (h *LocationHandler) InitRoutes(router *gin.Engine) *gin.Engine {

	locations := router.Group("/locations")
	{
		locations.Use(h.auth.checkAuthHeaderMiddleware)
		locations.GET("/:pointId",
			errorHandlerWrap(h.locationPoint),
		)
		locations.POST("",
			h.auth.authMiddleware,
			errorHandlerWrap(h.create),
		)
		locations.PUT("/:pointId",
			h.auth.authMiddleware,
			errorHandlerWrap(h.update),
		)
		locations.DELETE("/:pointId",
			h.auth.authMiddleware,
			errorHandlerWrap(h.delete),
		)
	}

	return router
}

func (h *LocationHandler) locationPoint(c *gin.Context) error {
	pointID, err := ParamID(c.Copy(), pointIDParam)
	if err != nil {
		return err
	}

	location, err := h.usecase.Location(pointID)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, location.Map())
	return nil
}

func (h *LocationHandler) create(c *gin.Context) error {
	var newLocation *domain.Location
	if err := c.BindJSON(&newLocation); err != nil {
		return NewErrBind(err)
	}

	location, err := h.usecase.Create(*newLocation.Latitude, *newLocation.Longitude)
	if err != nil {
		return err
	}

	c.JSON(http.StatusCreated, location.Map())
	return nil
}

func (h *LocationHandler) update(c *gin.Context) error {
	pointID, err := ParamID(c.Copy(), pointIDParam)
	if err != nil {
		return err
	}

	var newLocation *domain.Location
	if err = c.BindJSON(&newLocation); err != nil {
		return NewErrBind(err)
	}

	newLocation, err = h.usecase.Update(pointID, newLocation)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, newLocation.Map())
	return nil
}

func (h *LocationHandler) delete(c *gin.Context) error {
	pointID, err := ParamID(c.Copy(), pointIDParam)
	if err != nil {
		return err
	}

	err = h.usecase.Delete(pointID)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, nil)
	return nil
}
