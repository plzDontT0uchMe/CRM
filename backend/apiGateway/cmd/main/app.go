package main

import (
	"CRM/go/apiGateway/internal/handlers"
	"CRM/go/apiGateway/internal/middleware/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(cors.CORSMiddleware())

	r.POST("/api/auth", handlers.Authorization)
	r.POST("/api/reg", handlers.Registration)
	r.POST("/api/checkAuth", handlers.CheckAuthorization)

	r.Run(":3000")
}
