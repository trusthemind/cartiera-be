package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/trusthemind/go-cars-app/controllers"
	_ "github.com/trusthemind/go-cars-app/docs"
	"github.com/trusthemind/go-cars-app/initializers"
	"github.com/trusthemind/go-cars-app/middleware"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
	initializers.SyncDB()

}

// *SWAGGER SETTUP
// @title Cars Sales App API
// @version 0.6
// @description This is documentation for Cars Sales App API for all user operations
// @host localhost:3000
// @schemes http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	router := gin.Default()
	router.Static("/assets", "/assets")
	router.MaxMultipartMemory = 8 << 20

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
		payment_intent.GET("/all", middleware.RequireAuth, controllers.GetCustomerIntents)
		payment_intent.GET("/:id", middleware.RequireAuth, controllers.PaymentIntentByID)
		payment_intent.POST("/create", middleware.RequireAuth, controllers.CreatePaymentIntent)
		payment_intent.POST("/cancel", middleware.RequireAuth, controllers.CanceledPaymentIntent)
	}
	// !TEST
	router.GET("/auth/validate", middleware.RequireAuth, controllers.Validate)
	router.POST("/vincode/check", controllers.CheckVin)
	// router.POST("/posts/create", middleware.RequireAuth, controllers.CreatePost)

	// Run Server
	router.Run()
}
