package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ErrorMsg string

func (e ErrorMsg) String() string {
	return string(e)
}

const (
	BadJsonFormat              ErrorMsg = "[Err] wrong json format"
	UnSupportedMediaFormatType ErrorMsg = "[Err] invalid image type supported are jpeg and png"
	InternalServerError        ErrorMsg = "[Err] internal server error"
)

func imageUpload(ctx *gin.Context) {

	payload := ImageUpload{}

	if err := ctx.BindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Stdout: BadJsonFormat.String(), Error: err.Error()})
		return
	}
	fmt.Println(payload)
	fileName := ""
	switch payload.Format {
	case "image/png":
		fileName = "image.png"
	case "image/jpeg", "image/jpg":
		fileName = "image.jpeg"
	default:
		ctx.JSON(http.StatusUnsupportedMediaType, Response{Stdout: UnSupportedMediaFormatType.String(), Error: UnSupportedMediaFormatType.String()})
		return
	}

	if err := os.WriteFile(fileName, payload.RawImage, 0666); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Stdout: InternalServerError.String(), Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, Response{Stdout: "fake response uploaded"})
}

func imageGet(ctx *gin.Context) {
	payload := ImageGet{}
	if err := ctx.BindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Error: err.Error(), Stdout: BadJsonFormat.String()})
		return
	}
	img := payload.Uuid // demo for image
	fmt.Println(img)
	ctx.JSON(http.StatusOK, Response{
		Stdout: "fake",
		Image: ImageGetResp{
			RawImage: []byte(img),
		},
	})
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(cors.Default())
	r.POST("/image/upload", imageUpload)
	r.GET("/image/get", imageGet)
	r.GET("/image/healthz", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "image looks healthy",
		})
	})
	r.GET("/image/docs", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"upload": "/image/upload",
			"get":    "/image/get",
			"TODO":   "about payloads",
		})
	})

	_ = r.Run(":8090")
}
