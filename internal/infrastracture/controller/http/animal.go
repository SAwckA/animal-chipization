package http

import (
	"animal-chipization/internal/domain"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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
		animal.Use(h.middleware.ckeckAuthHeaderMiddleware)
		animal.GET(fmt.Sprintf("/:%s", animalIDParam),
			errorHandlerWrap(h.getAnimal),
		)

		animal.GET("/search",
			errorHandlerWrap(h.searchAnimal),
		)

		animal.POST("",
			h.middleware.authMiddleware,
			errorHandlerWrap(h.createAnimal),
		)

		animal.PUT(fmt.Sprintf("/:%s", animalIDParam),
			h.middleware.authMiddleware,
			errorHandlerWrap(h.updateAnimal),
		)

		animal.DELETE(fmt.Sprintf("/:%s", animalIDParam),
			h.middleware.authMiddleware,
			errorHandlerWrap(h.deleteAnimal),
		)

		types := animal.Group(fmt.Sprintf(":%s/types", animalIDParam))
		{
			types.POST(fmt.Sprintf(":%s", typeParam),
				h.middleware.authMiddleware,
				errorHandlerWrap(h.addAnimalType),
			)

			types.PUT("",
				h.middleware.authMiddleware,
				errorHandlerWrap(h.editAnimalType),
			)

			types.DELETE(fmt.Sprintf(":%s", typeParam),
				h.middleware.authMiddleware,
				errorHandlerWrap(h.deleteAnimalType),
			)
		}
	}

	return router
}

func (h *AnimalHandler) getAnimal(c *gin.Context) error {

	animalID, err := validateID(c.Copy(), animalIDParam)

	if err != nil {
		return err
	}

	animal, err := h.usecase.GetAnimal(animalID)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, animal.Response())
	return nil
}

func (h *AnimalHandler) searchAnimal(c *gin.Context) error {

	var input domain.AnimalSearchParams

	if err := c.BindQuery(&input); err != nil {
		return NewErrBind(err)
	}

	animalsList, err := h.usecase.SearchAnimal(&input)
	if err != nil {
		return err
	}

	if animalsList == nil {
		c.JSON(http.StatusOK, nil)
		return nil
	}
	resp := make([]map[string]interface{}, 0)

	for _, v := range *animalsList {
		resp = append(resp, v.Response())
	}
	c.JSON(http.StatusOK, resp)
	return nil
}

func (h *AnimalHandler) createAnimal(c *gin.Context) error {

	var input domain.AnimalCreateParams
	if err := c.BindJSON(&input); err != nil {
		return NewErrBind(err)
	}

	animal, err := h.usecase.CreateAnimal(input)
	if err != nil {
		return err
	}

	c.JSON(http.StatusCreated, animal.Response())
	return nil

}

func (h *AnimalHandler) updateAnimal(c *gin.Context) error {

	var input domain.AnimalUpdateParams

	animalID, err := validateID(c.Copy(), animalIDParam)
	if err != nil {
		return err
	}

	if err = c.BindJSON(&input); err != nil {
		return NewErrBind(err)
	}

	animal, err := h.usecase.UpdateAnimal(animalID, input)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, animal.Response())
	return nil
}

func (h *AnimalHandler) deleteAnimal(c *gin.Context) error {

	animalID, err := validateID(c.Copy(), animalIDParam)
	if err != nil {
		return NewErrBind(err)
	}

	err = h.usecase.DeleteAnimal(animalID)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, nil)
	return nil
}

func (h *AnimalHandler) addAnimalType(c *gin.Context) error {

	animalID, err := validateID(c.Copy(), animalIDParam)
	if err != nil {
		return err
	}

	animalTypeID, err := validateID(c.Copy(), typeParam)
	if err != nil {
		return err
	}

	animal, err := h.usecase.AddAnimalType(animalID, animalTypeID)
	if err != nil {
		return err
	}

	c.JSON(http.StatusCreated, animal.Response())
	return nil
}

func (h *AnimalHandler) editAnimalType(c *gin.Context) error {
	var input domain.AnimalEditTypeParams

	animalID, err := validateID(c.Copy(), animalIDParam)
	if err != nil {
		return err
	}

	if err := c.BindJSON(&input); err != nil {
		return NewErrBind(err)
	}

	animal, err := h.usecase.EditAnimalType(animalID, input)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, animal.Response())
	return nil
}

func (h *AnimalHandler) deleteAnimalType(c *gin.Context) error {

	animalID, err := validateID(c.Copy(), animalIDParam)
	if err != nil {
		return err
	}

	typeID, err := validateID(c.Copy(), typeParam)
	if err != nil {
		return err
	}

	animal, err := h.usecase.DeleteAnimalType(animalID, typeID)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, animal.Response())
	return nil
}
