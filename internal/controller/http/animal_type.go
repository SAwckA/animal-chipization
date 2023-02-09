package http

import (
	"animal-chipization/internal/domain"
	"animal-chipization/internal/errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const typeParam = "typeId"

type animalTypeUsecase interface {
	GetType(typeID int) *domain.AnimalType
	CreateType(typeName string) (*domain.AnimalType, error)
	UpdateType(typeID int, typeName string) (*domain.AnimalType, error)
	DeleteType(typeID int) error
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
		animalTypes.POST("/",
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
	typeID, err := validateID(c, typeParam)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	animalType := h.usecase.GetType(typeID)

	if animalType == nil {
		newErrorResponse(c, http.StatusNotFound, "Animal type not found", nil)
		return
	}

	c.JSON(http.StatusOK, animalType)
}

type animalTypeRequest struct {
	TypeName *string `json:"type"`
}

func (h *AnimalTypeHandler) createAnimalType(c *gin.Context) {
	var input animalTypeRequest

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	if input.TypeName == nil {
		newErrorResponse(c, http.StatusBadRequest, "Null type data", nil)
		return
	}

	if contentsOnlySpaces(*input.TypeName) {
		newErrorResponse(c, http.StatusBadRequest, "Contents only spaces", nil)
		return
	}

	animalType, err := h.usecase.CreateType(*input.TypeName)

	if err != nil {
		newErrorResponse(c, http.StatusConflict, "This type already exist", nil)
		return
	}

	c.JSON(http.StatusOK, animalType)
}

func (h *AnimalTypeHandler) updateAnimalType(c *gin.Context) {
	typeID, err := validateID(c, typeParam)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var input animalTypeRequest

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if input.TypeName == nil {
		newErrorResponse(c, http.StatusBadRequest, "Null type data", nil)
		return
	}

	if contentsOnlySpaces(*input.TypeName) {
		newErrorResponse(c, http.StatusBadRequest, "Contents only spaces", nil)
		return
	}

	newAnimalType, err := h.usecase.UpdateType(typeID, *input.TypeName)

	if err == errors.ErrNotFound {
		newErrorResponse(c, http.StatusNotFound, "animal type not found", nil)
		return
	}

	if err != nil {
		newErrorResponse(c, http.StatusConflict, "Animal type already exist", nil)
		return
	}

	c.JSON(http.StatusOK, newAnimalType)

}

func (h *AnimalTypeHandler) deleteAnimalType(c *gin.Context) {
	typeID, err := validateID(c, typeParam)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err = h.usecase.DeleteType(typeID)

	if err == errors.ErrNotFound {
		newErrorResponse(c, http.StatusNotFound, "Animal type not found", nil)
		return
	}

	if err != nil {
		newErrorResponse(c, http.StatusConflict, "Animal has this type", nil)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func validateID(c *gin.Context, name string) (int, error) {
	paramString := c.Param(name)

	if paramString == "" || paramString == "null" {
		return 0, errors.ErrInvalidID
	}

	res, err := strconv.Atoi(paramString)

	if err != nil {
		return 0, errors.ErrInvalidID
	}

	if res <= 0 {
		return 0, errors.ErrInvalidID
	}

	return res, nil
}
