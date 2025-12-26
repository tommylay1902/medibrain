package documentmeta

import "net/http"

type Route struct {
	Mux *http.ServeMux
}

func NewRoutes() *Route {
	mux := http.NewServeMux()
	handler := NewHandler()

	mux.HandleFunc("GET /documentmeta", handler.List)

	route := &Route{
		Mux: mux,
	}

	return route
}
