package controllers

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"

	"github.com/trusthemind/go-cars-app/initializers"
	"github.com/trusthemind/go-cars-app/models"
)

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

func GetAllEngines(c *gin.Context) {
	engines := []models.Engine{}
	result := initializers.DB.Find(&engines)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to get engines"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"engines": result})
}

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

func UpdateEngineInfo(c *gin.Context) {
	var RequestBody models.Engine

	if c.Bind(&RequestBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	var engineID = c.Param("id")
	engine := models.Engine{}

	result := initializers.DB.First(&engine, engineID)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to get engine"})
		return

	}

	if result.RowsAffected > 0 {
		reflectV := reflect.ValueOf(&RequestBody).Elem()
		reflectT := reflect.TypeOf(&RequestBody).Elem()

		for i := 0; i < reflectV.NumField(); i++ {
			field := reflectT.Field(i)
			fieldName := field.Name
			requestFieldValue := reflectV.FieldByName(fieldName).Interface()
			existingFieldValue := reflect.ValueOf(&engine).Elem().FieldByName(fieldName).Interface()
			if requestFieldValue != existingFieldValue {
				reflect.ValueOf(&engine).Elem().FieldByName(fieldName).Set(reflect.ValueOf(requestFieldValue))
			}
		}

	}

	c.JSON(http.StatusOK, gin.H{"message": "Engine updated successfully"})
}
