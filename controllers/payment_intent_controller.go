package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
	"github.com/trusthemind/go-cars-app/helpers"
	"github.com/trusthemind/go-cars-app/initializers"
	"github.com/trusthemind/go-cars-app/models"
)

func CreatePaymentIntent(c *gin.Context) {
	var RequestBody struct {
		// stripe usend only integers so first its need to * 100 to convert to int and reduce a float numbers
		Amount        float32 `json:"amount" gorm:"not null" biling:"required"`
		PaymentMethod string  `json:"payment_method" biling:"required" gorm:"not null"`
	}

	if c.Bind(&RequestBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	token, err := c.Request.Cookie("Authorization")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization cookies not found"})
		return
	}

	claims, err := helpers.ExtractClaims(token.Value, []byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}

	var user models.User
	var id = claims["sub"].(float64)
	var amount = helpers.ConvertAndRound(RequestBody.Amount, true)
	initializers.DB.First(&user, "ID = ?", id)

	params := &stripe.PaymentIntentParams{
		Amount:        stripe.Int64(amount),
		Currency:      stripe.String(string(stripe.CurrencyUSD)),
		Customer:      stripe.String(user.CustomerID),
		PaymentMethod: stripe.String(RequestBody.PaymentMethod),
	}

	result, err := paymentintent.New(params)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create a payment intent"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}
