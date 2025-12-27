package api

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/tommylay1902/medibrain/internal/api/domain/documentmeta"
)

type Mux struct {
	Mux *http.ServeMux
}

func NewMux(db *sqlx.DB) *Mux {
	documentMetaMux := documentmeta.NewRoutes(db)

	mainMux := http.NewServeMux()
	mainMux.Handle("/api/v1/", http.StripPrefix("/api/v1", documentMetaMux.Mux))

	return &Mux{
		Mux: mainMux,
	}
}
