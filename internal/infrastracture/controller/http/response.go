package http

import (
	"animal-chipization/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type errorResponse struct {
	Msg string `json:"msg"`
}

func NewErrBind(e error) error {
	return &domain.ApplicationError{
		OriginalError: e,
		SimplifiedErr: domain.ErrInvalidInput,
		Description:   "Invalid data",
	}
}

// Ошибка с любым статусом
func newErrorResponse(c *gin.Context, statusCode int, msg string) {
	c.AbortWithStatusJSON(statusCode, errorResponse{Msg: msg})
}

// Alias для newErrorResponse(c, http.StatusConflict, msg, nil)
func conflictResponse(c *gin.Context, msg string) {
	newErrorResponse(c, http.StatusConflict, msg)
}

// Alias для newErrorResponse(c, http.StatusBadRequest, msg, nil)
func badRequest(c *gin.Context, msg string) {
	newErrorResponse(c, http.StatusBadRequest, msg)
}

// Alias для newErrorResponse(c, http.StatusNotFound, msg, nil)
func notFoundResponse(c *gin.Context, msg string) {
	newErrorResponse(c, http.StatusNotFound, msg)
}

// Alias для newErrorResponse(c, http.StatusForbidden, msg, nil)
func forbiddenResponse(c *gin.Context, msg string) {
	newErrorResponse(c, http.StatusForbidden, msg)
}

// Alias для newErrorResponse(c, http.StatusUnauthorized, msg, nil)
func unauthorizedResponse(c *gin.Context, msg string) {
	newErrorResponse(c, http.StatusUnauthorized, msg)
}

// Необработаная ошибка
// Alias для newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
func internalError(c *gin.Context, err error) {
	if err == nil {
		newErrorResponse(c, http.StatusInternalServerError, "unexpected err")
		return
	}

	newErrorResponse(c, http.StatusInternalServerError, err.Error())

	logrus.Errorf("%s | %s | %s",
		err.(*domain.ApplicationError).OriginalError,
		err.(*domain.ApplicationError).SimplifiedErr,
		err.(*domain.ApplicationError).Description,
	)

	defer func() {
		_ = recover()
		logrus.Errorf(err.Error())
	}()
}
