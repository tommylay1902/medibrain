package metadata

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Handler struct {
	service *MetadataService
}

func NewHandler(service *MetadataService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(w http.ResponseWriter, req *http.Request) {
	pdfs, err := h.service.List()
	if err != nil {
		fmt.Println(err)
		fmt.Println("error getting pdfs list")
		w.WriteHeader(500)
		return
	}
	result, err := json.Marshal(pdfs)
	if err != nil {
		fmt.Println("error marshaling object")
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

func (h *Handler) CreateMetadata(w http.ResponseWriter, req *http.Request) {
}
