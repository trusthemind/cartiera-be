package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/trusthemind/go-cars-app/helpers"
	"github.com/trusthemind/go-cars-app/initializers"
	"github.com/trusthemind/go-cars-app/models"
)

// @Tags			Cars
// @Summary		Cars CRUD
// @Description	Create a car for sale
// @Accept			multipart/form-data
// @Produce		json
// @Param			data		formData	object	true	"Car"
// @Param			upload[]	formData	array	true	"Photos"
// @Success		200			{object}	models.Message
// @Failure		400			{object}	models.Error
// @Failure		401			{object}	models.Error
// @Router			/cars/create [post]
func CreateCar(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["upload[]"]
	data := c.Request.FormValue("data")

	var RequestData models.Car
	json.Unmarshal([]byte(data), &RequestData)

	if c.Bind(&RequestData) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header not found"})
		return
	}

	// Remove the "Bearer " prefix from the header value
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	claims, err := helpers.ExtractClaims(tokenString, []byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to get credentials"})
		return
	}
	var id = claims["sub"].(float64)

	var user models.User
	initializers.DB.First(&user, "ID = ?", id)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No uploaded files provided"})
		return
	}
	result, ok, err := helpers.SavePhotoToTable(c, files)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	car := models.Car{
		OwnerID:      int(user.ID),
		EngineID:     RequestData.EngineID,
		Brand:        RequestData.Brand,
		CarModel:     RequestData.CarModel,
		Year:         RequestData.Year,
		Price:        RequestData.Price,
		Photos:       result,
		Status:       RequestData.Status,
		VinCode:      RequestData.VinCode,
		Kilometers:   RequestData.Kilometers,
		Placement:    RequestData.Placement,
		OwnersNumber: RequestData.OwnersNumber,
		OwnerComment: RequestData.OwnerComment,
	}
	create_result := initializers.DB.Create(&car)

	if create_result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create car", "Asd": create_result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Car has been created successfully"})
}

// @Tags			Cars
// @Summary		Cars CRUD
// @Description	Get all cars
// @Accept			json
// @Produce		json
// @Success		200	{object}	[]models.Car
// @Failure		400	{object}	models.Error
// @Router			/cars/all [get]
func GetAllCars(c *gin.Context) {
	var cars []models.Car
	result := initializers.DB.Find(&cars)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get cars"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": cars})
}

// @Tags			Cars
// @Summary		Cars CRUD
// @Description	Get owned Cars
// @Accept			json
// @Produce		json
// @Success		200	{object}	[]models.Car
// @Failure		401	{object}	models.Error
// @Router			/cars/my [get]
func GetOwnedCars(c *gin.Context) {
	var cars []models.Car

	var err error
	token, err := c.Request.Cookie("Authorization")
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to get credentials"})
		return
	}

	claims, err := helpers.ExtractClaims(token.Value, []byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
	}

	var id = claims["sub"].(float64)

	log.Print(id)
	result := initializers.DB.Where("owner_id = ?", id).Find(&cars)
	log.Print(result)

	c.JSON(http.StatusOK, cars)
}

// @Tags			Cars
// @Summary		Cars CRUD
// @Description	Delete car by ID
// @Accept			json
// @Produce		json
// @Params			car_id path string "Car ID"
// @Success		200	{object}	[]models.Message
// @Failure		400	{object}	models.Error
// @Failure		401	{object}	models.Error
// @Router			/cars/delete/:id [delete]
func DeleteCarByID(c *gin.Context) {
	var car models.Car
	carID := c.Param("id")

	token, err := c.Request.Cookie("Authorization")
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to get credentials"})
		return
	}

	claims, err := helpers.ExtractClaims(token.Value, []byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to get user ID"})
		return
	}

	result := initializers.DB.Where("ID = ?", carID).Where("owner_id = ?", userID).Delete(&car)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "This car is not found"})
		return
	} else if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Car has been successfully deleted"})
}

// @Tags			Cars
// @Summary		Cars CRUD
// @Description	Update car by ID
// @Produce		json
// @Params			car_id path string "Car ID"
// @Success		200	{object}	[]models.Message
// @Failure		400	{object}	models.Error
// @Failure		401	{object}	models.Error
// @Router			/cars/update/:id [put]
func UpdateCarByID(c *gin.Context) {
	var requestBody map[string]interface{}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	carID := c.Param("id")
	car := models.Car{}

	token, err := c.Request.Cookie("Authorization")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to get credentials"})
		return
	}

	claims, err := helpers.ExtractClaims(token.Value, []byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to get user ID"})
		return
	}

	dbFields := map[string]interface{}{}
	for key, value := range requestBody {
		if key != "owner_id" && key != "vin_code" {
			dbFields[key] = value
		}
	}

	result := initializers.DB.Model(&car).Where("ID = ?", carID).Where("owner_id = ?", userID).Updates(dbFields)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "This car is not found"})
		return
	} else if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Car has been successfully updated"})
}
