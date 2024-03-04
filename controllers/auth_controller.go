package controllers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"github.com/trusthemind/go-cars-app/initializers"
	"github.com/trusthemind/go-cars-app/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(c *gin.Context) {
	stripe.Key = os.Getenv("STRIPE_KEY")
	var RequestBody struct {
		gorm.Model
		Name     string `gorm:"not null"`
		Email    string `gorm:"not null;unique"`
		Password string `gorm:"not null;min:8"`
	}

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
		"admin": user.IsAdmin,
		"exp":  time.Now().Add(time.Hour * 24).Unix(), //24 hours expire
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create token"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24, "", "", true, true)

	c.JSON(http.StatusOK, gin.H{})
}

func Logout(c *gin.Context) {
	token, err := c.Request.Cookie("Authorization")
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