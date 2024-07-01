package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name        string         `json:"name" gorm:"not null"`
	Email       string         `json:"email" gorm:"not null;unique"`
	Password    string         `json:"-" gorm:"not null"`
	Avatar      string         `json:"avatar" gorm:"not null"`
	CustomerID  string         `json:"customer_id" gorm:"not null;unique"`
	IsAdmin     bool           `json:"is_admin" gorm:"default:false"`
	PhoneNumber string         `json:"phone_number" gorm:"not null"`
	Telegram    string         `json:"tg_teg" gorm:"not null"`
	Liked       datatypes.JSON `json:"liked" gorm:"type:jsonb"`
}

type RequestRegistration struct {
	Name        string `json:"name" gorm:"not null" binding:"required"`
	Email       string `json:"email" gorm:"not null;unique" binding:"required"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password" gorm:"not null;min:8" binding:"required"`
}

type AdminRequestRegistration struct {
	*RequestRegistration
	IsCustomer bool `json:"is_customer" binding:"required"`
}

type RequestLogin struct {
	Email    string `json:"email" gorm:"not null;unique" binding:"required"`
	Password string `json:"password" gorm:"not null;min:8" binding:"required"`
}
