package tags

var TagsSchema = `CREATE TABLE tags(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name TEXT 
	)`
