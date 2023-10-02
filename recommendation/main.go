package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/rs/cors"
)

type ErrorMsg string

func (e ErrorMsg) String() string {
	return string(e)
}

const (
	BadJsonFormat       ErrorMsg = "[Err] wrong json format"
	InternalServerError ErrorMsg = "[Err] internal server error"
)

func GetRecommendations(w http.ResponseWriter, r *http.Request) (int, error) {

	if r.Method != http.MethodGet {
		return http.StatusMethodNotAllowed, apiError{Status: http.StatusMethodNotAllowed, Err: "GET method is allowed"}
	}

	payload := RecommendGet{}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return http.StatusBadRequest, apiError{Status: http.StatusBadRequest, Err: err.Error()}
	}
	log.Println(payload)

	return writeJson(w, http.StatusOK, Response{
		Stdout: "fake recommendation",
		Recommendations: Recommendations{
			Crops: []string{"fake01", "fake02"},
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
			"get":  "/recommend/get",
			"TODO": "about payloads",
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

	return writeJson(w, http.StatusOK, Health{Msg: "recommend looks healthy"})
}

func main() {

	http.HandleFunc("/recommend/get", makeHTTPHandler(GetRecommendations))
	http.HandleFunc("/recommend", makeHTTPHandler(Docs))
	http.HandleFunc("/recommend/healthz", makeHTTPHandler(Health))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},              // Allow all origins
		AllowedMethods: []string{"GET", "OPTIONS"}, // Allow GET, OPTIONS methods
		AllowedHeaders: []string{"Authorization"},  // Allow Authorization header
		// AllowCredentials: true,
		Debug: true,
	})

	s := &http.Server{
		Addr:           ":8100",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        c.Handler(http.DefaultServeMux),
	}

	log.Printf("Started to serve the recommend server on port {%v}\n", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}
