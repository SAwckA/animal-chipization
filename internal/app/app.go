package app

import (
	"animal-chipization/internal/infrastracture/controller"
	"animal-chipization/internal/infrastracture/controller/http"

	"animal-chipization/internal/infrastracture/repository"
	psql "animal-chipization/internal/infrastracture/repository/postgresql"
	"animal-chipization/internal/usecase"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func Run() error {
	// Точка инициализации слоёв приложения
	// Flow инициализации repository(data storage) -> usecase(логика приложения) -> handler(controller)

	// mongoClient, err := repository.NewMongoClient("mongodb://dev:changeme@localhost:27017")

	// if err != nil {
	// 	return err
	// }

	// mongoRepository := mongodb.NewMongoRepository(mongoClient)
	_ = godotenv.Load()

	psqlDB, err := repository.NewPostgresDB(os.Getenv("DB_HOST"), "5432", "dev", "animal-chipization", "changeme", "disable")

	if err != nil {
		return err
	}

	accountRepository := psql.NewAccountRepository(psqlDB)
	locationRepository := psql.NewLocationRepository(psqlDB)
	animalTypeRepository := psql.NewAnimalTypeRepository(psqlDB)
	animalRepository := psql.NewAnimalRepository(psqlDB)
	visitedLocationRepository := psql.NewVisitedLocationRepository(psqlDB)

	accountUsecase := usecase.NewAccountUsecase(accountRepository)
	locationUsecase := usecase.NewLocationUsecase(locationRepository)
	animalTypeUsecase := usecase.NewAnimalTypeUsecase(animalTypeRepository)
	animalUsecase := usecase.NewAnimalUsecase(animalRepository, animalTypeRepository)
	visitedLocationUsecase := usecase.NewVisitedLocationUsecase(visitedLocationRepository, locationRepository, animalRepository)

	middleware := http.NewAuthMiddleware(accountUsecase)

	accountHandler := http.NewAccountHandler(accountUsecase, middleware)
	registerHandler := http.NewRegisterHandler(accountUsecase, middleware)
	locationHandler := http.NewLocationHandler(locationUsecase, middleware)
	animalTypeHandler := http.NewAnimalTypeHandler(animalTypeUsecase, middleware)
	animalHandler := http.NewAnimalHandler(animalUsecase, middleware)
	visitedLocationHandler := http.NewVisitedLocationsHandler(visitedLocationUsecase, middleware)

	router := gin.New()

	router = accountHandler.InitRoutes(router)
	router = registerHandler.InitRoutes(router)
	router = locationHandler.InitRoutes(router)
	router = animalTypeHandler.InitRoutes(router)
	router = animalHandler.InitRoutes(router)
	router = visitedLocationHandler.InitRoutes(router)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("exclude_whitespace", http.ExcludeWhitespace)
		//v.RegisterValidation("default", httpController.DefaultValue, true)
	}

	server := controller.NewHTTPServer("8000", router)

	return server.Run()
}
