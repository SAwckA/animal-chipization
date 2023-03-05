package http

import (
	"animal-chipization/internal/domain"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const visitedPointIDParam = "visitedPointId"

type visitedLocationUsecase interface {
	Create(animalID, pointID int) (*domain.VisitedLocation, error)
	Update(animalID int, location domain.UpdateVisitedLocationDTO) (*domain.VisitedLocation, error)
	Delete(animalID int, locationID int) error
	Search(animalID int, params domain.SearchVisitedLocationDTO) (*[]domain.VisitedLocation, error)
}

type VisitedLocationsHandler struct {
	usecase    visitedLocationUsecase
	middleware authMiddleware
}

func NewVisitedLocationsHandler(usecase visitedLocationUsecase, middleware authMiddleware) *VisitedLocationsHandler {
	return &VisitedLocationsHandler{
		usecase:    usecase,
		middleware: middleware,
	}
}

func (h *VisitedLocationsHandler) InitRoutes(router *gin.Engine) *gin.Engine {

	locations := router.Group(fmt.Sprintf("animals/:%s/locations", animalIDParam))
	{
		locations.Use(h.middleware.checkAuthHeaderMiddleware)
		locations.GET("",
			errorHandlerWrap(h.search),
		)
		locations.POST(fmt.Sprintf("/:%s", pointIDParam),
			h.middleware.authMiddleware,
			errorHandlerWrap(h.create),
		)
		locations.PUT("",
			h.middleware.authMiddleware,
			errorHandlerWrap(h.update),
		)
		locations.DELETE(fmt.Sprintf("/:%s", visitedPointIDParam),
			h.middleware.authMiddleware,
			errorHandlerWrap(h.delete),
		)
	}

	return router
}

func (h *VisitedLocationsHandler) create(c *gin.Context) error {
	animalID, err := validateID(c.Copy(), animalIDParam)
	if err != nil {
		return err
	}

	pointID, err := validateID(c.Copy(), pointIDParam)
	if err != nil {
		return err
	}

	visitedLocation, err := h.usecase.Create(animalID, pointID)
	if err != nil {
		return err
	}

	c.JSON(http.StatusCreated, visitedLocation.Response())
	return nil
}

func (h *VisitedLocationsHandler) update(c *gin.Context) error {

	animalID, err := validateID(c.Copy(), animalIDParam)
	if err != nil {
		return err
	}

	var input domain.UpdateVisitedLocationDTO
	if err = c.BindJSON(&input); err != nil {
		return NewErrBind(err)
	}

	location, err := h.usecase.Update(animalID, input)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, location.Response())
	return nil
}

func (h *VisitedLocationsHandler) delete(c *gin.Context) error {

	animalID, err := validateID(c.Copy(), animalIDParam)
	if err != nil {
		return err
	}

	locationID, err := validateID(c.Copy(), visitedPointIDParam)
	if err != nil {
		return err
	}

	err = h.usecase.Delete(animalID, locationID)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, nil)
	return nil
}

func (h *VisitedLocationsHandler) search(c *gin.Context) error {
	animalID, err := validateID(c.Copy(), animalIDParam)
	if err != nil {
		return err
	}

	var input domain.SearchVisitedLocationDTO
	if err := c.BindQuery(&input); err != nil {
		return NewErrBind(err)
	}

	locations, err := h.usecase.Search(animalID, input)
	if err != nil {
		return err
	}

	resp := make([]map[string]interface{}, 0)
	tmp := *locations
	for _, v := range tmp {
		resp = append(resp, v.Response())
	}

	c.JSON(http.StatusOK, resp)
	return nil
}
