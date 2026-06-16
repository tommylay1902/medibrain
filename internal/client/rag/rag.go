// Package rag interacts with all the clients that are involved with the RAG process
package rag

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	client "github.com/tommylay1902/medibrain/internal/client/pydocument"
)

type Rag struct {
	llm        *ollama.LLM
	qClient    *qdrant.Client
	splitter   *textsplitter.RecursiveCharacter
	pydocument *client.Pydocument
}

func NewRag(llm *ollama.LLM, pydocument *client.Pydocument) *Rag {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "qdrant",
		Port: 6334,
	})
	if err != nil {
		slog.Error("error connecting to qdrant client", slog.Any("err", err))
		panic(err)
	}

	splitter := textsplitter.NewRecursiveCharacter(textsplitter.WithChunkSize(1000), textsplitter.WithChunkOverlap(200))

	return &Rag{
		llm:        llm,
		qClient:    client,
		splitter:   &splitter,
		pydocument: pydocument,
	}
}

func (r *Rag) StoreDocument(pdfBytes []byte, header *multipart.FileHeader, fid string, title *string, uploadDate *string, creationDate *string, keywords string) error {
	slog.Info(header.Filename)
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

	splitResponse, err := r.pydocument.GetDocumentSplits(pdfBytes, header)
	if err != nil {
		slog.Error("error getting document splits", slog.Any("err", err))
		return err
	}

	points := make([]*qdrant.PointStruct, 0, len(splitResponse.Chunks))
	for _, chunk := range splitResponse.Chunks {
		vec, err := getEmbedding(chunk.Text)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		payload := qdrant.NewValueMap(map[string]any{
			"fid":          fid,
			"title":        docTitle,
			"page":         chunk.Page,
			"filename":     header.Filename,
			"content":      chunk.Text,
			"uploadDate":   docUploadDate,
			"creationDate": docCreationDate,
			"keywords":     keywords,
		})

		points = append(points, &qdrant.PointStruct{
			Id: qdrant.NewID(uuid.NewString()),
			Vectors: qdrant.NewVectorsMap(map[string]*qdrant.Vector{
				"dense":  qdrant.NewVector(vec...),
				"sparse": qdrant.NewVectorDocument(&qdrant.Document{Text: chunk.Text, Model: "qdrant/bm25"}),
			}),
			Payload: payload,
		})
	}

	_, err = r.qClient.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: "documents",
		Points:         points,
	})
	if err != nil {
		slog.Error("error upserting chunk into qdrant", slog.Any("err", err))
		return err
	}
	return nil
}

func (r *Rag) StoreNote(id *uuid.UUID, content string, title string) error {
	document := schema.Document{
		PageContent: content,
	}
	chunks, _ := r.splitter.SplitText(document.PageContent)
	points := make([]*qdrant.PointStruct, 0, len(chunks))
	for _, chunk := range chunks {
		vec, err := getEmbedding(chunk)
		if err != nil {
			slog.Error("error embedding chunk")
			continue
		}
		payload := qdrant.NewValueMap(map[string]any{
			"id":      id.String(),
			"content": chunk,
			"title":   title,
		})

		points = append(points, &qdrant.PointStruct{Id: qdrant.NewID(uuid.NewString()), Vectors: qdrant.NewVectors(vec...), Payload: payload})
	}
	_, err := r.qClient.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: "notes",
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
	Page     int64  `json:"page"`
	Filename string `json:"filename"`
}

func (r *Rag) GetChunksByQuery(query string) ([]Response, string) {
	vec, err := getEmbedding(query)
	if err != nil {
		panic(err)
	}
	results, err := r.qClient.Query(context.Background(), &qdrant.QueryPoints{
		CollectionName: "documents",
		Prefetch: []*qdrant.PrefetchQuery{
			{
				Query: qdrant.NewQuery(vec...),
				Using: qdrant.PtrOf("dense"),
				Limit: qdrant.PtrOf(uint64(20)),
			},
			{
				Query: qdrant.NewQueryDocument(&qdrant.Document{
					Text:  query,
					Model: "qdrant/bm25",
				}),
				Using: qdrant.PtrOf("sparse"),
				Limit: qdrant.PtrOf(uint64(20)),
			},
		},
		Query:       qdrant.NewQueryFusion(qdrant.Fusion_RRF),
		WithPayload: qdrant.NewWithPayload(true),
		Limit:       qdrant.PtrOf(uint64(10)),
	})
	if err != nil {
		slog.Error("error getting chunk by query", slog.Any("err", err))
		panic(err)
	}

	mergedChunks := make([]string, 0, len(results))
	responses := make([]Response, 0, len(results))

	for _, result := range results {
		payload := result.Payload
		var r Response
		if fidValue, exists := payload["fid"]; exists && fidValue != nil {
			r.Fid = fidValue.GetStringValue()
		}
		if contentValue, exists := payload["content"]; exists && contentValue != nil {
			r.Content = contentValue.GetStringValue()
		}
		if titleValue, exists := payload["title"]; exists && titleValue != nil {
			r.Title = titleValue.GetStringValue()
		}
		if keywordsValue, exists := payload["keywords"]; exists && keywordsValue != nil {
			r.Keywords = keywordsValue.GetStringValue()
		}
		if filenameValue, exists := payload["filename"]; exists && filenameValue != nil {
			r.Filename = filenameValue.GetStringValue()
		}
		if pageValue, exists := payload["page"]; exists && pageValue != nil {
			r.Page = pageValue.GetIntegerValue()
		}

		slog.Info(fmt.Sprintf(
			"fid: %s, filename: %s, page: %d, content: %s",
			r.Fid, r.Filename, r.Page, r.Content,
		))

		mergedChunks = append(mergedChunks, fmt.Sprintf(
			"fid: %s, filename: %s, page: %d, content: %s",
			r.Fid, r.Filename, r.Page, r.Content,
		))
		responses = append(responses, r)
	}

	template := prompts.NewPromptTemplate(
		`You are a medical professional for question-answering tasks for someone who needs answers quickly
    (answer with concise and good summaries of the provided context).
    Use the following pieces of retrieved context to formulate your answers.
    Remember you are a medical professional so you can't give guesses as answers.
    If you can't find the answer within the context, just say you don't know.

    Each piece of context below is formatted as:
    fid: <document id>, filename: <source file>, page: <page number>, content: <text>

    Whenever you state a fact drawn from the context, immediately follow that sentence
		with a citation using the fid value, in this exact format: [fid:<document id>, filename: <filename>, page: <page number>]

		Example: "Aspirin can reduce inflammation [fid:abc123, filename: patient.pdf, page: 1]."

    Do not omit the citation for any claim taken from the context.

    QUESTION: {{.question}}
    CONTEXT: {{.context}}
    `,
		[]string{"question", "context"},
	)

	chunks := strings.Join(mergedChunks, "\n\n")
	formattedPrompt, _ := template.Format(map[string]any{
		"question": query,
		"context":  chunks,
	})

	return responses, formattedPrompt
}

func (r *Rag) StreamResponse(ctx context.Context, prompt string, w http.ResponseWriter, rc *http.ResponseController) error {
	_, err := llms.GenerateFromSinglePrompt(ctx, r.llm,
		prompt, llms.WithTemperature(0.0),
		llms.WithStreamingFunc(
			func(ctx context.Context, chunk []byte) error {
				_, writeErr := fmt.Fprintf(w, "data: %s\n\n", string(chunk))
				if writeErr != nil {
					return writeErr
				}

				return rc.Flush()
			}),
	)
	if err != nil {
		slog.Error(err.Error())
	}
	return nil
}

func GenerateCollections(r *Rag) {
	exists, err := r.qClient.CollectionExists(context.Background(), "documents")
	if err != nil {
		slog.Error("error checking for document collection", slog.Any("err", err))
		panic(err)
	}
	if exists {
		err = r.qClient.DeleteCollection(context.Background(), "documents")
		if err != nil {
			slog.Error("error deleting document collection", slog.Any("err", err))
			panic(err)
		}
	}

	exists, err = r.qClient.CollectionExists(context.Background(), "audio_logs")
	if err != nil {
		slog.Error("error checking for audio logs collection", slog.Any("err", err))
		panic(err)
	}
	if exists {
		err = r.qClient.DeleteCollection(context.Background(), "audio_logs")
		if err != nil {
			slog.Error("error deleting audio logs collection", slog.Any("err", err))
			panic(err)
		}
	}

	exists, err = r.qClient.CollectionExists(context.Background(), "notes")
	if err != nil {
		slog.Error("error checking for notes collection", slog.Any("err", err))
		panic(err)
	}
	if exists {
		err = r.qClient.DeleteCollection(context.Background(), "notes")
		if err != nil {
			slog.Error("error deleting notes collection", slog.Any("err", err))
			panic(err)
		}
	}
	err = r.qClient.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: "documents",
		VectorsConfig: qdrant.NewVectorsConfigMap(
			map[string]*qdrant.VectorParams{
				"dense": {
					Size:     1024,
					Distance: qdrant.Distance_Cosine,
				},
			},
		),
		SparseVectorsConfig: qdrant.NewSparseVectorsConfig(
			map[string]*qdrant.SparseVectorParams{
				"sparse": {
					Modifier: qdrant.Modifier_Idf.Enum(),
				},
			},
		),
	})
	if err != nil {
		slog.Error("error creating documents collection", slog.Any("err", err))
		panic(err)
	}

	err = r.qClient.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: "audio_logs",
		VectorsConfig: qdrant.NewVectorsConfigMap(
			map[string]*qdrant.VectorParams{
				"dense": {
					Size:     1024,
					Distance: qdrant.Distance_Cosine,
				},
			},
		),
		SparseVectorsConfig: qdrant.NewSparseVectorsConfig(
			map[string]*qdrant.SparseVectorParams{
				"sparse": {Modifier: qdrant.Modifier_Idf.Enum()},
			},
		),
	})
	if err != nil {
		fmt.Println("error creating audio logs collection")
		slog.Error("error creating audio logs collection", slog.Any("err", err))
		panic(err)
	}

	err = r.qClient.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: "notes",
		VectorsConfig: qdrant.NewVectorsConfigMap(
			map[string]*qdrant.VectorParams{
				"dense": {
					Size:     1024,
					Distance: qdrant.Distance_Cosine,
				},
			},
		),
		SparseVectorsConfig: qdrant.NewSparseVectorsConfig(
			map[string]*qdrant.SparseVectorParams{
				"sparse": {Modifier: qdrant.Modifier_Idf.Enum()},
			},
		),
	})
	if err != nil {
		slog.Error("error creating notes collection", slog.Any("err", err))
		panic(err)
	}
}

type Options struct {
	Dimensions int `json:"dimensions"`
}
type EmbeddingRequest struct {
	Model   string  `json:"model"`
	Prompt  string  `json:"prompt"`
	Options Options `json:"options"`
}

type EmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

func getEmbedding(text string) ([]float32, error) {
	reqBody := EmbeddingRequest{Model: "qwen3-embedding:0.6b", Options: Options{Dimensions: 1024}, Prompt: text}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		slog.Error("failed to marshal request", slog.Any("err", err))
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(
		"http://ollama:11434/api/embeddings",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		slog.Error("failed to call Ollama API", slog.Any("err", err))
		return nil, fmt.Errorf("failed to call Ollama API: %w", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("failed to read response", slog.Any("err", err))
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		slog.Error("error hit in calling ollama", slog.Any("err", string(body)))
		return nil, fmt.Errorf("ollama API error (%d): %s", resp.StatusCode, string(body))
	}

	var embeddingResp EmbeddingResponse
	if err := json.Unmarshal(body, &embeddingResp); err != nil {
		slog.Error("failed to parse response", slog.Any("err", err))
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	slog.Info(fmt.Sprintf("%d", (len(embeddingResp.Embedding))))
	return embeddingResp.Embedding, nil
}
