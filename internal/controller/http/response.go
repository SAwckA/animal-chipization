package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Msg string `json:"msg"`
}

// Ошибка с любым статусом
func newErrorResponse(c *gin.Context, statusCode int, msg string, err error) {
	if err != nil {
		logrus.Errorf(err.Error())
	}
	c.AbortWithStatusJSON(statusCode, errorResponse{Msg: msg})
}

// Alias для newErrorResponse(c, http.StatusConflict, msg, nil)
func conflictResponse(c *gin.Context, msg string) {
	newErrorResponse(c, http.StatusConflict, msg, nil)
}

// Alias для newErrorResponse(c, http.StatusBadRequest, msg, nil)
func badRequest(c *gin.Context, msg string) {
	newErrorResponse(c, http.StatusBadRequest, msg, nil)
}

// Alias для newErrorResponse(c, http.StatusNotFound, msg, nil)
func notFoundResponse(c *gin.Context, msg string) {
	newErrorResponse(c, http.StatusNotFound, msg, nil)
}

// Alias для newErrorResponse(c, http.StatusForbidden, msg, nil)
func forbiddenResponse(c *gin.Context, msg string) {
	newErrorResponse(c, http.StatusForbidden, msg, nil)
}

// Необработаная ошибка
// Alias для newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
func unreachableError(c *gin.Context, err error) {
	if err == nil {
		newErrorResponse(c, http.StatusInternalServerError, "unexpected err", nil)
		return
	}
	newErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
}
