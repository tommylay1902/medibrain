CREATE TABLE metadata(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	thumbnail_fid TEXT NOT NULL,
	pdf_fid TEXT NOT NULL,
	modification_date TIMESTAMP,
	creation_date TIMESTAMP,
	keywords TEXT NOT NULL,
	title TEXT, 
	author TEXT,
	subject TEXT
);



