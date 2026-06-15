package job

import (
	"net/http"
)

type Route struct {
	Mux *http.ServeMux
}

func NewJobRoutes(js *JobService) *Route {
	mux := http.NewServeMux()
	handler := NewJobHandler(js)
	mux.HandleFunc("GET /{id}", handler.GetJobStatus)
	route := &Route{
		Mux: mux,
	}

	return route
}
