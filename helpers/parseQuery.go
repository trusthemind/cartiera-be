package helpers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)
func ParseQueryParam(c *gin.Context, param string, defaultValue int) (int, error) {
	paramStr := c.DefaultQuery(param, strconv.Itoa(defaultValue))
	value, err := strconv.Atoi(paramStr)
	if err != nil || value <= 0 {
		return defaultValue, err
	}
	return value, nil
}