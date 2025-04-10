package usecase

import (
	"bytes"
	"os"
	"time"

	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/model"
)

type UserGenerateReq struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     int    `json:"role"`
}

type UserResp struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"` // this only use for testing purpose to simplify when needed
}

type ValidateUserReq struct {
	Email    string
	Password string
}

type SimulateLoanReq struct {
	PrincipalAmount float64 `json:"principal_amount"`
	Rate            uint32  `json:"rate"`
	LoanDuration    uint32  `json:"loan_duration"`
}

type SimulateLoanResp struct {
	TotalRepays        string `json:"total_repays"`
	WeeklyInstallments string `json:"weekly_installments"`
}

type ProposeLoanReq struct {
	UserID          uint64
	PrincipalAmount float64 `json:"principal_amount"`
	Rate            uint32  `json:"rate"`
	LoanDuration    uint32  `json:"loan_duration"`
}

type GetLoanResp struct {
	ID                   uint64    `json:"id"`
	UserID               uint64    `json:"user_id"`
	PrincipalAmount      string    `json:"principal_amount"`
	Rate                 uint32    `json:"rate"`
	LoanDuration         string    `json:"loan_duration"`
	Status               string    `json:"status"`
	ProposedDate         time.Time `json:"proposed_date"`
	PictureProofFilePath string    `json:"picture_proof_filepath,omitempty"`
	ApproverUID          uint64    `json:"approver_uid,omitempty"`
	ApprovalDate         time.Time `json:"approval_date,omitempty"`
	DisburserUID         uint64    `json:"disburser_uid,omitempty"`
	DisbursedDate        time.Time `json:"disbursement_date,omitempty"`
}

type PromoteLoanToApprovedReq struct {
	LoanID       uint64
	ApproverID   uint64
	PictureProof *bytes.Buffer
}

type FileResp struct {
	File     *os.File
	Filename string
}

type LendSimulateReq struct {
	UserID uint64
	LoanID uint64  `json:"loan_id"`
	Amount float64 `json:"amount"`
}

type LendSimulateResp struct {
	ROI    float64 `json:"roi"`
	Profit string  `json:"profit"`
}

type InvestReq struct {
	Lender   *model.User
	LoanID   uint64
	Amount   float64
	UserSign *bytes.Buffer
}

type PromoteLoanToDisburseReq struct {
	LoanID      uint64
	DisburserID uint64
	UserSign    *bytes.Buffer
}

type GetAgreementLetterReq struct {
	User   *model.User
	LendID uint64
	LoanID uint64
}

type GetListLender struct {
	UserID uint64
	LoanID uint64
}

type GetLendResp struct {
	ID                uint64    `json:"id"`
	LoanID            uint64    `json:"loan_id"`
	UserID            uint64    `json:"user_id"`
	Amount            string    `json:"amount"`
	UserSignPath      string    `json:"user_sign_path,omitempty"`
	AgreementFilePath string    `json:"agreement_file_path"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
