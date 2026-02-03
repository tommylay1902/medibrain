package metadata

import (
	"net/http"
)

type Route struct {
	Mux *http.ServeMux
}

func NewRoutes(dms *DocumentMetaService) *Route {
	mux := http.NewServeMux()
	handler := NewHandler(dms)

	mux.HandleFunc("GET /", handler.List)
	// mux.HandleFunc("POST /metadata", handler.GetDocumentMetaData)
	route := &Route{
		Mux: mux,
	}

	return route
}
