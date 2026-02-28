package note

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/tommylay1902/medibrain/internal/database"
)

type NoteRepo struct {
	uow database.UnitOfWorkFactory
}

func NewNoteRepo(uow database.UnitOfWorkFactory) *NoteRepo {
	return &NoteRepo{
		uow: uow,
	}
}

func (nr *NoteRepo) List(ctx context.Context) (NoteList, error) {
	db := nr.uow.GetDB(ctx)

	ext, ok := db.(sqlx.ExtContext)
	if !ok {
		return nil, errors.New("invalid database connection")
	}

	rows, err := ext.QueryxContext(ctx, "SELECT * FROM note")
	if err != nil {
		return nil, err
	}

	var notes []*Note
	for rows.Next() {
		var note Note
		err := rows.StructScan(&note)
		if err != nil {
			return nil, err
		}
		notes = append(notes, &note)
	}

	return notes, nil
}

func (nr *NoteRepo) ListWithKeywords(ctx context.Context) ([]*NoteWithTags, error) {
	db := nr.uow.GetDB(ctx)
	ext, ok := db.(sqlx.ExtContext)
	if !ok {
		return nil, errors.New("invalid database connection")
	}

	query := `
		SELECT n.id AS id, 
			n.creation_date AS creation_date, 
			n.modification_date as modification_date,
			n.content as content, 
			t.name as tag

		FROM note AS n 
		INNER JOIN note_tag as nt
		ON n.id = nt.note_id
		INNER JOIN tag AS t 
		ON t.id = nt.tag_id
		 `

	rows, err := ext.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var joinResult []NoteJoinTag
	for rows.Next() {
		var noteJoinTag NoteJoinTag
		err := rows.StructScan(&noteJoinTag)
		if err != nil {
			return nil, err
		}

		joinResult = append(joinResult, noteJoinTag)
	}

	if len(joinResult) == 0 {
		return nil, nil
	}

	noteMap := make(map[uuid.UUID]*NoteWithTags)

	for _, row := range joinResult {
		note, exists := noteMap[*row.ID]
		if !exists {
			noteMap[*row.ID] = &NoteWithTags{
				ID:               row.ID,
				CreationDate:     row.CreationDate,
				ModificationDate: row.ModificationDate,
				Content:          row.Content,
				Tags:             []*string{&row.Tag},
			}
		} else {
			note.Tags = append(note.Tags, &row.Tag)
		}
	}

	result := make([]*NoteWithTags, 0, len(noteMap))
	for _, note := range noteMap {
		result = append(result, note)
	}
	return result, nil
}

func (nr *NoteRepo) CreateNote(ctx context.Context, note *Note) (*uuid.UUID, error) {
	db := nr.uow.GetDB(ctx)
	ext, ok := db.(sqlx.ExtContext)
	if !ok {
		return nil, errors.New("invalid database connection")
	}

	query := `
    INSERT INTO note (
        creation_date, 
        modification_date, 
        content,
				title
    ) VALUES (
        $1, 
        $2, 
        $3,
				$4
    )
		RETURNING id
	`
	var uuid uuid.UUID
	row, err := ext.QueryContext(ctx, query, note.CreationDate, note.ModificationDate, note.Content, note.Title)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		if err = row.Scan(&uuid); err != nil {
			return nil, err
		}
	}
	return &uuid, err
}

func (nr *NoteRepo) CreateTagBatch(ctx context.Context, tags []string) ([]Tag, error) {
	db := nr.uow.GetDB(ctx)

	ext, ok := db.(sqlx.ExtContext)
	if !ok {
		return nil, errors.New("invalid database connection")
	}

	query := `
        INSERT INTO tag (id, name)
        SELECT gen_random_uuid(), name FROM UNNEST($1::text[]) AS name
        ON CONFLICT (name) DO UPDATE 
            SET name = EXCLUDED.name
        RETURNING id, name
    `

	rows, err := ext.QueryxContext(ctx, query, pq.Array(tags))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultTags []Tag
	for rows.Next() {
		var tag Tag
		if err := rows.StructScan(&tag); err != nil {
			return nil, err
		}

		resultTags = append(resultTags, tag)
	}

	return resultTags, rows.Err()
}

func (nr *NoteRepo) LinkNoteWithTags(ctx context.Context, noteId uuid.UUID, tagIds []*uuid.UUID) error {
	db := nr.uow.GetDB(ctx)

	ext, ok := db.(sqlx.ExtContext)
	if !ok {
		return errors.New("invalid database connection")
	}

	query := `
        INSERT INTO note_tag (note_id, tag_id) 
				SELECT $1, UNNEST($2::uuid[])
    `
	_, err := ext.ExecContext(ctx, query, noteId, pq.Array(tagIds))

	return err
}

func (nr *NoteRepo) ListTags(ctx context.Context) (TagList, error) {
	db := nr.uow.GetDB(ctx)

	ext, ok := db.(sqlx.ExtContext)

	if !ok {
		return nil, errors.New("invalid database connection")
	}

	query := `SELECT * FROM tag`
	rows, err := ext.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var tags TagList
	for rows.Next() {
		var tag Tag
		if err := rows.StructScan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, &tag)
	}
	return tags, nil
}

func (nr *NoteRepo) CreateTag(ctx context.Context, tag Tag) (*Tag, error) {
	db := nr.uow.GetDB(ctx)

	ext, ok := db.(sqlx.ExtContext)

	if !ok {
		return nil, errors.New("invalid database connection")
	}

	query := `
		INSERT INTO tag (name) values($1)
		RETURNING id, name
	`
	result, err := ext.QueryxContext(ctx, query, tag.Name)
	if err != nil {
		fmt.Println("error querying db")
		return nil, err
	}

	var createdTag Tag
	validRow := result.Next()
	if !validRow {
		return nil, errors.New("no tag not created succesfully")
	}

	err = result.StructScan(&createdTag)
	if err != nil {
		fmt.Println("error scanning tag")
		return nil, err
	}

	return &createdTag, nil
}
