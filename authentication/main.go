package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ErrorMsg string

func (e ErrorMsg) String() string {
	return string(e)
}

const (
	BadJsonFormat       ErrorMsg = "[Err] wrong json format"
	InternalServerError ErrorMsg = "[Err] internal server error"
)

func SignUp(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	switch req.Method {
	case http.MethodPost:
		log.Printf("Method [%s]: /account/signup\tTRIGGERED\n", req.Method)
	default:
		log.Printf("Method [%s]: /account/signup\t%v\n", req.Method, http.StatusBadRequest)
		http.Error(w, "Bad Method used use POST", http.StatusBadRequest)
		return
	}
	if req.Header.Get("Content-type") != "application/json" {
		http.Error(w, "Bad Content-type require json", http.StatusBadRequest)
		return
	}

	account := AccountSignUp{}
	if err := json.NewDecoder(req.Body).Decode(&account); err != nil {
		http.Error(w, BadJsonFormat.String()+"\nReason: "+err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(account)

	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(Response{
		Stdout: "signup successful",
		Account: AccountSignInRes{
			Uuid:        "abcd23e23",
			AccessToken: "32qwe32413212211(dummy)",
		},
	}); err != nil {
		http.Error(w, InternalServerError.String()+"\nReason: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Method [%s]: /account/signup\t%d\n", req.Method, http.StatusOK)
}

func SignIn(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	switch req.Method {
	case http.MethodPost:
		log.Printf("Method [%s]: /account/signin\tTRIGGERED\n", req.Method)
	default:
		log.Printf("Method [%s]: /account/signin\t%v\n", req.Method, http.StatusBadRequest)
		http.Error(w, "Bad Method used use POST", http.StatusBadRequest)
		return
	}
	if req.Header.Get("Content-type") != "application/json" {
		http.Error(w, "Bad Content-type require json", http.StatusBadRequest)
		return
	}

	account := AccountSignIn{}
	if err := json.NewDecoder(req.Body).Decode(&account); err != nil {
		http.Error(w, BadJsonFormat.String()+"\nReason: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Println(account)

	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(Response{
		Stdout: "logged in",
		Account: AccountSignInRes{
			Uuid:        "abcd23e23",
			AccessToken: "32qwe32413212211(dummy)",
		},
	}); err != nil {
		http.Error(w, InternalServerError.String()+"\nReason: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Method [%s]: /account/signin\t%d\n", req.Method, http.StatusOK)
}

func main() {

	// r.Use(cors.Default())
	// r.POST("/account/signin", SignIn)
	// r.POST("/account/signup", SignUp)
	// r.GET("/account/healthz", func(ctx *gin.Context) {
	// 	ctx.JSON(http.StatusOK, gin.H{
	// 		"message": "auth looks healthy",
	// 	})
	// })
	// r.GET("/account/docs", func(ctx *gin.Context) {
	// 	ctx.JSON(http.StatusOK, gin.H{
	// 		"signin": "/account/signin",
	// 		"signup": "/account/signup",
	// 		"TODO":   "about payloads",
	// 	})
	// })
	//
	s := &http.Server{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	http.HandleFunc("/account/signin", SignIn)
	http.HandleFunc("/account/signup", SignUp)

	http.HandleFunc("/account", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		switch r.Method {
		case http.MethodGet:
			log.Printf("Method [%s]: /account\tTRIGGERED\n", r.Method)
		default:
			log.Printf("Method [%s]: /account\t%v\n", r.Method, http.StatusBadRequest)
			http.Error(w, "Bad Method used use GET", http.StatusBadRequest)
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
			http.Error(w, "Bad Method used use GET", http.StatusBadRequest)
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
