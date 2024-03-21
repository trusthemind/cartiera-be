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

type VINRequest struct {
	VIN string `json:"vin_code" binding:"required"`
}

type VINResponse struct {
	VIN          string `json:"vin"`
	Country      string `json:"country"`
	Manufacturer string `json:"manufacturer"`
	Region       string `json:"region"`
	WMI          string `json:"wmi"`
	VDS          string `json:"vds"`
	VIS          string `json:"vis"`
	Years        []int  `json:"years"`
}
