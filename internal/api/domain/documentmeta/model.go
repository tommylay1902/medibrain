package documentmeta

import (
	"time"

	"github.com/google/uuid"
)

type DocumentMeta struct {
	Id                   *uuid.UUID `json:"id"`
	Fid                  *uint32    `json:"fid"`
	DateUserUploaded     time.Time  `json:"date_user_uploaded"`
	DateDocumentUploaded time.Time  `json:"date_document_uploaded"`
	Tags                 []string   `json:"tags"`
	Title                string     `json:"string"`
	Author               string     `json:"author"`
	Subject              string     `json:"string"`
	// RelatedTo
}
