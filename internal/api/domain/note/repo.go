package note

import (
	"github.com/jmoiron/sqlx"
)

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

func (nr *NoteRepo) ListWithKeywords() (*NoteWithTags, error) {
	var joinResult []NoteJoinTag
	err := nr.db.Select(&joinResult, "SELECT n.id AS id, n.creation_date AS creation_date, n.modification_date as modification_date n.content as content, nk.keyword as keyword FROM note AS n INNER JOIN note_keyword AS nk ON n.id = nk.note_id ")
	if err != nil {
		return nil, err
	}

	if len(joinResult) == 0 {
		return nil, nil
	}

	noteKeywords := &NoteWithTags{
		ID:               joinResult[0].ID,
		CreationDate:     joinResult[0].CreationDate,
		ModificationDate: joinResult[0].ModificationDate,
		Content:          joinResult[0].Content,
		Keywords:         make([]*string, 0, len(joinResult)),
	}

	for i := range joinResult {
		noteKeywords.Keywords = append(noteKeywords.Keywords, &joinResult[i].Keyword)
	}

	return noteKeywords, nil
}

// func (nr *NoteRepo) CreateNoteWithKeywords(notesWithKeywords *NoteWithKeywords) (uuid.UUID, error) {
// }
