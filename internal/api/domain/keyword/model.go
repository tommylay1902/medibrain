package keyword

import "github.com/google/uuid"

var KeywordSchema = `CREATE TABLE keyword(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name TEXT 
	)`

type Keyword struct {
	ID   *uuid.UUID `json:"id" db:"id"`
	Name string     `json:"name" db:"name"`
}

type Keywords []Keyword
