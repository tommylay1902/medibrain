package documentmeta

import (
	"github.com/google/uuid"
)

var DocumentMetaSchema = `
	CREATE TABLE document_meta(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	thumbnail_fid TEXT NOT NULL,
	pdf_fid TEXT NOT NULL,
	modification_date TIMESTAMP,
	creation_date TIMESTAMP,
	title TEXT, 
	author TEXT,
	subject TEXT
	)`

type DocumentMeta struct {
	ID               *uuid.UUID `json:"id,omitempty" db:"id"`
	ThumbnailFid     string     `json:"thumbnailFid" db:"thumbnail_fid"`
	PdfFid           string     `json:"pdfFid" db:"pdf_fid"`
	ModificationDate *string    `json:"modificationDate" db:"modification_date" `
	CreationDate     *string    `json:"creationDate" db:"creation_date"`
	// Tags                 []string   `json:"tags"`
	Title   *string `json:"title" db:"title"`
	Author  *string `json:"author" db:"author"`
	Subject *string `json:"subject" db:"subject"`
	// RelatedTo
}

type DocumentMetaList []*DocumentMeta
