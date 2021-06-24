package util

import (
	"encoding/json"
	"errors"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, v interface{}) error {

	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(http.StatusOK)
	if err := enc.Encode(v); err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return errors.New("fail to send json")
	}

	return nil
}

func ErrorJSON(w http.ResponseWriter, err error) {
	type jsonError struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
	}

	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(http.StatusOK)
	if err := enc.Encode(jsonError{Message: err.Error(), Status: http.StatusBadRequest}); err != nil {
		w.WriteHeader(http.StatusBadGateway)
	}
}
