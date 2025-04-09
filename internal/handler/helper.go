package handler

import (
	"encoding/json"
	"net/http"

	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/usecase"
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

func validateUserAuth(r *http.Request, userUC usecase.UserUsecase) uint64 {
	username, password, ok := r.BasicAuth()
	if !ok {
		return 0
	}

	return userUC.ValidateUser(r.Context(), usecase.ValidateUserReq{Name: username, Password: password})
}
