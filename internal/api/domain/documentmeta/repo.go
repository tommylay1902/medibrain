package documentmeta

import "github.com/jmoiron/sqlx"

type DocumentMetaRepo struct {
	db *sqlx.DB
}

func (dmr *DocumentMetaRepo) List() (DocumentMetaList, error) {
	var results DocumentMetaList

	err := dmr.db.Select(&results, "SELECT * FROM document_meta ")
	if err != nil {
		return nil, err
	}

	return results, nil
}
