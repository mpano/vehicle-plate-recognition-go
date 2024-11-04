package main

import (
	"vehicle-plate-recognition/core/environment"
	"vehicle-plate-recognition/handlers"

	//_ "vehicle-plate-recognition/docs"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"os"
	"vehicle-plate-recognition/store/postgres"
)

// @title Product Service API
// @version 1.0
// @description  product services API.
// @host localhost:4000
// @schemes http https

func main() {
	// Load environment variables from .env.example file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env.example file")
	}

	// Initialize environment variables
	envs := environment.New(
		os.Getenv("PORT"),
		os.Getenv("DB_URL"),
	)

	defer func() {
		if err := recover(); err != nil {
			println("Survived a panic")
		}
	}()

	// Set up the database and services
	db := postgres.New(envs)
	db.Migrate()

	if err != nil {
		println(err.Error())
		log.Fatal(err)
		return
	}

	// Set up Gin router
	r := gin.Default()

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Home routes
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Define vehicle routes
	handlers.RegisterVehicleRoutes(r, db)

	// Start the server
	err = r.Run(":" + envs.Port)
	if err != nil {
		log.Fatalf("Couldn't start the application: %s", err)
	}
}
