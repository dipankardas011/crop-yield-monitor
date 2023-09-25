package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
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

func imageUpload(w http.ResponseWriter, r *http.Request) (int, error) {

	if r.Method != http.MethodPost {
		return http.StatusMethodNotAllowed, apiError{Status: http.StatusMethodNotAllowed, Err: "POST method is allowed"}
	}

	payload := ImageUpload{}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return http.StatusBadRequest, apiError{Status: http.StatusBadRequest, Err: err.Error()}
	}

	log.Println(payload)
	fileName := ""
	switch payload.Format {
	case "image/png":
		fileName = "image.png"
	case "image/jpeg":
		fileName = "image.jpeg"
	default:
		return http.StatusUnsupportedMediaType, apiError{Status: http.StatusUnsupportedMediaType, Err: UnSupportedMediaFormatType.String()}
	}

	if err := os.WriteFile(fileName, payload.RawImage, 0666); err != nil {
		return http.StatusInternalServerError, err
	}

	return writeJson(w, http.StatusOK, Response{Stdout: "fake response uploaded"})
}

func imageGet(w http.ResponseWriter, r *http.Request) (int, error) {

	if r.Method != http.MethodGet {
		return http.StatusMethodNotAllowed, apiError{Status: http.StatusMethodNotAllowed, Err: "GET method is allowed"}
	}

	payload := ImageGet{}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return http.StatusBadRequest, apiError{Status: http.StatusBadRequest, Err: err.Error()}
	}

	img := payload.Uuid // demo for image
	fmt.Println(img)

	return writeJson(w, http.StatusOK, Response{
		Stdout: "fake",
		Image: ImageGetResp{
			RawImage: []byte(img),
		},
	})
}

func Docs(w http.ResponseWriter, r *http.Request) (int, error) {

	if r.Method != http.MethodGet {
		return http.StatusMethodNotAllowed, apiError{Status: http.StatusMethodNotAllowed, Err: "GET method is allowed"}
	}

	docs := struct {
		Loc map[string]string
	}{
		Loc: map[string]string{
			"upload": "/image/upload",
			"get":    "/image/get",
			"TODO":   "about payloads",
		},
	}

	return writeJson(w, http.StatusOK, docs)
}

func Health(w http.ResponseWriter, r *http.Request) (int, error) {

	if r.Method != http.MethodGet {
		return http.StatusMethodNotAllowed, apiError{Status: http.StatusMethodNotAllowed, Err: "GET method is allowed"}
	}

	type Health struct {
		Msg string `json:"message"`
	}

	return writeJson(w, http.StatusOK, Health{Msg: "image looks healthy"})
}

func main() {

	http.HandleFunc("/image/upload", makeHTTPHandler(imageUpload))
	http.HandleFunc("/image/get", makeHTTPHandler(imageGet))
	http.HandleFunc("/image", makeHTTPHandler(Docs))
	http.HandleFunc("/image/healthz", makeHTTPHandler(Health))

	s := &http.Server{
		Addr:           ":8090",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Started to serve the image server on port {%v}\n", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}
