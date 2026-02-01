package documentpipeline

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tommylay1902/medibrain/internal/api/domain/documentmeta"
)

type DocumentPipelineHandler struct {
	service *DocumentPipelineService
}

func NewHandler(service *DocumentPipelineService) *DocumentPipelineHandler {
	return &DocumentPipelineHandler{
		service: service,
	}
}

// TODO: return json error responses instead of just text
func (dph *DocumentPipelineHandler) UploadDocumentPipelineWithEdit(w http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(2 << 20)
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing multipart form: %v", err), http.StatusInternalServerError)
		return
	}

	defer func() {
		if req.MultipartForm != nil {
			req.MultipartForm.RemoveAll()
		}
	}()

	file, header, err := req.FormFile("fileInput")
	if err != nil {
		http.Error(w, fmt.Sprintf("internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	defer file.Close()
	maxSize := int64(2 << 20)
	if header.Size > maxSize {
		http.Error(w, fmt.Sprintf("file is too large: ~%.2f MB (max allowed: %.2f MB)", float64(header.Size)/(1024*1024), float64(maxSize)/(1024*1024)), http.StatusBadRequest)
		return
	}

	pdfBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	apiKey := req.Header.Get("X-API-KEY")
	var updateDM documentmeta.DocumentMeta

	metadataJSON := req.FormValue("metadata")

	if metadataJSON != "" {
		err := json.Unmarshal([]byte(metadataJSON), &updateDM)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid metadata JSON: %v", err), http.StatusBadRequest)
			return
		}
	}

	dm, err := dph.service.UploadDocumentPipelineWithEdit(pdfBytes, header, apiKey, &updateDM)
	if err != nil {
		http.Error(w, fmt.Sprintf("internal server err: %v", err), http.StatusInternalServerError)
	}

	w.WriteHeader(201)
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(dm)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error uploading document: %v", err), http.StatusInternalServerError)
		return
	}
}

func (dph *DocumentPipelineHandler) UploadDocumentPipeline(w http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing multipart form: %v", err), http.StatusInternalServerError)
		return
	}

	defer func() {
		if req.MultipartForm != nil {
			req.MultipartForm.RemoveAll()
		}
	}()
	file, header, err := req.FormFile("fileInput")
	if err != nil {
		http.Error(w, fmt.Sprintf("internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	defer file.Close()
	maxSize := int64(2 << 20)
	if header.Size > maxSize {
		http.Error(w, fmt.Sprintf("file is too large: ~%.2f MB (max allowed: %.2f MB)", float64(header.Size)/(1024*1024), float64(maxSize)/(1024*1024)), http.StatusBadRequest)
		return
	}

	pdfBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	apiKey := req.Header.Get("X-API-KEY")
	response, err := dph.service.UploadDocumentPipeline(pdfBytes, header, apiKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error uploading document: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(201)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error uploading documument: %v", err), http.StatusInternalServerError)
		return
	}
}
