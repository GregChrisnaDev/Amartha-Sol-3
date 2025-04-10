package usecase

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/GregChrisnaDev/Amartha-Sol-3/common/pdfgenerator"
	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/model"
	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/repository"
	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/storage"
	"github.com/google/uuid"
)

type loanUsecase struct {
	userRepo      repository.UserRepository
	loanRepo      repository.LoanRepository
	lendRepo      repository.LendRepository
	storageClient storage.Client
	pdfGenerator  pdfgenerator.Client
}

type LoanUsecase interface {
	ProposeLoan(ctx context.Context, params ProposeLoanReq) error
	GetLoanByLoanUID(ctx context.Context, userId uint64) ([]GetLoanResp, error)
	ApproveLoan(ctx context.Context, params PromoteLoanToApprovedReq) error
	GetProofPicture(ctx context.Context, loanId uint64) (FileResp, error)
	GetListApprovedLoan(ctx context.Context) ([]GetLoanResp, error)
	DisbursedLoan(ctx context.Context, params PromoteLoanToDisburseReq) error
	GetAgreementLetter(ctx context.Context, params GetAgreementLetterReq) (FileResp, error)
	GetListLender(ctx context.Context, params GetListLender) ([]GetLendResp, error)
}

func InitLoanUC(userRepo repository.UserRepository, loanRepo repository.LoanRepository, lendRepo repository.LendRepository, storageClient storage.Client, pdfGenerator pdfgenerator.Client) LoanUsecase {
	return &loanUsecase{
		userRepo:      userRepo,
		loanRepo:      loanRepo,
		lendRepo:      lendRepo,
		storageClient: storageClient,
		pdfGenerator:  pdfGenerator,
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

func (u *loanUsecase) GetProofPicture(ctx context.Context, loanId uint64) (FileResp, error) {
	loan, err := u.loanRepo.GetLoanByID(ctx, loanId)
	if err != nil {
		return FileResp{}, err
	}

	fileDetail, err := u.storageClient.DownloadFile(storage.PICTURE_PROOF_DIR + loan.PictureProofFilePath)
	if err != nil {
		return FileResp{}, err
	}

	return FileResp{
		Filename: loan.PictureProofFilePath,
		File:     fileDetail.File,
	}, nil
}

func (u *loanUsecase) GetListApprovedLoan(ctx context.Context) ([]GetLoanResp, error) {
	loans, err := u.loanRepo.GetListByStatus(ctx, model.Approved)
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

func (u *loanUsecase) DisbursedLoan(ctx context.Context, params PromoteLoanToDisburseReq) error {
	if exist, err := u.loanRepo.LoanExist(ctx, params.LoanID, model.Invested); err != nil {
		return err
	} else if !exist {
		return errors.New("loan id not exist")
	}

	imageFileName := uuid.New().String() + ".jpg"
	if err := u.storageClient.UploadImage(params.UserSign.Bytes(), storage.USER_SIGN_DIR, imageFileName); err != nil {
		return err
	}

	loaner, err := u.userRepo.GetByLoanId(ctx, params.LoanID)
	if err != nil {
		return err
	}

	lends, err := u.lendRepo.GetByLoanId(ctx, params.LoanID)
	if err != nil {
		return err
	}

	for _, lend := range lends {
		lender, err := u.userRepo.GetByID(ctx, lend.UserID)
		if err != nil {
			return err
		}

		err = u.pdfGenerator.GenerateAgreementLetter(pdfgenerator.AgreementLetterPDF{
			NameLender:    lender.Name,
			NameLoaner:    loaner.Name,
			AddressLender: lender.Address,
			AddressLoaner: loaner.Address,
			SignLender:    storage.USER_SIGN_DIR + lend.UserSignPath,
			SignLoaner:    storage.USER_SIGN_DIR + imageFileName,
			Filename:      lend.AgreementFilePath,
		})
		if err != nil {
			return err
		}
	}

	return u.loanRepo.PromoteLoanToDisbursed(ctx, model.Loan{
		ID:           params.LoanID,
		DisburserUID: params.DisburserID,
		UserSignPath: imageFileName,
	})
}

func (u *loanUsecase) GetAgreementLetter(ctx context.Context, params GetAgreementLetterReq) (FileResp, error) {
	filePath, err := u.loanRepo.GetAgreementFilePath(ctx, params.LendID, params.LoanID, params.User.ID)
	if err != nil {
		return FileResp{}, err
	}

	fileDetail, err := u.storageClient.DownloadFile(storage.AGREEMENT_LETTER_DIR + filePath)
	if err != nil {
		return FileResp{}, err
	}

	return FileResp{
		Filename: filePath,
		File:     fileDetail.File,
	}, nil
}

func (u *loanUsecase) GetListLender(ctx context.Context, params GetListLender) ([]GetLendResp, error) {
	lends, err := u.lendRepo.GetListLenderByLoanerID(ctx, params.LoanID, params.UserID)
	if err != nil {
		return []GetLendResp{}, err
	}

	var resp []GetLendResp
	for _, lend := range lends {
		resp = append(resp, GetLendResp{
			ID:                lend.ID,
			LoanID:            lend.LoanID,
			UserID:            lend.UserID,
			Amount:            convertToCurrency(lend.Amount),
			AgreementFilePath: lend.AgreementFilePath,
			CreatedAt:         lend.CreatedAt,
			UpdatedAt:         lend.UpdatedAt,
		})
	}

	return resp, nil
}
