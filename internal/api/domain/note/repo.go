package note

import "github.com/jmoiron/sqlx"

type NoteRepo struct {
	db *sqlx.DB
}

func NewNoteRepo(db *sqlx.DB) *NoteRepo {
	return &NoteRepo{
		db: db,
	}
}
