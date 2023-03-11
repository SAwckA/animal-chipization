package http

import (
	"animal-chipization/internal/domain"
	"encoding/base64"
	"strings"

	"github.com/gin-gonic/gin"
)

const accountCtx = "account"

type authUsecase interface {
	Login(email, password string) (*domain.Account, error)
}

type AuthMiddleware struct {
	usecase authUsecase
}

func NewAuthMiddleware(usecase authUsecase) *AuthMiddleware {
	return &AuthMiddleware{usecase: usecase}
}

// blockAuthHeader Обработчик аутентификации, отвечает за аутентификацию
// 		*ЗАПРЕЩАЕТ для авторизованных пользователей
func (m *AuthMiddleware) blockAuthHeader(c *gin.Context) {
	if authHeader := c.GetHeader("Authorization"); len(authHeader) > 0 {
		forbiddenResponse(c, "Forbidden for authorized users")
		return
	}
	c.Next()
}

// checkAuthHeaderMiddleware Обработчик аутентификации, отвечает за аутентификацию
// 		*НЕ Обязательная аутентификация
func (m *AuthMiddleware) checkAuthHeaderMiddleware(c *gin.Context) {
	if authHeader := c.GetHeader("Authorization"); len(authHeader) > 0 {
		m.authMiddleware(c)
		return
	}
	c.Next()
}

// authMiddleware Обработчик аутентификации, отвечает за аутентификацию
// 		*Обязательная аутентификация
func (m *AuthMiddleware) authMiddleware(c *gin.Context) {
	email, password, ok := getCredentials(c.Copy())
	if !ok {
		unauthorizedResponse(c, "no credentials")
		return
	}

	account, err := m.usecase.Login(email, password)
	if err != nil {
		unauthorizedResponse(c, err.Error())
		return
	}

	c.Set(accountCtx, account)
	c.Next()
}

// getCredentials Нужен для получения авторизационных данных из заголовка запроса
func getCredentials(cCp *gin.Context) (string, string, bool) {
	if token := strings.Split(cCp.GetHeader("Authorization"), " "); len(token) == 2 || token[0] == "Basic" {
		rawToken := token[1]

		t, err := base64.StdEncoding.DecodeString(rawToken)
		if err != nil {
			return "", "", false
		}

		credentials := strings.Split(string(t), ":")

		return credentials[0], credentials[1], true
	}

	return "", "", false
}
