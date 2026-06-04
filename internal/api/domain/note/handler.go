package note

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/lib/pq"
	"github.com/tommylay1902/medibrain/internal/client/rag"
)

type NoteHandler struct {
	noteService *NoteService
	ragClient   *rag.Rag
}

func NewNoteHandler(noteService *NoteService) *NoteHandler {
	return &NoteHandler{
		noteService: noteService,
	}
}

func (nh *NoteHandler) List(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	notes, err := nh.noteService.List(ctx)
	if err != nil {
		slog.Error("error listing note", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(notes)
	if err != nil {
		slog.Error("error encoding notes", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		return
	}
}

func (nh *NoteHandler) ListWithKeywords(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	notesWithKeyword, err := nh.noteService.ListWithKeywords(ctx)
	if err != nil {
		slog.Error("error listing notes with keywords", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(notesWithKeyword)
	if err != nil {
		slog.Error("error encoding notes with keywords", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		return
	}
}

func (nh *NoteHandler) ListTags(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	tags, err := nh.noteService.ListTag(ctx)
	if err != nil {
		slog.Error("error getting tags", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(tags)
	if err != nil {
		slog.Error("error encoding tags", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		return
	}
}

type CreateNoteBody struct {
	Note `json:"note"`
	Tags []string `json:"tags"`
}

func (nh *NoteHandler) CreateNote(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var body CreateNoteBody
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		slog.Error("can't parse body", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("Bad request error"), http.StatusBadRequest)
		return
	}

	noteId, err := nh.noteService.CreateNoteWithTags(ctx, &body.Note, body.Tags)
	if err != nil {
		slog.Error("error creating note with tags", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"id":      noteId,
		"message": "Note created successfully",
		"status":  "success",
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Error("error creating note with tags", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		return
	}
}

func (nh *NoteHandler) CreateTag(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var body Tag
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		slog.Error("parse body err", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("Bad request"), http.StatusBadRequest)
		return
	}

	tag, err := nh.noteService.CreateTag(ctx, body)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			slog.Error("unique constraint error", slog.Any("err", err))
			http.Error(w, fmt.Sprintf("tag already exists"), http.StatusConflict)
			return
		}
		slog.Error("error creating tag", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(tag)
	if err != nil {
		slog.Error("error encoding tag", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		return
	}
}

func (nh *NoteHandler) ChunkAndUploadNote(w http.ResponseWriter, req *http.Request) {
	var note Note

	err := json.NewDecoder(req.Body).Decode(&note)
	if err != nil {
		slog.Error("Error parsing request body")
		http.Error(w, "Couldn't parse request body", http.StatusBadRequest)
	}
	err = nh.noteService.StoreNote(note)
	if err != nil {
		slog.Error("error indexing note", slog.Any("err", err))
		http.Error(w, "error indexing note", http.StatusInternalServerError)
		return
	}
}
