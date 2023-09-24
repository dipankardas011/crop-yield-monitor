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
func SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	switch r.Method {
	case http.MethodPost:
		log.Printf("Method [%s]: /account/signup\tTRIGGERED\n", r.Method)
	default:
		log.Printf("Method [%s]: /account/signup\t%v\n", r.Method, http.StatusBadRequest)
		http.Error(w, "Bad Method use POST", http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-type") != "application/json" {
		log.Println("Content-type not json")
		http.Error(w, "Bad Content-type require json", http.StatusBadRequest)
		return
	}

	account := AccountSignUp{}
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		log.Println("failed to decode the json body")
		http.Error(w, BadJsonFormat.String()+"\nReason: "+err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(account)

	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(Response{
		Stdout: "signup successful",
		Account: AccountSignInRes{
			Uuid: "abcd23e23",
		},
	}); err != nil {
		log.Println("unable to encode the response")
		http.Error(w, InternalServerError.String()+"\nReason: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Method [%s]: /account/signup\t%d\n", r.Method, http.StatusOK)
}

// SignIn HTTP("POST")
func SignIn(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	switch r.Method {
	case http.MethodPost:
		log.Printf("Method [%s]: /account/signin\tTRIGGERED\n", r.Method)
	default:
		log.Printf("Method [%s]: /account/signin\t%v\n", r.Method, http.StatusBadRequest)
		http.Error(w, "Bad Method use POST", http.StatusBadRequest)
		return
	}
	if r.Header.Get("Content-type") != "application/json" {
		log.Println("Content-type not json")
		http.Error(w, "Bad Content-type require json", http.StatusBadRequest)
		return
	}

	account := AccountSignIn{}
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		log.Println("failed to decode the json body")
		http.Error(w, BadJsonFormat.String()+"\nReason: "+err.Error(), http.StatusBadRequest)
		return
	}

	////////////////// JWT

	if account.Password != "1234" || account.UserName != "dipankar" {
		log.Println("wrong password or username")
		http.Error(w, "Wrong password or username", http.StatusUnauthorized)
		return
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
		log.Println("signing token failed", err)
		http.Error(w, InternalServerError.String()+"\nReason: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	////////////////// JWT

	log.Println(account)

	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(Response{
		Stdout: "logged in do refer to cache for more getting the tokens",
		Account: AccountSignInRes{
			Uuid: "abcd23e23",
		},
	}); err != nil {
		log.Println("unable to encode the response")
		http.Error(w, InternalServerError.String()+"\nReason: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Method [%s]: /account/signin\t%d\n", r.Method, http.StatusOK)
}

// Welcome HTTP("GET")
func Welcome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	switch r.Method {
	case http.MethodGet:
		log.Printf("Method [%s]: /account/welcome\tTRIGGERED\n", r.Method)
	default:
		log.Printf("Method [%s]: /account/welcome\t%v\n", r.Method, http.StatusBadRequest)
		http.Error(w, "Bad Method use GET", http.StatusBadRequest)
		return
	}
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
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
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Write([]byte(fmt.Sprintf("{\"msg\":\"Welcome %s\"}", claims.Username)))
}

// Refresh HTTP("POST")
func Refresh(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	switch r.Method {
	case http.MethodPost:
		log.Printf("Method [%s]: /account/renew\tTRIGGERED\n", r.Method)
	default:
		log.Printf("Method [%s]: /account/renew\t%v\n", r.Method, http.StatusBadRequest)
		http.Error(w, "Bad Method use POST", http.StatusBadRequest)
		return
	}
	// (BEGIN) The code until this point is the same as the first part of the `Welcome` route
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tknStr := c.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)

	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the new token as the users `token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	w.Write([]byte(fmt.Sprintf("{\"msg\":\"Refreshed token for user=%s\"}", claims.Username)))
}

func Logout(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	switch r.Method {
	case http.MethodPost:
		log.Printf("Method [%s]: /account/renew\tTRIGGERED\n", r.Method)
	default:
		log.Printf("Method [%s]: /account/renew\t%v\n", r.Method, http.StatusBadRequest)
		http.Error(w, "Bad Method use POST", http.StatusBadRequest)
		return
	}

	// immediately clear the token cookie
	_, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	claims := &Claims{}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Expires: time.Now(),
	})

	w.Write([]byte(fmt.Sprintf("{\"msg\":\"logout success\"}", claims.Username)))
}

func main() {

	jwtKey = []byte(generateRandomToken(20))

	s := &http.Server{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	http.HandleFunc("/account/signin", SignIn)
	http.HandleFunc("/account/signup", SignUp)
	http.HandleFunc("/account/logout", Logout)
	http.HandleFunc("/account/welcome", Welcome)
	http.HandleFunc("/account/renew", Refresh)
	fmt.Println(`
POST /account/signin
POST /account/signup
POST /account/logout
GET /account/welcome
POST /account/renew`)

	http.HandleFunc("/account", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		switch r.Method {
		case http.MethodGet:
			log.Printf("Method [%s]: /account\tTRIGGERED\n", r.Method)
		default:
			log.Printf("Method [%s]: /account\t%v\n", r.Method, http.StatusBadRequest)
			http.Error(w, "Bad Method use GET", http.StatusBadRequest)
			return
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

		if err := json.NewEncoder(w).Encode(docs); err != nil {
			http.Error(w, InternalServerError.String()+"\nReason: "+err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Method [%s]: /account\t%d\n", r.Method, http.StatusOK)
	})

	http.HandleFunc("/account/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		switch r.Method {
		case http.MethodGet:
			log.Printf("Method [%s]: /account\tTRIGGERED\n", r.Method)
		default:
			log.Printf("Method [%s]: /account\t%v\n", r.Method, http.StatusBadRequest)
			http.Error(w, "Bad Method use GET", http.StatusBadRequest)
			return
		}

		type Health struct {
			Msg string `json:"message"`
		}
		if err := json.NewEncoder(w).Encode(Health{Msg: "auth looks healthy"}); err != nil {
			http.Error(w, InternalServerError.String()+"\nReason: "+err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Method [%s]: /account/healthz\t%d\n", r.Method, http.StatusOK)
	})

	log.Printf("Started to serve the authorization server on port {%v}\n", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}
