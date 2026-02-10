package note

import (
	"time"

	"github.com/google/uuid"
)

type Note struct {
	ID               *uuid.UUID `json:"id" db:"id"`
	CreationDate     *time.Time `json:"creation_date" db:"creation_date"`
	ModificationDate *time.Time `json:"modification_date" db:"modification_date"`
	Content          string     `json:"content" db:"content"`
}

type NoteKeyword struct {
	ID     *uuid.UUID `json:"id" db:"id"`
	NoteID uuid.UUID  `json:"note_id" db:"note_id"`
}

type (
	NoteList        []*Note
	NoteKeywordList []*NoteKeyword
)
