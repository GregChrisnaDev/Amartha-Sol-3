package usecase

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/model"
	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/repository"
	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/storage"
	"github.com/google/uuid"
)

type loanUsecase struct {
	loanRepo      repository.LoanRepository
	storageClient storage.StorageClient
}

type LoanUsecase interface {
	ProposeLoan(ctx context.Context, params ProposeLoanReq) error
	GetLoanByLoanUID(ctx context.Context, userId uint64) ([]GetLoanResp, error)
	ApproveLoan(ctx context.Context, params PromoteLoanToApprovedReq) error
	GetProofPicture(ctx context.Context, loanId uint64) (GetProofPictureResp, error)
}

func InitLoanUC(loanRepo repository.LoanRepository, storageClient storage.StorageClient) LoanUsecase {
	return &loanUsecase{
		loanRepo:      loanRepo,
		storageClient: storageClient,
	}
}

func (u *loanUsecase) ProposeLoan(ctx context.Context, params ProposeLoanReq) error {
	// TODO: add locking
	// TODO: defer unlock

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

	return err
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

func (u *loanUsecase) ApproveLoan(ctx context.Context, params PromoteLoanToApprovedReq) error {
	// TODO: add locking
	// TODO: defer unlock

	if exist, err := u.loanRepo.LoanExist(ctx, params.LoanID, model.Proposed); err != nil {
		return err
	} else if !exist {
		return errors.New("loan id not exist")
	}

	imageFileName := uuid.New().String() + ".jpg"
	if err := u.storageClient.UploadImage(params.PictureProof.Bytes(), storage.PICTURE_PROOF_DIR, imageFileName); err != nil {
		return err
	}

	return u.loanRepo.PromoteLoanToApproved(ctx, model.Loan{
		ID:                   params.LoanID,
		ApproverUID:          params.ApproverID,
		PictureProofFilePath: imageFileName,
	})
}

func (u *loanUsecase) GetProofPicture(ctx context.Context, loanId uint64) (GetProofPictureResp, error) {
	loan, err := u.loanRepo.GetLoanByID(ctx, loanId)
	if err != nil {
		return GetProofPictureResp{}, err
	}

	fileDetail, err := u.storageClient.DownloadFile(storage.PICTURE_PROOF_DIR + loan.PictureProofFilePath)
	if err != nil {
		return GetProofPictureResp{}, err
	}

	return GetProofPictureResp{
		Filename: loan.PictureProofFilePath,
		Image:    fileDetail.File,
	}, nil
}
