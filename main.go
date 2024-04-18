package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/timorodr/server/controllers"
	"github.com/timorodr/server/initializers"
	"github.com/timorodr/server/middleware"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
}
func main() {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(cors.Default())

	r.GET("/", controllers.HelloWorldHandler)
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.Run() // listen and serve on 0.0.0.0:8080
}
