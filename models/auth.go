package models

import "gorm.io/gorm"

type User struct {
	gorm.Model `json:"-"`
	Name       string `json:"name" gorm:"not null"`
	Email      string `json:"email" gorm:"not null;unique"`
	Password   string `json:"-" gorm:"not null"`
	Avatar     string `json:"avatar" gorm:"not null"`
	CustomerID string `json:"customer_id" gorm:"not null;unique"`
	IsAdmin    bool   `json:"is_admin" gorm:"default:false"`
	// Liked    []Car
}

type RequestRegistration struct {
	Name     string `json:"name" gorm:"not null" binding:"required"`
	Email    string `json:"email" gorm:"not null;unique" binding:"required"`
	Password string `json:"password" gorm:"not null;min:8" binding:"required"`
}

type AdminRequestRegistration struct {
	*RequestRegistration
	IsCustomer bool `json:"is_customer" binding:"required"`
}

type RequestLogin struct {
	Email    string `json:"email" gorm:"not null;unique" binding:"required"`
	Password string `json:"password" gorm:"not null;min:8" binding:"required"`
}
