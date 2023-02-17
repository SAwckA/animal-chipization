package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const visitedPointIDParam = "visitedPointId"

type visitedLocationUsecase interface {
}

type VisitedLocationsHandler struct {
	usecase    visitedLocationUsecase
	middleware authMiddleware
}

func NewVisitedLocationsHandler(usecase visitedLocationUsecase, middleware authMiddleware) *VisitedLocationsHandler {
	return &VisitedLocationsHandler{
		usecase:    usecase,
		middleware: middleware,
	}
}

func (h *VisitedLocationsHandler) InitRoutes(router *gin.Engine) *gin.Engine {

	locations := router.Group(fmt.Sprintf("animal/:%s/locations", animalIDParam))
	{
		locations.GET("/")
		locations.POST(fmt.Sprintf("/:%s", pointIDParam))
		locations.PUT("/")
		locations.DELETE(fmt.Sprintf("/:%s", visitedPointIDParam))
	}

	return router
}
