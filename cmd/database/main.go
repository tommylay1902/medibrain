package main

import "github.com/tommylay1902/medibrain/internal/database"

func main() {
	db := database.NewDB()
	database.CreateSchema(db)
}
