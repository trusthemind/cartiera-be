package controllers

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/trusthemind/go-cars-app/helpers"
	"github.com/trusthemind/go-cars-app/initializers"
	"github.com/trusthemind/go-cars-app/models"

)

// @Tags			Cars
// @Summary		Cars
// @Description	Create a car for sale
// @Accept			json
// @Produce		json
// @Param			request	body		models.Car	true	"Car"
// @Success		200		{object}	models.Message
// @Failure		400		{object}	models.Error
// @Failure		401		{object}	models.Error
// @Router			/cars/create [post]
func CreateCar(c *gin.Context) {
	// var RequestBody struct {
	// 	Brand        string `json:"brand" gorm:"not null"`
	// 	CarModel     string `json:"model" gorm:"not null"`
	// 	Year         int32  `json:"year" gorm:"not null"`
	// 	Price        int64  `json:"price" gorm:"not null"`
	// 	Status       string `json:"status" gorm:"not null"`
	// 	VinCode      string `json:"vin_code" gorm:"not null"`
	// 	Kilometers   int64  `json:"kilometers" default:"0"`
	// 	Placement    string `json:"placement" gorm:"not null"`
	// 	OwnersNumber int32  `json:"owners" default:"0"`
	// 	OwnerComment string `json:"comment"`
	// }
	var RequestBody models.Car

	if c.Bind(&RequestBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	token, err := c.Request.Cookie("Authorization")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization cookie not found"})
		return
	}

	claims, err := helpers.ExtractClaims(token.Value, []byte(os.Getenv("SECRET_KEY")))

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

	car := models.Car{
		OwnerID:      int(user.ID),
		EngineID:     RequestBody.EngineID,
		Brand:        RequestBody.Brand,
		CarModel:     RequestBody.CarModel,
		Year:         RequestBody.Year,
		Price:        RequestBody.Price,
		Status:       RequestBody.Status,
		VinCode:      RequestBody.VinCode,
		Kilometers:   RequestBody.Kilometers,
		Placement:    RequestBody.Placement,
		OwnersNumber: RequestBody.OwnersNumber,
		OwnerComment: RequestBody.OwnerComment,
	}
	result := initializers.DB.Create(&car)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create car"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Car has been created successfully"})
}

// @Tags			Cars
// @Summary		Cars
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
	c.JSON(http.StatusOK, gin.H{"cars": cars})
}

// @Tags			Cars
// @Summary		Cars
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
