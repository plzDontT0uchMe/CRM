package main

import (
	"CRM/go/apiGateway/internal/config"
	"CRM/go/apiGateway/internal/handlers"
	"CRM/go/apiGateway/internal/middleware/cors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
)

func main() {
	cfg := config.GetConfig()
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		fmt.Println(err)
		panic("error reading env")
	}
	config.SetConfig(&cfg)

	fmt.Println(cfg)

	r := gin.Default()

	r.Use(cors.CORSMiddleware())

	r.POST("/api/auth", handlers.Authorization)
	r.POST("/api/reg", handlers.Registration)
	r.POST("/api/checkAuth", handlers.CheckAuthorization)

	r.Run(cfg.Host + ":" + cfg.Port)
}
