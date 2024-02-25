package controllers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-zoox/fetch"
)
type VINResponse struct {
    VIN          string   `json:"vin"`
    Country      string   `json:"country"`
    Manufacturer string   `json:"manufacturer"`
    Region       string   `json:"region"`
    WMI          string   `json:"wmi"`
    VDS          string   `json:"vds"`
    VIS          string   `json:"vis"`
    Years        []int    `json:"years"`
}
func CheckVin(c *gin.Context) {
	var RequestBody struct {
		VIN string `json:"vin_code" binding:"required"`
	}

	if c.Bind(&RequestBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	str := os.Getenv("NINJA_URL") + "/vinlookup?vin=" + RequestBody.VIN
	key := os.Getenv("NINJA_KEY")
	response, err := fetch.Get(str, &fetch.Config{Headers: map[string]string{"X-Api-Key": key}})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get data from VIN code"})
		return
	}

	var vinResponse VINResponse
	err = json.Unmarshal(response.Body, &vinResponse);	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse VIN response"})
		return
	}

	c.JSON(http.StatusOK, vinResponse)
}
