package models

type Error struct {
	Error string `json:"error"`
}

type Message struct {
	Message string `json:"message"`
}

type RequestRegistration struct {
	Name     string `json:"name" gorm:"not null" binding:"required"`
	Email    string `json:"email" gorm:"not null;unique" binding:"required"`
	Password string `json:"password" gorm:"not null;min:8" binding:"required"`
}

type RequestLogin struct {
	Email    string `json:"email" gorm:"not null;unique" binding:"required"`
	Password string `json:"password" gorm:"not null;min:8" binding:"required"`
}
