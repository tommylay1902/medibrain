package note

import "net/http"

type Route struct {
	Mux *http.ServeMux
}

func NewNoteRoutes(ns *NoteService) *Route {
	mux := http.NewServeMux()
	handler := NewNoteHandler(ns)
	mux.HandleFunc("GET /", handler.ListWithKeywords)
	mux.HandleFunc("GET /tag", handler.ListTags)
	mux.HandleFunc("POST /", handler.CreateNote)
	mux.HandleFunc("POST /tag", handler.CreateTag)
	route := &Route{
		Mux: mux,
	}
	return route
}
