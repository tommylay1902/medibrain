package documentmeta

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Handler struct {
	service *DocumentMetaService
}

func NewHandler(service *DocumentMetaService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(w http.ResponseWriter, req *http.Request) {
	pdfs, err := h.service.List()
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

func (h *Handler) CreateDocumentMeta(w http.ResponseWriter, req *http.Request) {
}
