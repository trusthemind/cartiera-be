package models


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
