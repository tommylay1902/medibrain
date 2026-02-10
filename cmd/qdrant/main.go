package main

import (
	"github.com/tommylay1902/medibrain/internal/client/rag"
)

func main() {
	qdrant := rag.NewRag()
	rag.GenerateCollections(qdrant)
}
