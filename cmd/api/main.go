package main

import (
	"github.com/tommylay1902/medibrain/internal/api"
	"github.com/tommylay1902/medibrain/internal/api/domain/document"
	"github.com/tommylay1902/medibrain/internal/api/domain/metadata"
	seaweedclient "github.com/tommylay1902/medibrain/internal/client/seaweed"
	"github.com/tommylay1902/medibrain/internal/client/stirling"
	"github.com/tommylay1902/medibrain/internal/database"
)

func main() {
	db := database.NewDB()
	dmr := metadata.NewRepo(db)
	dms := metadata.NewService(dmr)

	sc := stirling.NewClient()
	swc := seaweedclient.NewClient()
	dps := document.NewService(dmr, swc, sc, dms)
	mux := api.NewMux(dms, dps)
	server := api.NewServer(":8080", mux)
	server.StartServer()
}
