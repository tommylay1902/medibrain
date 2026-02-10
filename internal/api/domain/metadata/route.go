package metadata

import (
	"net/http"
)

type Route struct {
	Mux *http.ServeMux
}

func NewRoutes(dms *MetadataService) *Route {
	mux := http.NewServeMux()
	handler := NewHandler(dms)

	mux.HandleFunc("GET /", handler.List)
	// mux.HandleFunc("POST /metadata", handler.GetMetadataData)
	route := &Route{
		Mux: mux,
	}

	return route
}
