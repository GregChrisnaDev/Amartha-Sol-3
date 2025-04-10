package usecase

import (
	"context"
	"errors"
	"log"

	"github.com/GregChrisnaDev/Amartha-Sol-3/common/pdfgenerator"
	"github.com/GregChrisnaDev/Amartha-Sol-3/common/postgres"
	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/model"
	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/repository"
	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/storage"
	"github.com/google/uuid"
)

type lendUsecase struct {
	userRepo      repository.UserRepository
	loanRepo      repository.LoanRepository
	lendRepo      repository.LendRepository
	dbTransaction postgres.DBTransaction
	storageClient storage.Client
	pdfGenerator  pdfgenerator.Client
}

type LendUsecase interface {
	Simulate(ctx context.Context, params LendSimulateReq) (LendSimulateResp, error)
	Invest(ctx context.Context, params InvestReq) error
	GetAgreementLetter(ctx context.Context, params GetAgreementLetterReq) (FileResp, error)
	GetListLend(ctx context.Context, params uint64) ([]GetLendResp, error)
}

func InitLendUC(userRepo repository.UserRepository, loanRepo repository.LoanRepository, lendRepo repository.LendRepository, dbTransaction postgres.DBTransaction, storageClient storage.Client, pdfGenerator pdfgenerator.Client) LendUsecase {
	return &lendUsecase{
		userRepo:      userRepo,
		loanRepo:      loanRepo,
		lendRepo:      lendRepo,
		dbTransaction: dbTransaction,
		storageClient: storageClient,
		pdfGenerator:  pdfGenerator,
	}
}

func (u *lendUsecase) Simulate(ctx context.Context, params LendSimulateReq) (LendSimulateResp, error) {
	loan, err := u.loanRepo.GetByIDStatus(ctx, params.LoanID, model.Approved)
	if err != nil {
		return LendSimulateResp{}, err
	}

	if params.UserID == loan.UserID {
		log.Println("[Simulate] forbidden user")
		return LendSimulateResp{}, errors.New("forbidden")
	}

	if params.Amount > loan.PrincipalAmount {
		log.Println("[Simulate] invalid amount")
		return LendSimulateResp{}, errors.New("invalid amount")
	}

	ratePerWeek := float64(loan.Rate) / 52
	totalInterest := loan.PrincipalAmount * ratePerWeek * float64(loan.LoanDuration) / 100

	profit := params.Amount / loan.PrincipalAmount * totalInterest
	roi := profit / params.Amount * 100

	return LendSimulateResp{
		ROI:    roi,
		Profit: convertToCurrency(profit),
	}, nil
}

func (u *lendUsecase) Invest(ctx context.Context, params InvestReq) error {
	// TODO: Add Locking
	// TODO: Defer Locking

	loan, err := u.loanRepo.GetByIDStatus(ctx, params.LoanID, model.Approved)
	if err != nil {
		return err
	}

	if params.Lender.ID == loan.UserID {
		log.Println("[Invest] forbidden user")
		return errors.New("forbidden")
	}

	lendsByLoanId, err := u.lendRepo.GetByLoanId(ctx, loan.ID)
	if err != nil {
		return err
	}

	var totalFund float64
	for _, lend := range lendsByLoanId {
		if lend.UserID != params.Lender.ID {
			totalFund += lend.Amount
		}
	}

	if params.Amount > loan.PrincipalAmount-totalFund {
		log.Println("[Invest] invalid amount")
		return errors.New("invalid amount")
	}

	if lendsByUID, err := u.lendRepo.GetByUidLoanId(ctx, params.LoanID, params.Lender.ID); err != nil {
		if err.Error() != "not found" {
			return err
		}
	} else {
		// Update data
		return u.dbTransaction.Execute(ctx, func(ctx context.Context) error {
			if err := u.lendRepo.Update(ctx, model.Lend{ID: lendsByUID.ID, Amount: params.Amount}); err != nil {
				return err
			}
			if params.Amount == loan.PrincipalAmount-totalFund {
				return u.loanRepo.PromoteLoanToInvested(ctx, loan.ID)
			}
			return nil
		})
	}

	imageFileName := uuid.New().String() + ".jpg"
	if err := u.storageClient.UploadImage(params.UserSign.Bytes(), storage.USER_SIGN_DIR, imageFileName); err != nil {
		return err
	}

	loaner, err := u.userRepo.GetByLoanId(ctx, params.LoanID)
	if err != nil {
		return err
	}

	fileName := uuid.New().String() + ".pdf"
	err = u.pdfGenerator.GenerateAgreementLetter(pdfgenerator.AgreementLetterPDF{
		NameLender:    params.Lender.Name,
		NameLoaner:    loaner.Name,
		AddressLender: params.Lender.Address,
		AddressLoaner: loaner.Address,
		SignLender:    storage.USER_SIGN_DIR + imageFileName,
		Filename:      fileName,
	})
	if err != nil {
		return err
	}

	return u.dbTransaction.Execute(ctx, func(ctx context.Context) error {
		if err := u.lendRepo.Add(ctx, model.Lend{LoanID: params.LoanID, UserID: params.Lender.ID, UserSignPath: imageFileName, AgreementFilePath: fileName, Amount: params.Amount}); err != nil {
			return err
		}

		if params.Amount == loan.PrincipalAmount-totalFund {
			return u.loanRepo.PromoteLoanToInvested(ctx, loan.ID)
		}
		return nil
	})

	//TODO: Send Email
}

func (u *lendUsecase) GetAgreementLetter(ctx context.Context, params GetAgreementLetterReq) (FileResp, error) {
	lend, err := u.lendRepo.GetByUidLoanId(ctx, params.LoanID, params.User.ID)
	if err != nil {
		return FileResp{}, err
	}

	fileDetail, err := u.storageClient.DownloadFile(storage.AGREEMENT_LETTER_DIR + lend.AgreementFilePath)
	if err != nil {
		return FileResp{}, err
	}

	return FileResp{
		Filename: lend.AgreementFilePath,
		File:     fileDetail.File,
	}, nil
}

func (u *lendUsecase) GetListLend(ctx context.Context, params uint64) ([]GetLendResp, error) {
	lends, err := u.lendRepo.GetByUID(ctx, params)
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
