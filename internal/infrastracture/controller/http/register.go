package http

import (
	"animal-chipization/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type registerUsecase interface {
	Register(params domain.RegistrationParams) (*domain.Account, error)
}

type RegisterHandler struct {
	usecase    registerUsecase
	middleware AuthMiddleware
}

func NewRegisterHandler(usecase registerUsecase, middleware *AuthMiddleware) *RegisterHandler {
	return &RegisterHandler{
		usecase:    usecase,
		middleware: *middleware,
	}
}

func (h *RegisterHandler) InitRoutes(router *gin.Engine) *gin.Engine {

	router.POST("registration",
		h.middleware.blockAuthHeader,
		errorHandlerWrap(h.CreateAccount),
	)

	return router
}

func (h *RegisterHandler) CreateAccount(c *gin.Context) error {
	var input domain.RegistrationParams
	if err := c.BindJSON(&input); err != nil {
		return NewErrBind(err)
	}

	account, err := h.usecase.Register(input)
	if err != nil {
		return err
	}

	c.JSON(http.StatusCreated, account.Map())
	return nil
}
