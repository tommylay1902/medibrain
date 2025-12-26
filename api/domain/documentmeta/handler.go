package documentmeta

import (
	"net/http"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) List(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(documentMeta)
}
