package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trusthemind/go-cars-app/initializers"
	"github.com/trusthemind/go-cars-app/models"
)

func CreateEngine(c *gin.Context) {
	var RequestBody struct {
		Name        string  `json:"name" binding:"required"`
		Fuel        string  `json:"fuel" binding:"required"`
		Cilinders   int32   `json:"ciliders"`
		Consumption float32 `json:"consumption"`
	}

	if c.Bind(&RequestBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	engine := models.Engine{
		Name:        RequestBody.Name,
		Fuel:        RequestBody.Fuel,
		Cilinders:   RequestBody.Cilinders,
		Consumption: RequestBody.Consumption,
	}
	result := initializers.DB.Create(&engine)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create engine"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Engine created successfully"})
}
