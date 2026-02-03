package api

import (
	"net/http"
	"strings"

	"github.com/tommylay1902/medibrain/internal/api/domain/document"
	"github.com/tommylay1902/medibrain/internal/api/domain/metadata"
)

type Mux struct {
	Mux *http.ServeMux
}

// NewMux creates and configures the main router for the API application.
// It organizes routes by domain, where each domain defines its specific routes
// (e.g., metadata domain defines "GET /metadata"). These domain-specific
// routes are then initialized here and mounted under the "/api/v1/" prefix.
//
// Returns the configured main router with all routes ready for server use.
func NewMux(dms *metadata.DocumentMetaService, dps *document.DocumentPipelineService) *Mux {
	mainMux := http.NewServeMux()

	// Create API v1 subrouter
	apiV1 := http.NewServeMux()

	// Get domain muxes
	documentMetaMux := metadata.NewRoutes(dms)
	documentPipelineMux := document.NewRoutes(dps)

	// Mount with prefixes
	mountSubrouter(apiV1, "metadata", documentMetaMux.Mux)
	mountSubrouter(apiV1, "document", documentPipelineMux.Mux)
	apiV1Handler := applyMiddleware(http.StripPrefix("/api/v1", apiV1), CorsMiddleware)
	mainMux.Handle("/api/v1/", apiV1Handler)

	// mainMux.HandleFunc("/", rootHandler)

	return &Mux{Mux: mainMux}
}

func applyMiddleware(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func mountSubrouter(parent *http.ServeMux, prefix string, child *http.ServeMux) {
	if child == nil {
		return
	}

	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		strippedPath := strings.TrimPrefix(r.URL.Path, prefix)
		if strippedPath == "" {
			strippedPath = "/"
		}

		r2 := *r
		r2.URL = &(*r.URL)
		r2.URL.Path = strippedPath

		child.ServeHTTP(w, &r2)
	})

	parent.Handle(prefix, handler)
	parent.Handle(prefix+"/", handler) // Also handle with trailing slash
	parent.Handle(prefix+"/*", handler)
}
