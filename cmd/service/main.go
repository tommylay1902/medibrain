package main

import (
	"github.com/tommylay1902/medibrain/internal/api"
	"github.com/tommylay1902/medibrain/internal/database"
)

func main() {
	db := database.NewDB()
	mux := api.NewMux(db)

	server := api.NewServer(":8080", mux)
	server.StartServer()
}
