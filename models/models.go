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
}

type UserPost struct {
	gorm.Model
	UserID      float64 `gorm:"not null"`
	UserName    string  `gorm:"not null"`
	Title       string  `gorm:"not null"`
	Description string  `gorm:"not null"`
}
