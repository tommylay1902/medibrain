package documentmeta

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type DocumentMetaRepo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *DocumentMetaRepo {
	return &DocumentMetaRepo{db: db}
}

func (dmr *DocumentMetaRepo) List() (DocumentMetaList, error) {
	var results DocumentMetaList

	err := dmr.db.Select(&results, "SELECT * FROM document_meta ")
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (dmr *DocumentMetaRepo) Create(meta *DocumentMeta) error {
	_, err := dmr.db.NamedQuery("INSERT INTO document_meta(fid, title, author, subject) VALUES(:fid,:title, :author, :subject)", meta)
	if err != nil {
		fmt.Println("error ")
		return err
	}
	return nil
}
