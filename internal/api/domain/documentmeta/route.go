package documentmeta

import (
	"net/http"

	"github.com/jmoiron/sqlx"
)

type Route struct {
	Mux *http.ServeMux
}

func NewRoutes(db *sqlx.DB) *Route {
	mux := http.NewServeMux()
	handler := NewHandler(db)

	mux.HandleFunc("GET /documentmeta", handler.List)

	route := &Route{
		Mux: mux,
	}

	return route
}
