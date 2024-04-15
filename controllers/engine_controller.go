package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/trusthemind/go-cars-app/initializers"
	"github.com/trusthemind/go-cars-app/models"
)

// @Tags			Engine
// @Summary		Engine CRUD
// @Description	Create new Engine
// @Accept			json
// @Produce		json
// @Param			request	body		models.Engine	true	"Engine Info"
// @Success		200		{object}	models.Engine
// @Failure		400		{object}	models.Error
// @Router			/engine/create [post]
func CreateEngine(c *gin.Context) {
	var RequestBody models.Engine

	if c.Bind(&RequestBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	engine := models.Engine{
		Name:        RequestBody.Name,
		Brand:       RequestBody.Brand,
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

// @Tags			Engine
// @Summary		Engine CRUD
// @Description	Get All engines
// @Accept			json
// @Produce		json
// @Param			brand	query		string	false	"Engine Brand"
// @Success		200		{array}		models.Engine
// @Failure		404		{object}	models.Error
// @Router			/engine [get]
func GetAllEngines(c *gin.Context) {
	brand := c.Query("brand")
	engines := []models.Engine{}
	var result *gorm.DB
	if brand != "" {
		result = initializers.DB.Where("Brand = ?", brand).Find(&engines)
	} else {
		result = initializers.DB.Find(&engines)
	}

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to get engines"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": engines})
}

// @Tags			Engine
// @Summary		Engine CRUD
// @Description	Delete engine by ID
// @Accept			json
// @Produce		json
// @Params			engine_id path string "Engine ID"
// @Success		200	{object}	[]models.Message
// @Failure		400	{object}	models.Error
// @Router			/engine/delete/:id [delete]
func DeleteEngineByID(c *gin.Context) {
	var engine_id = c.Param("id")
	engine := models.Engine{}

	result := initializers.DB.Where("ID =?", engine_id).Delete(&engine)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete engine"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Engine deleted successfully"})
}

// @Tags			Engine
// @Summary		Engine CRUD
// @Description	Update engine info
// @Accept			json
// @Produce		json
// @Params			engine_id path string "Engine ID"
// @Params			request body models.Engine true "Field"
// @Success		200	{object}	[]models.Engine
// @Failure		400	{object}	models.Error
// @Failure		404	{object}	models.Error
// @Router			/engine/update/:id [put]
func UpdateEngineInfo(c *gin.Context) {
    var RequestBody map[string]interface{}
    var engineID = c.Param("id")

    if err := c.BindJSON(&RequestBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
        return
    }

    engine := models.Engine{}
    if result := initializers.DB.First(&engine, "ID = ?", engineID); result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Engine not found"})
        return
    }

    dbFields := map[string]interface{}{}
    for key, value := range RequestBody {
        if key != "name" {
            dbFields[key] = value
        }
    }

    if err := initializers.DB.Model(&engine).Updates(dbFields).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update engine"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Engine updated successfully"})
}
