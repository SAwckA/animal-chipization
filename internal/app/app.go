package app

import (
	"animal-chipization/internal/controller"
	httpController "animal-chipization/internal/controller/http"
	"animal-chipization/internal/infrastracture/repository"
	psql "animal-chipization/internal/infrastracture/repository/postgresql"
	"animal-chipization/internal/usecase"

	"github.com/gin-gonic/gin"
)

func Run() error {
	// Точка инициализации слоёв приложения
	// Flow инициализации repository(data storage) -> usecase(логика приложения) -> handler(controller)

	// mongoClient, err := repository.NewMongoClient("mongodb://dev:changeme@localhost:27017")

	// if err != nil {
	// 	return err
	// }

	// mongoRepository := mongodb.NewMongoRepository(mongoClient)

	psqlDB, err := repository.NewPostgresDB("localhost", "5432", "dev", "animal-chipization", "changeme", "disable")

	if err != nil {
		return err
	}

	accountRepository := psql.NewAccountRepository(psqlDB)
	locationRepository := psql.NewLocationRepository(psqlDB)

	accountUsecase := usecase.NewAccountUsecase(accountRepository)
	locationUsecase := usecase.NewLocationUsecase(locationRepository)

	middlerware := httpController.NewMiddleware(accountUsecase)

	accountHandler := httpController.NewAccountHandler(accountUsecase, middlerware)
	registerHandler := httpController.NewRegisterHandler(accountUsecase, middlerware)
	locationHandler := httpController.NewLocationHandler(locationUsecase, middlerware)

	router := gin.New()

	router = accountHandler.InitRoutes(router)
	router = registerHandler.InitRoutes(router)
	router = locationHandler.InitRoutes(router)

	server := controller.NewHTTPServer("8000", router)

	return server.Run()
}
