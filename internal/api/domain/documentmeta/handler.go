package documentmeta

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type Handler struct {
	repo DocumentMetaRepo
}

func NewHandler(db *sqlx.DB) *Handler {
	repo := DocumentMetaRepo{db: db}
	return &Handler{repo: repo}
}

func (h *Handler) List(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	pdfs, err := h.repo.List()
	if err != nil {
		fmt.Println("error getting pdfs list")
		w.WriteHeader(500)
		return
	}
	result, err := json.Marshal(pdfs)
	if err != nil {
		fmt.Println("error mashaling object")
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	_, err = w.Write(result)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println("error writing result")
		w.WriteHeader(500)
		return
	}
}
