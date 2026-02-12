package api

import (
	"net/http"
	"strings"

	"github.com/tommylay1902/medibrain/internal/api/domain/document"
	"github.com/tommylay1902/medibrain/internal/api/domain/metadata"
	"github.com/tommylay1902/medibrain/internal/api/domain/note"
)

type Mux struct {
	Mux *http.ServeMux
}

// NewMux creates and configures the main router for the API application.
// It organizes routes by domain, where each domain defines its specific routes
// This function will add the base domain path to each route but each domain
// will implement more specific paths if needed
// (e.g., metadata domain defines "GET /"). These domain-specific
// routes are then initialized here and mounted under the "/api/v1/" prefix.
//
// Returns the configured main router with all routes ready for server use.
func NewMux(dms *metadata.MetadataService, dps *document.DocumentPipelineService, ns *note.NoteService) *Mux {
	mainMux := http.NewServeMux()

	// Create API v1 subrouter
	apiV1 := http.NewServeMux()

	// Get domain muxes
	metadataMux := metadata.NewRoutes(dms)
	documentPipelineMux := document.NewRoutes(dps)
	noteMux := note.NewNoteRoutes(ns)

	// Mount with prefixes
	mountSubrouter(apiV1, "metadata", metadataMux.Mux)
	mountSubrouter(apiV1, "document", documentPipelineMux.Mux)
	mountSubrouter(apiV1, "note", noteMux.Mux)
	apiV1Handler := applyMiddleware(http.StripPrefix("/api/v1", apiV1), CorsMiddleware)
	mainMux.Handle("/api/v1/", apiV1Handler)

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

		urlCopy := *r.URL
		urlCopy.Path = strippedPath
		r2.URL = &urlCopy
		r2.URL.Path = strippedPath
		r2.Header = r.Header.Clone()

		child.ServeHTTP(w, &r2)
	})

	parent.Handle(prefix, handler)
	parent.Handle(prefix+"/", handler)
}
