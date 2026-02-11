package note

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type NoteHandler struct {
	noteService *NoteService
}

func NewNoteHandler(noteService *NoteService) *NoteHandler {
	return &NoteHandler{noteService: noteService}
}

func (nh *NoteHandler) List(w http.ResponseWriter, req *http.Request) {
	notes, err := nh.noteService.List()
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(notes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		return
	}
}

func (nh *NoteHandler) ListWithKeywords(w http.ResponseWriter, req *http.Request) {
	notesWithKeyword, err := nh.noteService.ListWithKeywords()
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(notesWithKeyword)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		return
	}
}
