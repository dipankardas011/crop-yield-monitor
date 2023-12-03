package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	AUTH_SVR_URL        string   = "http://auth:8080/account/token"
	IMG_SVR_URL         string   = "http://image:8090/image/get"
	ML_SVR_URL          string   = "http://ml:9090/start"
	BadJsonFormat       ErrorMsg = "[Err] wrong json format"
	InternalServerError ErrorMsg = "[Err] internal server error"
)

type ErrorMsg string

func (e ErrorMsg) String() string {
	return string(e)
}

func checkAuthenticUser(r *http.Request) (int, string, error) {
	// it should pass the token extracted from parent functions which call this

	payload := AuthResponse{}
	stat, err := makeHTTPCall(r, http.MethodGet, AUTH_SVR_URL, nil, &payload)
	if err != nil {
		return stat, "", err
	}

	return http.StatusOK, payload.Stdout, nil
}

func getImage(r *http.Request) (int, *Image, error) {

	payload := ImgResponse{}

	stat, err := makeHTTPCall(r, http.MethodGet, AUTH_SVR_URL, nil, &payload)
	if err != nil {
		return stat, nil, err
	}

	return http.StatusOK, &payload.Image, nil
}

func GetRecommendations(w http.ResponseWriter, r *http.Request) (int, error) {

	if r.Method != http.MethodGet {
		return http.StatusMethodNotAllowed, apiError{Status: http.StatusMethodNotAllowed, Err: "GET method is allowed"}
	}

	status, username, err := checkAuthenticUser(r)
	if err != nil {
		return status, err
	}

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
			recommend = &Recommendations{Status: RecommendationNotReady}

			if err := dbClient.WriteRecommendations(username, *recommend); err != nil {
				return http.StatusInternalServerError, err
			}
		} else {
			return http.StatusInternalServerError, err
		}
	}

	// NOTE: redundant check for readability
	if recommend.Status == RecommendationNotReady && recommend.Status != RecommendationScheduled {
		// call the ML will be avoided (DUPLICATION of trigger can happen) for that Flag is there
		// WARN: Responsibility of the ML developer to handle these
		// NOTE: Need to decide whether the ML part will require another auth call or we simply pass the username as json body for it to handle rest
		// WARN: making assumption its PUT request
		// TODO: the ML will recieve the Auth Token and the username

		// get the image
		stat, rawImg, err := getImage(r)
		if err != nil {
			return stat, err
		}

		rawImgPayload, err := json.Marshal(rawImg)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		// Create a new PUT request
		req, err := http.NewRequest("PUT", ML_SVR_URL+"?username="+username, bytes.NewBuffer(rawImgPayload))
		if err != nil {
			return http.StatusInternalServerError, err
		}

		userCookie, err := r.Cookie("user_token")
		if err != nil {
			if err == http.ErrNoCookie {
				return http.StatusUnauthorized, apiError{Err: "Missing Cookie", Status: http.StatusUnauthorized}
			}
			return http.StatusInternalServerError, apiError{Err: err.Error(), Status: http.StatusInternalServerError}
		}

		// Set content type to application/json
		req.Header.Set("Content-Type", "application/json")
		// sharing the authorization link
		req.Header.Set("Authorization", "Bearer "+userCookie.Value)

		// Use http.DefaultClient to send the request without waiting for the response

		// TODO: once the ml is ready do add it
		_, err = http.DefaultClient.Do(req)
		// if err != nil {
		// 	return http.StatusInternalServerError, err
		// }
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

// NOTE: THIS SERVER HAS NOT BEEN TESTED!!!
func main() {

	RECOMMEND_SVR_URL = os.Getenv("DB_URL")
	PASS = os.Getenv("DB_PASSWORD")

	http.HandleFunc("/recommend/get", makeHTTPHandler(GetRecommendations)) // User-facing
	http.HandleFunc("/recommend/docs", makeHTTPHandler(Docs))              // User-facing
	http.HandleFunc("/recommend/healthz", makeHTTPHandler(Health))         // User-facing
	http.HandleFunc("/recommend/db/read", makeHTTPHandler(DatabaseAccessRead))
	http.HandleFunc("/recommend/db/write", makeHTTPHandler(DatabaseAccessWrite))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},       // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},      // Allow GET, POST, and OPTIONS methods
		AllowedHeaders:   []string{"Authorization", "Set-Cookie"}, // Allow Authorization header
		AllowCredentials: true,
		Debug:            true,
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
