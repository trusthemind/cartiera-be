package helpers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)
func GetPhoto(c *gin.Context) {
    photoName := c.Param("photoName")
    photoPath := fmt.Sprintf("uploads/%s", photoName)

    // Check if file exists and open
    _, err := os.Stat(photoPath)
    if os.IsNotExist(err) {
        c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
        return
    }

    // Serve the file
    c.File(photoPath)
}
