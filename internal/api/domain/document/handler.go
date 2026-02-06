package document

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/ledongthuc/pdf"
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
	var updateDM metadata.DocumentMeta

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

func (dph *DocumentPipelineHandler) ChunkAndUploadText(w http.ResponseWriter, req *http.Request) {
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
	textBody, err := dph.service.stirlingClient.GetTextFromPdf(pdfBytes, header, apiKey)
	if err != nil || textBody == nil {
		http.Error(w, fmt.Sprintf("internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	dm, err := dph.service.stirlingClient.GetMetaData(pdfBytes, header, apiKey)

	fid := req.FormValue("fid")
	if err = dph.service.ragClient.StoreDocument(*textBody, fid, dm.Title, dm.CreationDate, dm.ModificationDate); err != nil {
		fmt.Println(dm.Title)
		http.Error(w, fmt.Sprintf("internal server error: %v", err), http.StatusInternalServerError)
		return
	}
}

type SearchBody struct {
	Search string `json:"search"`
}

func (dph *DocumentPipelineHandler) GetSearchQuery(w http.ResponseWriter, req *http.Request) {
	var searchBody SearchBody
	if err := json.NewDecoder(req.Body).Decode(&searchBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	results := dph.service.ragClient.GetChunksByQuery(searchBody.Search)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func extractTextFromPDF(pdfData []byte) (string, error) {
	text, err := extractTextWithGoPDF(pdfData)
	if err == nil && len(strings.TrimSpace(text)) > 0 {
		return text, nil
	}

	return text, nil
}

func extractTextWithGoPDF(pdfData []byte) (string, error) {
	reader := bytes.NewReader(pdfData)
	pdfReader, err := pdf.NewReader(reader, int64(reader.Len()))
	if err != nil {
		return "", err
	}

	var textBuilder strings.Builder
	numPages := pdfReader.NumPage()

	for i := range numPages {
		page := pdfReader.Page(i + 1)

		content, err := page.GetPlainText(nil)
		if err != nil {
			continue
		}

		textBuilder.WriteString(content)
		textBuilder.WriteString("\n")
	}

	return textBuilder.String(), nil
}

func isLikelyBase64(s string) bool {
	if len(s) < 4 {
		return false
	}

	validBase64Chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
	base64Count := 0
	totalChars := 0

	for _, r := range s {
		if r == '\n' || r == '\r' || r == ' ' {
			continue
		}
		totalChars++
		if strings.ContainsRune(validBase64Chars, r) {
			base64Count++
		}
	}

	if totalChars == 0 {
		return false
	}

	return float64(base64Count)/float64(totalChars) > 0.95
}

func isValidUTF8(data []byte) bool {
	return utf8.Valid(data)
}
