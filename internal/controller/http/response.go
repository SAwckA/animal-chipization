package http

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Msg string `json:"msg"`
}

//Ответ об ошибке
func newErrorResponse(c *gin.Context, statusCode int, msg string, err error) {
	if err != nil {
		logrus.Errorf(err.Error())
	}
	c.AbortWithStatusJSON(statusCode, errorResponse{Msg: msg})
}
