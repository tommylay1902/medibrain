package document

import "net/http"

type Route struct {
	Mux *http.ServeMux
}

func NewRoutes(dps *DocumentPipelineService) *Route {
	mux := http.NewServeMux()
	handler := NewHandler(dps)

	mux.HandleFunc("POST /upload", handler.UploadDocumentPipeline)
	mux.HandleFunc("POST /upload-with-edit", handler.UploadDocumentPipelineWithEdit)
	// mux.HandleFunc("POST /upload-chunks", handler.UploadChunks)
	mux.HandleFunc("POST /query", handler.GetSearchQuery)
	mux.HandleFunc("POST /chunk", handler.ChunkAndUploadText)
	route := &Route{
		Mux: mux,
	}

	return route
}
