package tags

import "github.com/google/uuid"

var TagsSchema = `CREATE TABLE tags(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name TEXT 
	)`

type Tags struct {
	ID   *uuid.UUID `json:"id" db:"id"`
	Name string     `json:"name" db:"name"`
}
