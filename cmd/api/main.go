package main

import (
	"github.com/tommylay1902/medibrain/internal/api"
	"github.com/tommylay1902/medibrain/internal/api/domain/document"
	"github.com/tommylay1902/medibrain/internal/api/domain/metadata"
	"github.com/tommylay1902/medibrain/internal/api/domain/note"
	"github.com/tommylay1902/medibrain/internal/client/rag"
	seaweedclient "github.com/tommylay1902/medibrain/internal/client/seaweed"
	"github.com/tommylay1902/medibrain/internal/client/stirling"
	"github.com/tommylay1902/medibrain/internal/database"
)

func main() {
	rag := rag.NewRag()
	db := database.NewDB()
	dmr := metadata.NewRepo(db)
	nr := note.NewNoteRepo(db)

	dms := metadata.NewService(dmr)
	ns := note.NewNoteService(nr)

	sc := stirling.NewClient()
	swc := seaweedclient.NewClient()

	dps := document.NewService(dmr, swc, sc, dms, rag)

	mux := api.NewMux(dms, dps, ns)

	server := api.NewServer(":8080", mux)
	server.StartServer()
}
