package note

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/tommylay1902/medibrain/internal/database"
)

type NoteService struct {
	repo *NoteRepo
	uow  database.UnitOfWorkFactory
}

func NewNoteService(repo *NoteRepo, uow database.UnitOfWorkFactory) *NoteService {
	return &NoteService{
		repo: repo,
		uow:  uow,
	}
}

func (ns *NoteService) List(ctx context.Context) (NoteList, error) {
	notes, err := ns.repo.List(ctx)

	return notes, err
}

func (ns *NoteService) ListWithKeywords(ctx context.Context) ([]*NoteWithTags, error) {
	result, err := ns.repo.ListWithKeywords(ctx)
	return result, err
}

func (ns *NoteService) CreateNoteWithTags(ctx context.Context, note *Note, tags []string) error {
	if len(tags) > 7 {
		return errors.New("a note can only have a maxium of 7 tags")
	}

	now := time.Now()

	if note.CreationDate == nil {
		note.CreationDate = &now
	}

	if note.ModificationDate == nil {
		note.ModificationDate = &now
	}

	uow, ctx, err := ns.uow.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := uow.Rollback(); rbErr != nil {
				log.Printf("rollback failed: %v, original error: %v", rbErr, err)
			}
		}
	}()

	noteID, err := ns.repo.CreateNote(ctx, note)
	if err != nil {
		return fmt.Errorf("error creating notes: %v", err)
	}

	results, err := ns.repo.CreateTagBatch(ctx, tags)
	if err != nil {
		return fmt.Errorf("error creating tag batch: %v", err)
	}

	tagIDs := make([]*uuid.UUID, len(results))
	for i := range results {
		tagIDs[i] = results[i].ID
	}

	err = ns.repo.LinkNoteWithTags(ctx, *noteID, tagIDs)
	if err != nil {
		return fmt.Errorf("error Linking Notes with Tags: %v", err)
	}

	return uow.Commit()
}

func (ns *NoteService) CreateTag(ctx context.Context, tag Tag) (*Tag, error) {
	id, err := ns.repo.CreateTag(ctx, tag)
	return id, err
}

func (ns *NoteService) ListTag(ctx context.Context) (TagList, error) {
	return ns.repo.ListTags(ctx)
}
