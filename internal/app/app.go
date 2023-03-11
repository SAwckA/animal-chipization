package app

import (
	"animal-chipization/config"
	"animal-chipization/internal/infrastracture/controller"
	"animal-chipization/internal/infrastracture/controller/http"
	"animal-chipization/internal/infrastracture/repository"
	psql "animal-chipization/internal/infrastracture/repository/postgresql"
	"animal-chipization/internal/usecase"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	_ = godotenv.Load()

	appConfig := config.LoadConfig()

	psqlDB, err := repository.NewPostgresDB(appConfig.PostgresConfig)

	if err != nil {
		log.Fatalf("cant connect to database, cause: %s", err.Error())
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

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router = accountHandler.InitRoutes(router)
	router = registerHandler.InitRoutes(router)
	router = locationHandler.InitRoutes(router)
	router = animalTypeHandler.InitRoutes(router)
	router = animalHandler.InitRoutes(router)
	router = visitedLocationHandler.InitRoutes(router)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("exclude_whitespace", http.ExcludeWhitespace)
		_ = v.RegisterValidation("allowed_strings", http.AllowedStrings)
	}

	server := controller.NewHTTPServer(appConfig.HttpConfig.Port, router)

	logrus.Infof("HTTP SERVER IS STARTING AT PORT: %s", appConfig.HttpConfig.Port)
	go func() {
		_ = server.Run()
	}()

	quit := make(chan os.Signal)

	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("SHUTDOWN HTTP SERVER...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Error during shutdown http server:", err)
	}

	select {
	case <-ctx.Done():
		log.Println("timeout of 3 seconds.")
	}

	log.Println("Server exiting")
}
