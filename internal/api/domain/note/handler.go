package note

import "net/http"

type NoteHandler struct {
	noteService *NoteService
}

func NewNoteHandler(noteService *NoteService) *NoteHandler {
	return &NoteHandler{noteService: noteService}
}

func (nh *NoteHandler) List(w http.ResponseWriter, req *http.Request) {
}
