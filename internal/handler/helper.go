package handler

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := APIResponse{
		Code:    statusCode,
		Message: message,
		Data:    data,
	}
	json.NewEncoder(w).Encode(resp)
}
