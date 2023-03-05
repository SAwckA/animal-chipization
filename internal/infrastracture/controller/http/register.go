package http

import (
	"animal-chipization/internal/domain"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type registerUsecase interface {
	Register(params domain.RegistrationDTO) (*domain.Account, error)
}

type RegisterHandler struct {
	usecase    registerUsecase
	middleware Middleware
}

func NewRegisterHandler(usecase registerUsecase, middleware *Middleware) *RegisterHandler {
	return &RegisterHandler{
		usecase:    usecase,
		middleware: *middleware,
	}
}

func (h *RegisterHandler) InitRoutes(router *gin.Engine) *gin.Engine {

	router.POST("registration", h.middleware.blockAuthHeader, h.CreateAccount)

	return router
}

func (h *RegisterHandler) CreateAccount(c *gin.Context) {
	var input domain.RegistrationDTO

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	account, err := h.usecase.Register(input)

	switch {
	case errors.Unwrap(err) == domain.ErrAlreadyExist:
		newErrorResponse(c, http.StatusConflict, err.Error(), nil)
	case err != nil:
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
	default:
		c.JSON(http.StatusCreated, account.Response())
	}
}
