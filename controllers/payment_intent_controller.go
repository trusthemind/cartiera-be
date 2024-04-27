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

// @Tags			Payment Intent
// @Summary		Payment Intent Operation
// @Description	Create Payment Intent
// @Accept			json
// @Produce		json
// @Param			request	body		models.PaymentIntent	true	"Email, Password"
// @Success		200		{object}	models.Message
// @Failure		400		{object}	models.Error
// @Failure		401		{object}	models.Error
// @Router			/paymnet_intent/create [post]
func CreatePaymentIntent(c *gin.Context) {
	var RequestBody models.PaymentIntentCreateRequest

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
	initializers.DB.First(&user, "ID = ?", id)
	var amount = helpers.ConvertAndRound(RequestBody.Amount, true)

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
	confirm_params := &stripe.PaymentIntentConfirmParams{
		PaymentMethod: stripe.String(RequestBody.PaymentMethod),
	}

	result, _ = paymentintent.Confirm(result.ID, confirm_params)

	if result.Status == "success" {
		c.JSON(http.StatusOK, gin.H{"message": "Payment intent is successfully confirmed"})
		return
	}

	paymentIntent := models.PaymentIntent{
		StripeID:     result.ID,
		CustomerID:   result.Customer.ID,
		ClientSecret: result.ClientSecret,
		Amount:       result.Amount,
		CanceledAt:   result.CanceledAt,
		Currency:     result.Currency,
		Status:       string(result.Status),
	}
	save_result := initializers.DB.Create(&paymentIntent)

	if save_result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to save a payment intent"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment intent is successfully created"})
}

// @Tags			Payment Intent
// @Summary		Payment Intent Operation
// @Description	Get Customers Payment Intents
// @Accept			json
// @Produce		json
// @Success		200	{object}	[]models.PaymentIntentList
// @Failure		400	{object}	models.Error
// @Failure		401	{object}	models.Error
// @Router			/paymnet_intent/all [get]
func GetCustomerIntents(c *gin.Context) {
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
	initializers.DB.First(&user, "ID =?", id)

	params := &stripe.PaymentIntentListParams{
		Customer: stripe.String(user.CustomerID),
	}

	var paymentIntents []*stripe.PaymentIntent
	var responceArray []*models.PaymentIntent
	var responceItem *models.PaymentIntent
	interator := paymentintent.List(params)

	for interator.Next() {
		item := interator.PaymentIntent()
		responceItem = &models.PaymentIntent{
			StripeID:   item.ID,
			Status:     string(item.Status),
			Currency:   item.Currency,
			CustomerID: item.Customer.ID,
			CanceledAt: item.CanceledAt,
			Amount:     item.Amount,
		}
		paymentIntents = append(paymentIntents, item)
		responceArray = append(responceArray, responceItem)
	}

	if err := interator.Err(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"length": len(responceArray), "data": responceArray})
}

// @Tags			Payment Intent
// @Summary		Payment Intent Operation
// @Description	Get Payment Intent by ID
// @Accept			json
// @Produce		json
// @Param			payment_id	path		string	true	"Payment Intent ID"
// @Success		200			{object}	models.PaymentIntent
// @Failure		400			{object}	models.Error
// @Router			/paymnet_intent/:id [get]
func PaymentIntentByID(c *gin.Context) {
	var payment_id = c.Param("id")

	params := &stripe.PaymentIntentParams{}

	result, err := paymentintent.Get(payment_id, params)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment intent ID"})
	}

	c.JSON(http.StatusOK, result)
}

// @Tags			Payment Intent
// @Summary		Payment Intent Operation
// @Description	Cancel Payment Intent
// @Accept			json
// @Produce		json
// @Param			request	body		string	true	"id"
// @Success		200		{object}	models.Message
// @Failure		400		{object}	models.Error
// @Router			/paymnet_intent/cancel [post]
func CanceledPaymentIntent(c *gin.Context) {
	var RequestBody struct {
		ID string `json:"payment_intent_id" gorm:"not null" binding:"required"`
	}

	if c.Bind(&RequestBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	params := &stripe.PaymentIntentCancelParams{}

	_, err := paymentintent.Cancel(RequestBody.ID, params)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment intent has been canceled"})
}
