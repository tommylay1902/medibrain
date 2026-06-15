package main

import (
	"log/slog"

	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tommylay1902/medibrain/internal/client/rag"
)

func main() {
	llm, err := ollama.New(
		ollama.WithModel("llama3.2:1b"),
		ollama.WithServerURL("http://ollama:11434"),
	)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
	qdrant := rag.NewRag(llm)
	rag.GenerateCollections(qdrant)
}
