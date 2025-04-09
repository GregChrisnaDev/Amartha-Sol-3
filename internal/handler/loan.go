package handler

import (
	"encoding/json"
	"net/http"

	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/usecase"
)

type loanHandler struct {
	userUC usecase.UserUsecase
	loanUC usecase.LoanUsecase
}

type LoanHandler interface {
	ProposeLoanHandler(w http.ResponseWriter, r *http.Request)
	GetLoanByUIDHandler(w http.ResponseWriter, r *http.Request)
}

func InitLoanHandler(userUC usecase.UserUsecase, loanUC usecase.LoanUsecase) LoanHandler {
	return &loanHandler{
		userUC: userUC,
		loanUC: loanUC,
	}
}

func (h *loanHandler) ProposeLoanHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// validate auth
	userId := validateUserAuth(r, h.userUC)
	if userId == 0 {
		writeJSON(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	var req usecase.ProposeLoanReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	req.UserID = userId

	if req.PrincipalAmount <= 0 || req.LoanDuration <= 0 || req.Rate <= 0 {
		writeJSON(w, http.StatusBadRequest, "Invalid Parameter", nil)
		return
	}

	err := h.loanUC.ProposeLoan(ctx, req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, "Failed", nil)
		return
	}

	writeJSON(w, http.StatusOK, "Success", nil)
}

func (h *loanHandler) GetLoanByUIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// validate auth
	userId := validateUserAuth(r, h.userUC)
	if userId == 0 {
		writeJSON(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	resp, err := h.loanUC.GetLoanByLoanUID(ctx, userId)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, "Failed", nil)
		return
	}

	writeJSON(w, http.StatusOK, "Success", resp)
}
