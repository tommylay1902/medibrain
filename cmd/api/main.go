package main

import (
	"log/slog"
	"os"

	"github.com/hibiken/asynq"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tommylay1902/medibrain/internal/api"
	"github.com/tommylay1902/medibrain/internal/api/domain/document"
	"github.com/tommylay1902/medibrain/internal/api/domain/job"
	"github.com/tommylay1902/medibrain/internal/api/domain/metadata"
	"github.com/tommylay1902/medibrain/internal/api/domain/note"
	"github.com/tommylay1902/medibrain/internal/client/rag"
	seaweedclient "github.com/tommylay1902/medibrain/internal/client/seaweed"
	"github.com/tommylay1902/medibrain/internal/client/stirling"
	"github.com/tommylay1902/medibrain/internal/database"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	slog.SetDefault(logger)

	llm, err := ollama.New(
		ollama.WithModel("llama3.2:1b"),
		ollama.WithServerURL("http://ollama:11434"),
	)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
	rag := rag.NewRag(llm)
	db := database.NewDB()
	uowFactory := database.NewUnitOfWorkFactory(db)
	dmr := metadata.NewRepo(db)
	nr := note.NewNoteRepo(uowFactory)

	dms := metadata.NewService(dmr)
	ns := note.NewNoteService(nr, uowFactory, rag)
	sc := stirling.NewClient()
	swc := seaweedclient.NewClient()

	// start up asynq service
	redisOpt := asynq.RedisClientOpt{Addr: "redis-cache:6379"}
	asynqClient := asynq.NewClient(redisOpt)
	asynqInspector := asynq.NewInspector(redisOpt)

	defer asynqClient.Close()
	defer asynqInspector.Close()
	js := job.NewJobService(asynqInspector)

	dps := document.NewService(dmr, swc, sc, dms, rag, asynqClient)

	mux := api.NewMux(dms, dps, ns, js)

	server := api.NewServer(":8080", mux)
	server.StartServer()
}
