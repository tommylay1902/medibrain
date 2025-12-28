-- From: ../internal/api/domain/documentmeta/model.go
CREATE TABLE document_meta(
	id UUID PRIMARY KEY gen_random_uuid(),
	fid BIGINT,
	date_user_uploaded TIMESTAMP,
	date_document_uploaded TIMESTAMP,
	title TEXT, 
	author TEXT,
	subject TEXT
	)
;

-- From: ../internal/api/domain/tags/model.go
CREATE TABLE tags(
	id UUID PRIMARY KEY gen_random_uuid(),
	name STRING
	)
;

