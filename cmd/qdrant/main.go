package main

import (
	"log/slog"

	"github.com/tmc/langchaingo/llms/ollama"
	client "github.com/tommylay1902/medibrain/internal/client/pydocument"
	"github.com/tommylay1902/medibrain/internal/client/rag"
)

func main() {
	llm, err := ollama.New(
		ollama.WithModel("llama3.2:9b"),
		ollama.WithServerURL("http://ollama:11434"),
	)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	pydc := client.NewPydocument()
	qdrant := rag.NewRag(llm, pydc)
	rag.GenerateCollections(qdrant)
}
