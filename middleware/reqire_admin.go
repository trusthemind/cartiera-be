package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/trusthemind/go-cars-app/helpers"
	"github.com/trusthemind/go-cars-app/initializers"
	"github.com/trusthemind/go-cars-app/models"
)

func RequireAdmin(c *gin.Context) {
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

	var user models.User
	initializers.DB.First(&user, claims["sub"])

	if user.ID == 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	} else if user.IsAdmin == false {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	// if user

	c.Set("user", user)
	c.Next()

	fmt.Println(claims["foo"], claims["nbf"])
}
