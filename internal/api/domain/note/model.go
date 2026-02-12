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
	Tag              string     `json:"name" db:"tag"`
}

type NoteWithTags struct {
	ID               *uuid.UUID `json:"id" db:"id"`
	CreationDate     *time.Time `json:"creation_date" db:"creation_date"`
	ModificationDate *time.Time `json:"modification_date" db:"modification_date"`
	Content          string     `json:"content" db:"content"`
	Tags             []*string  `json:"tags"`
}

type Tag struct {
	ID   *uuid.UUID `json:"id" db:"id"`
	Name string     `json:"name" db:"name"`
}

type (
	NoteList []*Note
	TagList  []*Tag
)
