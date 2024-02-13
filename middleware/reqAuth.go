package middleware

import "github.com/gin-gonic/gin"

func RequireAuth(c *gin.Context) {
	// прокидування мідвари на контролер
	c.Next()
}