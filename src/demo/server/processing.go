package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func uploadImage(c *gin.Context) {
	var json struct {
		Image  []byte `json:"image"`
		Format string `json:"format"`
	}

	if err := c.BindJSON(&json); err != nil {
		return
	}

	fileName := "image.jpeg"
	if json.Format == "image/png" {
		fileName = "image.png"
	}

	// Write the image data to a file
	err := os.WriteFile(fileName, json.Image, 0666)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "image uploaded"})
}
