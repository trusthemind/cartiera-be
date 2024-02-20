package main

import (
	"github.com/gin-gonic/gin"
	"github.com/trusthemind/go-auth/controllers"
	"github.com/trusthemind/go-auth/initializers"
	"github.com/trusthemind/go-auth/middleware"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
	initializers.SyncDB()

}

func main() {
	router := gin.Default()

	// AUTH
	router.POST("/auth/registration", controllers.Register)
	router.POST("/auth/login", controllers.Login)
	router.POST("/auth/logout", controllers.Logout)
	//using middleware for request
	router.GET("/auth/validate", middleware.RequireAuth, controllers.Validate)
	router.POST("/auth/validate", middleware.RequireAuth, controllers.Validate)

	// Run Server
	router.Run()
}
