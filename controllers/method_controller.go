package controllers

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentmethod"
	"github.com/trusthemind/go-cars-app/helpers"
	"github.com/trusthemind/go-cars-app/initializers"
	"github.com/trusthemind/go-cars-app/models"
)

type BillingDetails struct {
	Country string `json:"country" binding:"required" gorm:"not null"`
	City    string `json:"city" binding:"required" gorm:"not null"`
	ALine   string `json:"address_line" binding:"required" gorm:"not null"`
}

type PaymentMethod struct {
	Number         string         `json:"number" gorm:"not null" binding:"required" min:"16"`
	Exp            string         `json:"exp" binding:"required" gorm:"not null" min:"4" max:"5"`
	CVC            string         `json:"cvc" binding:"required" gorm:"not null" min:"3" max:"3"`
	BillingDetails BillingDetails `json:"billing_details" binding:"required"`
	Phone          string         `json:"phone" binding:"required" gorm:"not null"`
}

func CreatePaymentMethod(c *gin.Context) {
	var RequestBody PaymentMethod

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
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
	}

	var user models.User
	var id = claims["sub"].(float64)
	initializers.DB.First(&user, "ID = ?", id)

	exp := strings.Split(RequestBody.Exp, "/")
	params := &stripe.PaymentMethodParams{Type: stripe.String(string(stripe.PaymentMethodTypeCard)),
		Card: &stripe.PaymentMethodCardParams{
			Number:   stripe.String(RequestBody.Number),
			ExpMonth: stripe.String(exp[0]),
			ExpYear:  stripe.String("20" + exp[1]),
			CVC:      stripe.String(RequestBody.CVC),
		},
		BillingDetails: &stripe.BillingDetailsParams{
			Name:  stripe.String(user.Email),
			Email: stripe.String(user.Email),
			Phone: stripe.String(RequestBody.Phone),
			Address: &stripe.AddressParams{
				Country: stripe.String(RequestBody.BillingDetails.Country),
				City:    stripe.String(RequestBody.BillingDetails.City),
				Line1:   stripe.String(RequestBody.BillingDetails.ALine),
			},
		},
	}

	result, err := paymentmethod.New(params)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create a payment method"})
		return
	}

	// attachParams := &stripe.PaymentMethodAttachParams{
	// Customer: stripe.String(user.CustomerID),
	// }

	// attachResult, err := paymentmethod.Attach("", attachParams)
	c.JSON(http.StatusOK, gin.H{"message": result})
	// c.JSON(http.StatusOK, gin.H{"message": "Payment method has created successfully"})

}
