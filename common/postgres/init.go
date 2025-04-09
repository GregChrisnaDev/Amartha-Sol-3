package postgres

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Init() *sqlx.DB {
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

	return db
}
