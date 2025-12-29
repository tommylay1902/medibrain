package documentmeta

import (
	"net/http"
)

type Route struct {
	Mux *http.ServeMux
}

func NewRoutes(dms *DocumentMetaService) *Route {
	mux := http.NewServeMux()
	handler := NewHandler(dms)

	mux.HandleFunc("GET /documentmeta", handler.List)
	// mux.HandleFunc("POST /documentmeta", handler.GetDocumentMetaData)
	route := &Route{
		Mux: mux,
	}

	return route
}
