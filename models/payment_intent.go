package models

import "github.com/stripe/stripe-go"

type PaymentIntentCreateRequest struct {
	// stripe usend only integers so first its need to * 100 to convert to int and reduce a float numbers
	Amount        float32 `json:"amount" gorm:"not null" biling:"required"`
	PaymentMethod string  `json:"payment_method" biling:"required" gorm:"not null"`
}

type PaymentIntentList struct {
	Length int `json:"length"`
	Data []stripe.PaymentIntent
}
