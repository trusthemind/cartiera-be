package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/trusthemind/go-cars-app/helpers"
	"github.com/trusthemind/go-cars-app/initializers"
	"github.com/trusthemind/go-cars-app/models"
	"gorm.io/gorm"
)

func CreatePost(c *gin.Context) {
	var RequestBody struct {
		gorm.Model
		Title       string `gorm:"not null"`
		Description string `gorm:"not null"`
	}

	if c.Bind(&RequestBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	token, err := c.Request.Cookie("Authorization")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization cookie not found"})
		return
	}

	claims, err := helpers.ExtractClaims(token.Value, []byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to extract claims"})
	}

	var name = claims["name"].(string)
	var id = claims["sub"].(float64)

	post := models.UserPost{
		UserID:      id,
		UserName:    name,
		Title:       RequestBody.Title,
		Description: RequestBody.Description,
	}
	result := initializers.DB.Create(&post)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
