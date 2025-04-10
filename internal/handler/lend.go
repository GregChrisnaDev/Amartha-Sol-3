package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/model"
	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/usecase"
)

type lendHandler struct {
	userUC usecase.UserUsecase
	lendUC usecase.LendUsecase
}

type LendHandler interface {
	SimulateHandler(w http.ResponseWriter, r *http.Request)
	GetListLenderHandler(w http.ResponseWriter, r *http.Request)
	InvestHandler(w http.ResponseWriter, r *http.Request)
	GetAgreementLetterHandler(w http.ResponseWriter, r *http.Request)
}

func InitLendHandler(userUC usecase.UserUsecase, lendUC usecase.LendUsecase) LendHandler {
	return &lendHandler{
		userUC: userUC,
		lendUC: lendUC,
	}
}

func (h *lendHandler) SimulateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// validate auth
	user := validateUserAuth(r, h.userUC, 0)
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	var req usecase.LendSimulateReq
	req.UserID = user.ID
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if req.LoanID <= 0 || req.Amount <= 0 {
		writeJSON(w, http.StatusBadRequest, "Invalid Parameter", nil)
		return
	}

	resp, err := h.lendUC.Simulate(ctx, req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, "Failed", nil)
		return
	}

	writeJSON(w, http.StatusOK, "Success", resp)
}

func (h *lendHandler) InvestHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error

	// validate auth
	user := validateUserAuth(r, h.userUC, model.Customer)
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	var req usecase.InvestReq
	req.Lender = user
	req.Amount, err = strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, "Invalid Parameter", nil)
		return
	}

	req.LoanID, err = strconv.ParseUint(r.FormValue("loan_id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, "Invalid Parameter", nil)
		return
	}

	req.UserSign, err = convertImageToBuffer(r, "user_sign")
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, "Failed parse image", nil)
		return
	}

	err = h.lendUC.Invest(ctx, req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, "Failed", nil)
		return
	}

	writeJSON(w, http.StatusOK, "Success", nil)
}

func (h *lendHandler) GetAgreementLetterHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// validate auth
	user := validateUserAuth(r, h.userUC, model.Customer)
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	loanId, err := strconv.ParseUint(r.URL.Query().Get("loan_id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, "Invalid Parameter", nil)
		return
	}

	resp, err := h.lendUC.GetAgreementLetter(ctx, usecase.GetAgreementLetterReq{
		User:   user,
		LoanID: loanId,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, "Failed", nil)
		return
	}
	defer resp.File.Close()

	stat, _ := resp.File.Stat()

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename=\""+resp.Filename+"\"")
	http.ServeContent(w, r, resp.Filename, stat.ModTime(), resp.File)
}

func (h *lendHandler) GetListLenderHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// validate auth
	user := validateUserAuth(r, h.userUC, model.Customer)
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	resp, err := h.lendUC.GetListLend(ctx, user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, "Failed", nil)
		return
	}

	writeJSON(w, http.StatusOK, "Success", resp)
}
