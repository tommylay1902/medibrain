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
	// if err != nil {
	// 	http.Error(w, fmt.Sprintf("Error parsing form: %v", err), http.StatusBadRequest)
	// 	return
	// }
	// file, header, err := req.FormFile("file") // "file" is the form field name
	// if err != nil {
	// 	http.Error(w, fmt.Sprintf("Error retrieving file: %v", err), http.StatusBadRequest)
	// 	return
	// }
	// defer file.Close()
	err := dph.service.UploadDocumentPipeline(req)
	// err dph.service.UploadDocumentPipeline(req.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error uploading document: %v", err), http.StatusInternalServerError)
		return
	}
	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte("File uploaded successfully"))
}
