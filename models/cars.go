package models

import (
	"gorm.io/gorm"

)

type Car struct {
	gorm.Model
	OwnerID      int `json:"owner_id" gorm:"not null"`
	EngineID     int
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

// !Electrical version of engine Double check
type Engine struct {
	gorm.Model
	Brand       string  `json:"brand" gorm:"not null"`
	Name        string  `json:"name" gorm:"not null"`
	Fuel        string  `json:"fuel" gorm:"default:0"`
	Cilinders   int32   `json:"ciliders" gorm:"default:0"`
	Consumption float32 `json:"consumption" gorm:"default:0"`
}
