package http

import (
	"animal-chipization/internal/domain"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

const animalIDParam = "animalId"

type animalUsecase interface {
	Animal(id int) (*domain.Animal, error)
	Search(params *domain.AnimalSearchParams) ([]domain.Animal, error)
	Create(params *domain.AnimalCreateParams) (*domain.Animal, error)
	Update(id int, params *domain.AnimalUpdateParams) (*domain.Animal, error)
	Delete(id int) error

	AddAnimalType(animalID, typeID int) (*domain.Animal, error)
	EditAnimalType(animalID int, params *domain.AnimalEditTypeParams) (*domain.Animal, error)
	DeleteAnimalType(animalID, typeID int) (*domain.Animal, error)
}

type AnimalHandler struct {
	usecase animalUsecase
	auth    authMiddleware
}

func NewAnimalHandler(usecase animalUsecase, auth authMiddleware) *AnimalHandler {
	return &AnimalHandler{usecase: usecase, auth: auth}
}

func (h *AnimalHandler) InitRoutes(router *gin.Engine) *gin.Engine {

	animal := router.Group("animals")
	{
		animal.Use(h.auth.checkAuthHeaderMiddleware)
		animal.GET(fmt.Sprintf("/:%s", animalIDParam),
			errorHandlerWrap(h.animal),
		)

		animal.GET("/search",
			errorHandlerWrap(h.search),
		)

		animal.POST("",
			h.auth.authMiddleware,
			errorHandlerWrap(h.create),
		)

		animal.PUT(fmt.Sprintf("/:%s", animalIDParam),
			h.auth.authMiddleware,
			errorHandlerWrap(h.update),
		)

		animal.DELETE(fmt.Sprintf("/:%s", animalIDParam),
			h.auth.authMiddleware,
			errorHandlerWrap(h.delete),
		)

		types := animal.Group(fmt.Sprintf(":%s/types", animalIDParam))
		{
			types.POST(fmt.Sprintf(":%s", typeParam),
				h.auth.authMiddleware,
				errorHandlerWrap(h.addAnimalType),
			)

			types.PUT("",
				h.auth.authMiddleware,
				errorHandlerWrap(h.editAnimalType),
			)

			types.DELETE(fmt.Sprintf(":%s", typeParam),
				h.auth.authMiddleware,
				errorHandlerWrap(h.deleteAnimalType),
			)
		}
	}

	return router
}

func (h *AnimalHandler) animal(c *gin.Context) error {
	animalID, err := ParamID(c.Copy(), animalIDParam)
	if err != nil {
		return err
	}

	animal, err := h.usecase.Animal(animalID)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, animal.Map())
	return nil
}

func (h *AnimalHandler) search(c *gin.Context) error {
	var input domain.AnimalSearchParams
	if err := c.BindQuery(&input); err != nil {
		return NewErrBind(err)
	}

	animalsList, err := h.usecase.Search(&input)
	if err != nil {
		return err
	}

	resp := make([]map[string]interface{}, 0)

	if animalsList != nil {
		for _, v := range animalsList {
			resp = append(resp, v.Map())
		}
	}

	c.JSON(http.StatusOK, resp)
	return nil
}

func (h *AnimalHandler) create(c *gin.Context) error {
	var input *domain.AnimalCreateParams
	if err := c.BindJSON(&input); err != nil {
		return NewErrBind(err)
	}

	animal, err := h.usecase.Create(input)
	if err != nil {
		return err
	}

	c.JSON(http.StatusCreated, animal.Map())
	return nil

}

func (h *AnimalHandler) update(c *gin.Context) error {
	animalID, err := ParamID(c.Copy(), animalIDParam)
	if err != nil {
		return err
	}

	var input domain.AnimalUpdateParams
	if err = c.BindJSON(&input); err != nil {
		return NewErrBind(err)
	}

	animal, err := h.usecase.Update(animalID, &input)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, animal.Map())
	return nil
}

func (h *AnimalHandler) delete(c *gin.Context) error {
	animalID, err := ParamID(c.Copy(), animalIDParam)
	if err != nil {
		return NewErrBind(err)
	}

	err = h.usecase.Delete(animalID)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, nil)
	return nil
}

func (h *AnimalHandler) addAnimalType(c *gin.Context) error {
	animalID, err := ParamID(c.Copy(), animalIDParam)
	if err != nil {
		return err
	}

	animalTypeID, err := ParamID(c.Copy(), typeParam)
	if err != nil {
		return err
	}

	animal, err := h.usecase.AddAnimalType(animalID, animalTypeID)
	if err != nil {
		return err
	}

	c.JSON(http.StatusCreated, animal.Map())
	return nil
}

func (h *AnimalHandler) editAnimalType(c *gin.Context) error {
	animalID, err := ParamID(c.Copy(), animalIDParam)
	if err != nil {
		return err
	}

	var input *domain.AnimalEditTypeParams
	if err = c.BindJSON(&input); err != nil {
		return NewErrBind(err)
	}

	animal, err := h.usecase.EditAnimalType(animalID, input)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, animal.Map())
	return nil
}

func (h *AnimalHandler) deleteAnimalType(c *gin.Context) error {
	animalID, err := ParamID(c.Copy(), animalIDParam)
	if err != nil {
		return err
	}

	typeID, err := ParamID(c.Copy(), typeParam)
	if err != nil {
		return err
	}

	animal, err := h.usecase.DeleteAnimalType(animalID, typeID)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, animal.Map())
	return nil
}
