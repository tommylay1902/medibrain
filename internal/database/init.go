package database

import (
	_ "embed"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDB() *sqlx.DB {
	db, err := sqlx.Connect("postgres", "host=medibrain-db port=5432 user=root password=1234 dbname=medibrain sslmode=disable")
	if err != nil {
		slog.Error("UPDATE")
		slog.Error(err.Error())
		panic(err)
	}

	fmt.Println("succesfully connected to db")
	return db
}
