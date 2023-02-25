package http

import (
	"animal-chipization/internal/domain"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const visitedPointIDParam = "visitedPointId"

type visitedLocationUsecase interface {
	Create(animalID, pointID int) (*domain.VisitedLocation, error)
	Update(animalID int, location domain.UpdateVisitedLocationDTO) (*domain.VisitedLocation, error)
	Delete(animalID int, locatoinID int) error
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
		locations.GET("",
			h.middleware.ckeckAuthHeaderMiddleware,
			h.search,
		)
		locations.POST(fmt.Sprintf("/:%s", pointIDParam),
			h.middleware.authMiddleware,
			h.create,
		)
		locations.PUT("",
			h.middleware.authMiddleware,
			h.update,
		)
		locations.DELETE(fmt.Sprintf("/:%s", visitedPointIDParam),
			h.middleware.authMiddleware,
			h.delete,
		)
	}

	return router
}

func (h *VisitedLocationsHandler) create(c *gin.Context) {
	animalID, err := validateID(c.Copy(), animalIDParam)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	pointID, err := validateID(c.Copy(), pointIDParam)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	visitedLocation, err := h.usecase.Create(animalID, pointID)

	switch errors.Unwrap(err) {
	case domain.ErrNotFound:
		notFoundResponse(c, err.Error())

	case domain.ErrInvalidInput, domain.ErrAlreadyExist:
		badRequest(c, err.Error())

	case domain.ErrUnknown:
		unreachableError(c, err)

	case nil:
		c.JSON(http.StatusCreated, visitedLocation.Response())

	default:
		unreachableError(c, err)
	}

}

func (h *VisitedLocationsHandler) update(c *gin.Context) {

	animalID, err := validateID(c.Copy(), animalIDParam)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var input domain.UpdateVisitedLocationDTO
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	location, err := h.usecase.Update(animalID, input)

	switch errors.Unwrap(err) {
	case domain.ErrNotFound:
		notFoundResponse(c, err.Error())

	case domain.ErrInvalidInput, domain.ErrAlreadyExist:
		badRequest(c, err.Error())

	case domain.ErrUnknown:
		unreachableError(c, err)

	case nil:
		c.JSON(http.StatusOK, location.Response())

	default:
		unreachableError(c, err)
	}

}

func (h *VisitedLocationsHandler) delete(c *gin.Context) {

	animalID, err := validateID(c.Copy(), animalIDParam)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	locationID, err := validateID(c.Copy(), visitedPointIDParam)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err = h.usecase.Delete(animalID, locationID)

	switch errors.Unwrap(err) {
	case domain.ErrNotFound:
		notFoundResponse(c, err.Error())

	case domain.ErrInvalidInput, domain.ErrAlreadyExist:
		badRequest(c, err.Error())

	case domain.ErrUnknown:
		unreachableError(c, err)

	case nil:
		c.JSON(http.StatusOK, nil)

	default:
		unreachableError(c, err)
	}
}

func (h *VisitedLocationsHandler) search(c *gin.Context) {
	animalID, err := validateID(c.Copy(), animalIDParam)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var input domain.SearchVisitedLocationDTO
	if err := c.BindQuery(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	locations, err := h.usecase.Search(animalID, input)

	switch errors.Unwrap(err) {
	case domain.ErrNotFound:
		notFoundResponse(c, err.Error())

	case domain.ErrInvalidInput, domain.ErrAlreadyExist:
		badRequest(c, err.Error())

	case domain.ErrUnknown:
		unreachableError(c, err)

	case nil:
		resp := make([]map[string]interface{}, 0)
		tmp := *locations
		for _, v := range tmp {
			resp = append(resp, v.Response())
		}

		c.JSON(http.StatusOK, resp)

	default:
		unreachableError(c, err)
	}

}
