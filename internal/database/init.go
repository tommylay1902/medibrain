package database

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDB() *sqlx.DB {
	// Use the service name "db" as host
	db, err := sqlx.Connect("postgres", "host=localhost port=5432 user=root password=1234 dbname=medibrain sslmode=disable")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println("succesfully conncted to db")
	return db
}

func CreateSchema(db *sqlx.DB) {
	sql, err := os.ReadFile("internal/database/migrations/test2.sql")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	db.MustExec(string(sql))
}
