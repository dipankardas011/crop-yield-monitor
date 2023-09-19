package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SignUp(ctx *gin.Context) {
	account := AccountSignUp{}
	if err := ctx.BindJSON(&account); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Bad Request format, [%v]", err),
		})
	}
	fmt.Println(account)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "account created successfully",
	})
}

func SignIn(ctx *gin.Context) {
	account := AccountSignIn{}
	if err := ctx.BindJSON(&account); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Bad Request format, [%v]", err),
		})
	}
	fmt.Println(account)

	ctx.JSON(http.StatusOK, gin.H{
		"usertoken": "32qwe32413212211(dummy)",
		"message":   "signin successful",
	})
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(cors.Default())
	r.POST("/account/signin", SignIn)
	r.POST("/account/signup", SignUp)
	r.GET("/account/healthz", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "auth looks healthy",
		})
	})
	r.GET("/account/docs", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"signin": "/account/signin",
			"signup": "/account/signup",
			"TODO":   "about payloads",
		})
	})

	_ = r.Run(":8080")
}
