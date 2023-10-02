package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/rs/cors"
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

	if response.StatusCode >= 300 {
		return response.StatusCode, "", apiError{Status: response.StatusCode, Err: payload.Error}
	}

	return http.StatusOK, payload.Stdout, nil
}

func imageUpload(w http.ResponseWriter, r *http.Request) (int, error) {

	if r.Method != http.MethodPost {
		return http.StatusMethodNotAllowed, apiError{Status: http.StatusMethodNotAllowed, Err: "POST method is allowed"}
	}

	payload := Image{}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return http.StatusBadRequest, apiError{Status: http.StatusBadRequest, Err: err.Error()}
	}

	status, username, err := checkAuthenticUser(r)
	if err != nil {
		return status, err
	}

	if payload.Format != "image/png" && payload.Format != "image/jpeg" {
		return http.StatusUnsupportedMediaType, apiError{Status: http.StatusUnsupportedMediaType, Err: UnSupportedMediaFormatType.String()}
	}

	if err := dbClient.WriteImage(username, payload); err != nil {
		return http.StatusInternalServerError, err
	}

	return writeJson(w, http.StatusOK, Response{Stdout: "fake response uploaded"})
}

func imageGet(w http.ResponseWriter, r *http.Request) (int, error) {

	if r.Method != http.MethodGet {
		return http.StatusMethodNotAllowed, apiError{Status: http.StatusMethodNotAllowed, Err: "GET method is allowed"}
	}

	status, username, err := checkAuthenticUser(r)
	if err != nil {
		return status, err
	}

	img, err := dbClient.ReadImage(username)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return writeJson(w, http.StatusOK, Response{
		Stdout: "fake",
		Image:  *img,
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
			"[POST] upload a new image":      "/image/upload",
			"[GET] get the already uploaded": "/image/get",
			"[GET] health of the server":     "/image/healthz",
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

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},                      // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"}, // Allow GET, POST, and OPTIONS methods
		AllowedHeaders:   []string{"Authorization"},          // Allow Authorization header
		AllowCredentials: true,
		Debug:            true,
	})

	s := &http.Server{
		Addr:           ":8090",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        c.Handler(http.DefaultServeMux),
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
