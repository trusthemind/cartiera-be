package main

import (
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/trusthemind/go-cars-app/controllers"
	admin_controllers "github.com/trusthemind/go-cars-app/controllers/admin"
	_ "github.com/trusthemind/go-cars-app/docs"
	"github.com/trusthemind/go-cars-app/helpers"
	"github.com/trusthemind/go-cars-app/initializers"
	"github.com/trusthemind/go-cars-app/middleware"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
	initializers.SyncDB()
}

// *SWAGGER SETUP
//	@title			Cars Sales App API
//	@version		0.6
//	@description	This is documentation for Cars Sales App API for all user operations
//	@host			localhost:3000
//	@schemes		http

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {
	router := gin.Default()
	router.Static("/uploads", "/uploads")
	router.MaxMultipartMemory = 8 << 20

	corsConfig := cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://car-sales-app-v2.up.railway.app", },
		AllowMethods:     []string{"PUT", "POST", "DELETE", "GET"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Cookie"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000" || strings.HasPrefix(origin, "http://localhost:3000/")
		},
		
	})

	router.Use(corsConfig)

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
		engine.GET("", controllers.GetAllEngines)
		engine.POST("/create", middleware.RequireAuth, controllers.CreateEngine)
		engine.PUT("/update/:id", middleware.RequireAuth, controllers.UpdateEngineInfo)
		engine.DELETE("/delete/:id", controllers.DeleteEngineByID)
	}

	// *CAR
	cars := router.Group("/cars")
	{
		cars.GET("/all", controllers.GetAllCars)
		cars.POST("/create", middleware.RequireAuth, controllers.CreateCar)
		cars.GET("/my", middleware.RequireAuth, controllers.GetOwnedCars)
		cars.PUT("/update/:id", middleware.RequireAuth, controllers.UpdateCarByID)
		cars.DELETE("/delete/:id", middleware.RequireAuth, controllers.DeleteCarByID)
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

	admin := router.Group("/admin")
	{
		admin.GET("/users", middleware.RequireAdmin, middleware.RequireAuth, admin_controllers.GetAllUsers)
		admin.PUT("/users/update/:id", middleware.RequireAdmin, middleware.RequireAuth, admin_controllers.UpdateUserByID)
		admin.POST("/new-user", middleware.RequireAdmin, middleware.RequireAuth, admin_controllers.CreateNewUser)
		admin.DELETE("/delete/:id", middleware.RequireAdmin, middleware.RequireAuth, admin_controllers.DeleteUserbyID)
	}
	router.GET("photos/:photoName", helpers.GetPhoto)
	// !TEST
	router.POST("/auth/validate", middleware.RequireAuth, controllers.Validate)
	router.POST("/vincode/check", controllers.CheckVin)
	// router.POST("/posts/create", middleware.RequireAuth, controllers.CreatePost)

	// Run Server
	router.Run()
}
