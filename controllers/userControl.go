package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
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
	var user models.User
	initializers.DB.First(&user, "email = ?", RequstBody.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	//compare the req password with existing password using bcrypt
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(RequstBody.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Password"})
		return
	}

	// create a token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 4).Unix(), //4 hours expire
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create token"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*4, "", "", true, true)

	c.JSON(http.StatusOK, gin.H{})
}

func Logout(c *gin.Context) {
	token, err := c.Request.Cookie("Authorization")
	fmt.Print("aaaaaaaaaaaaaa", err, token)
	if err != nil {
		// Handle missing cookie error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization cookie not found"})
		return
	}
	if token != nil {
		c.SetCookie("Authorization", "", 0, "", "", true, true)

		c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
	}
}
