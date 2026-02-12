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
	ctx := req.Context()
	notes, err := nh.noteService.List(ctx)
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
	ctx := req.Context()
	notesWithKeyword, err := nh.noteService.ListWithKeywords(ctx)
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

type CreateNoteBody struct {
	Note Note     `json:"note"`
	Tags []string `json:"tags"`
}

func (nh *NoteHandler) CreateNote(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var body CreateNoteBody
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Bad request error; %v", err), http.StatusBadRequest)
		return
	}
	err = nh.noteService.CreateNoteWithTags(ctx, &body.Note, body.Tags)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
