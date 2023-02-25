package http

import (
	"animal-chipization/internal/domain"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const animalIDParam = "animalId"

type animalUsecase interface {
	GetAnimal(animalID int) (*domain.Animal, error)
	SearchAnimal(params *domain.AnimalSearchParams) (*[]domain.Animal, error)
	CreateAnimal(params domain.AnimalCreateParams) (*domain.Animal, error)
	UpdateAnimal(animalID int, params domain.AnimalUpdateParams) (*domain.Animal, error)
	DeleteAnimal(animalID int) error

	AddAnimalType(animalID, typeID int) (*domain.Animal, error)
	EditAnimalType(animalID int, params domain.AnimalEditTypeParams) (*domain.Animal, error)
	DeleteAnimalType(animalID, typeID int) (*domain.Animal, error)
}

type AnimalHandler struct {
	usecase    animalUsecase
	middleware authMiddleware
}

func NewAnimalHandler(usecase animalUsecase, middleware authMiddleware) *AnimalHandler {
	return &AnimalHandler{usecase: usecase, middleware: middleware}
}

func (h *AnimalHandler) InitRoutes(router *gin.Engine) *gin.Engine {

	animal := router.Group("animals")
	{
		animal.GET(fmt.Sprintf("/:%s", animalIDParam),
			h.middleware.ckeckAuthHeaderMiddleware,
			h.getAnimal,
		)

		animal.GET("/search",
			h.middleware.ckeckAuthHeaderMiddleware,
			h.searchAnimal,
		)

		animal.POST("",
			h.middleware.authMiddleware,
			h.createAnimal,
		)

		animal.PUT(fmt.Sprintf("/:%s", animalIDParam),
			h.middleware.authMiddleware,
			h.updateAnimal,
		)

		animal.DELETE(fmt.Sprintf("/:%s", animalIDParam),
			h.middleware.authMiddleware,
			h.deleteAnimal,
		)

		types := animal.Group(fmt.Sprintf(":%s/types", animalIDParam))
		{
			types.POST(fmt.Sprintf(":%s", typeParam),
				h.middleware.authMiddleware,
				h.addAnimalType,
			)

			types.PUT("",
				h.middleware.authMiddleware,
				h.editAnimalType,
			)

			types.DELETE(fmt.Sprintf(":%s", typeParam),
				h.middleware.authMiddleware,
				h.deleteAnimalType,
			)
		}
	}

	return router
}

func (h *AnimalHandler) getAnimal(c *gin.Context) {

	animalID, err := validateID(c.Copy(), animalIDParam)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	animal, err := h.usecase.GetAnimal(animalID)

	switch errors.Unwrap(err) {

	case domain.ErrInvalidInput:
		badRequest(c, err.Error())

	case domain.ErrNotFound:
		notFoundResponse(c, err.Error())

	case domain.ErrConflict:
		conflictResponse(c, err.Error())

	case domain.ErrUnknown:
		unreachableError(c, err)

	default:
		c.JSON(http.StatusOK, animal.Response())

	}
}

func (h *AnimalHandler) searchAnimal(c *gin.Context) {

	input, err := parseSearchAnimalParams(c)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	animalsList, err := h.usecase.SearchAnimal(input)

	switch errors.Unwrap(err) {

	case domain.ErrInvalidInput:
		badRequest(c, err.Error())

	case domain.ErrNotFound:
		notFoundResponse(c, err.Error())

	case domain.ErrConflict:
		conflictResponse(c, err.Error())

	case domain.ErrUnknown:
		unreachableError(c, err)

	default:
		if animalsList == nil {
			c.JSON(http.StatusOK, nil)
			return
		}
		resp := make([]map[string]interface{}, 0)

		for _, v := range *animalsList {
			resp = append(resp, v.Response())
		}
		c.JSON(http.StatusOK, resp)
	}

}

func (h *AnimalHandler) createAnimal(c *gin.Context) {

	var input domain.AnimalCreateParams
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid body", nil)
		return
	}

	animal, err := h.usecase.CreateAnimal(input)

	switch errors.Unwrap(err) {

	case domain.ErrInvalidInput:
		badRequest(c, err.Error())

	case domain.ErrNotFound:
		notFoundResponse(c, err.Error())

	case domain.ErrConflict:
		conflictResponse(c, err.Error())

	case domain.ErrUnknown:
		unreachableError(c, err)

	default:
		c.JSON(http.StatusCreated, animal.Response())
	}

}

func (h *AnimalHandler) updateAnimal(c *gin.Context) {

	var input domain.AnimalUpdateParams

	animalID, err := validateID(c.Copy(), animalIDParam)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	if err = c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid body", nil)
		return
	}

	animal, err := h.usecase.UpdateAnimal(animalID, input)

	switch errors.Unwrap(err) {

	case domain.ErrInvalidInput:
		badRequest(c, err.Error())

	case domain.ErrNotFound:
		notFoundResponse(c, err.Error())

	case domain.ErrConflict:
		conflictResponse(c, err.Error())

	case domain.ErrUnknown:
		unreachableError(c, err)

	default:
		c.JSON(http.StatusOK, animal.Response())
	}
}

func (h *AnimalHandler) deleteAnimal(c *gin.Context) {

	animalID, err := validateID(c.Copy(), animalIDParam)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err = h.usecase.DeleteAnimal(animalID)

	switch errors.Unwrap(err) {

	case domain.ErrInvalidInput:
		badRequest(c, err.Error())

	case domain.ErrNotFound:
		notFoundResponse(c, err.Error())

	case domain.ErrConflict:
		conflictResponse(c, err.Error())

	case domain.ErrUnknown:
		unreachableError(c, err)

	default:
		c.JSON(http.StatusOK, nil)
	}
}

func (h *AnimalHandler) addAnimalType(c *gin.Context) {

	animalID, err := validateID(c.Copy(), animalIDParam)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	animalTypeID, err := validateID(c.Copy(), typeParam)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	animal, err := h.usecase.AddAnimalType(animalID, animalTypeID)

	switch errors.Unwrap(err) {

	case domain.ErrInvalidInput:
		badRequest(c, err.Error())

	case domain.ErrNotFound:
		notFoundResponse(c, err.Error())

	case domain.ErrConflict:
		conflictResponse(c, err.Error())

	case domain.ErrUnknown:
		unreachableError(c, err)

	default:
		c.JSON(http.StatusCreated, animal.Response())
	}
}

func (h *AnimalHandler) editAnimalType(c *gin.Context) {
	var input domain.AnimalEditTypeParams

	animalID, err := validateID(c.Copy(), animalIDParam)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid body", nil)
		return
	}

	animal, err := h.usecase.EditAnimalType(animalID, input)

	switch errors.Unwrap(err) {

	case domain.ErrInvalidInput:
		badRequest(c, err.Error())

	case domain.ErrNotFound:
		notFoundResponse(c, err.Error())

	case domain.ErrConflict:
		conflictResponse(c, err.Error())

	case domain.ErrUnknown:
		unreachableError(c, err)

	default:
		c.JSON(http.StatusOK, animal.Response())
	}

}

func (h *AnimalHandler) deleteAnimalType(c *gin.Context) {

	animalID, err := validateID(c.Copy(), animalIDParam)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	typeID, err := validateID(c.Copy(), typeParam)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	animal, err := h.usecase.DeleteAnimalType(animalID, typeID)

	switch errors.Unwrap(err) {

	case domain.ErrInvalidInput:
		badRequest(c, err.Error())

	case domain.ErrNotFound:
		notFoundResponse(c, err.Error())

	case domain.ErrConflict:
		conflictResponse(c, err.Error())

	case domain.ErrUnknown:
		unreachableError(c, err)

	default:
		c.JSON(http.StatusOK, animal.Response())
	}

}

// FIXME:
// УНИЧТОЖИТЬ
// FIXME:
func parseSearchAnimalParams(c *gin.Context) (*domain.AnimalSearchParams, error) {
	var animalSearchParams domain.AnimalSearchParams

	startDateTimeQuery := c.Query("startDateTime")
	endDateTimeQuery := c.Query("endDateTime")

	chipperIDQuery := c.Query("chipperId")
	chippedLocationIDQuery := c.Query("chippedLocationId")

	fromQuery := c.Query("from")
	sizeQuery := c.Query("size")

	lifeStatus := c.Query("lifeStatus")
	gender := c.Query("gender")

	if startDateTimeQuery == "" || startDateTimeQuery == "null" {
		animalSearchParams.StartDateTime = nil

	} else {
		startDateTime, err := time.Parse(time.RFC3339, startDateTimeQuery)
		animalSearchParams.StartDateTime = &startDateTime

		if err != nil {
			return nil, err
		}
	}

	if endDateTimeQuery == "" || endDateTimeQuery == "null" {
		animalSearchParams.EndDateTime = nil

	} else {
		endDateTime, err := time.Parse(time.RFC3339, endDateTimeQuery)
		animalSearchParams.EndDateTime = &endDateTime

		if err != nil {
			return nil, err
		}
	}

	if chipperIDQuery == "" || chipperIDQuery == "null" {
		animalSearchParams.ChipperID = nil

	} else {
		chipperID, err := strconv.Atoi(chipperIDQuery)
		if chipperID <= 0 || err != nil {
			return nil, err
		}
		animalSearchParams.ChipperID = &chipperID
	}

	if chippedLocationIDQuery == "" || chippedLocationIDQuery == "null" {
		animalSearchParams.ChippedLocationID = nil

	} else {
		chippedLocationID, err := strconv.Atoi(chippedLocationIDQuery)
		if chippedLocationID <= 0 || err != nil {
			return nil, err
		}
		animalSearchParams.ChippedLocationID = &chippedLocationID
	}

	if lifeStatus == "" || lifeStatus == "null" {
		animalSearchParams.LifeStatus = nil

	} else {
		if lifeStatus == "ALIVE" || lifeStatus == "DEAD" {
			animalSearchParams.LifeStatus = &lifeStatus
		} else {
			return nil, errors.New("invalid lifeStatus")
		}
	}

	if gender == "" || gender == "null" {
		animalSearchParams.Gender = nil

	} else {
		if gender == "MALE" || gender == "FEMALE" || gender == "OTHER" {
			animalSearchParams.Gender = &gender
		} else {
			return nil, errors.New("invalid gender")
		}
	}

	if fromQuery == "" || fromQuery == "null" {
		animalSearchParams.From = 0

	} else {
		from, err := strconv.Atoi(fromQuery)
		if err != nil {
			return nil, errors.New("invalid query from")
		}
		if from < 0 {
			return nil, errors.New("invalid from < 0")
		}

		animalSearchParams.From = from
	}

	if sizeQuery == "" || sizeQuery == "null" {
		animalSearchParams.Size = 10

	} else {
		size, err := strconv.Atoi(sizeQuery)
		if err != nil {
			return nil, errors.New("invalid query size")
		}
		if size <= 0 {
			return nil, errors.New("invalide size <= 0")
		}
		animalSearchParams.Size = size
	}
	return &animalSearchParams, nil
}
