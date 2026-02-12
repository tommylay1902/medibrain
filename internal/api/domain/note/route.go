package note

import "net/http"

type Route struct {
	Mux *http.ServeMux
}

func NewNoteRoutes(ns *NoteService) *Route {
	mux := http.NewServeMux()
	handler := NewNoteHandler(ns)
	mux.HandleFunc("GET /", handler.List)
	mux.HandleFunc("GET /keywords", handler.ListWithKeywords)
	mux.HandleFunc("POST /", handler.CreateNote)
	route := &Route{
		Mux: mux,
	}
	return route
}
