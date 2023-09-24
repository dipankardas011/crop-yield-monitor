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

func getRecommendations(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	switch r.Method {
	case http.MethodGet:
		log.Printf("Method [%s]: /recommend/get\tTRIGGERED\n", r.Method)
	default:
		log.Printf("Method [%s]: /recommend/get\t%v\n", r.Method, http.StatusBadRequest)
		http.Error(w, "Bad Method use GET", http.StatusBadRequest)
		return
	}
	if r.Header.Get("Content-type") != "application/json" {
		log.Println("Content-type not json")
		http.Error(w, "Bad Content-type require json", http.StatusBadRequest)
		return
	}

	payload := RecommendGet{}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Println("failed to decode the json body")
		http.Error(w, BadJsonFormat.String()+"\nReason: "+err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(payload)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(Response{Stdout: "fake recommendation", Recommendations: Recommendations{
		Crops: []string{"fake01", "fake02"},
	}}); err != nil {
		log.Println("unable to encode the response")
		http.Error(w, InternalServerError.String()+"\nReason: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Method [%s]: /recommend/get\t%d\n", r.Method, http.StatusOK)
}

func main() {

	s := &http.Server{
		Addr:           ":8100",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	http.HandleFunc("/recommend/get", getRecommendations)

	http.HandleFunc("/recommend", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		switch r.Method {
		case http.MethodGet:
			log.Printf("Method [%s]: /recommend\tTRIGGERED\n", r.Method)
		default:
			log.Printf("Method [%s]: /recommend\t%v\n", r.Method, http.StatusBadRequest)
			http.Error(w, "Bad Method use GET", http.StatusBadRequest)
			return
		}
		docs := struct {
			Loc map[string]string
		}{
			Loc: map[string]string{
				"get":  "/recommend/get",
				"TODO": "about payloads",
			},
		}

		if err := json.NewEncoder(w).Encode(docs); err != nil {
			http.Error(w, InternalServerError.String()+"\nReason: "+err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Method [%s]: /recommend\t%d\n", r.Method, http.StatusOK)
	})

	http.HandleFunc("/recommend/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		switch r.Method {
		case http.MethodGet:
			log.Printf("Method [%s]: /recommend/healthz\tTRIGGERED\n", r.Method)
		default:
			log.Printf("Method [%s]: /recommend/healthz\t%v\n", r.Method, http.StatusBadRequest)
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

		log.Printf("Method [%s]: /recommend/healthz\t%d\n", r.Method, http.StatusOK)
	})

	log.Printf("Started to serve the recommend server on port {%v}\n", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}
