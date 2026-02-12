package main

import (
	"github.com/tommylay1902/medibrain/internal/api"
	"github.com/tommylay1902/medibrain/internal/api/domain/document"
	"github.com/tommylay1902/medibrain/internal/api/domain/metadata"
	"github.com/tommylay1902/medibrain/internal/api/domain/note"
	"github.com/tommylay1902/medibrain/internal/api/domain/tag"
	"github.com/tommylay1902/medibrain/internal/client/rag"
	seaweedclient "github.com/tommylay1902/medibrain/internal/client/seaweed"
	"github.com/tommylay1902/medibrain/internal/client/stirling"
	"github.com/tommylay1902/medibrain/internal/database"
)

func main() {
	rag := rag.NewRag()
	db := database.NewDB()
	uowFactory := database.NewUnitOfWorkFactory(db)
	dmr := metadata.NewRepo(db)
	nr := note.NewNoteRepo(uowFactory)
	tr := tag.NewTagRepo(db)

	dms := metadata.NewService(dmr)
	ns := note.NewNoteService(nr, uowFactory)
	ts := tag.NewTagService(tr)
	sc := stirling.NewClient()
	swc := seaweedclient.NewClient()

	dps := document.NewService(dmr, swc, sc, dms, rag)

	mux := api.NewMux(dms, dps, ns, ts)

	server := api.NewServer(":8080", mux)
	server.StartServer()
}
