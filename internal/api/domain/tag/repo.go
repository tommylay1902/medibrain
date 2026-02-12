package tag

import "github.com/jmoiron/sqlx"

type TagRepo struct {
	db *sqlx.DB
}

func NewTagRepo(db *sqlx.DB) *TagRepo {
	return &TagRepo{db: db}
}

func (tr *TagRepo) List() (TagList, error) {
	var tags TagList
	err := tr.db.Select(&tags, "SELECT * FROM TAG")
	return tags, err
}

func (tr *TagRepo) Create(tag *Tag) error {
	_, err := tr.db.NamedQuery("INSERT INTO tag(name) VALUES(:name)", tag)
	return err
}
