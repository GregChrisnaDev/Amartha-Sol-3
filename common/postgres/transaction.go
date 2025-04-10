package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Transaction is an interface that models the standard transaction in
// `database/sql`.
//
// To ensure `TxFn` funcs cannot commit or rollback a transaction (which is
// handled by `WithTransaction`), those methods are not included here.
type Transaction interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	NamedExecContext(ctx context.Context, query string, args interface{}) (sql.Result, error)
}

// A Txfn is a function that will be called with an initialized `Transaction` object
// that can be used for executing statements and queries against a database.
type TxFn func(Transaction) error

// WithTransaction creates a new transaction and handles rollback/commit based on the
// error object returned by the `TxFn`
func WithTransaction(ctx context.Context, db *sqlx.DB, fn TxFn) (err error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			_ = tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}

type DBTransactionFunc func(ctx context.Context) error

type DBTransaction interface {
	Execute(ctx context.Context, fun DBTransactionFunc) (err error)
}

type dbTransaction struct {
	db *sqlx.DB
}

func NewDBTransaction(db *sqlx.DB) DBTransaction {
	return dbTransaction{
		db: db,
	}
}

func (r dbTransaction) Execute(ctx context.Context, fun DBTransactionFunc) (err error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			_ = tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit()
		}
	}()

	return fun(r.injectTx(ctx, tx))
}

type txKey struct{}

func (r dbTransaction) injectTx(ctx context.Context, tx Transaction) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func extractTx(ctx context.Context) Transaction {
	if tx, ok := ctx.Value(txKey{}).(Transaction); ok {
		return tx
	}
	return nil
}
