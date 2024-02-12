package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trusthemind/go-auth/initializers"
	"github.com/trusthemind/go-auth/models"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var RequstBody struct {
		Name     string
		Email    string
		Password string
	}
	if c.Bind(&RequstBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(RequstBody.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Name:     RequstBody.Name,
		Email:    RequstBody.Email,
		Password: string(hash)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func Login(c *gin.Context) {
	var RequstBody struct {
		Email    string
		Password string
	}
	if c.Bind(&RequstBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}
}
