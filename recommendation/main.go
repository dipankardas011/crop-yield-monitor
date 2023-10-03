package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/cors"
)

var (
	dbClient *RecommendDBClient
)

const (
	AUTH_SVR_URL                 = "http://auth:8080/account/token"
	ML_SVR_URL                   = "http://ml:9090/start"
	BadJsonFormat       ErrorMsg = "[Err] wrong json format"
	InternalServerError ErrorMsg = "[Err] internal server error"
)

type ErrorMsg string

func (e ErrorMsg) String() string {
	return string(e)
}

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

func GetRecommendations(w http.ResponseWriter, r *http.Request) (int, error) {

	if r.Method != http.MethodGet {
		return http.StatusMethodNotAllowed, apiError{Status: http.StatusMethodNotAllowed, Err: "GET method is allowed"}
	}

	status, username, err := checkAuthenticUser(r)
	if err != nil {
		return status, err
	}

	fmt.Println(username)

	// first check with the key in the nosql db
	// if we didn't find it we will trigger the ml workload
	// and pass the status to the ml
	// ml will first write to the db that its in pending
	// TODO: first write for the noReady state when the record is not present
	// second if the record was present but with Status notReady we need to wait for ML

	recommend, err := dbClient.ReadRecommendations(username)
	if err != nil {
		if err == redis.Nil {
			// need to write to the db
			recommend := Recommendations{Status: RecommendationPending}

			if err := dbClient.WriteRecommendations(username, recommend); err != nil {
				return http.StatusInternalServerError, err
			}
		} else {
			return http.StatusInternalServerError, err
		}
	}

	// NOTE: redundant check for readability
	if recommend.Status == RecommendationPending && recommend.Status != RecommendationScheduled {
		// call the ML will be avoided (DUPLICATION of trigger can happen) for that Flag is there
		// WARN: Responsibility of the ML developer to handle these
		// NOTE: Need to decide whether the ML part will require another auth call or we simply pass the username as json body for it to handle rest
		// WARN: making assumption its PUT request
		// TODO: the ML will recieve the Auth Token and the username

		// Create a new PUT request
		req, err := http.NewRequest("PUT", ML_SVR_URL+"?username="+username, nil)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		// sharing the authorization link
		req.Header.Set("Authorization", r.Header.Get("Authorization"))

		// Use http.DefaultClient to send the request without waiting for the response
		_, err = http.DefaultClient.Do(req)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	return writeJson(w, http.StatusOK, Response{
		Stdout:          "recommendation for username " + username,
		Recommendations: *recommend,
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
			"[GET] to get the recommendations":           "/recommend/get",
			"[GET] to get the health":                    "/recommend/healthz",
			"[GET] to get read recommendation of data":   "/recommend/db/read",
			"[POST] to get write recommendation of data": "/recommend/db/write",
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

func DatabaseAccessWrite(w http.ResponseWriter, r *http.Request) (int, error) {

	if r.Method != http.MethodPost {
		return http.StatusMethodNotAllowed, apiError{Status: http.StatusMethodNotAllowed, Err: "POST method is allowed"}
	}

	status, username, err := checkAuthenticUser(r)
	if err != nil {
		return status, err
	}
	// decode from the r.Body
	recommend := Recommendations{}

	if err := json.NewDecoder(r.Body).Decode(&recommend); err != nil {
		return http.StatusBadRequest, apiError{Status: http.StatusBadRequest, Err: BadJsonFormat.String() + "\nReason: " + err.Error()}
	}

	if err := dbClient.WriteRecommendations(username, recommend); err != nil {
		return http.StatusInternalServerError, err
	}

	return writeJson(w, http.StatusOK, Response{Stdout: "Written to the database of username" + username})
}

func DatabaseAccessRead(w http.ResponseWriter, r *http.Request) (int, error) {

	if r.Method != http.MethodGet {
		return http.StatusMethodNotAllowed, apiError{Status: http.StatusMethodNotAllowed, Err: "GET method is allowed"}
	}

	status, username, err := checkAuthenticUser(r)
	if err != nil {
		return status, err
	}

	recommend, err := dbClient.ReadRecommendations(username)
	if err != nil {
		if err == redis.Nil {
			return http.StatusServiceUnavailable, fmt.Errorf("No Data present")
		} else {
			return http.StatusInternalServerError, err
		}
	}
	return writeJson(w, http.StatusOK, recommend)
}

func main() {

	RECOMMEND_SVR_URL = os.Getenv("DB_URL")
	PASS = os.Getenv("DB_PASSWORD")

	http.HandleFunc("/recommend/get", makeHTTPHandler(GetRecommendations))
	http.HandleFunc("/recommend", makeHTTPHandler(Docs))
	http.HandleFunc("/recommend/healthz", makeHTTPHandler(Health))
	http.HandleFunc("/recommend/db/read", makeHTTPHandler(DatabaseAccessRead))
	http.HandleFunc("/recommend/db/write", makeHTTPHandler(DatabaseAccessWrite))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},                      // Allow all origins
		AllowedMethods: []string{"GET", "POST", "OPTIONS"}, // Allow GET, OPTIONS methods
		AllowedHeaders: []string{"Authorization"},          // Allow Authorization header
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

	dbClient = &RecommendDBClient{}
	if err := dbClient.NewClient(); err != nil {
		panic(err)
	}
	log.Printf("Started to serve the recommend server on port {%v}\n", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}
