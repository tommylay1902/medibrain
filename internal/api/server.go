package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	Server *http.Server
	mux    *Mux
}

func NewServer(port string, mux *Mux) *Server {
	return &Server{
		Server: &http.Server{
			Addr:    port,
			Handler: mux.Mux,
		},
		mux: mux,
	}
}

func (s *Server) StartServer() {
	serverErrors := make(chan error, 1)
	go func() {
		log.Println("http server listening on port 8080")
		if err := s.Server.ListenAndServe(); err != nil {
			serverErrors <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("Fatal server error: %v", err)
	case <-stop:
		log.Println("Shutdown signal received, shutting down server gracefully...")

		// Create a context with a timeout for the shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Attempt graceful shutdown
		if err := s.Server.Shutdown(ctx); err != nil {
			log.Fatalf("Server shutdown failed: %v", err)
		}
		log.Println("Server gracefully stopped")
	}
}
