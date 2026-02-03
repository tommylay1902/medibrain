package document

import "github.com/jmoiron/sqlx"

type DocumentPipelineRepo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *DocumentPipelineRepo {
	return &DocumentPipelineRepo{db: db}
}
