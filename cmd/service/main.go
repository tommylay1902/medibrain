package main

import (
	"github.com/tommylay1902/medibrain/internal/api"
	"github.com/tommylay1902/medibrain/internal/api/domain/documentmeta"
	"github.com/tommylay1902/medibrain/internal/database"
)

func main() {
	db := database.NewDB()
	dmr := documentmeta.NewRepo(db)
	dms := documentmeta.NewService(dmr)
	mux := api.NewMux(dms)

	server := api.NewServer(":8080", mux)
	server.StartServer()
}
