package database

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

//go:embed migrations/schemas.sql
var sqlContent string

func NewDB() *sqlx.DB {
	db, err := sqlx.Connect("postgres", "host=localhost port=5432 user=root password=1234 dbname=medibrain sslmode=disable")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println("succesfully connected to db")
	return db
}

func CreateSchema(db *sqlx.DB) {
	fmt.Println(sqlContent)
	sql, err := os.ReadFile(sqlContent)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	db.MustExec(string(sql))
}
