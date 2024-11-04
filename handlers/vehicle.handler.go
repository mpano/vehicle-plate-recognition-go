package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"vehicle-plate-recognition/services"
	"vehicle-plate-recognition/store/postgres"
)

func RegisterVehicleRoutes(router *gin.Engine, db *postgres.Database) {
	vehicleService := services.NewVehicleService(db)
	router.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
			return
		}

		// Process the file
		plate, err := vehicleService.ProcessVehicleImage(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Vehicle record saved", "plate": plate})
	})
}
