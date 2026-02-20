package note

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/lib/pq"
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

func (nh *NoteHandler) ListTags(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	tags, err := nh.noteService.ListTag(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(tags)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		return
	}
}

type CreateNoteBody struct {
	Note
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

func (nh *NoteHandler) CreateTag(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var body Tag
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Bad request: %v", err), http.StatusBadRequest)
		return
	}

	tag, err := nh.noteService.CreateTag(ctx, body)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			fmt.Println("test err")
			http.Error(w, fmt.Sprintf("tag already exists: %v", pqErr), http.StatusConflict)
			return
		}
		http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(tag)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		return
	}
}
