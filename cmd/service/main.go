package main

import (
	"github.com/tommylay1902/medibrain/api"
)

func main() {
	mux := api.NewMux()
	server := api.NewServer(":8080", mux)
	server.StartServer()
}
