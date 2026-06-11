// could be used for scaling purposes later
package main

import (
	"log/slog"

	"github.com/hibiken/asynq"
	"github.com/tommylay1902/medibrain/internal/client/rag"
	"github.com/tommylay1902/medibrain/internal/task"
)

func main() {
	rag := rag.NewRag()
	handler := &task.IngestHandler{
		Rag: rag,
	}

	redisOpt := asynq.RedisClientOpt{Addr: "localhost:6379"}
	srv := asynq.NewServer(redisOpt, asynq.Config{Concurrency: 4})
	mux := asynq.NewServeMux()
	mux.HandleFunc(task.TypeIngestDocument, handler.Handle)

	if err := srv.Run(mux); err != nil {
		slog.Error("asynq worker stopped", "err", err)
	}
	slog.Info("succesfully started asynq worker server")
}
