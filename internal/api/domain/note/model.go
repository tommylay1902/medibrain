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

const NoteSchema = `
	CREATE TABLE note(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	creation_date TEXT NOT NULL,
	modification_date TEXT NOT NULL,
	content TEXT
	)`

const NoteKeywordsSchema = `
	CREATE TABLE note_keyword(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	note_id UUID NOT NULL
	)`
