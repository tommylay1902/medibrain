package document

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

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
	w.WriteHeader(201)

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

	w.WriteHeader(201)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Error("Error uploading document", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("Error uploading documument: %v", err), http.StatusInternalServerError)
		return
	}
}

// TODO: need to revise this looks like service logic is leaking into handler logic
func (dph *DocumentPipelineHandler) ChunkAndUploadText(w http.ResponseWriter, req *http.Request) {
	slog.Info("entering chunking")
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
		slog.Error("err parsing fileInput", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	defer file.Close()
	maxSize := int64(2 << 20)
	if header.Size > maxSize {
		slog.Error("file is too large")
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

	fid := req.FormValue("fid")
	if strings.TrimSpace(fid) == "" {
		slog.Error("need fid")
		http.Error(w, fmt.Sprintf("Need fid: %v", err), http.StatusBadRequest)
		return
	}

	jobId, err := dph.service.ChunkAndUploadText(pdfBytes, header, apiKey, fid)
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
	var searchBody SearchBody
	if err := json.NewDecoder(req.Body).Decode(&searchBody); err != nil {
		slog.Error("request body parsing err", slog.Any("err", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	results := dph.service.ragClient.GetChunksByQuery(searchBody.Search)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		slog.Error("err getting chunk by query", slog.Any("err", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// func extractTextWithGoPDF(pdfData []byte) (string, error) {
// 	reader := bytes.NewReader(pdfData)
// 	pdfReader, err := pdf.NewReader(reader, int64(reader.Len()))
// 	if err != nil {
// 		return "", err
// 	}
//
// 	var textBuilder strings.Builder
// 	numPages := pdfReader.NumPage()
//
// 	for i := range numPages {
// 		page := pdfReader.Page(i + 1)
//
// 		content, err := page.GetPlainText(nil)
// 		if err != nil {
// 			continue
// 		}
//
// 		textBuilder.WriteString(content)
// 		textBuilder.WriteString("\n")
// 	}
//
// 	return textBuilder.String(), nil
// }

// func extractTextFromPDF(pdfData []byte) (string, error) {
// 	text, err := extractTextWithGoPDF(pdfData)
// 	if err == nil && len(strings.TrimSpace(text)) > 0 {
// 		return text, nil
// 	}
//
// 	return text, nil
// }

// func isLikelyBase64(s string) bool {
// 	if len(s) < 4 {
// 		return false
// 	}
//
// 	validBase64Chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
// 	base64Count := 0
// 	totalChars := 0
//
// 	for _, r := range s {
// 		if r == '\n' || r == '\r' || r == ' ' {
// 			continue
// 		}
// 		totalChars++
// 		if strings.ContainsRune(validBase64Chars, r) {
// 			base64Count++
// 		}
// 	}
//
// 	if totalChars == 0 {
// 		return false
// 	}
//
// 	return float64(base64Count)/float64(totalChars) > 0.95
// }

// func isValidUTF8(data []byte) bool {
// 	return utf8.Valid(data)
// }
