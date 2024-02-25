package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `form:"name" gorm:"not null"`
	Email    string `form:"email" gorm:"not null;unique"`
	Password string `form:"password" gorm:"not null"`
	Avatar   string `form:"avatar" gorm:"not null"`
	IsAdmin  bool   `form:"is_admin" default:"false"`
	Liked    []Car
}

type Car struct {
	gorm.Model
	Owner        User
	Engine       Engine
	OwerComment  string `form:"comment" gorm:"default:null"`
	Year         int32  `form:"year" gorm:"not null"`
	OwnersNumber int32  `form:"owners_number" gorm:"not null"`
	Price        int64  `form:"price" gorm:"not null"`
	Kilometers   int64  `form:"kilometers" gorm:"not null"`
	Brand        string `form:"brand" gorm:"not null"`
	Status       string `form:"status" gorm:"not null"`
	CarModel     string `form:"model" gorm:"not null"`
	VinCode      string `form:"vincode" gorm:"not null;unique"`
	Placement    string `form:"placement" gorm:"not null"`
}

type Engine struct {
	gorm.Model
	Name        string
	Fuel        string
	Cilinders   int32
	Consumption float32
}

type UserPost struct {
	gorm.Model
	UserID      float64 `gorm:"not null"`
	UserName    string  `gorm:"not null"`
	Title       string  `gorm:"not null"`
	Description string  `gorm:"not null"`
}
