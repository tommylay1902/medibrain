package tag

import "net/http"

type Route struct {
	Mux *http.ServeMux
}

func NewTagRoutes(ts *TagService) *Route {
	mux := http.NewServeMux()
	handler := NewTagHandler(ts)
	mux.HandleFunc("GET /", handler.List)
	route := &Route{
		Mux: mux,
	}
	return route
}
