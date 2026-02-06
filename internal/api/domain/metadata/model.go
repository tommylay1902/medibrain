package metadata

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
	keywords TEXT NOT NULL,
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
	Keywords         string     `json:"keywords" db:"keywords"`
	Title            *string    `json:"title" db:"title"`
	Author           *string    `json:"author" db:"author"`
	Subject          *string    `json:"subject" db:"subject"`
}

type DocumentMetaList []*DocumentMeta
