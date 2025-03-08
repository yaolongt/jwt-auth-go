package routes

import (
	"go-jwt/controllers"
	"go-jwt/middleware"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRoutes() {
	middleware := middleware.NewMiddleware()
	user := controllers.NewUserController()

	r := gin.Default()

	r.Use(gin.Recovery())

	cfg := cors.DefaultConfig()
	cfg.AllowHeaders = append(cfg.AllowHeaders, "Authorization")
	cfg.AllowAllOrigins = true
	r.Use(cors.New(cfg))

	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "helloworld",
	// 	})
	// })

	r.POST("/signup", user.SignUp)
	r.POST("/login", user.Login)
	r.GET("/validate", middleware.Auth, user.Validate)

	r.Run(":" + os.Getenv("PORT"))
}
