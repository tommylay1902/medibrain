-- From: ../internal/api/domain/documentmeta/model.go
CREATE TABLE document_meta(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	fid BIGINT,
	upload_date TIMESTAMP,
	creation_date TIMESTAMP,
	title TEXT, 
	author TEXT,
	subject TEXT
	)
;

-- From: ../internal/api/domain/tags/model.go
CREATE TABLE tags(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name TEXT 
	)
;

