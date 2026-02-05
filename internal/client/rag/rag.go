package rag

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/qdrant/go-client/qdrant"
	"github.com/tmc/langchaingo/llms/huggingface"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
)

type embedder struct {
	llm *huggingface.LLM
}

func newEmbedder() *embedder {
	// TODO: need to fix the pathing for loading env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error loading .env")
		wd, _ := os.Getwd()
		fmt.Printf("Current working directory: %s\n", wd)
		panic(err)
	}
	llm, err := huggingface.New(
		huggingface.WithModel("sentence-transformers/all-MiniLM-L6-v2"),
		huggingface.WithToken(os.Getenv("HF_TOKEN")),
		huggingface.WithURL("https://router.huggingface.co/hf-inference"),
	)
	if err != nil {
		fmt.Println("error getting llm client")
		panic(err)
	}

	return &embedder{
		llm: llm,
	}
}

func (e *embedder) GenerateEmbedding(ctx context.Context, texts []string) ([][]float32, error) {
	vectors, err := e.llm.CreateEmbedding(
		ctx,
		texts,
		"sentence-transformers/all-MiniLM-L6-v2/pipeline/feature-extraction",
		"",
	)
	if err != nil {
		return nil, err
	}

	return vectors, nil
}

type Rag struct {
	qClient  *qdrant.Client
	embedder *embedder
	splitter *textsplitter.RecursiveCharacter
}

func NewRag() *Rag {
	embedder := newEmbedder()
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
	})
	if err != nil {
		fmt.Println("error connecting to qdrant client")
		panic(err)
	}

	splitter := textsplitter.NewRecursiveCharacter(textsplitter.WithChunkSize(800), textsplitter.WithChunkOverlap(40))

	return &Rag{
		qClient:  client,
		embedder: embedder,
		splitter: &splitter,
	}
}

func (r *Rag) StoreDocument(doc string, fid string) []string {
	document := schema.Document{
		PageContent: doc,
	}
	chunks, _ := r.splitter.SplitText(document.PageContent)
	vec, err := r.embedder.GenerateEmbedding(context.Background(), chunks)
	if err != nil {
		fmt.Println("error embedding chunks")
		panic(err)
	}
	points := make([]*qdrant.PointStruct, 0, len(chunks))
	for i := range vec {
		payload := qdrant.NewValueMap(map[string]any{
			"fid":     fid,
			"content": chunks[i],
		})
		points = append(points, &qdrant.PointStruct{Id: qdrant.NewID(uuid.NewString()), Vectors: qdrant.NewVectors(vec[i]...), Payload: payload})
	}
	_, err = r.qClient.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: "documents",
		Points:         points,
	})
	if err != nil {
		fmt.Println(err)
	}

	return chunks
}

type Response struct {
	Content string `json:"content"`
	Fid     string `json:"fid"`
}

func (r *Rag) GetChunksByQuery(query string) []Response {
	chunks, err := r.embedder.GenerateEmbedding(context.Background(), []string{query})
	if err != nil {
		panic(err)
	}
	results, err := r.qClient.Query(context.Background(), &qdrant.QueryPoints{
		CollectionName: "documents",
		Query:          qdrant.NewQuery(chunks[0]...),
		WithPayload:    qdrant.NewWithPayload(true),
	})
	if err != nil {
		panic(err)
	}

	responses := make([]Response, 0, len(results))

	for _, result := range results {
		payload := result.Payload

		var r Response
		if fidValue, exists := payload["fid"]; exists && fidValue != nil {
			fid := fidValue.GetStringValue() // Use qdrant.Value methods
			r.Fid = fid
		}

		if contentValue, exists := payload["content"]; exists && contentValue != nil {
			content := contentValue.GetStringValue()
			r.Content = content
		}
		responses = append(responses, r)
	}

	return responses
}

// Access fields
func GenerateCollections(r *Rag) {
	exists, err := r.qClient.CollectionExists(context.Background(), "documents")
	if err != nil {
		fmt.Println("error checking for document collection")
		panic(err)
	}
	if exists {
		err := r.qClient.DeleteCollection(context.Background(), "documents")
		if err != nil {
			fmt.Println("error deleting document collection")
			panic(err)
		}
	}

	exists, err = r.qClient.CollectionExists(context.Background(), "audio_logs")
	if err != nil {
		fmt.Println("error checking for audio logs collection")
		panic(err)
	}
	if exists {
		err := r.qClient.DeleteCollection(context.Background(), "audio_logs")
		if err != nil {
			fmt.Println("error deleting audio logs collection")
			panic(err)
		}
	}

	exists, err = r.qClient.CollectionExists(context.Background(), "notes")
	if err != nil {
		fmt.Println("error checking for notes collection")
		panic(err)
	}
	if exists {
		err := r.qClient.DeleteCollection(context.Background(), "notes")
		if err != nil {
			fmt.Println("error deleting notes collection")
			panic(err)
		}
	}
	err = r.qClient.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: "documents",
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     384,
			Distance: qdrant.Distance_Cosine,
		}),
	})
	if err != nil {
		fmt.Println("error creating documents collection")
		panic(err)
	}

	err = r.qClient.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: "audio_logs",
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     384,
			Distance: qdrant.Distance_Cosine,
		}),
	})
	if err != nil {
		fmt.Println("error creating audio logs collection")
		panic(err)
	}

	err = r.qClient.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: "notes",
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     384,
			Distance: qdrant.Distance_Cosine,
		}),
	})
	if err != nil {
		fmt.Println("error creating notes collection")
		panic(err)
	}
}
