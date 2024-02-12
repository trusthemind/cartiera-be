package main

import (
	"github.com/gin-gonic/gin"
	"github.com/trusthemind/go-auth/controllers"
	"github.com/trusthemind/go-auth/initializers"
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

	// Run Server
	r.Run()
}
