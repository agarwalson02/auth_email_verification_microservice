package main

import (
	"log"

	"api_gateway/internal/client"
	"api_gateway/internal/delivery/http"
	"api_gateway/internal/service"
	"api_gateway/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func main() {

	authClient, err := client.NewAuthClient("localhost:5000")
	if err != nil {
		log.Fatal(err)
	}

	authService := service.NewAuthService(authClient)
	handler := http.NewHandler(authService)

	r := gin.Default()

	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)

	protected := r.Group("/")
	protected.Use(middleware.AuthHeaderMiddleware())
	{
		protected.POST("/send-email", handler.SendEmail)
		protected.POST("/logout", handler.Logout)
	}



	log.Println("API Gateway running on port 8080")
	r.Run(":8080")
}
