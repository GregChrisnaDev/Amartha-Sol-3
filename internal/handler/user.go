package handler

import (
	"encoding/json"
	"net/http"

	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/usecase"
)

type userHandler struct {
	userUC usecase.UserUsecase
}

type UserHandler interface {
	UserGenerateHandler(w http.ResponseWriter, r *http.Request)
	GetAllUserHandler(w http.ResponseWriter, r *http.Request)
}

func InitUserHandler(userUC usecase.UserUsecase) UserHandler {
	return &userHandler{
		userUC: userUC,
	}
}

func (h *userHandler) UserGenerateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req usecase.UserGenerateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, "Invalid request body", err.Error(), nil)
		return
	}

	if req.Name == "" || req.Address == "" || req.Email == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, "Missing parameter", "invalid parameter", nil)
		return
	}

	if req.Role != 1 && req.Role != 2 {
		writeJSON(w, http.StatusBadRequest, "Invalid role", "invalid role", nil)
		return
	}

	resp, err := h.userUC.GenerateUser(ctx, req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, "Failed", err.Error(), nil)
		return
	}

	writeJSON(w, http.StatusOK, "Success", "", resp)
}

func (h *userHandler) GetAllUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp, err := h.userUC.GetAllUser(ctx)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, "Failed", err.Error(), resp)
		return
	}

	writeJSON(w, http.StatusOK, "Success", "", resp)
}
