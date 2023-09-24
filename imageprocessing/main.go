package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type ErrorMsg string

func (e ErrorMsg) String() string {
	return string(e)
}

const (
	BadJsonFormat              ErrorMsg = "[Err] wrong json format"
	UnSupportedMediaFormatType ErrorMsg = "[Err] invalid image type supported are jpeg and png"
	InternalServerError        ErrorMsg = "[Err] internal server error"
)

func imageUpload(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	switch r.Method {
	case http.MethodPost:
		log.Printf("Method [%s]: /image/upload\tTRIGGERED\n", r.Method)
	default:
		log.Printf("Method [%s]: /image/upload\t%v\n", r.Method, http.StatusBadRequest)
		http.Error(w, "Bad Method use POST", http.StatusBadRequest)
		return
	}
	if r.Header.Get("Content-type") != "application/json" {
		log.Println("Content-type not json")
		http.Error(w, "Bad Content-type require json", http.StatusBadRequest)
		return
	}

	payload := ImageUpload{}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Println("failed to decode the json body")
		http.Error(w, BadJsonFormat.String()+"\nReason: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Println(payload)
	fileName := ""
	switch payload.Format {
	case "image/png":
		fileName = "image.png"
	case "image/jpeg":
		fileName = "image.jpeg"
	default:
		log.Println("image format type missmatch")
		http.Error(w, UnSupportedMediaFormatType.String(), http.StatusUnsupportedMediaType)
		return
	}

	if err := os.WriteFile(fileName, payload.RawImage, 0666); err != nil {
		log.Println("unable to write the image file")
		http.Error(w, InternalServerError.String()+"\nReason: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(Response{
		Stdout: "fake response uploaded",
	}); err != nil {
		log.Println("unable to encode the response")
		http.Error(w, InternalServerError.String()+"\nReason: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Method [%s]: /image/upload\t%d\n", r.Method, http.StatusOK)
}

func imageGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	switch r.Method {
	case http.MethodGet:
		log.Printf("Method [%s]: /image/get\tTRIGGERED\n", r.Method)
	default:
		log.Printf("Method [%s]: /image/get\t%v\n", r.Method, http.StatusBadRequest)
		http.Error(w, "Bad Method use GET", http.StatusBadRequest)
		return
	}
	if r.Header.Get("Content-type") != "application/json" {
		log.Println("Content-type not json")
		http.Error(w, "Bad Content-type require json", http.StatusBadRequest)
		return
	}

	payload := ImageGet{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Println("failed to decode the json body")
		http.Error(w, BadJsonFormat.String()+"\nReason: "+err.Error(), http.StatusBadRequest)
		return
	}
	img := payload.Uuid // demo for image
	fmt.Println(img)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(Response{
		Stdout: "fake",
		Image: ImageGetResp{
			RawImage: []byte(img),
		},
	}); err != nil {
		log.Println("unable to encode the response")
		http.Error(w, InternalServerError.String()+"\nReason: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Method [%s]: /image/get\t%d\n", r.Method, http.StatusOK)
}

func main() {

	s := &http.Server{
		Addr:           ":8090",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	http.HandleFunc("/image/upload", imageUpload)
	http.HandleFunc("/image/get", imageGet)

	http.HandleFunc("/image", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		switch r.Method {
		case http.MethodGet:
			log.Printf("Method [%s]: /image\tTRIGGERED\n", r.Method)
		default:
			log.Printf("Method [%s]: /image\t%v\n", r.Method, http.StatusBadRequest)
			http.Error(w, "Bad Method use GET", http.StatusBadRequest)
			return
		}
		docs := struct {
			Loc map[string]string
		}{
			Loc: map[string]string{
				"upload": "/image/upload",
				"get":    "/image/get",
				"TODO":   "about payloads",
			},
		}

		if err := json.NewEncoder(w).Encode(docs); err != nil {
			http.Error(w, InternalServerError.String()+"\nReason: "+err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Method [%s]: /image\t%d\n", r.Method, http.StatusOK)
	})

	http.HandleFunc("/image/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		switch r.Method {
		case http.MethodGet:
			log.Printf("Method [%s]: /image/healthz\tTRIGGERED\n", r.Method)
		default:
			log.Printf("Method [%s]: /image/healthz\t%v\n", r.Method, http.StatusBadRequest)
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

		log.Printf("Method [%s]: /image/healthz\t%d\n", r.Method, http.StatusOK)
	})

	log.Printf("Started to serve the image server on port {%v}\n", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}
