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

type NoteJoinTag struct {
	ID               *uuid.UUID `json:"id" db:"id"`
	CreationDate     *time.Time `json:"creation_date" db:"creation_date"`
	ModificationDate *time.Time `json:"modification_date" db:"modification_date"`
	Content          string     `json:"content" db:"content"`
	Keyword          string     `json:"keyword"`
}

type NoteWithTags struct {
	ID               *uuid.UUID `json:"id" db:"id"`
	CreationDate     *time.Time `json:"creation_date" db:"creation_date"`
	ModificationDate *time.Time `json:"modification_date" db:"modification_date"`
	Content          string     `json:"content" db:"content"`
	Keywords         []*string  `json:"keywords"`
}

type NoteList []*Note
