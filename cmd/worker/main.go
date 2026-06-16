// could be used for scaling purposes later
package main

import (
	"log/slog"

	"github.com/hibiken/asynq"
	"github.com/tmc/langchaingo/llms/ollama"
	client "github.com/tommylay1902/medibrain/internal/client/pydocument"
	"github.com/tommylay1902/medibrain/internal/client/rag"
	"github.com/tommylay1902/medibrain/internal/task"
)

func main() {
	llm, err := ollama.New(
		ollama.WithModel("medgemma:4b"),
		ollama.WithServerURL("http://ollama:11434"),
	)
	if err != nil {

		slog.Error(err.Error())
		panic(err)
	}

	pydc := client.NewPydocument()
	rag := rag.NewRag(llm, pydc)
	handler := &task.IngestHandler{
		Rag: rag,
	}

	redisOpt := asynq.RedisClientOpt{Addr: "redis-cache:6379"}
	srv := asynq.NewServer(redisOpt, asynq.Config{Concurrency: 4, Queues: map[string]int{
		"ingest": 4,
	}})
	mux := asynq.NewServeMux()
	mux.HandleFunc(task.TypeIngestDocument, handler.Handle)

	if err := srv.Run(mux); err != nil {
		slog.Error("asynq worker stopped", "err", err)
	}
	slog.Info("succesfully started asynq worker server")
}
