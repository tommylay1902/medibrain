package documentpipeline

import "net/http"

type Route struct {
	Mux *http.ServeMux
}

func NewRoutes(dps *DocumentPipelineService) *Route {
	mux := http.NewServeMux()
	handler := NewHandler(dps)

	mux.HandleFunc("POST /upload", handler.UploadDocumentPipeline)
	// mux.HandleFunc("POST /documentmeta", handler.GetDocumentMetaData)
	route := &Route{
		Mux: mux,
	}

	return route
}
