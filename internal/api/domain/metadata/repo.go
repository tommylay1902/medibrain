package metadata

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type MetadataRepo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *MetadataRepo {
	return &MetadataRepo{db: db}
}

func (dmr *MetadataRepo) List() (MetadataList, error) {
	var results MetadataList

	err := dmr.db.Select(&results, "SELECT * FROM document_meta ")
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (dmr *MetadataRepo) Create(meta *Metadata) error {
	_, err := dmr.db.NamedQuery("INSERT INTO metadata(thumbnail_fid, pdf_fid, keywords, title, author, subject) VALUES(:thumbnail_fid, :pdf_fid, :keywords, :title, :author, :subject)", meta)
	if err != nil {
		fmt.Println("error ")
		return err
	}
	return nil
}
