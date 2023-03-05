package app

import (
	"animal-chipization/internal/controller"

	httpController "animal-chipization/internal/controller/http"
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
	visitedLocatoinRepository := psql.NewVisitedLocationRepository(psqlDB)

	accountUsecase := usecase.NewAccountUsecase(accountRepository)
	registerAccountUsecase := usecase.NewRegisterAccountUsecase(accountRepository)
	locationUsecase := usecase.NewLocationUsecase(locationRepository)
	animalTypeUsecase := usecase.NewAnimalTypeUsecase(animalTypeRepository)
	animalUsecase := usecase.NewAnimalUsecase(animalRepository, animalTypeRepository)
	visitedLocationUsecase := usecase.NewVisitedLocationUsecase(visitedLocatoinRepository, locationRepository, animalRepository)

	middlerware := httpController.NewMiddleware(registerAccountUsecase)

	accountHandler := httpController.NewAccountHandler(accountUsecase, middlerware)
	registerHandler := httpController.NewRegisterHandler(registerAccountUsecase, middlerware)
	locationHandler := httpController.NewLocationHandler(locationUsecase, middlerware)
	animalTypeHandler := httpController.NewAnimalTypeHandler(animalTypeUsecase, middlerware)
	animalHandler := httpController.NewAnimalHandler(animalUsecase, middlerware)
	visitedLocationHandler := httpController.NewVisitedLocationsHandler(visitedLocationUsecase, middlerware)

	router := gin.New()

	router = accountHandler.InitRoutes(router)
	router = registerHandler.InitRoutes(router)
	router = locationHandler.InitRoutes(router)
	router = animalTypeHandler.InitRoutes(router)
	router = animalHandler.InitRoutes(router)
	router = visitedLocationHandler.InitRoutes(router)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("exclude_whitespace", httpController.ExcludeWhitespace)
		//v.RegisterValidation("default", httpController.DefaultValue, true)
	}

	server := controller.NewHTTPServer("8000", router)

	return server.Run()
}
