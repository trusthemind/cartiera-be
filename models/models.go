package models

import (
	// "mime/multipart"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name       string `form:"name" gorm:"not null"`
	Email      string `form:"email" gorm:"not null;unique"`
	Password   string `form:"password" gorm:"not null"`
	Avatar     string `form:"avatar" gorm:"not null"`
	CustomerID string `form:"customer_id" gorm:"not null;unique"`
	IsAdmin    bool   `form:"is_admin" default:"false"`
	// Liked    []Car
}

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

// fix user input datas for example not null and all of this

type Engine struct {
	gorm.Model
	Name        string
	Fuel        string
	Cilinders   int32
	Consumption float32
}
