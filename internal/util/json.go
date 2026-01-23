package util

import (
	"encoding/json"
	"kasir-api/internal/models"
	"net/http"
)

func EncodeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func DecodeJSON(r *http.Request, payload any) error {
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(payload); err != nil {
		return models.NewError(http.StatusInternalServerError, "Failed to decode JSON payload")
	}

	return nil
}
