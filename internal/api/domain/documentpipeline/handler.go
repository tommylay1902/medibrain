package documentpipeline

import (
	"fmt"
	"net/http"
)

type DocumentPipelineHandler struct {
	service *DocumentPipelineService
}

func NewHandler(service *DocumentPipelineService) *DocumentPipelineHandler {
	return &DocumentPipelineHandler{
		service: service,
	}
}

func (dph *DocumentPipelineHandler) UploadDocumentPipeline(w http.ResponseWriter, req *http.Request) {
	err := dph.service.UploadDocumentPipeline(req)
	// err dph.service.UploadDocumentPipeline(req.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error uploading document: %v", err), http.StatusInternalServerError)
		return
	}

	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte("File uploaded successfully"))
}
