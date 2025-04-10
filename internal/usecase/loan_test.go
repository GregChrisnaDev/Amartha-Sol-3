package usecase_test

import (
	"bytes"
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/model"
	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/storage"
	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/usecase"
	mock_cache "github.com/GregChrisnaDev/Amartha-Sol-3/tests/mock/common/cache"
	mock_pdfgenerator "github.com/GregChrisnaDev/Amartha-Sol-3/tests/mock/common/pdfgenerator"
	mock_repository "github.com/GregChrisnaDev/Amartha-Sol-3/tests/mock/repository"
	mock_storage "github.com/GregChrisnaDev/Amartha-Sol-3/tests/mock/storage"
)

type LoanTestSuite struct {
	suite.Suite

	ctrl          *gomock.Controller
	userRepo      *mock_repository.MockUserRepository
	loanRepo      *mock_repository.MockLoanRepository
	lendRepo      *mock_repository.MockLendRepository
	storageClient *mock_storage.MockClient
	pdfGenerator  *mock_pdfgenerator.MockClient
	cacheClient   *mock_cache.MockRedisLock

	loanUC usecase.LoanUsecase
}

func TestLoanTestSuite(t *testing.T) {
	suite.Run(t, new(LoanTestSuite))
}

func (s *LoanTestSuite) BeforeTest(string, string) {
	s.ctrl = gomock.NewController(s.T())
	s.userRepo = mock_repository.NewMockUserRepository(s.ctrl)
	s.lendRepo = mock_repository.NewMockLendRepository(s.ctrl)
	s.loanRepo = mock_repository.NewMockLoanRepository(s.ctrl)
	s.pdfGenerator = mock_pdfgenerator.NewMockClient(s.ctrl)
	s.storageClient = mock_storage.NewMockClient(s.ctrl)
	s.cacheClient = mock_cache.NewMockRedisLock(s.ctrl)

	s.loanUC = usecase.InitLoanUC(s.userRepo, s.loanRepo, s.lendRepo, s.storageClient, s.pdfGenerator, s.cacheClient)
}

func (s *LoanTestSuite) AfterTest(string, string) {
	defer s.ctrl.Finish()
}

func (s *LoanTestSuite) TestLoanUsecase_ProposeLoan() {
	testCases := []struct {
		name     string
		params   usecase.ProposeLoanReq
		mockFunc func()
		err      error
	}{
		{
			name: "success to propose loan",
			params: usecase.ProposeLoanReq{
				UserID:          1,
				PrincipalAmount: 10000000,
				Rate:            10,
				LoanDuration:    50,
			},
			mockFunc: func() {
				s.cacheClient.EXPECT().Acquire(gomock.Any(), "propose_loan:1", 1000*time.Millisecond, 200*time.Millisecond).Return(true)
				s.cacheClient.EXPECT().Release(gomock.Any(), "propose_loan:1").Return(nil)
				s.loanRepo.EXPECT().PendingLoanExist(gomock.Any(), uint64(1)).Return(false, nil)
				s.loanRepo.EXPECT().ProposeLoan(gomock.Any(), model.Loan{
					UserID:          1,
					PrincipalAmount: 10000000,
					Rate:            10,
					LoanDuration:    50,
				}).Return(nil)
			},
			err: nil,
		},
		{
			name: "error lock transaction",
			params: usecase.ProposeLoanReq{
				UserID:          1,
				PrincipalAmount: 10000000,
				Rate:            10,
				LoanDuration:    50,
			},
			mockFunc: func() {
				s.cacheClient.EXPECT().Acquire(gomock.Any(), "propose_loan:1", 1000*time.Millisecond, 200*time.Millisecond).Return(false)
			},
			err: errors.New("failed to lock"),
		},
		{
			name: "error user has pending loan",
			params: usecase.ProposeLoanReq{
				UserID:          1,
				PrincipalAmount: 10000000,
				Rate:            10,
				LoanDuration:    50,
			},
			mockFunc: func() {
				s.cacheClient.EXPECT().Acquire(gomock.Any(), "propose_loan:1", 1000*time.Millisecond, 200*time.Millisecond).Return(true)
				s.cacheClient.EXPECT().Release(gomock.Any(), "propose_loan:1").Return(nil)
				s.loanRepo.EXPECT().PendingLoanExist(gomock.Any(), uint64(1)).Return(true, nil)
			},
			err: errors.New("user has pending loan"),
		},
		{
			name: "error while get pending loan",
			params: usecase.ProposeLoanReq{
				UserID:          1,
				PrincipalAmount: 10000000,
				Rate:            10,
				LoanDuration:    50,
			},
			mockFunc: func() {
				s.cacheClient.EXPECT().Acquire(gomock.Any(), "propose_loan:1", 1000*time.Millisecond, 200*time.Millisecond).Return(true)
				s.cacheClient.EXPECT().Release(gomock.Any(), "propose_loan:1").Return(nil)
				s.loanRepo.EXPECT().PendingLoanExist(gomock.Any(), uint64(1)).Return(false, errors.New("error get"))
			},
			err: errors.New("error get"),
		},
		{
			name: "error while store loan to db",
			params: usecase.ProposeLoanReq{
				UserID:          1,
				PrincipalAmount: 10000000,
				Rate:            10,
				LoanDuration:    50,
			},
			mockFunc: func() {
				s.cacheClient.EXPECT().Acquire(gomock.Any(), "propose_loan:1", 1000*time.Millisecond, 200*time.Millisecond).Return(true)
				s.cacheClient.EXPECT().Release(gomock.Any(), "propose_loan:1").Return(nil)
				s.loanRepo.EXPECT().PendingLoanExist(gomock.Any(), uint64(1)).Return(false, nil)
				s.loanRepo.EXPECT().ProposeLoan(gomock.Any(), model.Loan{
					UserID:          1,
					PrincipalAmount: 10000000,
					Rate:            10,
					LoanDuration:    50,
				}).Return(errors.New("error insert"))
			},
			err: errors.New("error insert"),
		},
	}

	for _, test := range testCases {
		s.Run(test.name, func() {
			test.mockFunc()
			err := s.loanUC.ProposeLoan(context.Background(), test.params)
			s.Equal(err, test.err)
		})
	}
}

func (s *LoanTestSuite) TestLoanUsecase_ApproveLoan() {
	testCases := []struct {
		name     string
		params   usecase.PromoteLoanToApprovedReq
		mockFunc func()
		err      error
	}{
		{
			name: "success to approve loan",
			params: usecase.PromoteLoanToApprovedReq{
				LoanID:       1,
				ApproverID:   2,
				PictureProof: bytes.NewBuffer([]byte{}),
			},
			mockFunc: func() {
				s.cacheClient.EXPECT().Acquire(gomock.Any(), "approve_loan:1", 1000*time.Millisecond, 200*time.Millisecond).Return(true)
				s.cacheClient.EXPECT().Release(gomock.Any(), "approve_loan:1").Return(nil)
				s.loanRepo.EXPECT().LoanExist(gomock.Any(), uint64(1), model.Proposed).Return(true, nil)
				s.storageClient.EXPECT().UploadImage(bytes.NewBuffer([]byte{}).Bytes(), storage.PICTURE_PROOF_DIR, gomock.Any()).Return(nil)
				s.loanRepo.EXPECT().PromoteLoanToApproved(gomock.Any(), gomock.Any()).Return(nil)
			},
			err: nil,
		},
		{
			name: "error lock transaction",
			params: usecase.PromoteLoanToApprovedReq{
				LoanID:       1,
				ApproverID:   2,
				PictureProof: bytes.NewBuffer([]byte{}),
			},
			mockFunc: func() {
				s.cacheClient.EXPECT().Acquire(gomock.Any(), "approve_loan:1", 1000*time.Millisecond, 200*time.Millisecond).Return(false)
			},
			err: errors.New("failed to lock"),
		},
		{
			name: "error loan not exist",
			params: usecase.PromoteLoanToApprovedReq{
				LoanID:       1,
				ApproverID:   2,
				PictureProof: bytes.NewBuffer([]byte{}),
			},
			mockFunc: func() {
				s.cacheClient.EXPECT().Acquire(gomock.Any(), "approve_loan:1", 1000*time.Millisecond, 200*time.Millisecond).Return(true)
				s.cacheClient.EXPECT().Release(gomock.Any(), "approve_loan:1").Return(nil)
				s.loanRepo.EXPECT().LoanExist(gomock.Any(), uint64(1), model.Proposed).Return(false, nil)
			},
			err: errors.New("loan id not exist"),
		},
		{
			name: "error while check loan",
			params: usecase.PromoteLoanToApprovedReq{
				LoanID:       1,
				ApproverID:   2,
				PictureProof: bytes.NewBuffer([]byte{}),
			},
			mockFunc: func() {
				s.cacheClient.EXPECT().Acquire(gomock.Any(), "approve_loan:1", 1000*time.Millisecond, 200*time.Millisecond).Return(true)
				s.cacheClient.EXPECT().Release(gomock.Any(), "approve_loan:1").Return(nil)
				s.loanRepo.EXPECT().LoanExist(gomock.Any(), uint64(1), model.Proposed).Return(false, errors.New("error get"))
			},
			err: errors.New("error get"),
		},
		{
			name: "error while upload image",
			params: usecase.PromoteLoanToApprovedReq{
				LoanID:       1,
				ApproverID:   2,
				PictureProof: bytes.NewBuffer([]byte{}),
			},
			mockFunc: func() {
				s.cacheClient.EXPECT().Acquire(gomock.Any(), "approve_loan:1", 1000*time.Millisecond, 200*time.Millisecond).Return(true)
				s.cacheClient.EXPECT().Release(gomock.Any(), "approve_loan:1").Return(nil)
				s.loanRepo.EXPECT().LoanExist(gomock.Any(), uint64(1), model.Proposed).Return(true, nil)
				s.storageClient.EXPECT().UploadImage(bytes.NewBuffer([]byte{}).Bytes(), storage.PICTURE_PROOF_DIR, gomock.Any()).Return(errors.New("error upload"))
			},
			err: errors.New("error upload"),
		},
		{
			name: "error while promote loan to db",
			params: usecase.PromoteLoanToApprovedReq{
				LoanID:       1,
				ApproverID:   2,
				PictureProof: bytes.NewBuffer([]byte{}),
			},
			mockFunc: func() {
				s.cacheClient.EXPECT().Acquire(gomock.Any(), "approve_loan:1", 1000*time.Millisecond, 200*time.Millisecond).Return(true)
				s.cacheClient.EXPECT().Release(gomock.Any(), "approve_loan:1").Return(nil)
				s.loanRepo.EXPECT().LoanExist(gomock.Any(), uint64(1), model.Proposed).Return(true, nil)
				s.storageClient.EXPECT().UploadImage(bytes.NewBuffer([]byte{}).Bytes(), storage.PICTURE_PROOF_DIR, gomock.Any()).Return(nil)
				s.loanRepo.EXPECT().PromoteLoanToApproved(gomock.Any(), gomock.Any()).Return(errors.New("error update"))
			},
			err: errors.New("error update"),
		},
	}

	for _, test := range testCases {
		s.Run(test.name, func() {
			test.mockFunc()
			err := s.loanUC.ApproveLoan(context.Background(), test.params)
			s.Equal(err, test.err)
		})
	}
}

func (s *LoanTestSuite) TestLoanUsecase_DisburseLoan() {
	testCases := []struct {
		name     string
		params   usecase.PromoteLoanToDisburseReq
		mockFunc func()
		err      error
	}{
		{
			name: "success to disburse loan",
			params: usecase.PromoteLoanToDisburseReq{
				LoanID:      1,
				DisburserID: 2,
				UserSign:    bytes.NewBuffer([]byte{}),
			},
			mockFunc: func() {
				s.cacheClient.EXPECT().Acquire(gomock.Any(), "disburse_loan:1", 2000*time.Millisecond, 500*time.Millisecond).Return(true)
				s.cacheClient.EXPECT().Release(gomock.Any(), "disburse_loan:1").Return(nil)
				s.loanRepo.EXPECT().GetByIDStatus(gomock.Any(), uint64(1), model.Invested).Return(model.Loan{
					ID:              1,
					PrincipalAmount: 10000000,
					Rate:            10,
					LoanDuration:    50,
					UserSignPath:    "sign3",
				}, nil)
				s.storageClient.EXPECT().UploadImage(bytes.NewBuffer([]byte{}).Bytes(), storage.USER_SIGN_DIR, gomock.Any()).Return(nil)
				s.userRepo.EXPECT().GetByLoanId(gomock.Any(), uint64(1)).Return(model.User{
					Name:    "test",
					Address: "addr",
				}, nil)
				s.lendRepo.EXPECT().GetByLoanId(gomock.Any(), uint64(1)).Return([]model.Lend{
					{
						UserID:            2,
						UserSignPath:      "sign",
						AgreementFilePath: "file",
						Amount:            5000000,
					},
					{
						UserID:            3,
						UserSignPath:      "sign2",
						AgreementFilePath: "file2",
						Amount:            5000000,
					},
				}, nil)
				s.userRepo.EXPECT().GetByID(gomock.Any(), uint64(2)).Return(model.User{
					Name:    "test2",
					Address: "addr2",
				}, nil)

				s.userRepo.EXPECT().GetByID(gomock.Any(), uint64(3)).Return(model.User{
					Name:    "test3",
					Address: "addr3",
				}, nil)
				s.pdfGenerator.EXPECT().GenerateAgreementLetter(gomock.Any()).Times(2).Return(nil)
				s.loanRepo.EXPECT().PromoteLoanToDisbursed(gomock.Any(), gomock.Any()).Return(nil)
			},
			err: nil,
		},
		{
			name: "error while lock transaction",
			params: usecase.PromoteLoanToDisburseReq{
				LoanID:      1,
				DisburserID: 2,
				UserSign:    bytes.NewBuffer([]byte{}),
			},
			mockFunc: func() {
				s.cacheClient.EXPECT().Acquire(gomock.Any(), "disburse_loan:1", 2000*time.Millisecond, 500*time.Millisecond).Return(false)
			},
			err: errors.New("failed to lock"),
		},
		{
			name: "error while get loan data",
			params: usecase.PromoteLoanToDisburseReq{
				LoanID:      1,
				DisburserID: 2,
				UserSign:    bytes.NewBuffer([]byte{}),
			},
			mockFunc: func() {
				s.cacheClient.EXPECT().Acquire(gomock.Any(), "disburse_loan:1", 2000*time.Millisecond, 500*time.Millisecond).Return(true)
				s.cacheClient.EXPECT().Release(gomock.Any(), "disburse_loan:1").Return(nil)
				s.loanRepo.EXPECT().GetByIDStatus(gomock.Any(), uint64(1), model.Invested).Return(model.Loan{}, errors.New("error get"))
			},
			err: errors.New("error get"),
		},
		{
			name: "error while get loaner data",
			params: usecase.PromoteLoanToDisburseReq{
				LoanID:      1,
				DisburserID: 2,
				UserSign:    bytes.NewBuffer([]byte{}),
			},
			mockFunc: func() {
				s.cacheClient.EXPECT().Acquire(gomock.Any(), "disburse_loan:1", 2000*time.Millisecond, 500*time.Millisecond).Return(true)
				s.cacheClient.EXPECT().Release(gomock.Any(), "disburse_loan:1").Return(nil)
				s.loanRepo.EXPECT().GetByIDStatus(gomock.Any(), uint64(1), model.Invested).Return(model.Loan{
					ID:              1,
					PrincipalAmount: 10000000,
					Rate:            10,
					LoanDuration:    50,
					UserSignPath:    "sign3",
				}, nil)
				s.storageClient.EXPECT().UploadImage(bytes.NewBuffer([]byte{}).Bytes(), storage.USER_SIGN_DIR, gomock.Any()).Return(errors.New("error upload"))
			},
			err: errors.New("error upload"),
		},
		{
			name: "error while get loaner data",
			params: usecase.PromoteLoanToDisburseReq{
				LoanID:      1,
				DisburserID: 2,
				UserSign:    bytes.NewBuffer([]byte{}),
			},
			mockFunc: func() {
				s.cacheClient.EXPECT().Acquire(gomock.Any(), "disburse_loan:1", 2000*time.Millisecond, 500*time.Millisecond).Return(true)
				s.cacheClient.EXPECT().Release(gomock.Any(), "disburse_loan:1").Return(nil)
				s.loanRepo.EXPECT().GetByIDStatus(gomock.Any(), uint64(1), model.Invested).Return(model.Loan{
					ID:              1,
					PrincipalAmount: 10000000,
					Rate:            10,
					LoanDuration:    50,
					UserSignPath:    "sign3",
				}, nil)
				s.storageClient.EXPECT().UploadImage(bytes.NewBuffer([]byte{}).Bytes(), storage.USER_SIGN_DIR, gomock.Any()).Return(nil)
				s.userRepo.EXPECT().GetByLoanId(gomock.Any(), uint64(1)).Return(model.User{}, errors.New("error get"))
			},
			err: errors.New("error get"),
		},
		{
			name: "error while get lend data",
			params: usecase.PromoteLoanToDisburseReq{
				LoanID:      1,
				DisburserID: 2,
				UserSign:    bytes.NewBuffer([]byte{}),
			},
			mockFunc: func() {
				s.cacheClient.EXPECT().Acquire(gomock.Any(), "disburse_loan:1", 2000*time.Millisecond, 500*time.Millisecond).Return(true)
				s.cacheClient.EXPECT().Release(gomock.Any(), "disburse_loan:1").Return(nil)
				s.loanRepo.EXPECT().GetByIDStatus(gomock.Any(), uint64(1), model.Invested).Return(model.Loan{
					ID:              1,
					PrincipalAmount: 10000000,
					Rate:            10,
					LoanDuration:    50,
					UserSignPath:    "sign3",
				}, nil)
				s.storageClient.EXPECT().UploadImage(bytes.NewBuffer([]byte{}).Bytes(), storage.USER_SIGN_DIR, gomock.Any()).Return(nil)
				s.userRepo.EXPECT().GetByLoanId(gomock.Any(), uint64(1)).Return(model.User{
					Name:    "test",
					Address: "addr",
				}, nil)
				s.lendRepo.EXPECT().GetByLoanId(gomock.Any(), uint64(1)).Return([]model.Lend{}, errors.New("error get"))
			},
			err: errors.New("error get"),
		},
		{
			name: "error while get user data lender",
			params: usecase.PromoteLoanToDisburseReq{
				LoanID:      1,
				DisburserID: 2,
				UserSign:    bytes.NewBuffer([]byte{}),
			},
			mockFunc: func() {
				s.cacheClient.EXPECT().Acquire(gomock.Any(), "disburse_loan:1", 2000*time.Millisecond, 500*time.Millisecond).Return(true)
				s.cacheClient.EXPECT().Release(gomock.Any(), "disburse_loan:1").Return(nil)
				s.loanRepo.EXPECT().GetByIDStatus(gomock.Any(), uint64(1), model.Invested).Return(model.Loan{
					ID:              1,
					PrincipalAmount: 10000000,
					Rate:            10,
					LoanDuration:    50,
					UserSignPath:    "sign3",
				}, nil)
				s.storageClient.EXPECT().UploadImage(bytes.NewBuffer([]byte{}).Bytes(), storage.USER_SIGN_DIR, gomock.Any()).Return(nil)
				s.userRepo.EXPECT().GetByLoanId(gomock.Any(), uint64(1)).Return(model.User{
					Name:    "test",
					Address: "addr",
				}, nil)
				s.lendRepo.EXPECT().GetByLoanId(gomock.Any(), uint64(1)).Return([]model.Lend{
					{
						UserID:            2,
						UserSignPath:      "sign",
						AgreementFilePath: "file",
						Amount:            5000000,
					},
				}, nil)
				s.userRepo.EXPECT().GetByID(gomock.Any(), uint64(2)).Return(model.User{}, errors.New("error get"))
			},
			err: errors.New("error get"),
		},
		{
			name: "error while update data loan to db",
			params: usecase.PromoteLoanToDisburseReq{
				LoanID:      1,
				DisburserID: 2,
				UserSign:    bytes.NewBuffer([]byte{}),
			},
			mockFunc: func() {
				s.cacheClient.EXPECT().Acquire(gomock.Any(), "disburse_loan:1", 2000*time.Millisecond, 500*time.Millisecond).Return(true)
				s.cacheClient.EXPECT().Release(gomock.Any(), "disburse_loan:1").Return(nil)
				s.loanRepo.EXPECT().GetByIDStatus(gomock.Any(), uint64(1), model.Invested).Return(model.Loan{
					ID:              1,
					PrincipalAmount: 10000000,
					Rate:            10,
					LoanDuration:    50,
					UserSignPath:    "sign3",
				}, nil)
				s.storageClient.EXPECT().UploadImage(bytes.NewBuffer([]byte{}).Bytes(), storage.USER_SIGN_DIR, gomock.Any()).Return(nil)
				s.userRepo.EXPECT().GetByLoanId(gomock.Any(), uint64(1)).Return(model.User{
					Name:    "test",
					Address: "addr",
				}, nil)
				s.lendRepo.EXPECT().GetByLoanId(gomock.Any(), uint64(1)).Return([]model.Lend{
					{
						UserID:            2,
						UserSignPath:      "sign",
						AgreementFilePath: "file",
						Amount:            5000000,
					},
				}, nil)
				s.userRepo.EXPECT().GetByID(gomock.Any(), uint64(2)).Return(model.User{
					Name:    "test2",
					Address: "addr2",
				}, nil)
				s.pdfGenerator.EXPECT().GenerateAgreementLetter(gomock.Any()).Return(errors.New("error generate"))
			},
			err: errors.New("error generate"),
		},
		{
			name: "error while update data loan to db",
			params: usecase.PromoteLoanToDisburseReq{
				LoanID:      1,
				DisburserID: 2,
				UserSign:    bytes.NewBuffer([]byte{}),
			},
			mockFunc: func() {
				s.cacheClient.EXPECT().Acquire(gomock.Any(), "disburse_loan:1", 2000*time.Millisecond, 500*time.Millisecond).Return(true)
				s.cacheClient.EXPECT().Release(gomock.Any(), "disburse_loan:1").Return(nil)
				s.loanRepo.EXPECT().GetByIDStatus(gomock.Any(), uint64(1), model.Invested).Return(model.Loan{
					ID:              1,
					PrincipalAmount: 10000000,
					Rate:            10,
					LoanDuration:    50,
					UserSignPath:    "sign3",
				}, nil)
				s.storageClient.EXPECT().UploadImage(bytes.NewBuffer([]byte{}).Bytes(), storage.USER_SIGN_DIR, gomock.Any()).Return(nil)
				s.userRepo.EXPECT().GetByLoanId(gomock.Any(), uint64(1)).Return(model.User{
					Name:    "test",
					Address: "addr",
				}, nil)
				s.lendRepo.EXPECT().GetByLoanId(gomock.Any(), uint64(1)).Return([]model.Lend{
					{
						UserID:            2,
						UserSignPath:      "sign",
						AgreementFilePath: "file",
						Amount:            5000000,
					},
					{
						UserID:            3,
						UserSignPath:      "sign2",
						AgreementFilePath: "file2",
						Amount:            5000000,
					},
				}, nil)
				s.userRepo.EXPECT().GetByID(gomock.Any(), uint64(2)).Return(model.User{
					Name:    "test2",
					Address: "addr2",
				}, nil)

				s.userRepo.EXPECT().GetByID(gomock.Any(), uint64(3)).Return(model.User{
					Name:    "test3",
					Address: "addr3",
				}, nil)
				s.pdfGenerator.EXPECT().GenerateAgreementLetter(gomock.Any()).Times(2).Return(nil)
				s.loanRepo.EXPECT().PromoteLoanToDisbursed(gomock.Any(), gomock.Any()).Return(errors.New("error update"))
			},
			err: errors.New("error update"),
		},
	}

	for _, test := range testCases {
		s.Run(test.name, func() {
			test.mockFunc()
			err := s.loanUC.DisbursedLoan(context.Background(), test.params)
			s.Equal(err, test.err)
		})
	}
}

func (s *LoanTestSuite) TestLoanUsecase_GetAgreementLetter() {
	dummyFile := &os.File{}
	testCases := []struct {
		name     string
		params   usecase.GetAgreementLetterReq
		mockFunc func()
		wantResp usecase.FileResp
		err      error
	}{
		{
			name: "success to get agreement letter",
			params: usecase.GetAgreementLetterReq{
				LoanID: 1,
				User: &model.User{
					ID: 1,
				},
				LendID: 1,
			},
			mockFunc: func() {
				s.loanRepo.EXPECT().GetAgreementFilePath(gomock.Any(), uint64(1), uint64(1), uint64(1)).Return("filepath", nil)
				s.storageClient.EXPECT().DownloadFile(storage.AGREEMENT_LETTER_DIR+"filepath").Return(storage.DownloadFileResp{
					File: dummyFile,
				}, nil)
			},
			wantResp: usecase.FileResp{
				File:     dummyFile,
				Filename: "filepath",
			},
			err: nil,
		},
		{
			name: "error while get agreement path",
			params: usecase.GetAgreementLetterReq{
				LoanID: 1,
				User: &model.User{
					ID: 1,
				},
				LendID: 1,
			},
			mockFunc: func() {
				s.loanRepo.EXPECT().GetAgreementFilePath(gomock.Any(), uint64(1), uint64(1), uint64(1)).Return("", errors.New("error get"))
			},
			err: errors.New("error get"),
		},
		{
			name: "error while download file",
			params: usecase.GetAgreementLetterReq{
				LoanID: 1,
				User: &model.User{
					ID: 1,
				},
				LendID: 1,
			},
			mockFunc: func() {
				s.loanRepo.EXPECT().GetAgreementFilePath(gomock.Any(), uint64(1), uint64(1), uint64(1)).Return("filepath", nil)
				s.storageClient.EXPECT().DownloadFile(storage.AGREEMENT_LETTER_DIR+"filepath").Return(storage.DownloadFileResp{}, errors.New("error download"))
			},
			err: errors.New("error download"),
		},
	}

	for _, test := range testCases {
		s.Run(test.name, func() {
			test.mockFunc()
			resp, err := s.loanUC.GetAgreementLetter(context.Background(), test.params)
			s.Equal(err, test.err)
			s.Equal(resp, test.wantResp)
		})
	}
}
