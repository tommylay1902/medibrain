package tag

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type TagHandler struct {
	ts *TagService
}

func NewTagHandler(ts *TagService) *TagHandler {
	return &TagHandler{
		ts: ts,
	}
}

func (th *TagHandler) List(w http.ResponseWriter, req *http.Request) {
	tags, err := th.ts.List()
	if err != nil {
		http.Error(w, fmt.Sprintf("internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(tags)
	if err != nil {
		http.Error(w, fmt.Sprintf("interal server error: %v", err), http.StatusInternalServerError)
		return
	}
}

func (th *TagHandler) Create(w http.ResponseWriter, req *http.Request) {
	var tag Tag
	err := json.NewDecoder(req.Body).Decode(&tag)
	if err != nil {
		http.Error(w, fmt.Sprintf("Bad request error: %v", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = th.ts.Create(&tag)
	if err != nil {
		http.Error(w, fmt.Sprintf("interal server error: %v", err), http.StatusInternalServerError)
		return
	}
}
