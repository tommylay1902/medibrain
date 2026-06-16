package document

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/tommylay1902/medibrain/internal/api/domain/metadata"
)

type DocumentPipelineHandler struct {
	service *DocumentPipelineService
}

func NewHandler(service *DocumentPipelineService) *DocumentPipelineHandler {
	return &DocumentPipelineHandler{
		service: service,
	}
}

func (dph *DocumentPipelineHandler) UploadDocumentTest(w http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(2 << 20)
	if err != nil {
		slog.Error("error parsing multipart form", slog.Any("err", err))
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
		slog.Error("error getting fileInput", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	defer file.Close()
	maxSize := int64(2 << 20)
	if header.Size > maxSize {
		slog.Error("file size too large")
		http.Error(w, fmt.Sprintf("file is too large: ~%.2f MB (max allowed: %.2f MB)", float64(header.Size)/(1024*1024), float64(maxSize)/(1024*1024)), http.StatusBadRequest)
		return
	}

	pdfBytes, err := io.ReadAll(file)
	if err != nil {
		slog.Error("error reading file", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	apiKey := req.Header.Get("X-API-KEY")
	var updateDM metadata.Metadata

	metadataJSON := req.FormValue("metadata")

	if metadataJSON != "" {
		err = json.Unmarshal([]byte(metadataJSON), &updateDM)
		if err != nil {
			slog.Error("invalid metadata JSON", slog.Any("err", err))
			http.Error(w, fmt.Sprintf("invalid metadata JSON: %v", err), http.StatusBadRequest)
			return
		}
	}

	dm, err := dph.service.UploadDocumentPipelineWithEdit(pdfBytes, header, apiKey, &updateDM)
	if err != nil {
		slog.Error("error with upload document service", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("internal server err: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)

	err = json.NewEncoder(w).Encode(dm)
	if err != nil {
		slog.Error("error uploading document", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("Error uploading document: %v", err), http.StatusInternalServerError)
		return
	}

	jobId, err := dph.service.ChunkAndUploadText(pdfBytes, header, apiKey, dm.PdfFid)
	if err != nil {
		slog.Error("error chunking", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(202)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"jobId": jobId,
	}); err != nil {
		slog.Error("Error writing jobId json result", slog.Any("err", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// TODO: return json error responses instead of just text
func (dph *DocumentPipelineHandler) UploadDocumentPipelineWithEdit(w http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(2 << 20)
	if err != nil {
		slog.Error("error parsing multipart form", slog.Any("err", err))
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
		slog.Error("error getting fileInput", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	defer file.Close()
	maxSize := int64(2 << 20)
	if header.Size > maxSize {
		slog.Error("file size too large")
		http.Error(w, fmt.Sprintf("file is too large: ~%.2f MB (max allowed: %.2f MB)", float64(header.Size)/(1024*1024), float64(maxSize)/(1024*1024)), http.StatusBadRequest)
		return
	}

	pdfBytes, err := io.ReadAll(file)
	if err != nil {
		slog.Error("error reading file", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	apiKey := req.Header.Get("X-API-KEY")
	var updateDM metadata.Metadata

	metadataJSON := req.FormValue("metadata")

	if metadataJSON != "" {
		err = json.Unmarshal([]byte(metadataJSON), &updateDM)
		if err != nil {
			slog.Error("invalid metadata JSON", slog.Any("err", err))
			http.Error(w, fmt.Sprintf("invalid metadata JSON: %v", err), http.StatusBadRequest)
			return
		}
	}

	dm, err := dph.service.UploadDocumentPipelineWithEdit(pdfBytes, header, apiKey, &updateDM)
	if err != nil {
		slog.Error("error with upload document service", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("internal server err: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(202)

	err = json.NewEncoder(w).Encode(dm)
	if err != nil {
		slog.Error("error uploading document", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("Error uploading document: %v", err), http.StatusInternalServerError)
		return
	}
}

func (dph *DocumentPipelineHandler) UploadDocumentPipeline(w http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(10 << 20)
	if err != nil {
		slog.Error("error parsing multipart form", slog.Any("err", err))
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
		slog.Error("error parsing fileInput", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	defer file.Close()
	maxSize := int64(2 << 20)
	if header.Size > maxSize {
		slog.Error("file size is too large")
		http.Error(w, fmt.Sprintf("file is too large: ~%.2f MB (max allowed: %.2f MB)", float64(header.Size)/(1024*1024), float64(maxSize)/(1024*1024)), http.StatusBadRequest)
		return
	}

	pdfBytes, err := io.ReadAll(file)
	if err != nil {
		slog.Error("error reading file", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	apiKey := req.Header.Get("X-API-KEY")
	response, err := dph.service.UploadDocumentPipeline(pdfBytes, header, apiKey)
	if err != nil {
		slog.Error("Error with upload document service", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("Error uploading document: %v", err), http.StatusInternalServerError)
		return
	}

	jobId, err := dph.service.ChunkAndUploadText(pdfBytes, header, apiKey, response.PdfFid)
	if err != nil {
		slog.Error("error chunking", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(202)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"jobId": jobId,
	}); err != nil {
		slog.Error("Error writing jobId json result", slog.Any("err", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type SearchBody struct {
	Search string `json:"search"`
}

func (dph *DocumentPipelineHandler) GetSearchQuery(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	search := req.URL.Query().Get("search")
	rc := http.NewResponseController(w)

	ctx := req.Context()

	_, prompt := dph.service.ragClient.GetChunksByQuery(search)
	err := dph.service.ragClient.StreamResponse(ctx, prompt, w, rc)
	if err != nil {
		slog.Error("stream error", slog.Any("err", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
