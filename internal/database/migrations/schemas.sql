-- From: ../internal/api/domain/tags/model.go
CREATE TABLE tags(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name TEXT 
	)
;

-- From: ../internal/api/domain/metadata/model.go

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
	)
;

