package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/model"
	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/usecase"
)

type loanHandler struct {
	userUC usecase.UserUsecase
	loanUC usecase.LoanUsecase
}

type LoanHandler interface {
	ProposeLoanHandler(w http.ResponseWriter, r *http.Request)
	GetLoanByUIDHandler(w http.ResponseWriter, r *http.Request)
	ApproveLoanHandler(w http.ResponseWriter, r *http.Request)
	GetProofPictureHandler(w http.ResponseWriter, r *http.Request)
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
	user := validateUserAuth(r, h.userUC, model.Customer)
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	var req usecase.ProposeLoanReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	req.UserID = user.ID

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
	user := validateUserAuth(r, h.userUC, model.Customer)
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	resp, err := h.loanUC.GetLoanByLoanUID(ctx, user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, "Failed", nil)
		return
	}

	writeJSON(w, http.StatusOK, "Success", resp)
}

func (h *loanHandler) ApproveLoanHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// validate auth
	user := validateUserAuth(r, h.userUC, model.Employee)
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Limit request body size (10MB)
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		writeJSON(w, http.StatusUnauthorized, "Failed to parse multipart form", nil)
		return
	}

	loanId, err := strconv.ParseUint(r.FormValue("loan_id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, "Invalid Parameter", nil)
		return
	}

	imageBuf, err := convertImageToBuffer(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, "Failed parse image", nil)
		return
	}

	req := usecase.PromoteLoanToApprovedReq{
		LoanID:       loanId,
		ApproverID:   user.ID,
		PictureProof: imageBuf,
	}

	err = h.loanUC.ApproveLoan(ctx, req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, "Failed", nil)
		return
	}

	writeJSON(w, http.StatusOK, "Success", nil)
}

func (h *loanHandler) GetProofPictureHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// validate auth
	user := validateUserAuth(r, h.userUC, 0)
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	loanId, err := strconv.ParseUint(r.URL.Query().Get("loan_id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, "Invalid Parameter", nil)
		return
	}

	resp, err := h.loanUC.GetProofPicture(ctx, loanId)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, "Failed", nil)
		return
	}
	defer resp.Image.Close()

	stat, _ := resp.Image.Stat()

	w.Header().Set("Content-Type", "image/jpg")
	w.Header().Set("Content-Disposition", "inline; filename=\""+resp.Filename+"\"")
	http.ServeContent(w, r, resp.Filename, stat.ModTime(), resp.Image)
}
