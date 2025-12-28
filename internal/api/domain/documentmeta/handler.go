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

//
// func (h *Handler) GetDocumentMetaData(w http.ResponseWriter, req *http.Request) {
// 	req.Body = http.MaxBytesReader(w, req.Body, int64(10<<20))
// 	err := req.ParseMultipartForm(int64(10 << 20))
// 	if err != nil {
// 		fmt.Println("error parsing multipartform")
// 		return
// 	}
//
// 	file, handler, err := req.FormFile("inputFile")
// 	defer file.Close()
// 	file.Read()
// }

func (h *Handler) CreateDocumentMeta(w http.ResponseWriter, req *http.Request) {
}
