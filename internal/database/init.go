package database

import "github.com/jmoiron/sqlx"

func NewDB() *sqlx.DB {
	db, err := sqlx.Connect()
}
