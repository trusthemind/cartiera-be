package controllers

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/trusthemind/go-cars-app/helpers"
	"github.com/trusthemind/go-cars-app/initializers"
	"github.com/trusthemind/go-cars-app/models"

)

// @Tags			Authorization
// @Summary		Registration
// @Description	Register a new user
// @Accept			json
// @Produce		json
// @Param			request	body		models.RequestRegistration	true	"Name, Email, Password"
// @Success		200		{object}	models.Message
// @Failure		400		{object}	models.Error
// @Router			/auth/registration [post]
func Register(c *gin.Context) {
	stripe.Key = os.Getenv("STRIPE_KEY")
	var RequestBody models.RequestRegistration

	if c.Bind(&RequestBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(RequestBody.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to hash password"})
		return
	}
	params := &stripe.CustomerParams{
		Email: &RequestBody.Email,
		Name:  &RequestBody.Name,
	}

	customer, err := customer.New(params)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create customer"})
		return
	}
	log.Print(customer.ID)

	user := models.User{
		CustomerID: customer.ID,
		Name:       RequestBody.Name,
		Email:      RequestBody.Email,
		Password:   string(hash)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
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

// @Tags			Authorization
// @Summary		Login
// @Description	Login to app
// @Accept			json
// @Produce		json
// @Param			request	body		models.RequestLogin	true	"Email, Password"
// @Success		200		{object}	models.Message
// @Failure		400		{object}	models.Error
// @Router			/auth/login [post]
func Login(c *gin.Context) {
	var RequestBody struct {
		gorm.Model
		Email    string `gorm:"not null"`
		Password string `gorm:"not null"`
	}
	if c.Bind(&RequestBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}
	var user models.User
	initializers.DB.First(&user, "email = ?", RequestBody.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	//compare the req password with existing password using bcrypt
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(RequestBody.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Password"})
		return
	}

	// create a token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"name": user.Name,
		"exp":  time.Now().Add(time.Hour * 24).Unix(), //24 hours expire
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create token"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24, "", "", true, true)

	c.JSON(http.StatusOK, gin.H{"username": user.Name, "token": tokenString})
}

func Logout(c *gin.Context) {
	token, err := c.Request.Cookie("Authorization")
	if err != nil {
		// Handle missing cookie error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization cookie not found"})
		return
	}
	if token != nil {
		c.SetCookie("Authorization", "", time.Now().Hour()-1, "", "", false, true)

		c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
	}
}
