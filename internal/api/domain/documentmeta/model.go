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
	ID                   *uuid.UUID `json:"id"`
	Fid                  *uint32    `json:"fid"`
	DateUserUploaded     time.Time  `json:"date_user_uploaded"`
	DateDocumentUploaded time.Time  `json:"date_document_uploaded"`
	// Tags                 []string   `json:"tags"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Subject string `json:"subject"`
	// RelatedTo
}

type DocumentMetaList []*DocumentMeta
