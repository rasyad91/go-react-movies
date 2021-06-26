package util

import (
	"encoding/json"
	"errors"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, wrap string, v interface{}) error {
	wrapper := make(map[string]interface{})
	wrapper[wrap] = v

	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(http.StatusOK)
	if err := enc.Encode(wrapper); err != nil {
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
	wrapper := make(map[string]interface{})
	wrapper["error"] = jsonError{
		Message: err.Error(),
		Status:  http.StatusBadRequest,
	}

	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(http.StatusBadRequest)
	if err := enc.Encode(wrapper); err != nil {
		w.WriteHeader(http.StatusBadGateway)
	}
}
