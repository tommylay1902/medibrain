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

func (nr *NoteRepo) List() (NoteList, error) {
	var notes []*Note
	err := nr.db.Select(&notes, "SELECT * FROM note")
	if err != nil {
		return nil, err
	}
	return notes, nil
}
