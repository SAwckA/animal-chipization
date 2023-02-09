package http

import (
	"animal-chipization/internal/domain"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type registerUsecase interface {
	CreateUser(firstName, lastName, email, password string) (*domain.Account, error)
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

func (h *RegisterHandler) CreateAccount(ctx *gin.Context) {
	firstName := ctx.Query("firstName")
	lastName := ctx.Query("lastName")
	email := ctx.Query("email")
	password := ctx.Query("password")

	if len(firstName) < 1 || strings.Contains(firstName, " ") {
		newErrorResponse(ctx, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	if len(lastName) < 1 || strings.Contains(lastName, " ") {
		newErrorResponse(ctx, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	if len(email) < 1 || strings.Contains(email, " ") || !validateEmail(email) {
		newErrorResponse(ctx, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	if len(password) < 1 || strings.Contains(password, " ") {
		newErrorResponse(ctx, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	account, err := h.usecase.CreateUser(firstName, lastName, email, password)

	if err != nil {
		// FIXME: Обработка ошибки индекса email
		newErrorResponse(ctx, http.StatusConflict, "User with this email exist", err)
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id":        account.ID,
		"firstName": account.FirstName,
		"lastName":  account.LastName,
		"email":     account.Email,
	})

}

func validateEmail(email string) bool {
	// FIXME: Валидация по regex
	if strings.Contains(email, "@") && strings.Contains(email, ".") {
		return true
	}
	return false
}
