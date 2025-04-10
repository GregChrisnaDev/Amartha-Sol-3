package model

import "time"

type Lend struct {
	ID                uint64
	LoanID            uint64
	UserID            uint64
	Amount            float64
	UserSignPath      string
	AgreementFilePath string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
