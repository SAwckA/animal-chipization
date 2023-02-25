package http

import (
	"animal-chipization/internal/domain"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const accountCtx = "account"

type middlewareUsecase interface {
	Login(email, password string) (*domain.Account, error)
}

type Middleware struct {
	usecase middlewareUsecase
}

func NewMiddleware(usecase middlewareUsecase) *Middleware {
	return &Middleware{usecase: usecase}
}

// blockAuthHeader Обработчик аутентификации, отвечает за аутентификацию
// 		*ЗАПРЕЩАЕТ для авторизованных пользователей
func (m *Middleware) blockAuthHeader(ctx *gin.Context) {
	if authHeader := ctx.GetHeader("Authorization"); len(authHeader) > 0 {
		newErrorResponse(ctx, http.StatusForbidden, "Forbidden for authorized users", nil)
		return
	}
	ctx.Next()
}

// ckeckAuthHeaderMiddleware Обработчик аутентификации, отвечает за аутентификацию
// 		*НЕ Обязательная аутентификация
func (m *Middleware) ckeckAuthHeaderMiddleware(ctx *gin.Context) {
	if authHeader := ctx.GetHeader("Authorization"); len(authHeader) > 0 {
		m.authMiddleware(ctx)
		return
	}
	ctx.Next()
}

// authMiddleware Обработчик аутентификации, отвечает за аутентификацию
// 		*Обязательная аутентификация
func (m *Middleware) authMiddleware(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")

	token, err := validateHeader(authHeader)

	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	account, err := authorize(string(token), m.usecase.Login)

	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	ctx.Set(accountCtx, account)

	ctx.Next()
}

func validateHeader(header string) (string, error) {
	splitedHeaderToken := strings.Split(header, " ")

	if len(splitedHeaderToken) != 2 || splitedHeaderToken[0] != "Basic" {
		return "", errors.New("invalid token")
	}

	encodedToken := splitedHeaderToken[1]

	token, err := base64.StdEncoding.DecodeString(encodedToken)

	if err != nil {
		return "", errors.New("invalid token")
	}

	return string(token), err
}

func authorize(token string, loginFunc func(email string, password string) (*domain.Account, error)) (*domain.Account, error) {

	if splitedAuthString := strings.Split(token, ":"); len(splitedAuthString) == 2 {
		login, password := splitedAuthString[0], splitedAuthString[1]

		account, err := loginFunc(login, password)

		if err != nil {
			return nil, errors.New("invalid credentials")
		}

		return account, err
	}

	return nil, errors.New("invalid token")
}
