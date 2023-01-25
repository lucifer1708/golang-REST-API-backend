package main

import (
	"fmt"
	"go-backend/controllers"
	"go-backend/middlewares"
	"go-backend/models"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println(os.Getenv("API_SECRET"))
	models.ConnectDB()
	r := gin.Default()
	public := r.Group("/api")
	public.POST("/register", controllers.Register)
	public.POST("/login", controllers.Login)
	public.GET("/users", controllers.AllUser)
	public.GET("/user/:id", controllers.GetUser)
	public.DELETE("/user/:id", controllers.DeleteUser)
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
