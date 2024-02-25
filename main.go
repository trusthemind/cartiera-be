package main

import (
	"github.com/gin-gonic/gin"
	"github.com/trusthemind/go-cars-app/controllers"
	"github.com/trusthemind/go-cars-app/initializers"
	"github.com/trusthemind/go-cars-app/middleware"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
	initializers.SyncDB()

}

func main() {
	router := gin.Default()
	router.Static("/assets", "/assets")
	router.MaxMultipartMemory = 8 << 20

	// *AUTH
	auth := router.Group("/auth")
	{
		auth.POST("/registration", controllers.Register)
		auth.POST("/login", controllers.Login)
		auth.POST("/logout", controllers.Logout)
	}

	// *ENGINE
	engine := router.Group("/engine")
	{
		engine.POST("/create", controllers.CreateEngine)
	}

	// !TEST
	//using middleware for request
	router.GET("/auth/validate", middleware.RequireAuth, controllers.Validate)
	router.POST("/vincode/check", controllers.CheckVin)
	// router.POST("/posts/create", middleware.RequireAuth, controllers.CreatePost)

	// Run Server
	router.Run()
}
