package api

import (
	"net/http"

	"github.com/tommylay1902/medibrain/internal/api/domain/documentmeta"
	"github.com/tommylay1902/medibrain/internal/api/domain/documentpipeline"
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
func NewMux(dms *documentmeta.DocumentMetaService, dps *documentpipeline.DocumentPipelineService) *Mux {
	mainMux := http.NewServeMux()

	// Create API v1 subrouter
	apiV1 := http.NewServeMux()

	// Get domain muxes
	documentMetaMux := documentmeta.NewRoutes(dms)
	documentPipelineMux := documentpipeline.NewRoutes(dps)

	// Mount with prefixes
	mountSubrouter(apiV1, "/documentmeta", documentMetaMux.Mux)
	mountSubrouter(apiV1, "/documentpipeline", documentPipelineMux.Mux)

	mainMux.Handle("/api/v1/",
		// applyMiddleware(
		http.StripPrefix("/api/v1", apiV1),
		// loggingMiddleware,
		// corsMiddleware,

		// )
	)

	// mainMux.HandleFunc("/", rootHandler)

	return &Mux{Mux: mainMux}
}

// func applyMiddleware(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
// 	for i := len(middlewares) - 1; i >= 0; i-- {
// 		h = middlewares[i](h)
// 	}
// 	return h
// }

func mountSubrouter(parent *http.ServeMux, prefix string, child *http.ServeMux) {
	if child == nil {
		return
	}
	parent.Handle(prefix+"/", http.StripPrefix(prefix, child))
}
