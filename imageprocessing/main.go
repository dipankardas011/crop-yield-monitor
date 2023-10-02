package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	AUTH_SVR_URL = "http://auth:8080/account/token"
)

var (
	dbClient *ImageDBClient
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

func checkAuthenticUser(r *http.Request) (int, string, error) {
	// it should pass the token extracted from parent functions which call this

	// Get the JWT token from the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return http.StatusUnauthorized, "", apiError{Err: "Missing Authorization header", Status: http.StatusUnauthorized}
	}

	request, error := http.NewRequest(http.MethodGet, AUTH_SVR_URL, nil)
	if error != nil {
		return http.StatusInternalServerError, "", error
	}

	request.Header.Set("Authorization", r.Header.Get("Authorization"))

	client := &http.Client{}

	response, error := client.Do(request)
	if error != nil {
		return http.StatusInternalServerError, "", error
	}

	responseBody, error := io.ReadAll(response.Body)
	if error != nil {
		return http.StatusInternalServerError, "", error
	}

	payload := AuthResponse{}

	if err := json.Unmarshal(responseBody, &payload); err != nil {
		return http.StatusInternalServerError, "", apiError{Status: http.StatusInternalServerError, Err: err.Error()}
	}
	fmt.Println(payload.String())

	if response.StatusCode >= 300 {
		return response.StatusCode, "", apiError{Status: response.StatusCode, Err: payload.Error}
	}

	if len(payload.Error) > 0 {
		return response.StatusCode, "", apiError{Status: response.StatusCode, Err: payload.Error}
	}

	return http.StatusOK, payload.Stdout, nil
}

func imageUpload(w http.ResponseWriter, r *http.Request) (int, error) {

	if r.Method != http.MethodPost {
		return http.StatusMethodNotAllowed, apiError{Status: http.StatusMethodNotAllowed, Err: "POST method is allowed"}
	}

	payload := ImageUpload{}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return http.StatusBadRequest, apiError{Status: http.StatusBadRequest, Err: err.Error()}
	}

	status, username, err := checkAuthenticUser(r)
	if err != nil {
		return status, err
	}

	log.Println(payload)
	fileName := ""
	switch payload.Format {
	case "image/png":
		fileName = username + ".png"
	case "image/jpeg":
		fileName = username + ".jpeg"
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

	// get the username from the token from the auth server
	// img := payload.Uuid // demo for image

	status, _, err := checkAuthenticUser(r)
	if err != nil {
		return status, err
	}

	img := "<DUMMY>"
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

	dbClient = &ImageDBClient{}
	if err := dbClient.NewClient(); err != nil {
		panic(err)
	}

	log.Printf("Started to serve the image server on port {%v}\n", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}
