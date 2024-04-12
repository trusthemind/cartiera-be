package models

import (
	"gorm.io/gorm"
)

type PaymentIntentCreateRequest struct {
	// stripe usend only integers so first its need to * 100 to convert to int and reduce a float numbers
	Amount        float32 `json:"amount" gorm:"not null" biling:"required"`
	PaymentMethod string  `json:"payment_method" biling:"required" gorm:"not null"`
}

type PaymentIntentList struct {
	Length int             `json:"length"`
	Data   []PaymentIntent `json:"data"`
}

type PaymentIntentCard struct {
	ID string `json:"method_id" gorm:"not null"`
}
type PaymentIntent struct {
	gorm.Model
	StripeID     string  `json:"_id" gorm:"not null"`
	Status       string  `json:"status" gorm:"not null"`
	Currency     string  `json:"currency" gorm:"not null"`
	CustomerID   string  `json:"customer" gorm:"not null"`
	CanceledAt   int64   `json:"canceled" gorm:"not null"`
	Amount       float32 `json:"amount" gorm:"not null"`
	ClientSecret string  `json:"secret" gorm:"not null"`
}
