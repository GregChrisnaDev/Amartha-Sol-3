package model

import "time"

type Lend struct {
	ID                uint64
	LoanID            uint64
	UserID            uint64
	Amount            float64
	AgreementFilePath string
	CreatedDate       time.Time
	UpdatedDate       time.Time
}
