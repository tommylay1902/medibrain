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

CREATE TABLE note(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	creation_date TIMESTAMP NOT NULL DEFAULT NOW(),
	modification_date TIMESTAMP NOT NULL DEFAULT NOW(),
	content TEXT
);

CREATE TABLE tag(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT UNIQUE
);

CREATE TABLE note_tag(
  note_id UUID REFERENCES note,
  tag_id UUID REFERENCES tag,
  PRIMARY KEY(note_id, tag_id)
);

