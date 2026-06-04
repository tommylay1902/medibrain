package metadata

import (
	"github.com/google/uuid"
)

type Metadata struct {
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

type MetadataList []*Metadata
