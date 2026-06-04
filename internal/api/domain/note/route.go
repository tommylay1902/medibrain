package note

import (
	"net/http"
)

type Route struct {
	Mux *http.ServeMux
}

func NewNoteRoutes(ns *NoteService) *Route {
	mux := http.NewServeMux()
	handler := NewNoteHandler(ns)
	mux.HandleFunc("GET /", handler.ListWithKeywords)
	mux.HandleFunc("GET /tag", handler.ListTags)
	// mux.HandleFunc("POST /note", handler.CreateNote)
	mux.HandleFunc("POST /", handler.CreateNote)
	mux.HandleFunc("POST /tag", handler.CreateTag)
	mux.HandleFunc("POST /chunk", handler.ChunkAndUploadNote)
	route := &Route{
		Mux: mux,
	}
	return route
}
