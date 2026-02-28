package rag

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	splitter *textsplitter.RecursiveCharacter
}

func NewRag() *Rag {
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
		splitter: &splitter,
	}
}

func (r *Rag) StoreDocument(doc string, fid string, title *string, uploadDate *string, creationDate *string, keywords string) error {
	docTitle := ""
	if title != nil {
		docTitle = *title
	}
	docUploadDate := ""
	if uploadDate != nil {
		docUploadDate = *uploadDate
	}

	docCreationDate := ""
	if creationDate != nil {
		docCreationDate = *creationDate
	}
	document := schema.Document{
		PageContent: doc,
	}
	chunks, _ := r.splitter.SplitText(document.PageContent)
	points := make([]*qdrant.PointStruct, 0, len(chunks))
	for _, chunk := range chunks {
		vec, err := getEmbedding(chunk)
		if err != nil {
			continue
		}
		payload := qdrant.NewValueMap(map[string]any{
			"fid":          fid,
			"title":        docTitle,
			"content":      chunk,
			"uploadDate":   docUploadDate,
			"creationDate": docCreationDate,
			"keywords":     keywords,
		})
		points = append(points, &qdrant.PointStruct{Id: qdrant.NewID(uuid.NewString()), Vectors: qdrant.NewVectors(vec...), Payload: payload})
	}
	_, err := r.qClient.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: "documents",
		Points:         points,
	})
	if err != nil {
		return err
	}

	return nil
}

type Response struct {
	Content  string `json:"content"`
	Fid      string `json:"fid"`
	Title    string `json:"title"`
	Keywords string `json:"keywords"`
}

func (r *Rag) GetChunksByQuery(query string) []Response {
	vec, err := getEmbedding(query)
	if err != nil {
		panic(err)
	}
	results, err := r.qClient.Query(context.Background(), &qdrant.QueryPoints{
		CollectionName: "documents",
		Query:          qdrant.NewQuery(vec...),
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

		if titleValue, exists := payload["title"]; exists && titleValue != nil {
			title := titleValue.GetStringValue()
			r.Title = title
		}

		if keywordsValue, exists := payload["keywords"]; exists && keywordsValue != nil {
			keywords := keywordsValue.GetStringValue()
			r.Keywords = keywords
		}

		responses = append(responses, r)
	}

	return responses
}

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

type EmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type EmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

func getEmbedding(text string) ([]float32, error) {
	reqBody := EmbeddingRequest{Model: "all-minilm:l6-v2", Prompt: text}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(
		"http://localhost:11434/api/embeddings",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call Ollama API: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error hit in calling")
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("error hit in calling ollama")
		return nil, fmt.Errorf("Ollama API error (%d): %s", resp.StatusCode, string(body))
	}

	// 4. Parse the embedding
	var embeddingResp EmbeddingResponse
	if err := json.Unmarshal(body, &embeddingResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return embeddingResp.Embedding, nil
}
