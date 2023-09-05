package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	// Use the default CORS middleware
	router.Use(cors.Default())

	router.POST("/upload", uploadImage)
	if err := router.Run("localhost:8080"); err != nil {
		panic(err)
	}
}
