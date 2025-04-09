package model

import "time"

type Loan struct {
	ID                   uint64
	UserID               uint64
	PrincipalAmount      float64
	Rate                 uint32
	LoanDuration         uint32
	Status               int8
	ProposedDate         time.Time
	PictureProofFilePath string
	ApproverUID          uint64
	ApprovalDate         time.Time
	DisburserUID         uint64
	DisbursedDate        time.Time
}

const (
	// Loan Status
	Proposed int8 = iota + 1
	Approved
	Invested
	Disbursed
)

var LoanStatusMapping = map[int8]string{
	Proposed:  "Proposed",
	Approved:  "Approved",
	Invested:  "Invested",
	Disbursed: "Disbursed",
}
