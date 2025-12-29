package api

import (
	"net/http"

	"github.com/tommylay1902/medibrain/internal/api/domain/documentmeta"
)

type Mux struct {
	Mux *http.ServeMux
}

// NewMux creates and configures the main router for the API application.
// It organizes routes by domain, where each domain defines its specific routes
// (e.g., documentmeta domain defines "GET /documentmeta"). These domain-specific
// routes are then initialized here and mounted under the "/api/v1/" prefix.
//
// Returns the configured main router with all routes ready for server use.
func NewMux(dms *documentmeta.DocumentMetaService) *Mux {
	documentMetaMux := documentmeta.NewRoutes(dms)

	mainMux := http.NewServeMux()
	mainMux.Handle("/api/v1/", http.StripPrefix("/api/v1", documentMetaMux.Mux))

	return &Mux{
		Mux: mainMux,
	}
}
