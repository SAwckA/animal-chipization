package http

import (
	"animal-chipization/internal/domain"
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
		animalTypes.Use(h.middleware.ckeckAuthHeaderMiddleware)
		animalTypes.GET(fmt.Sprintf("/:%s", typeParam),
			errorHandlerWrap(h.getAnimalType),
		)
		animalTypes.POST("",
			h.middleware.authMiddleware,
			errorHandlerWrap(h.createAnimalType),
		)
		animalTypes.PUT(fmt.Sprintf("/:%s", typeParam),
			h.middleware.authMiddleware,
			errorHandlerWrap(h.updateAnimalType),
		)
		animalTypes.DELETE(fmt.Sprintf("/:%s", typeParam),
			h.middleware.authMiddleware,
			errorHandlerWrap(h.deleteAnimalType),
		)
	}

	return router
}

func (h *AnimalTypeHandler) getAnimalType(c *gin.Context) error {
	typeID, err := validateID(c.Copy(), typeParam)
	if err != nil {
		return err
	}

	animalType, err := h.usecase.Get(typeID)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, animalType.Response())
	return nil
}

func (h *AnimalTypeHandler) createAnimalType(c *gin.Context) error {
	var input domain.AnimalTypeDTO

	if err := c.BindJSON(&input); err != nil {
		return NewErrBind(err)
	}

	animalType, err := h.usecase.Create(input.Type)
	if err != nil {
		return err
	}

	c.JSON(http.StatusCreated, animalType.Response())
	return nil
}

func (h *AnimalTypeHandler) updateAnimalType(c *gin.Context) error {
	var input domain.AnimalTypeDTO

	typeID, err := validateID(c.Copy(), typeParam)
	if err != nil {
		return err
	}

	if err := c.BindJSON(&input); err != nil {
		return NewErrBind(err)
	}

	animalType, err := h.usecase.Update(typeID, input.Type)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, animalType.Response())
	return nil
}

func (h *AnimalTypeHandler) deleteAnimalType(c *gin.Context) error {
	typeID, err := validateID(c.Copy(), typeParam)
	if err != nil {
		return err
	}

	err = h.usecase.Delete(typeID)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, nil)
	return nil
}
