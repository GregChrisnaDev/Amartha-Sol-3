package usecase_test

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/GregChrisnaDev/Amartha-Sol-3/common"
	"github.com/GregChrisnaDev/Amartha-Sol-3/common/postgres"
	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/model"
	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/storage"
	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/usecase"
	mock_cache "github.com/GregChrisnaDev/Amartha-Sol-3/tests/mock/common/cache"
	mock_mail "github.com/GregChrisnaDev/Amartha-Sol-3/tests/mock/common/mail"
	mock_pdfgenerator "github.com/GregChrisnaDev/Amartha-Sol-3/tests/mock/common/pdfgenerator"
	mock_postgres "github.com/GregChrisnaDev/Amartha-Sol-3/tests/mock/common/postgres"
	mock_repository "github.com/GregChrisnaDev/Amartha-Sol-3/tests/mock/repository"
	mock_storage "github.com/GregChrisnaDev/Amartha-Sol-3/tests/mock/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type LendTestSuite struct {
	suite.Suite

	ctrl          *gomock.Controller
	userRepo      *mock_repository.MockUserRepository
	loanRepo      *mock_repository.MockLoanRepository
	lendRepo      *mock_repository.MockLendRepository
	dbTransaction *mock_postgres.MockDBTransaction
	storageClient *mock_storage.MockClient
	pdfGenerator  *mock_pdfgenerator.MockClient
	mailClient    *mock_mail.MockClient
	cacheClient   *mock_cache.MockRedisLock

	lendUC usecase.LendUsecase
}

func TestLendTestSuite(t *testing.T) {
	suite.Run(t, new(LendTestSuite))
}

func (s *LendTestSuite) BeforeTest(string, string) {
	s.ctrl = gomock.NewController(s.T())
	s.userRepo = mock_repository.NewMockUserRepository(s.ctrl)
	s.lendRepo = mock_repository.NewMockLendRepository(s.ctrl)
	s.loanRepo = mock_repository.NewMockLoanRepository(s.ctrl)
	s.dbTransaction = mock_postgres.NewMockDBTransaction(s.ctrl)
	s.storageClient = mock_storage.NewMockClient(s.ctrl)
	s.pdfGenerator = mock_pdfgenerator.NewMockClient(s.ctrl)
	s.mailClient = mock_mail.NewMockClient(s.ctrl)
	s.cacheClient = mock_cache.NewMockRedisLock(s.ctrl)
	common.AsyncMakeSync()

	s.lendUC = usecase.InitLendUC(s.userRepo, s.loanRepo, s.lendRepo, s.dbTransaction, s.storageClient, s.pdfGenerator, s.mailClient, s.cacheClient)
}

func (s *LendTestSuite) AfterTest(string, string) {
	defer s.ctrl.Finish()
}

func (s *LendTestSuite) TestLendUsecase_Invest() {
	testCases := []struct {
		name     string
		params   usecase.InvestReq
		mockFunc func()
		err      error
	}{
		{
			name: "success to new invest loan",
			params: usecase.InvestReq{
				Lender: &model.User{
					ID:      2,
					Name:    "test",
					Address: "addr",
					Email:   "email",
				},
				LoanID:   1,
				Amount:   9000000,
				UserSign: &bytes.Buffer{},
			},
			mockFunc: func() {
				s.cacheClient.EXPECT().Acquire(gomock.Any(), "invest_loan:1", 1000*time.Millisecond, 200*time.Millisecond).Return(true)
				s.cacheClient.EXPECT().Release(gomock.Any(), "invest_loan:1").Return(nil)
				s.loanRepo.EXPECT().GetByIDStatus(gomock.Any(), uint64(1), model.Approved).Return(model.Loan{
					ID:              1,
					UserID:          1,
					PrincipalAmount: 10000000,
				}, nil)
				s.lendRepo.EXPECT().GetByLoanId(gomock.Any(), uint64(1)).Return([]model.Lend{
					{
						UserID: 3,
						Amount: 1000000,
					},
				}, nil)
				s.userRepo.EXPECT().GetByLoanId(gomock.Any(), uint64(1)).Return(model.User{
					Name:    "test",
					Address: "addr",
				}, nil)
				s.lendRepo.EXPECT().GetByUidLoanId(gomock.Any(), uint64(1), uint64(2)).Return(model.Lend{}, errors.New("not found"))
				s.storageClient.EXPECT().UploadImage(gomock.Any(), storage.USER_SIGN_DIR, gomock.Any()).Return(nil)
				s.pdfGenerator.EXPECT().GenerateAgreementLetter(gomock.Any()).Return(nil)
				s.dbTransaction.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fun postgres.DBTransactionFunc) error {
					return fun(ctx)
				})
				s.lendRepo.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil)
				s.loanRepo.EXPECT().PromoteLoanToInvested(gomock.Any(), uint64(1)).Return(nil)
				s.mailClient.EXPECT().SendMail(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			err: nil,
		},
		{
			name: "success to new invest loan without promote to invested",
			params: usecase.InvestReq{
				Lender: &model.User{
					ID:      2,
					Name:    "test",
					Address: "addr",
					Email:   "email",
				},
				LoanID:   1,
				Amount:   8000000,
				UserSign: &bytes.Buffer{},
			},
			mockFunc: func() {
				s.cacheClient.EXPECT().Acquire(gomock.Any(), "invest_loan:1", 1000*time.Millisecond, 200*time.Millisecond).Return(true)
				s.cacheClient.EXPECT().Release(gomock.Any(), "invest_loan:1").Return(nil)
				s.loanRepo.EXPECT().GetByIDStatus(gomock.Any(), uint64(1), model.Approved).Return(model.Loan{
					ID:              1,
					UserID:          1,
					PrincipalAmount: 10000000,
				}, nil)
				s.lendRepo.EXPECT().GetByLoanId(gomock.Any(), uint64(1)).Return([]model.Lend{
					{
						UserID: 3,
						Amount: 1000000,
					},
				}, nil)
				s.userRepo.EXPECT().GetByLoanId(gomock.Any(), uint64(1)).Return(model.User{
					Name:    "test",
					Address: "addr",
				}, nil)
				s.lendRepo.EXPECT().GetByUidLoanId(gomock.Any(), uint64(1), uint64(2)).Return(model.Lend{}, errors.New("not found"))
				s.storageClient.EXPECT().UploadImage(gomock.Any(), storage.USER_SIGN_DIR, gomock.Any()).Return(nil)
				s.pdfGenerator.EXPECT().GenerateAgreementLetter(gomock.Any()).Return(nil)
				s.dbTransaction.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fun postgres.DBTransactionFunc) error {
					return fun(ctx)
				})
				s.lendRepo.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil)
				s.mailClient.EXPECT().SendMail(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			err: nil,
		},
		{
			name: "success to update invest loan without promote invested",
			params: usecase.InvestReq{
				Lender: &model.User{
					ID:      2,
					Name:    "test",
					Address: "addr",
					Email:   "email",
				},
				LoanID:   1,
				Amount:   8000000,
				UserSign: &bytes.Buffer{},
			},
			mockFunc: func() {
				s.cacheClient.EXPECT().Acquire(gomock.Any(), "invest_loan:1", 1000*time.Millisecond, 200*time.Millisecond).Return(true)
				s.cacheClient.EXPECT().Release(gomock.Any(), "invest_loan:1").Return(nil)
				s.loanRepo.EXPECT().GetByIDStatus(gomock.Any(), uint64(1), model.Approved).Return(model.Loan{
					ID:              1,
					UserID:          1,
					PrincipalAmount: 10000000,
				}, nil)
				s.lendRepo.EXPECT().GetByLoanId(gomock.Any(), uint64(1)).Return([]model.Lend{
					{
						UserID: 3,
						Amount: 1000000,
					},
				}, nil)
				s.userRepo.EXPECT().GetByLoanId(gomock.Any(), uint64(1)).Return(model.User{
					Name:    "test",
					Address: "addr",
				}, nil)
				s.lendRepo.EXPECT().GetByUidLoanId(gomock.Any(), uint64(1), uint64(2)).Return(model.Lend{
					ID:                1,
					UserSignPath:      "sign",
					AgreementFilePath: "file",
				}, nil)
				s.pdfGenerator.EXPECT().GenerateAgreementLetter(gomock.Any()).Return(nil)
				s.dbTransaction.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fun postgres.DBTransactionFunc) error {
					return fun(ctx)
				})
				s.lendRepo.EXPECT().Update(gomock.Any(), model.Lend{ID: 1, Amount: 8000000}).Return(nil)
				s.mailClient.EXPECT().SendMail(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			err: nil,
		},
		{
			name: "success to update invest loan",
			params: usecase.InvestReq{
				Lender: &model.User{
					ID:      2,
					Name:    "test",
					Address: "addr",
					Email:   "email",
				},
				LoanID:   1,
				Amount:   9000000,
				UserSign: &bytes.Buffer{},
			},
			mockFunc: func() {
				s.cacheClient.EXPECT().Acquire(gomock.Any(), "invest_loan:1", 1000*time.Millisecond, 200*time.Millisecond).Return(true)
				s.cacheClient.EXPECT().Release(gomock.Any(), "invest_loan:1").Return(nil)
				s.loanRepo.EXPECT().GetByIDStatus(gomock.Any(), uint64(1), model.Approved).Return(model.Loan{
					ID:              1,
					UserID:          1,
					PrincipalAmount: 10000000,
				}, nil)
				s.lendRepo.EXPECT().GetByLoanId(gomock.Any(), uint64(1)).Return([]model.Lend{
					{
						UserID: 3,
						Amount: 1000000,
					},
				}, nil)
				s.userRepo.EXPECT().GetByLoanId(gomock.Any(), uint64(1)).Return(model.User{
					Name:    "test",
					Address: "addr",
				}, nil)
				s.lendRepo.EXPECT().GetByUidLoanId(gomock.Any(), uint64(1), uint64(2)).Return(model.Lend{
					ID:                1,
					UserSignPath:      "sign",
					AgreementFilePath: "file",
				}, nil)
				s.pdfGenerator.EXPECT().GenerateAgreementLetter(gomock.Any()).Return(nil)
				s.dbTransaction.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fun postgres.DBTransactionFunc) error {
					return fun(ctx)
				})
				s.lendRepo.EXPECT().Update(gomock.Any(), model.Lend{ID: 1, Amount: 9000000}).Return(nil)
				s.loanRepo.EXPECT().PromoteLoanToInvested(gomock.Any(), uint64(1)).Return(nil)
				s.mailClient.EXPECT().SendMail(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			err: nil,
		},
		{
			name: "error while get loan data",
			params: usecase.InvestReq{
				Lender: &model.User{
					ID:      2,
					Name:    "test",
					Address: "addr",
					Email:   "email",
				},
				LoanID:   1,
				Amount:   8000000,
				UserSign: &bytes.Buffer{},
			},
			mockFunc: func() {
				s.cacheClient.EXPECT().Acquire(gomock.Any(), "invest_loan:1", 1000*time.Millisecond, 200*time.Millisecond).Return(true)
				s.cacheClient.EXPECT().Release(gomock.Any(), "invest_loan:1").Return(nil)
				s.loanRepo.EXPECT().GetByIDStatus(gomock.Any(), uint64(1), model.Approved).Return(model.Loan{}, errors.New("error get"))

			},
			err: errors.New("error get"),
		},
		{
			name: "error while get lend data",
			params: usecase.InvestReq{
				Lender: &model.User{
					ID:      2,
					Name:    "test",
					Address: "addr",
					Email:   "email",
				},
				LoanID:   1,
				Amount:   8000000,
				UserSign: &bytes.Buffer{},
			},
			mockFunc: func() {
				s.cacheClient.EXPECT().Acquire(gomock.Any(), "invest_loan:1", 1000*time.Millisecond, 200*time.Millisecond).Return(true)
				s.cacheClient.EXPECT().Release(gomock.Any(), "invest_loan:1").Return(nil)
				s.loanRepo.EXPECT().GetByIDStatus(gomock.Any(), uint64(1), model.Approved).Return(model.Loan{
					ID:              1,
					UserID:          1,
					PrincipalAmount: 10000000,
				}, nil)
				s.lendRepo.EXPECT().GetByLoanId(gomock.Any(), uint64(1)).Return([]model.Lend{}, errors.New("error get"))
			},
			err: errors.New("error get"),
		},
		{
			name: "error while get user lender",
			params: usecase.InvestReq{
				Lender: &model.User{
					ID:      2,
					Name:    "test",
					Address: "addr",
					Email:   "email",
				},
				LoanID:   1,
				Amount:   8000000,
				UserSign: &bytes.Buffer{},
			},
			mockFunc: func() {
				s.cacheClient.EXPECT().Acquire(gomock.Any(), "invest_loan:1", 1000*time.Millisecond, 200*time.Millisecond).Return(true)
				s.cacheClient.EXPECT().Release(gomock.Any(), "invest_loan:1").Return(nil)
				s.loanRepo.EXPECT().GetByIDStatus(gomock.Any(), uint64(1), model.Approved).Return(model.Loan{
					ID:              1,
					UserID:          1,
					PrincipalAmount: 10000000,
				}, nil)
				s.lendRepo.EXPECT().GetByLoanId(gomock.Any(), uint64(1)).Return([]model.Lend{
					{
						UserID: 3,
						Amount: 1000000,
					},
				}, nil)
				s.userRepo.EXPECT().GetByLoanId(gomock.Any(), uint64(1)).Return(model.User{}, errors.New("error get"))

			},
			err: errors.New("error get"),
		},
		{
			name: "error while generate pdf",
			params: usecase.InvestReq{
				Lender: &model.User{
					ID:      2,
					Name:    "test",
					Address: "addr",
					Email:   "email",
				},
				LoanID:   1,
				Amount:   8000000,
				UserSign: &bytes.Buffer{},
			},
			mockFunc: func() {
				s.cacheClient.EXPECT().Acquire(gomock.Any(), "invest_loan:1", 1000*time.Millisecond, 200*time.Millisecond).Return(true)
				s.cacheClient.EXPECT().Release(gomock.Any(), "invest_loan:1").Return(nil)
				s.loanRepo.EXPECT().GetByIDStatus(gomock.Any(), uint64(1), model.Approved).Return(model.Loan{
					ID:              1,
					UserID:          1,
					PrincipalAmount: 10000000,
				}, nil)
				s.lendRepo.EXPECT().GetByLoanId(gomock.Any(), uint64(1)).Return([]model.Lend{
					{
						UserID: 3,
						Amount: 1000000,
					},
				}, nil)
				s.userRepo.EXPECT().GetByLoanId(gomock.Any(), uint64(1)).Return(model.User{
					Name:    "test",
					Address: "addr",
				}, nil)
				s.lendRepo.EXPECT().GetByUidLoanId(gomock.Any(), uint64(1), uint64(2)).Return(model.Lend{}, errors.New("error get"))
			},
			err: errors.New("error get"),
		},
		{
			name: "error while upload image",
			params: usecase.InvestReq{
				Lender: &model.User{
					ID:      2,
					Name:    "test",
					Address: "addr",
					Email:   "email",
				},
				LoanID:   1,
				Amount:   8000000,
				UserSign: &bytes.Buffer{},
			},
			mockFunc: func() {
				s.cacheClient.EXPECT().Acquire(gomock.Any(), "invest_loan:1", 1000*time.Millisecond, 200*time.Millisecond).Return(true)
				s.cacheClient.EXPECT().Release(gomock.Any(), "invest_loan:1").Return(nil)
				s.loanRepo.EXPECT().GetByIDStatus(gomock.Any(), uint64(1), model.Approved).Return(model.Loan{
					ID:              1,
					UserID:          1,
					PrincipalAmount: 10000000,
				}, nil)
				s.lendRepo.EXPECT().GetByLoanId(gomock.Any(), uint64(1)).Return([]model.Lend{
					{
						UserID: 3,
						Amount: 1000000,
					},
				}, nil)
				s.userRepo.EXPECT().GetByLoanId(gomock.Any(), uint64(1)).Return(model.User{
					Name:    "test",
					Address: "addr",
				}, nil)
				s.lendRepo.EXPECT().GetByUidLoanId(gomock.Any(), uint64(1), uint64(2)).Return(model.Lend{}, errors.New("not found"))
				s.storageClient.EXPECT().UploadImage(gomock.Any(), storage.USER_SIGN_DIR, gomock.Any()).Return(errors.New("error upload"))
			},
			err: errors.New("error upload"),
		},
		{
			name: "error while add lend",
			params: usecase.InvestReq{
				Lender: &model.User{
					ID:      2,
					Name:    "test",
					Address: "addr",
					Email:   "email",
				},
				LoanID:   1,
				Amount:   8000000,
				UserSign: &bytes.Buffer{},
			},
			mockFunc: func() {
				s.cacheClient.EXPECT().Acquire(gomock.Any(), "invest_loan:1", 1000*time.Millisecond, 200*time.Millisecond).Return(true)
				s.cacheClient.EXPECT().Release(gomock.Any(), "invest_loan:1").Return(nil)
				s.loanRepo.EXPECT().GetByIDStatus(gomock.Any(), uint64(1), model.Approved).Return(model.Loan{
					ID:              1,
					UserID:          1,
					PrincipalAmount: 10000000,
				}, nil)
				s.lendRepo.EXPECT().GetByLoanId(gomock.Any(), uint64(1)).Return([]model.Lend{
					{
						UserID: 3,
						Amount: 1000000,
					},
				}, nil)
				s.userRepo.EXPECT().GetByLoanId(gomock.Any(), uint64(1)).Return(model.User{
					Name:    "test",
					Address: "addr",
				}, nil)
				s.lendRepo.EXPECT().GetByUidLoanId(gomock.Any(), uint64(1), uint64(2)).Return(model.Lend{}, errors.New("not found"))
				s.storageClient.EXPECT().UploadImage(gomock.Any(), storage.USER_SIGN_DIR, gomock.Any()).Return(nil)
				s.pdfGenerator.EXPECT().GenerateAgreementLetter(gomock.Any()).Return(errors.New("error generate"))
			},
			err: errors.New("error generate"),
		},
		{
			name: "error while add lend",
			params: usecase.InvestReq{
				Lender: &model.User{
					ID:      2,
					Name:    "test",
					Address: "addr",
					Email:   "email",
				},
				LoanID:   1,
				Amount:   8000000,
				UserSign: &bytes.Buffer{},
			},
			mockFunc: func() {
				s.cacheClient.EXPECT().Acquire(gomock.Any(), "invest_loan:1", 1000*time.Millisecond, 200*time.Millisecond).Return(true)
				s.cacheClient.EXPECT().Release(gomock.Any(), "invest_loan:1").Return(nil)
				s.loanRepo.EXPECT().GetByIDStatus(gomock.Any(), uint64(1), model.Approved).Return(model.Loan{
					ID:              1,
					UserID:          1,
					PrincipalAmount: 10000000,
				}, nil)
				s.lendRepo.EXPECT().GetByLoanId(gomock.Any(), uint64(1)).Return([]model.Lend{
					{
						UserID: 3,
						Amount: 1000000,
					},
				}, nil)
				s.userRepo.EXPECT().GetByLoanId(gomock.Any(), uint64(1)).Return(model.User{
					Name:    "test",
					Address: "addr",
				}, nil)
				s.lendRepo.EXPECT().GetByUidLoanId(gomock.Any(), uint64(1), uint64(2)).Return(model.Lend{}, errors.New("not found"))
				s.storageClient.EXPECT().UploadImage(gomock.Any(), storage.USER_SIGN_DIR, gomock.Any()).Return(nil)
				s.pdfGenerator.EXPECT().GenerateAgreementLetter(gomock.Any()).Return(nil)
				s.dbTransaction.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fun postgres.DBTransactionFunc) error {
					return fun(ctx)
				})
				s.lendRepo.EXPECT().Add(gomock.Any(), gomock.Any()).Return(errors.New("error add"))
			},
			err: errors.New("error add"),
		},
	}

	for _, test := range testCases {
		s.Run(test.name, func() {
			test.mockFunc()
			err := s.lendUC.Invest(context.Background(), test.params)
			s.Equal(err, test.err)
		})
	}
}
