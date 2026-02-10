CREATE TABLE note(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	creation_date TEXT NOT NULL,
	modification_date TEXT NOT NULL,
	content TEXT
);

CREATE TABLE note_keyword(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	note_id UUID REFERENCES note
);


