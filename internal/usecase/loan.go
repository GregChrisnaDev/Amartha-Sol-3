package usecase

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/model"
	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/repository"
)

type loanUsecase struct {
	loanRepo repository.LoanRepository
}

type LoanUsecase interface {
	ProposeLoan(ctx context.Context, params ProposeLoanReq) error
	GetLoanByLoanID(ctx context.Context, loanId uint64) (GetLoanResp, error)
	GetLoanByLoanUID(ctx context.Context, userId uint64) ([]GetLoanResp, error)
}

func InitLoanUC(loanRepo repository.LoanRepository) LoanUsecase {
	return &loanUsecase{
		loanRepo: loanRepo,
	}
}

func (u *loanUsecase) ProposeLoan(ctx context.Context, params ProposeLoanReq) error {
	// redis locking

	if hasPendingLoan, err := u.loanRepo.PendingLoanExist(ctx, params.UserID); err != nil {
		return err
	} else if hasPendingLoan {
		return errors.New("user has pending loan")
	}

	err := u.loanRepo.ProposeLoan(ctx, model.Loan{
		UserID:          params.UserID,
		PrincipalAmount: params.PrincipalAmount,
		Rate:            params.Rate,
		LoanDuration:    params.LoanDuration,
	})

	// redis unlock

	return err
}

func (u *loanUsecase) GetLoanByLoanID(ctx context.Context, loanId uint64) (GetLoanResp, error) {
	loan, err := u.loanRepo.GetLoanByID(ctx, loanId)
	if err != nil {
		return GetLoanResp{}, err
	}

	return GetLoanResp{
		ID:                   loan.ID,
		UserID:               loan.UserID,
		PrincipalAmount:      convertToCurrency(loan.PrincipalAmount),
		Rate:                 loan.Rate,
		LoanDuration:         fmt.Sprintf("%s week", strconv.FormatUint(uint64(loan.LoanDuration), 10)),
		Status:               model.LoanStatusMapping[loan.Status],
		ProposedDate:         loan.ProposedDate,
		PictureProofFilePath: loan.PictureProofFilePath,
		ApproverUID:          loan.ApproverUID,
		ApprovalDate:         loan.ApprovalDate,
		DisburserUID:         loan.DisburserUID,
		DisbursedDate:        loan.DisbursedDate,
	}, nil
}

func (u *loanUsecase) GetLoanByLoanUID(ctx context.Context, userId uint64) ([]GetLoanResp, error) {
	loans, err := u.loanRepo.GetLoanByUID(ctx, userId)
	if err != nil {
		return []GetLoanResp{}, err
	}

	var resp []GetLoanResp
	for _, loan := range loans {
		resp = append(resp, GetLoanResp{
			ID:                   loan.ID,
			UserID:               loan.UserID,
			PrincipalAmount:      convertToCurrency(loan.PrincipalAmount),
			Rate:                 loan.Rate,
			LoanDuration:         fmt.Sprintf("%s week", strconv.FormatUint(uint64(loan.LoanDuration), 10)),
			Status:               model.LoanStatusMapping[loan.Status],
			ProposedDate:         loan.ProposedDate,
			PictureProofFilePath: loan.PictureProofFilePath,
			ApproverUID:          loan.ApproverUID,
			ApprovalDate:         loan.ApprovalDate,
			DisburserUID:         loan.DisburserUID,
			DisbursedDate:        loan.DisbursedDate,
		})
	}

	return resp, nil
}
