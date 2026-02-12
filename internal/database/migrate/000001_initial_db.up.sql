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
	creation_date TEXT NOT NULL,
	modification_date TEXT NOT NULL,
	content TEXT
);

CREATE TABLE tag(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT 
);

CREATE TABLE note_tag(
  note_id UUID REFERENCES note,
  tag_id UUID REFERENCES tag,
  PRIMARY KEY(note_id, tag_id)
);

