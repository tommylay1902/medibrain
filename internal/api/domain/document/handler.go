package document

import (
	"bytes"
	"encoding/base64"
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

func (dph *DocumentPipelineHandler) UploadChunks(w http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(10 << 20) // 32MB
	if err != nil {
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	var allText strings.Builder

	for _, files := range req.MultipartForm.File {
		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				fmt.Printf("Error opening file %s: %v\n", fileHeader.Filename, err)
				continue
			}

			content, err := io.ReadAll(file)
			file.Close()
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", fileHeader.Filename, err)
				continue
			}

			fmt.Printf("Processing file: %s (size: %d bytes)\n", fileHeader.Filename, len(content))

			var pdfData []byte
			if isLikelyBase64(string(content)) {
				fmt.Println("Detected base64 encoded file")
				clean := strings.Map(func(r rune) rune {
					if r == '\n' || r == '\r' || r == ' ' || r == '\t' {
						return -1
					}
					return r
				}, string(content))

				if len(clean)%4 != 0 {
					clean += strings.Repeat("=", 4-len(clean)%4)
				}

				decoded, err := base64.StdEncoding.DecodeString(clean)
				if err != nil {
					fmt.Printf("Failed to decode base64: %v\n", err)
					decoded, err = base64.URLEncoding.DecodeString(clean)
					if err != nil {
						fmt.Printf("Failed to decode base64 URL: %v\n", err)
						continue
					}
				}
				pdfData = decoded
			} else {
				pdfData = content
			}

			if len(pdfData) > 4 && string(pdfData[:4]) == "%PDF" {
				fmt.Printf("Extracting text from PDF: %s\n", fileHeader.Filename)
				text, err := extractTextFromPDF(pdfData)
				if err != nil {
					fmt.Printf("Error extracting PDF text: %v\n", err)
					continue
				}

				allText.WriteString(fmt.Sprintf("=== File: %s ===\n", fileHeader.Filename))
				allText.WriteString(text)
				allText.WriteString("\n---\n")
			} else {
				fmt.Printf("File %s is not a PDF (or base64 decoding failed)\n", fileHeader.Filename)
			}
		}
	}

	for key, values := range req.MultipartForm.Value {
		for _, value := range values {
			if len(value) > 100 {
				fmt.Printf("Processing field: %s (length: %d)\n", key, len(value))

				clean := strings.TrimSpace(value)
				if strings.HasPrefix(clean, "data:") {
					if idx := strings.Index(clean, "base64,"); idx != -1 {
						clean = clean[idx+7:]
					}
				}

				if len(clean)%4 != 0 {
					clean += strings.Repeat("=", 4-len(clean)%4)
				}

				decoded, err := base64.StdEncoding.DecodeString(clean)
				if err != nil {
					decoded, err = base64.URLEncoding.DecodeString(clean)
					if err != nil {
						continue
					}
				}

				if len(decoded) > 4 && string(decoded[:4]) == "%PDF" {
					fmt.Printf("Extracting text from base64 PDF in field: %s\n", key)
					text, err := extractTextFromPDF(decoded)
					if err != nil {
						fmt.Printf("Error extracting PDF text: %v\n", err)
						continue
					}

					allText.WriteString(fmt.Sprintf("=== Field: %s ===\n", key))
					allText.WriteString(text)
					allText.WriteString("\n---\n")
				}
			}
		}
	}
	fmt.Printf("Extracted text length: %d characters\n", allText.Len())

	if allText.Len() > 0 {
		fid := req.FormValue("fid")
		chunks := dph.service.ragClient.StoreDocument(allText.String(), fid)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Created %d chunks", len(chunks))))
	} else {
		http.Error(w, "No text content found in PDFs", http.StatusBadRequest)
	}
}

func (dph *DocumentPipelineHandler) GetSearchQuery(w http.ResponseWriter, req *http.Request) {
	results := dph.service.ragClient.GetChunksByQuery("with progression with femoral lesion though does not appearstructurally at risk")
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Extract text from PDF bytes
func extractTextFromPDF(pdfData []byte) (string, error) {
	// Create a reader from the byte slice

	// Method 1: Using github.com/ledongthuc/pdf (simpler)
	text, err := extractTextWithGoPDF(pdfData)
	if err == nil && len(strings.TrimSpace(text)) > 0 {
		return text, nil
	}

	return text, nil
}

// Simple PDF text extraction using github.com/ledongthuc/pdf
func extractTextWithGoPDF(pdfData []byte) (string, error) {
	reader := bytes.NewReader(pdfData)
	pdfReader, err := pdf.NewReader(reader, int64(reader.Len()))
	if err != nil {
		return "", err
	}

	var textBuilder strings.Builder
	numPages := pdfReader.NumPage()
	if err != nil {
		return "", err
	}

	for i := 0; i < numPages; i++ {
		page := pdfReader.Page(i + 1)
		// if page == nil {
		// 	continue
		// }

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
