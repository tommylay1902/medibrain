package note

import (
	"time"

	"github.com/google/uuid"
)

type Note struct {
	ID               *uuid.UUID
	CreationDate     *time.Time
	ModificationDate *time.Time
	Content          string
}

type NoteKeywords struct {
	ID     *uuid.UUID
	NoteID uuid.UUID
}
