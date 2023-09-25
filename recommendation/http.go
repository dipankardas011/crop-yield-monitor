package main

import (
	"log"
	"net/http"
	"time"
)

type apiFunc func(http.ResponseWriter, *http.Request) (int, error)

func makeHTTPHandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Printf("[%s] %s ⚡", r.Method, r.URL.Path)
		start := time.Now()

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Accept", "application/json; charset=utf-8")
		w.Header().Set("server", "authentication-server")

		statCode, err := f(w, r)
		if err != nil {
			log.Println(err)

			if e, ok := err.(apiError); ok {
				_, _ = writeJson(w, e.Status, Response{Error: e.Error()})
			} else {
				_, _ = writeJson(w, statCode, Response{Error: err.Error()})
			}
		}

		log.Printf("[%s] %s {%d} %v", r.Method, r.URL.Path, statCode, time.Since(start))
	}
}
