package documentmeta

import (
	"time"

	"github.com/google/uuid"
)

var DocumentMetaSchema = `CREATE TABLE document_meta(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	fid BIGINT,
	date_user_uploaded TIMESTAMP,
	date_document_uploaded TIMESTAMP,
	title TEXT, 
	author TEXT,
	subject TEXT
	)`

type DocumentMeta struct {
	ID                   *uuid.UUID `json:"id" db:"id"`
	Fid                  *uint32    `json:"fid" db:"fid"`
	DateUserUploaded     time.Time  `json:"date_user_uploaded" db:"date_user_uploaded" `
	DateDocumentUploaded time.Time  `json:"date_document_uploaded" db:"date_document_uploaded"`
	// Tags                 []string   `json:"tags"`
	Title   string `json:"title" db:"title"`
	Author  string `json:"author" db:"author"`
	Subject string `json:"subject" db:"subject"`
	// RelatedTo
}

type DocumentMetaList []*DocumentMeta
