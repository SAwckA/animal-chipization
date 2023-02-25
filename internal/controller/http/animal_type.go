package http

import (
	"animal-chipization/internal/domain"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const typeParam = "typeId"

type animalTypeUsecase interface {
	Get(id int) (*domain.AnimalType, error)
	Create(typeName string) (*domain.AnimalType, error)
	Update(id int, typeName string) (*domain.AnimalType, error)
	Delete(id int) error
}

type AnimalTypeHandler struct {
	usecase    animalTypeUsecase
	middleware authMiddleware
}

func NewAnimalTypeHandler(usecase animalTypeUsecase, middleware authMiddleware) *AnimalTypeHandler {
	return &AnimalTypeHandler{usecase: usecase, middleware: middleware}
}

func (h *AnimalTypeHandler) InitRoutes(router *gin.Engine) *gin.Engine {

	animalTypes := router.Group("animals/types")
	{
		animalTypes.GET(fmt.Sprintf("/:%s", typeParam),
			h.middleware.ckeckAuthHeaderMiddleware,
			h.getAnimalType,
		)
		animalTypes.POST("",
			h.middleware.authMiddleware,
			h.createAnimalType,
		)
		animalTypes.PUT(fmt.Sprintf("/:%s", typeParam),
			h.middleware.authMiddleware,
			h.updateAnimalType,
		)
		animalTypes.DELETE(fmt.Sprintf("/:%s", typeParam),
			h.middleware.authMiddleware,
			h.deleteAnimalType,
		)
	}

	return router
}

func (h *AnimalTypeHandler) getAnimalType(c *gin.Context) {
	var typeID domain.TypeId

	if err := c.BindUri(&typeID); err != nil {
		badRequest(c, err.Error())
		return
	}

	animalType, err := h.usecase.Get(typeID.ID)

	switch errors.Unwrap(err) {
	case domain.ErrNotFound:
		notFoundResponse(c, err.Error())

	case domain.ErrAlreadyExist:
		conflictResponse(c, err.Error())

	case domain.ErrLinked:
		badRequest(c, err.Error())

	case domain.ErrUnknown:
		unreachableError(c, err)

	default:
		c.JSON(http.StatusOK, animalType.Response())
	}
}

func (h *AnimalTypeHandler) createAnimalType(c *gin.Context) {
	var input domain.AnimalTypeDTO

	if err := c.BindJSON(&input); err != nil {
		badRequest(c, err.Error())
		return
	}

	animalType, err := h.usecase.Create(input.Type)

	switch errors.Unwrap(err) {
	case domain.ErrNotFound:
		notFoundResponse(c, err.Error())

	case domain.ErrAlreadyExist:
		conflictResponse(c, err.Error())

	case domain.ErrLinked:
		badRequest(c, err.Error())

	case domain.ErrUnknown:
		unreachableError(c, err)

	default:
		c.JSON(http.StatusCreated, animalType.Response())
	}
}

func (h *AnimalTypeHandler) updateAnimalType(c *gin.Context) {
	var typeID domain.TypeId
	var input domain.AnimalTypeDTO

	if err := c.BindUri(&typeID); err != nil {
		badRequest(c, err.Error())
		return
	}

	if err := c.BindJSON(&input); err != nil {
		badRequest(c, err.Error())
		return
	}

	animalType, err := h.usecase.Update(typeID.ID, input.Type)

	switch errors.Unwrap(err) {
	case domain.ErrNotFound:
		notFoundResponse(c, err.Error())

	case domain.ErrAlreadyExist:
		conflictResponse(c, err.Error())

	case domain.ErrLinked:
		badRequest(c, err.Error())

	case domain.ErrUnknown:
		unreachableError(c, err)

	default:
		c.JSON(http.StatusOK, animalType.Response())
	}
}

func (h *AnimalTypeHandler) deleteAnimalType(c *gin.Context) {
	var typeID domain.TypeId

	if err := c.BindUri(&typeID); err != nil {
		badRequest(c, err.Error())
		return
	}

	err := h.usecase.Delete(typeID.ID)

	switch errors.Unwrap(err) {
	case domain.ErrNotFound:
		notFoundResponse(c, err.Error())

	case domain.ErrAlreadyExist:
		conflictResponse(c, err.Error())

	case domain.ErrLinked:
		badRequest(c, err.Error())

	case domain.ErrUnknown:
		unreachableError(c, err)

	default:
		c.JSON(http.StatusOK, nil)
	}
}
