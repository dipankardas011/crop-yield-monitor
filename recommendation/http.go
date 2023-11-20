package main

import (
	"encoding/json"
	"io"
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

func makeHTTPCall(r *http.Request, httpMethod, url string, body io.Reader, res any) (int, error) {
	// it should pass the token extracted from parent functions which call this
	// Get the JWT token from the Authorization header
	//authHeader := r.Header.Get("Authorization")
	//if authHeader == "" {
	//	return http.StatusUnauthorized, apiError{Err: "Missing Authorization header", Status: http.StatusUnauthorized}
	//}

	userCookie, err := r.Cookie("user_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return http.StatusUnauthorized, apiError{Err: "Missing Cookie", Status: http.StatusUnauthorized}
		}
		return http.StatusInternalServerError, apiError{Err: err.Error(), Status: http.StatusInternalServerError}
	}

	request, error := http.NewRequest(httpMethod, url, body)
	if error != nil {
		return http.StatusInternalServerError, error
	}

	request.Header.Set("Authorization", "Bearer "+userCookie.Value)

	client := &http.Client{}

	response, error := client.Do(request)
	if error != nil {
		return http.StatusInternalServerError, error
	}

	responseBody, error := io.ReadAll(response.Body)
	if error != nil {
		return http.StatusInternalServerError, error
	}

	if err := json.Unmarshal(responseBody, &res); err != nil {
		return http.StatusInternalServerError, apiError{Status: http.StatusInternalServerError, Err: err.Error()}
	}

	if response.StatusCode >= 300 {
		err := ""
		switch o := res.(type) {
		case *Response:
			err = o.Error
		case *AuthResponse:
			err = o.Error
		case *ImgResponse:
			err = o.Error
		}
		return response.StatusCode, apiError{Status: response.StatusCode, Err: err}
	}

	return http.StatusOK, nil
}
