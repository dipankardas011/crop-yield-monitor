package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type ErrorMsg string

func (e ErrorMsg) String() string {
	return string(e)
}

const (
	BadJsonFormat       ErrorMsg = "[Err] wrong json format"
	InternalServerError ErrorMsg = "[Err] internal server error"
)

var (
	jwtKey    []byte
	sqlClient *DBClient
)

// SignUp HTTP("POST")
func SignUp(w http.ResponseWriter, r *http.Request) (int, error) {

	if r.Method != http.MethodPost {
		return http.StatusMethodNotAllowed, apiError{Err: "Bad Method, expected POST", Status: http.StatusMethodNotAllowed}
	}

	account := AccountSignUp{}
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		return http.StatusBadRequest, apiError{Status: http.StatusBadRequest, Err: BadJsonFormat.String() + "\nReason: " + err.Error()}
	}

	log.Println(account)

	if err := sqlClient.CreateUser(account); err != nil {
		return http.StatusBadRequest, apiError{Status: http.StatusBadRequest, Err: err.Error()}
	}

	return writeJson(w, http.StatusOK, Response{
		Stdout: "signup successful",
	})
}

// SignIn HTTP("POST")
func SignIn(w http.ResponseWriter, r *http.Request) (int, error) {

	if r.Method != http.MethodPost {
		return http.StatusMethodNotAllowed, apiError{Err: "Bad Method, expected POST", Status: http.StatusMethodNotAllowed}
	}

	account := AccountSignIn{}
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		return http.StatusBadRequest, apiError{Status: http.StatusBadRequest, Err: BadJsonFormat.String() + "\nReason: " + err.Error()}
	}

	dbStoredData, err := sqlClient.GetPasswordByUsername(account.UserName)
	if err != nil {
		if err != sql.ErrNoRows {
			return http.StatusUnauthorized, apiError{Err: err.Error(), Status: http.StatusUnauthorized}
		} else {
			return http.StatusUnauthorized, apiError{Err: "Invalid Username", Status: http.StatusUnauthorized}
		}
	}

	if err := func() error {
		getHash := genHash(account.Password + dbStoredData.salt)
		if dbStoredData.password != getHash {
			return errors.New("invalid password")
		}
		return nil
	}(); err != nil {
		return http.StatusUnauthorized, apiError{Err: err.Error(), Status: http.StatusUnauthorized}
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: account.UserName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return http.StatusInternalServerError, apiError{Err: err.Error(), Status: http.StatusInternalServerError}
	}

	return writeJson(w, http.StatusOK, Response{
		Stdout: "Login Successful",
		Account: struct {
			Token string `json:"token"`
		}{Token: tokenString},
	})
}

// Refresh HTTP("POST")
func Refresh(w http.ResponseWriter, r *http.Request) (int, error) {

	if r.Method != http.MethodPost {
		return http.StatusMethodNotAllowed, apiError{Err: "Bad Method, expected GET", Status: http.StatusMethodNotAllowed}
	}

	// Get the JWT token from the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return http.StatusUnauthorized, apiError{Err: "Missing Authorization header", Status: http.StatusUnauthorized}
	}

	// Extract the token from the header (assuming Bearer token format)
	tknStr := strings.TrimPrefix(authHeader, "Bearer ")

	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return http.StatusUnauthorized, err
		}
		return http.StatusBadRequest, err
	}
	if !tkn.Valid {
		return http.StatusUnauthorized, err
	}

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		return http.StatusNotAcceptable, apiError{Status: http.StatusNotAcceptable, Err: "a new token will only be issued if the old token is within 30 seconds after expiry"}
	}

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)

	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return http.StatusInternalServerError, apiError{Err: "Failed to generate new token", Status: http.StatusInternalServerError}
	}

	return writeJson(w, http.StatusOK, Response{
		Stdout: "Refreshed token for user=" + claims.Username,
		Account: struct {
			Token string `json:"token"`
		}{Token: tokenString},
	})
}

func Logout(w http.ResponseWriter, r *http.Request) (int, error) {

	if r.Method != http.MethodPost {
		return http.StatusMethodNotAllowed, apiError{Err: "Bad Method, expected GET", Status: http.StatusMethodNotAllowed}
	}

	// Get the JWT token from the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return http.StatusUnauthorized, apiError{Err: "Missing Authorization header", Status: http.StatusUnauthorized}
	}

	// Extract the token from the header (assuming Bearer token format)

	_ = strings.TrimPrefix(authHeader, "Bearer ")

	// FIXME: to enhance security, you can maintain a blacklist of invalidated tokens on the server-side. When a user logs out, you can add the current token to the blacklist.
	// 	For each request, you can check if the token in the Authorization header is not in the blacklist before processing the request. This extra step helps prevent the use of invalidated tokens even if they are somehow retained by the client.

	claims := &Claims{}

	return writeJson(w, http.StatusOK, Response{Stdout: "logout success of " + claims.Username})
}

func Docs(w http.ResponseWriter, r *http.Request) (int, error) {

	if r.Method != http.MethodGet {
		return http.StatusMethodNotAllowed, apiError{Err: "invalid method", Status: http.StatusMethodNotAllowed}
	}
	docs := struct {
		Loc map[string]string
	}{
		Loc: map[string]string{
			"[POST] signin":                    "/account/signin",
			"[POST] signup":                    "/account/signup",
			"[POST] logout":                    "/account/logout",
			"[GET] token renew":                "/account/renew",
			"[GET] authorization bearer token": "/account/token",
			"[GET] Health status":              "/account/healthz",
		},
	}

	return writeJson(w, http.StatusOK, Response{Account: docs})
}

func Health(w http.ResponseWriter, r *http.Request) (int, error) {

	if r.Method != http.MethodGet {
		return http.StatusMethodNotAllowed, apiError{Err: "invalid method", Status: http.StatusMethodNotAllowed}
	}

	return writeJson(w, http.StatusOK, Response{Stdout: "auth looks healthy"})
}

func IsAuthenticToken(w http.ResponseWriter, r *http.Request) (int, error) {

	if r.Method != http.MethodGet {
		return http.StatusMethodNotAllowed, apiError{Err: "Bad Method, expected GET", Status: http.StatusMethodNotAllowed}
	}

	// Get the JWT token from the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return http.StatusUnauthorized, apiError{Err: "Missing Authorization header", Status: http.StatusUnauthorized}
	}

	// Extract the token from the header (assuming Bearer token format)
	tknStr := strings.TrimPrefix(authHeader, "Bearer ")

	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return http.StatusUnauthorized, err
		}
		return http.StatusBadRequest, err
	}
	if !tkn.Valid {
		return http.StatusUnauthorized, err
	}

	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		return http.StatusNotAcceptable, apiError{Status: http.StatusNotAcceptable, Err: "a new token will only be issued if the old token is within 30 seconds of expiry"}
	}

	return writeJson(w, http.StatusOK, Response{Stdout: claims.Username})
}

func main() {

	DB_CONN_ADDR = os.Getenv("DB_URL")
	dbPassword = os.Getenv("DB_PASSWORD")
	jwtKey = []byte(generateRandomToken(20))

	http.HandleFunc("/account/signin", makeHTTPHandler(SignIn))
	http.HandleFunc("/account/signup", makeHTTPHandler(SignUp))
	http.HandleFunc("/account/logout", makeHTTPHandler(Logout))
	http.HandleFunc("/account/renew", makeHTTPHandler(Refresh))
	http.HandleFunc("/account/token", makeHTTPHandler(IsAuthenticToken))

	http.HandleFunc("/account", makeHTTPHandler(Docs))
	http.HandleFunc("/account/healthz", makeHTTPHandler(Health))

	s := &http.Server{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// create the mysql server client
	sqlClient = &DBClient{}

	if err := sqlClient.MySqlNewClient(); err != nil {
		panic(err)
	}

	log.Printf("Started to serve the authorization server on port {%v}\n", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}
