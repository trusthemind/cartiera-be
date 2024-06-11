package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/trusthemind/go-cars-app/helpers"
	"github.com/trusthemind/go-cars-app/initializers"
	"github.com/trusthemind/go-cars-app/models"
)

func UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse avatar file"})
		return
	}
	fileName := uuid.New().String() + filepath.Ext(file.Filename)

	if err := c.SaveUploadedFile(file, "uploads/"+fileName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save avatar file"})
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header not found"})
		return
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)


	claims, err := helpers.ExtractClaims(tokenString, []byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
	}

	var id = claims["sub"].(float64)

	avatarPath := "uploads/" + fileName
	if err := helpers.UpdateAvatar(uint(id), avatarPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update avatar path in the database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Avatar is updated successfully"})
}

func GetUserInfo(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header not found"})
		return
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	claims, err := helpers.ExtractClaims(tokenString, []byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to get credentials"})
		return
	}

	var id = claims["sub"].(float64)

	var user models.User
	initializers.DB.First(&user, "ID = ?", id)

	c.JSON(http.StatusOK, user)
}

func Validate(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["upload[]"]

	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No uploaded files provided"})
		return
	}

	result, ok, err := helpers.SavePhotoToTable(c, files)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": result, "type": reflect.TypeOf(result).String()})

	c.JSON(http.StatusOK, gin.H{"message": "Middleware is passed"})
}
