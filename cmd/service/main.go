package main

import (
	"github.com/tommylay1902/medibrain/internal/api"
	"github.com/tommylay1902/medibrain/internal/database"
)

func main() {
	db := database.NewDB()
	database.CreateSchema(db)
	mux := api.NewMux()

	server := api.NewServer(":8080", mux)
	server.StartServer()
}
