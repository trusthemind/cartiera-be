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

	// *USERS
	users := router.Group("/users")
	{
		users.PUT("/avatar/update", middleware.RequireAuth, controllers.UploadAvatar)
	}

	// *ENGINE
	engine := router.Group("/engine")
	{
		engine.POST("/create", controllers.CreateEngine)
	}

	// *CAR
	cars := router.Group("/cars")
	{
		cars.GET("/all", controllers.GetAllCars)
		cars.POST("/create", middleware.RequireAuth, controllers.CreateCar)
		cars.GET("/my", middleware.RequireAuth, controllers.GetOwnedCars)
	}

	payment_method := router.Group("/payment_method")
	{
		payment_method.POST("/create", middleware.RequireAuth, controllers.CreatePaymentMethod)
		payment_method.GET("/all", middleware.RequireAuth, controllers.GetAllPaymentMethod)

	}

	payment_intent := router.Group("/payment_intent")
	 {
		payment_intent.POST("/create", middleware.RequireAuth, controllers.CreatePaymentIntent)
	}
	// !TEST
	router.GET("/auth/validate", middleware.RequireAuth, controllers.Validate)
	router.POST("/vincode/check", controllers.CheckVin)
	// router.POST("/posts/create", middleware.RequireAuth, controllers.CreatePost)

	// Run Server
	router.Run()
}
