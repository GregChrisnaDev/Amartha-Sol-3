package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/GregChrisnaDev/Amartha-Sol-3/common/postgres"
	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/model"
)

type lendRepository struct {
	db postgres.DB
}

type LendRepository interface {
	Add(ctx context.Context, params model.Lend) error
	Update(ctx context.Context, params model.Lend) error
	GetByUID(ctx context.Context, userId uint64) ([]model.Lend, error)
	GetByLoanId(ctx context.Context, loanId uint64) ([]model.Lend, error)
	GetByUidLoanId(ctx context.Context, loanId, userId uint64) (model.Lend, error)
	GetListLenderByLoanerID(ctx context.Context, loanId, userId uint64) ([]model.Lend, error)
}

type lend struct {
	ID                uint64    `db:"id"`
	LoanID            uint64    `db:"loan_id"`
	UserID            uint64    `db:"user_id"`
	Amount            float64   `db:"amount"`
	UserSignPath      string    `db:"user_sign_path"`
	AgreementFilePath string    `db:"agreement_file_path"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}

func InitLendRepo(db postgres.DB) LendRepository {
	return &lendRepository{
		db: db,
	}
}

func (r *lendRepository) Add(ctx context.Context, params model.Lend) error {
	if _, err := r.db.ConnTx(ctx).ExecContext(ctx, "INSERT INTO lends(loan_id, user_id, amount, user_sign_path, agreement_file_path, created_at, updated_at) VALUES($1,$2,$3,$4,$5,NOW(),NOW())", params.LoanID, params.UserID, params.Amount, params.UserSignPath, params.AgreementFilePath); err != nil {
		log.Println("[LendRepo][Add] error while insert lend", err.Error())
		return err
	}

	return nil
}

func (r *lendRepository) Update(ctx context.Context, params model.Lend) error {
	if _, err := r.db.ConnTx(ctx).ExecContext(ctx, "UPDATE lends SET amount = $1, updated_at = NOW() WHERE id = $2", params.Amount, params.ID); err != nil {
		log.Println("[LendRepo][Update] error while update lend", err.Error())
		return err
	}

	return nil
}

func (r *lendRepository) GetByUID(ctx context.Context, userId uint64) ([]model.Lend, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, loan_id, user_id, amount, user_sign_path, agreement_file_path, created_at, updated_at FROM lends WHERE user_id = $1", userId)
	if err != nil {
		log.Println("[LendRepo][GetByUID] error get data", err.Error())
		return []model.Lend{}, err
	}

	var listLend []model.Lend
	for rows.Next() {
		var lend lend

		err = rows.Scan(&lend.ID, &lend.LoanID, &lend.UserID, &lend.Amount, &lend.UserSignPath, &lend.AgreementFilePath, &lend.CreatedAt, &lend.UpdatedAt)
		if err != nil {
			log.Println("[LendRepo][GetByUID] error while scan", err.Error())
			return []model.Lend{}, err
		}

		listLend = append(listLend, model.Lend{
			ID:                lend.ID,
			LoanID:            lend.LoanID,
			UserID:            lend.UserID,
			Amount:            lend.Amount,
			UserSignPath:      lend.UserSignPath,
			AgreementFilePath: lend.AgreementFilePath,
			CreatedAt:         lend.CreatedAt,
			UpdatedAt:         lend.UpdatedAt,
		})
	}

	return listLend, nil
}

func (r *lendRepository) GetByLoanId(ctx context.Context, loanId uint64) ([]model.Lend, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, loan_id, user_id, amount, user_sign_path, agreement_file_path, created_at, updated_at FROM lends WHERE loan_id = $1", loanId)
	if err != nil {
		log.Println("[LendRepo][GetByLoanId] error get data", err.Error())
		return []model.Lend{}, err
	}

	var listLend []model.Lend
	for rows.Next() {
		var lend lend

		err = rows.Scan(&lend.ID, &lend.LoanID, &lend.UserID, &lend.Amount, &lend.UserSignPath, &lend.AgreementFilePath, &lend.CreatedAt, &lend.UpdatedAt)
		if err != nil {
			log.Println("[LendRepo][GetByLoanId] error while scan", err.Error())
			return []model.Lend{}, err
		}

		listLend = append(listLend, model.Lend{
			ID:                lend.ID,
			LoanID:            lend.LoanID,
			UserID:            lend.UserID,
			Amount:            lend.Amount,
			UserSignPath:      lend.UserSignPath,
			AgreementFilePath: lend.AgreementFilePath,
			CreatedAt:         lend.CreatedAt,
			UpdatedAt:         lend.UpdatedAt,
		})
	}

	return listLend, nil
}

func (r *lendRepository) GetByUidLoanId(ctx context.Context, loanId, userId uint64) (model.Lend, error) {
	rows := r.db.QueryRowContext(ctx, "SELECT id, loan_id, user_id, amount, user_sign_path, agreement_file_path, created_at, updated_at FROM lends WHERE loan_id = $1 AND user_id = $2", loanId, userId)

	var lend lend
	if err := rows.Scan(&lend.ID, &lend.LoanID, &lend.UserID, &lend.Amount, &lend.UserSignPath, &lend.AgreementFilePath, &lend.CreatedAt, &lend.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Lend{}, errors.New("not found")
		}
		log.Println("[LendRepo][GetByUidLoanId] error while get lend", err.Error())
		return model.Lend{}, err
	}

	return model.Lend{
		ID:                lend.ID,
		LoanID:            lend.LoanID,
		UserID:            lend.UserID,
		Amount:            lend.Amount,
		UserSignPath:      lend.UserSignPath,
		AgreementFilePath: lend.AgreementFilePath,
		CreatedAt:         lend.CreatedAt,
		UpdatedAt:         lend.UpdatedAt,
	}, nil
}

func (r *lendRepository) GetListLenderByLoanerID(ctx context.Context, loanId, userId uint64) ([]model.Lend, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT le.id, le.loan_id, le.user_id, le.amount, le.agreement_file_path, le.created_at, le.updated_at FROM lends le JOIN loans lo ON le.loan_id = lo.id WHERE lo.id = $1 AND lo.user_id = $2", loanId, userId)
	if err != nil {
		log.Println("[LendRepo][GetListLenderByLoanerID] error get data", err.Error())
		return []model.Lend{}, err
	}

	var listLend []model.Lend
	for rows.Next() {
		var lend lend

		err = rows.Scan(&lend.ID, &lend.LoanID, &lend.UserID, &lend.Amount, &lend.AgreementFilePath, &lend.CreatedAt, &lend.UpdatedAt)
		if err != nil {
			log.Println("[LendRepo][GetListLenderByLoanerID] error while scan", err.Error())
			return []model.Lend{}, err
		}

		listLend = append(listLend, model.Lend{
			ID:                lend.ID,
			LoanID:            lend.LoanID,
			UserID:            lend.UserID,
			Amount:            lend.Amount,
			AgreementFilePath: lend.AgreementFilePath,
			CreatedAt:         lend.CreatedAt,
			UpdatedAt:         lend.UpdatedAt,
		})
	}

	return listLend, nil
}
