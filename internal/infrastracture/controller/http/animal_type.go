package http

import (
	"animal-chipization/internal/domain"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const typeParam = "typeId"

type animalTypeUsecase interface {
	AnimalType(id int) (*domain.AnimalType, error)
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
		animalTypes.Use(h.middleware.checkAuthHeaderMiddleware)
		animalTypes.GET(fmt.Sprintf("/:%s", typeParam),
			errorHandlerWrap(h.animalType),
		)
		animalTypes.POST("",
			h.middleware.authMiddleware,
			errorHandlerWrap(h.create),
		)
		animalTypes.PUT(fmt.Sprintf("/:%s", typeParam),
			h.middleware.authMiddleware,
			errorHandlerWrap(h.update),
		)
		animalTypes.DELETE(fmt.Sprintf("/:%s", typeParam),
			h.middleware.authMiddleware,
			errorHandlerWrap(h.delete),
		)
	}

	return router
}

func (h *AnimalTypeHandler) animalType(c *gin.Context) error {
	typeID, err := validateID(c.Copy(), typeParam)
	if err != nil {
		return err
	}

	animalType, err := h.usecase.AnimalType(typeID)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, animalType.Map())
	return nil
}

func (h *AnimalTypeHandler) create(c *gin.Context) error {
	var input domain.AnimalTypeCreate
	if err := c.BindJSON(&input); err != nil {
		return NewErrBind(err)
	}

	animalType, err := h.usecase.Create(input.Type)
	if err != nil {
		return err
	}

	c.JSON(http.StatusCreated, animalType.Map())
	return nil
}

func (h *AnimalTypeHandler) update(c *gin.Context) error {
	typeID, err := validateID(c.Copy(), typeParam)
	if err != nil {
		return err
	}

	var input domain.AnimalTypeCreate
	if err = c.BindJSON(&input); err != nil {
		return NewErrBind(err)
	}

	animalType, err := h.usecase.Update(typeID, input.Type)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, animalType.Map())
	return nil
}

func (h *AnimalTypeHandler) delete(c *gin.Context) error {
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
