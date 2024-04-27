package controllers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/trusthemind/go-cars-app/helpers"
	"github.com/trusthemind/go-cars-app/initializers"
	"github.com/trusthemind/go-cars-app/models"
)

func CreateDetail(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["upload[]"]
	data := c.Request.FormValue("data")

	var RequestData models.Detail
	json.Unmarshal([]byte(data), &RequestData)

	err := c.Bind(&RequestData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No uploaded files provided"})
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

	result, ok, err := helpers.SavePhotoToTable(c, files)
	if !ok {
        c.JSON(http.StatusBadRequest, gin.H{"error": err})
        return
    }
	
	detail:= models.Detail{
		Name: RequestData.Name,
		Price: RequestData.Price,
		Condiiton: RequestData.Condiiton,
		OwnerID: int(user.ID),
		Photos: result,
	}

	create_result := initializers.DB.Create(&detail)

	if create_result.Error != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create detail", "Asd": create_result.Error.Error()})
        return
    }

	c.JSON(http.StatusOK, gin.H{"message": "Detail has been created successfully"})
}

func GetAllDetails(c *gin.Context) {
	var details []models.Detail

	result := initializers.DB.Find(&details)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get details"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": details})
}
