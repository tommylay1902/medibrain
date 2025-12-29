package documentmeta

import (
	"time"

	"github.com/google/uuid"
)

var DocumentMetaSchema = `CREATE TABLE document_meta(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	fid BIGINT,
	upload_date TIMESTAMP,
	creation_date TIMESTAMP,
	title TEXT, 
	author TEXT,
	subject TEXT
	)`

type DocumentMeta struct {
	ID           *uuid.UUID `json:"id" db:"id"`
	Fid          *uint32    `json:"fid" db:"fid"`
	UploadDate   time.Time  `json:"upload_date" db:"upload_date" `
	CreationDate time.Time  `json:"creation_date" db:"creation_date"`
	// Tags                 []string   `json:"tags"`
	Title   string `json:"title" db:"title"`
	Author  string `json:"author" db:"author"`
	Subject string `json:"subject" db:"subject"`
	// RelatedTo
}

type DocumentMetaList []*DocumentMeta
