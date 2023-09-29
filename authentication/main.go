package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
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

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	log.Println(account)

	return writeJson(w, http.StatusOK, Response{
		Stdout:  "Login Successful",
		Account: "Check the cookie",
	})
}

// Refresh HTTP("POST")
func Refresh(w http.ResponseWriter, r *http.Request) (int, error) {

	if r.Method != http.MethodPost {
		return http.StatusMethodNotAllowed, apiError{Err: "Bad Method, expected GET", Status: http.StatusMethodNotAllowed}
	}
	// (BEGIN) The code until this point is the same as the first part of the `Welcome` route
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return http.StatusUnauthorized, err
		}
		return http.StatusBadRequest, err
	}
	tknStr := c.Value
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
		return http.StatusNotAcceptable, apiError{Status: http.StatusNotAcceptable, Err: "a new token will only be issued if the old token is within 30 seconds of expiry"}
	}

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)

	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Set the new token as the users `token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	return writeJson(w, http.StatusOK, Response{Stdout: "Refreshed token for user=" + claims.Username})
}

func Logout(w http.ResponseWriter, r *http.Request) (int, error) {

	if r.Method != http.MethodPost {
		return http.StatusMethodNotAllowed, apiError{Err: "Bad Method, expected GET", Status: http.StatusMethodNotAllowed}
	}

	// immediately clear the token cookie
	_, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return http.StatusUnauthorized, err
		}
		return http.StatusBadRequest, err
	}
	claims := &Claims{}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Expires: time.Now(),
	})

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
			"signin": "/account/signin",
			"signup": "/account/signup",
			"TODO":   "about payloads",
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

func IsValidToken(w http.ResponseWriter, r *http.Request) (int, error) {

	if r.Method != http.MethodGet {
		return http.StatusMethodNotAllowed, apiError{Err: "Bad Method, expected GET", Status: http.StatusMethodNotAllowed}
	}
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return http.StatusUnauthorized, err
		}
		// For any other type of error, return a bad request status
		return http.StatusBadRequest, err
	}

	// Get the JWT string from the cookie
	tknStr := c.Value

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

	return writeJson(w, http.StatusOK, Response{Stdout: "Welcome " + claims.Username})
}

func main() {

	jwtKey = []byte(generateRandomToken(20))

	http.HandleFunc("/account/signin", makeHTTPHandler(SignIn))
	http.HandleFunc("/account/signup", makeHTTPHandler(SignUp))
	http.HandleFunc("/account/logout", makeHTTPHandler(Logout))
	http.HandleFunc("/account/renew", makeHTTPHandler(Refresh))
	http.HandleFunc("/account/token/status", makeHTTPHandler(IsValidToken))

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
