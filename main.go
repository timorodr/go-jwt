package main

import (
	"github.com/gin-gonic/gin"
	"github.com/timorodr/server/controllers"
	"github.com/timorodr/server/initializers"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
}
func main() {
	r := gin.Default()

	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.Run() // listen and serve on 0.0.0.0:8080
}
