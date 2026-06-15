package task

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/hibiken/asynq"
	"github.com/tommylay1902/medibrain/internal/client/rag"
)

type (
	IngestPayload struct {
		Fid          string
		Text         string
		Title        *string
		UploadDate   *string
		CreationDate *string
		Keywords     string
	}
	IngestHandler struct {
		Rag *rag.Rag
	}
)

const TypeIngestDocument = "document:ingest"

func (h *IngestHandler) Handle(ctx context.Context, t *asynq.Task) error {
	var p IngestPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		slog.Error("error unmarshaling payload")
		return err
	}
	slog.Info("calling rag store document")
	if err := h.Rag.StoreDocument(
		p.Text,
		p.Fid,
		p.Title,
		p.UploadDate,
		p.CreationDate,
		p.Keywords,
	); err != nil {
		return fmt.Errorf("store document: %w", err)
	}

	slog.Info("rag indexing complete", "fid", p.Fid)

	return nil
}
