-- From: ../internal/api/domain/documentmeta/model.go

DROP TABLE IF EXISTS document_meta CASCADE;
	CREATE TABLE document_meta(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	fid TEXT,
	upload_date TIMESTAMP,
	creation_date TIMESTAMP,
	title TEXT, 
	author TEXT,
	subject TEXT
	)
;

-- From: ../internal/api/domain/keyword/model.go
DROP TABLE IF EXISTS keyword CASCADE;
CREATE TABLE keyword(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name TEXT 
	)
;

-- From: ../internal/api/domain/tags/model.go
DROP TABLE IF EXISTS tags CASCADE;
CREATE TABLE tags(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name TEXT 
	)
;

