package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/model"
	"github.com/jmoiron/sqlx"
)

type loanRepository struct {
	db *sqlx.DB
}

type LoanRepository interface {
	ProposeLoan(ctx context.Context, params model.Loan) error
	GetLoanByUID(ctx context.Context, userId uint64) ([]model.Loan, error)
	PendingLoanExist(ctx context.Context, userId uint64) (bool, error)
	GetLoanByID(ctx context.Context, loanId uint64) (model.Loan, error)
	PromoteLoanToApproved(ctx context.Context, params model.Loan) error
	LoanExist(ctx context.Context, loanId uint64, status int8) (bool, error)
}

type loan struct {
	ID                   uint64         `db:"id"`
	UserID               uint64         `db:"user_id"`
	PrincipalAmount      float64        `db:"principal_amount"`
	Rate                 uint32         `db:"rate"`
	LoanDuration         uint32         `db:"loan_duration"`
	Status               int8           `db:"status"`
	ProposedDate         time.Time      `db:"proposed_date"`
	PictureProofFilePath sql.NullString `db:"picture_proof_filepath"`
	ApproverUID          sql.NullInt64  `db:"approver_uid"`
	ApprovalDate         sql.NullTime   `db:"approval_date"`
	DisburserUID         sql.NullInt64  `db:"disburser_uid"`
	DisbursedDate        sql.NullTime   `db:"disbursement_date"`
}

func InitLoanRepo(db *sqlx.DB) LoanRepository {
	return &loanRepository{
		db: db,
	}
}

func (r *loanRepository) ProposeLoan(ctx context.Context, params model.Loan) error {
	if _, err := r.db.ExecContext(ctx, "INSERT INTO loans(user_id, principal_amount, rate, loan_duration, status, proposed_date) VALUES($1,$2,$3,$4,$5,NOW())", params.UserID, params.PrincipalAmount, params.Rate, params.LoanDuration, model.Proposed); err != nil {
		log.Println("[LoanRepo][ProposeLoan] error while insert user", err.Error())
		return err
	}

	return nil
}

func (r *loanRepository) GetLoanByUID(ctx context.Context, userId uint64) ([]model.Loan, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, user_id, principal_amount, rate, loan_duration, status, proposed_date, picture_proof_filepath, approver_uid, approval_date, disburser_uid, disbursement_date FROM loans WHERE user_id = $1", userId)
	if err != nil {
		log.Println("[LoanRepo][GetLoanByUID] error get data", err.Error())
		return []model.Loan{}, err
	}

	var listLoan []model.Loan
	for rows.Next() {
		var loan loan

		err = rows.Scan(&loan.ID, &loan.UserID, &loan.PrincipalAmount, &loan.Rate, &loan.LoanDuration, &loan.Status, &loan.ProposedDate, &loan.PictureProofFilePath, &loan.ApproverUID, &loan.ApprovalDate, &loan.DisburserUID, &loan.DisbursedDate)
		if err != nil {
			log.Println("[LoanRepo][GetLoanByUID] error while scan", err.Error())
			return []model.Loan{}, err
		}

		listLoan = append(listLoan, model.Loan{
			ID:                   loan.ID,
			UserID:               loan.UserID,
			PrincipalAmount:      loan.PrincipalAmount,
			Rate:                 loan.Rate,
			LoanDuration:         loan.LoanDuration,
			Status:               loan.Status,
			ProposedDate:         loan.ProposedDate,
			PictureProofFilePath: loan.PictureProofFilePath.String,
			ApproverUID:          uint64(loan.ApproverUID.Int64),
			ApprovalDate:         loan.ApprovalDate.Time,
			DisburserUID:         uint64(loan.DisburserUID.Int64),
			DisbursedDate:        loan.DisbursedDate.Time,
		})
	}

	return listLoan, nil
}

func (r *loanRepository) PendingLoanExist(ctx context.Context, userId uint64) (bool, error) {
	var exist bool
	err := r.db.QueryRowContext(ctx, "SELECT 1 FROM loans WHERE user_id = $1 AND status != $2", userId, model.Disbursed).Scan(&exist)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		log.Println("[LoanRepo][PendingLoanExist] error while scan", err.Error())
		return false, err
	}

	return exist, nil
}

func (r *loanRepository) LoanExist(ctx context.Context, loanId uint64, status int8) (bool, error) {
	var exist bool
	var err error
	args := []interface{}{loanId}
	query := "SELECT 1 FROM loans WHERE id = $1"
	if status != 0 {
		query += " AND status = $2"
		query, args, err = sqlx.In(query, loanId, status)
		if err != nil {
			return false, err
		}
	}

	query = r.db.Rebind(query)
	err = r.db.GetContext(ctx, &exist, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		log.Println("[LoanRepo][LoanExist] error while scan", err.Error())
		return false, err
	}

	return exist, nil
}

func (r *loanRepository) GetLoanByID(ctx context.Context, loanId uint64) (model.Loan, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, user_id, principal_amount, rate, loan_duration, status, proposed_date, picture_proof_filepath, approver_uid, approval_date, disburser_uid, disbursement_date FROM loans WHERE id = $1", loanId)

	var loan loan
	err := row.Scan(&loan.ID, &loan.UserID, &loan.PrincipalAmount, &loan.Rate, &loan.LoanDuration, &loan.Status, &loan.ProposedDate, &loan.PictureProofFilePath, &loan.ApproverUID, &loan.ApprovalDate, &loan.DisburserUID, &loan.DisbursedDate)
	if err != nil {
		log.Println("[LoanRepo][GetLoanByID] error while scan", err.Error())
		return model.Loan{}, err
	}

	return model.Loan{
		ID:                   loan.ID,
		UserID:               loan.UserID,
		PrincipalAmount:      loan.PrincipalAmount,
		Rate:                 loan.Rate,
		LoanDuration:         loan.LoanDuration,
		Status:               loan.Status,
		ProposedDate:         loan.ProposedDate,
		PictureProofFilePath: loan.PictureProofFilePath.String,
		ApproverUID:          uint64(loan.ApproverUID.Int64),
		ApprovalDate:         loan.ApprovalDate.Time,
		DisburserUID:         uint64(loan.DisburserUID.Int64),
		DisbursedDate:        loan.DisbursedDate.Time,
	}, nil
}

func (r *loanRepository) PromoteLoanToApproved(ctx context.Context, params model.Loan) error {
	if _, err := r.db.ExecContext(ctx, "UPDATE loans SET picture_proof_filepath=$1, approver_uid=$2, approval_date=NOW(), status=$3 WHERE id=$4", params.PictureProofFilePath, params.ApproverUID, model.Approved, params.ID); err != nil {
		log.Println("[LoanRepo][PromoteLoanToApproved] error while promote loan to approved", err.Error())
		return err
	}

	return nil
}
