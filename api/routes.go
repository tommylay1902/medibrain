package api

import (
	"net/http"

	"github.com/tommylay1902/medibrain/api/domain/documentmeta"
)

type Mux struct {
	Mux *http.ServeMux
}

func NewMux() *Mux {
	documentMetaMux := documentmeta.NewRoutes()

	mainMux := http.NewServeMux()
	mainMux.Handle("/api/v1/", http.StripPrefix("/api/v1", documentMetaMux.Mux))

	return &Mux{
		Mux: mainMux,
	}
}
