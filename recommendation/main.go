package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ErrorMsg string

func (e ErrorMsg) String() string {
	return string(e)
}

const (
	BadJsonFormat       ErrorMsg = "[Err] wrong json format"
	InternalServerError ErrorMsg = "[Err] internal server error"
)

func getRecommendations(ctx *gin.Context) {
	payload := RecommendGet{}

	if err := ctx.BindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Stdout: BadJsonFormat.String(), Error: err.Error()})
		return
	}
	fmt.Println(payload)

	ctx.JSON(http.StatusOK, Response{Stdout: "fake recommendation", Recommendations: Recommendations{
		Crops: []string{"fake01", "fake02"},
	}})
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(cors.Default())
	r.GET("/recommend/get", getRecommendations)
	r.GET("/recommend/healthz", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "recommend looks healthy",
		})
	})
	r.GET("/recommend/docs", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"get":  "/recommendations/get",
			"TODO": "about payloads",
		})
	})

	_ = r.Run(":8100")
}
