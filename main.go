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
	r := gin.Default()

	// AUTH
	r.POST("/auth/registration", controllers.Register)
	r.POST("/auth/login", controllers.Login)
	r.POST("/auth/logout", controllers.Logout)
	//using middleware for request
	r.GET("/auth/validate", middleware.RequireAuth, controllers.Validate)
	r.POST("/auth/validate", middleware.RequireAuth, controllers.Validate)

	// Run Server
	r.Run()
}
