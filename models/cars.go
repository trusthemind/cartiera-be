package models

import (
	"gorm.io/gorm"
)

type Car struct {
	gorm.Model
	// Engine       Engine
	OwnerID      int
	OwnerComment string `json:"comment" gorm:"default:null"`
	Year         int32  `json:"year" gorm:"not null"`
	OwnersNumber int32  `json:"owners_number" gorm:"not null"`
	Price        int64  `json:"price" gorm:"not null"`
	Kilometers   int64  `json:"kilometers" gorm:"not null"`
	Brand        string `json:"brand" gorm:"not null"`
	Status       string `json:"status" gorm:"not null"`
	CarModel     string `json:"model" gorm:"not null"`
	VinCode      string `json:"vin_code" gorm:"not null;unique"`
	Placement    string `json:"placement" gorm:"not null"`
}

type Engine struct {
	gorm.Model
	Name        string
	Fuel        string
	Cilinders   int32
	Consumption float32
}

// TODO finish a saving in DB with proper form fields
// type PaymentIntent struct {
// 	gorm.Model
// }