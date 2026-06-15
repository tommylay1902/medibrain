package textstream

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type TextStreamHandler struct{}

func NewTextStreamHandler() *TextStreamHandler {
	return &TextStreamHandler{}
}

func (tsh *TextStreamHandler) StreamLLMResponse(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	rc := http.NewResponseController(w)

	ctx := req.Context()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	log.Println("Client connected to SSE stream")

	for {
		select {
		case <-ctx.Done():
			log.Println("Client disconnected from SSE stream")
			return

		case t := <-ticker.C:
			eventData := fmt.Sprintf("The current server time is: %s", t.Format(time.RFC3339))
			_, err := fmt.Fprintf(w, "data: %s\n\n", eventData)
			if err != nil {
				log.Printf("Error writing to stream: %v\n", err)
				return
			}

			if err := rc.Flush(); err != nil {
				log.Printf("Error flushing stream: %v\n", err)
				return
			}
		}
	}
}
