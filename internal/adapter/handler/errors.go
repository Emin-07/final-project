package handler

import (
	"encoding/json"
	"errors"
	"net/http"
)

var noIdProvidedErr = errors.New("Не указан идентификатор")

func WriteJson(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func JsonError(w http.ResponseWriter, error string, status int) {
	WriteJson(w, map[string]string{"error": error}, status)
}
