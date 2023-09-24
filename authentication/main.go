package main

import (
	"encoding/json"
	"fmt"
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

var jwtKey []byte

// SignUp HTTP("POST")
func SignUp(w http.ResponseWriter, r *http.Request) error {

	if r.Method != http.MethodPost {
		return apiError{Err: "Bad Method, expected POST", Status: http.StatusMethodNotAllowed}
	}

	account := AccountSignUp{}
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		return apiError{Status: http.StatusBadRequest, Err: BadJsonFormat.String() + "\nReason: " + err.Error()}
	}

	defer log.Printf("Method [%s]: /account/signup\t%d\n", r.Method, http.StatusOK)
	fmt.Println(account)

	return writeJson(w, http.StatusOK, Response{
		Stdout: "signup successful",
		Account: AccountSignInRes{
			Uuid: "abcd23e23",
		},
	})
	// w.WriteHeader(http.StatusAccepted)
	// if err := json.NewEncoder(w).Encode(Response{
	// 	Stdout: "signup successful",
	// 	Account: AccountSignInRes{
	// 		Uuid: "abcd23e23",
	// 	},
	// }); err != nil {
	// 	log.Println("unable to encode the response")
	// 	http.Error(w, InternalServerError.String()+"\nReason: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }
}

// SignIn HTTP("POST")
func SignIn(w http.ResponseWriter, r *http.Request) error {

	if r.Method != http.MethodPost {
		return apiError{Err: "Bad Method, expected POST", Status: http.StatusMethodNotAllowed}
	}

	account := AccountSignIn{}
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		return apiError{Status: http.StatusBadRequest, Err: BadJsonFormat.String() + "\nReason: " + err.Error()}
	}

	if account.Password != "1234" || account.UserName != "dipankar" {
		return apiError{Err: "wrong password or username", Status: http.StatusUnauthorized}
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
		return apiError{Err: err.Error(), Status: http.StatusInternalServerError}
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
		Stdout: "logged in do refer to cache for more getting the tokens",
		Account: AccountSignInRes{
			Uuid: "abcd23e23",
		},
	})
}

// Welcome HTTP("GET")
func Welcome(w http.ResponseWriter, r *http.Request) error {

	if r.Method != http.MethodGet {
		return apiError{Err: "Bad Method, expected GET", Status: http.StatusMethodNotAllowed}
	}
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return apiError{Status: http.StatusUnauthorized, Err: err.Error()}
		}
		// For any other type of error, return a bad request status
		return apiError{Status: http.StatusBadRequest, Err: err.Error()}
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
			return apiError{Status: http.StatusUnauthorized, Err: err.Error()}
		}
		return apiError{Status: http.StatusBadRequest, Err: err.Error()}
	}
	if !tkn.Valid {
		return apiError{Status: http.StatusUnauthorized, Err: err.Error()}
	}
	return writeJson(w, http.StatusOK, struct{ Msg string }{Msg: "Welcome " + claims.Username})
}

// Refresh HTTP("POST")
func Refresh(w http.ResponseWriter, r *http.Request) error {

	if r.Method != http.MethodPost {
		return apiError{Err: "Bad Method, expected GET", Status: http.StatusMethodNotAllowed}
	}
	// (BEGIN) The code until this point is the same as the first part of the `Welcome` route
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return apiError{Status: http.StatusUnauthorized, Err: err.Error()}
		}
		return apiError{Status: http.StatusBadRequest, Err: err.Error()}
	}
	tknStr := c.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return apiError{Status: http.StatusUnauthorized, Err: err.Error()}
		}
		return apiError{Status: http.StatusBadRequest, Err: err.Error()}
	}
	if !tkn.Valid {
		return apiError{Status: http.StatusUnauthorized, Err: err.Error()}
	}

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		return apiError{Status: http.StatusBadRequest, Err: "a new token will only be issued if the old token is within 30 seconds of expiry"}
	}

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)

	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return apiError{Status: http.StatusInternalServerError, Err: err.Error()}
	}

	// Set the new token as the users `token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	return writeJson(w, http.StatusOK, struct{ Msg string }{Msg: "Refreshed token for user=" + claims.Username})
}

func Logout(w http.ResponseWriter, r *http.Request) error {

	if r.Method != http.MethodPost {
		return apiError{Err: "Bad Method, expected GET", Status: http.StatusMethodNotAllowed}
	}

	// immediately clear the token cookie
	_, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return apiError{Status: http.StatusUnauthorized, Err: err.Error()}
		}
		return apiError{Status: http.StatusBadRequest, Err: err.Error()}
	}
	claims := &Claims{}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Expires: time.Now(),
	})

	return writeJson(w, http.StatusOK, struct{ Msg string }{Msg: "logout success of " + claims.Username})
}

func Docs(w http.ResponseWriter, r *http.Request) error {

	if r.Method != http.MethodGet {
		return apiError{Err: "invalid method", Status: http.StatusMethodNotAllowed}
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

	return writeJson(w, http.StatusOK, docs)
}

func Health(w http.ResponseWriter, r *http.Request) error {

	if r.Method != http.MethodGet {
		return apiError{Err: "invalid method", Status: http.StatusMethodNotAllowed}
	}

	type Health struct {
		Msg string `json:"message"`
	}

	return writeJson(w, http.StatusOK, Health{Msg: "auth looks healthy"})
}

func main() {

	jwtKey = []byte(generateRandomToken(20))

	// mux := http.NewServeMux()

	http.HandleFunc("/account/signin", makeHTTPHandler(SignIn))
	http.HandleFunc("/account/signup", makeHTTPHandler(SignUp))
	http.HandleFunc("/account/logout", makeHTTPHandler(Logout))
	http.HandleFunc("/account/welcome", makeHTTPHandler(Welcome))
	http.HandleFunc("/account/renew", makeHTTPHandler(Refresh))
	fmt.Println(`
POST /account/signin
POST /account/signup
POST /account/logout
GET /account/welcome
POST /account/renew`)

	http.HandleFunc("/account", makeHTTPHandler(Docs))

	http.HandleFunc("/account/healthz", makeHTTPHandler(Health))

	// wrappedMux := NewLogger(mux)

	s := &http.Server{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Started to serve the authorization server on port {%v}\n", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}

func writeJson(w http.ResponseWriter, statusCode int, data any) error {
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(data)
}

type apiError struct {
	Status int
	Err    string
}

func (e apiError) Error() string {
	return e.Err
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Printf("[%s] %s ⚡", r.Method, r.URL.Path)
		start := time.Now()

		defer log.Printf("[%s] %s %v", r.Method, r.URL.Path, time.Since(start))

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Accept", "application/json; charset=utf-8")
		w.Header().Set("server", "authentication-server")

		if err := f(w, r); err != nil {
			log.Println(err)

			if e, ok := err.(apiError); ok {
				writeJson(w, e.Status, Response{Error: e.Error()})
			} else {
				writeJson(w, http.StatusInternalServerError, Response{Error: err.Error()})
			}
		}
	}
}
