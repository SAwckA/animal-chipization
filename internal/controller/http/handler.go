package http

import (
	"errors"

	"animal-chipization/internal/domain"
	"github.com/gin-gonic/gin"
)

// errorHandlerWrap Обработка ошибок контроллера
// Соотвествие упрощенной ошибки (domain.ApplicationError.(*)SimplifiedError) к HTTP коду запроса
func errorHandlerWrap(next func(c *gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {

		err := next(c)

		switch errors.Unwrap(err) {
		case domain.ErrInvalidInput, domain.ErrLinked:
			badRequest(c, err.Error())

		case domain.ErrBadDatabaseOut, domain.ErrAlreadyExist:
			conflictResponse(c, err.Error())

		case domain.ErrNotFound:
			notFoundResponse(c, err.Error())

		case domain.ErrForbidden:
			forbiddenResponse(c, err.Error())

		case nil:
			return

		default:
			internalError(c, err)
		}
	}
}
