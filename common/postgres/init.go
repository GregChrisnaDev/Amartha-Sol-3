package postgres

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	*sqlx.DB
}

func Init() DB {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?%s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_QUERY_STRING"),
	)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalln("[InitPostgres] error opening DB:", err.Error())
	}

	if err = db.Ping(); err != nil {
		log.Fatalln("[InitPostgres] error connecting to DB:", err.Error())
	}

	return DB{
		db,
	}
}

func (r DB) ConnTx(ctx context.Context) Transaction {
	if tx := extractTx(ctx); tx != nil {
		return tx
	}
	return r
}
