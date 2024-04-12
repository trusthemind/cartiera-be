package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-zoox/core-utils/fmt"

	"github.com/trusthemind/go-cars-app/initializers"
	"github.com/trusthemind/go-cars-app/models"
)

//	@Tags			Engine
//	@Summary		Engine CRUD
//	@Description	Create new Engine
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.EngineRequest	true	"Engine Info"
//	@Success		200		{object}	models.Engine
//	@Failure		400		{object}	models.Error
//	@Router			/engine/create [post]
func CreateEngine(c *gin.Context) {
	var RequestBody models.Engine

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

//	@Tags			Engine
//	@Summary		Engine CRUD
//	@Description	Get All engines
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]models.Engine
//	@Failure		404	{object}	models.Error
//	@Router			/engine/all [get]
func GetAllEngines(c *gin.Context) {
	engines := []models.Engine{}
	result := initializers.DB.Find(&engines)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to get engines"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": engines})
}

//	@Tags			Engine
//	@Summary		Engine CRUD
//	@Description	Delete engine by ID
//	@Accept			json
//	@Produce		json
//	@Params			engine_id path string "Engine ID"
//	@Success		200	{object}	[]models.Message
//	@Failure		400	{object}	models.Error
//	@Router			/engine/delete/:id [delete]
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

//	@Tags			Engine
//	@Summary		Engine CRUD
//	@Description	Update engine info
//	@Accept			json
//	@Produce		json
//	@Params			engine_id path string "Engine ID"
//	@Params			request body models.EngineRequest true "Field"
//	@Success		200	{object}	[]models.Engine
//	@Failure		400	{object}	models.Error
//	@Failure		404	{object}	models.Error
//	@Router			/engine/update/:id [put]
func UpdateEngineInfo(c *gin.Context) {
	var RequestBody models.Engine
	var engine_id = c.Param("id")

	if c.Bind(&RequestBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}
	engine := models.Engine{}

	result := initializers.DB.First(&engine, "ID = ?", engine_id)
	fmt.Println(engine_id, result)

	c.JSON(http.StatusOK, gin.H{"message": engine})
}
