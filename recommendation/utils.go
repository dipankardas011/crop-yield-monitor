package main

import (
	"encoding/json"
	"net/http"
)

func writeJson(w http.ResponseWriter, statusCode int, data any) (int, error) {
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
