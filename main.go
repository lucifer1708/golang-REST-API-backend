package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-backend/controllers"
	_ "go-backend/docs"
	"go-backend/middlewares"
	"go-backend/models"
	"os"
	"time"
)

// @title           Blog app in Go
// @version         1.0
// @description     A Blog app made to publish my blogs in Go using Gin framework.
// @contact.name   Sumit Dhiman
// @contact.url    https://sumitdhiman.live
// @contact.email  sd08012003@gmail.com
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:1323
// @BasePath  /api
func main() {
	fmt.Println(os.Getenv("API_SECRET"))
	models.ConnectDB()
	r := gin.Default()
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	public := r.Group("/api")
	public.POST("/register", controllers.Register)
	public.POST("/login", controllers.Login)
	public.GET("/users", controllers.AllUser)
	public.GET("/user/:id", controllers.GetUser)
	public.DELETE("/user/delete/:id", controllers.DeleteUser)
	protected := r.Group("/api/admin")
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.GET("/user", controllers.CurrentUser)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "DELETE", "POST"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Run(":1323")
}
