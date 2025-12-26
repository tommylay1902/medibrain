package database

import (
	"fmt"

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
