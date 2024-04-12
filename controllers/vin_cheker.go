package controllers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-zoox/fetch"

	"github.com/trusthemind/go-cars-app/models"
)

//	@Tags			VIN
//	@Summary		VIN
//	@Description	Use VIN-code for more details
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.VINRequest	true	"VIN-code"
//	@Success		200		{object}	models.VINResponse
//	@Failure		404		{object}	models.Error
//	@Router			/vincode/check [post]
func CheckVin(c *gin.Context) {
	var RequestBody models.VINRequest

	if c.Bind(&RequestBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	str := os.Getenv("NINJA_URL") + "/vinlookup?vin=" + RequestBody.VIN
	key := os.Getenv("NINJA_KEY")
	response, err := fetch.Get(str, &fetch.Config{Headers: map[string]string{"X-Api-Key": key}})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to get data from VIN code"})
		return
	}

	var vinResponse models.VINResponse
	err = json.Unmarshal(response.Body, &vinResponse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse VIN response"})
		return
	}

	c.JSON(http.StatusOK, vinResponse)
}
