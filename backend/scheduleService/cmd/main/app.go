package main

import "github.com/gin-gonic/gin"

func Hello(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello World",
	})
}

func main() {
	server := gin.Default()
	server.Use(func(context *gin.Context) {
		context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	})
	server.GET("/hello123", Hello)
	//r.POST("/proxy/reg", database.PostUser)
	server.Run(":3002")
}
