package usecase

import "time"

type UserGenerateReq struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Password string `json:"password"`
	Role     int    `json:"role"`
}

type UserResp struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Role     string `json:"role"`
	Password string `json:"password,omitempty"` // this only use for testing purpose to simplify when needed
}

type ValidateUserReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
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
