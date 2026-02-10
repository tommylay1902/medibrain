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
	// Get current working directory
	// wd, err := os.Getwd()
	// if err != nil {
	// 	panic(err)
	// }

	// Build absolute path
	// fullPath := filepath.Join(wd, "medibrain/internal/database/migrations/schemas.sql")
	fmt.Println(sqlContent)
	sql, err := os.ReadFile(sqlContent)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	db.MustExec(string(sql))
}
