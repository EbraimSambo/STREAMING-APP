package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// File holds the schema definition for the File entity.
type File struct {
	ent.Schema
}

// Fields of the File.
func (File) Fields() []ent.Field {
	return []ent.Field{
		field.Text("file").NotEmpty(),
		field.Bool("visibility").Default(false),
		field.Time("created_at").
			Default(time.Now),
		field.Time("deleted_at").
			Optional(),
	}
}

// Edges of the File.
func (File) Edges() []ent.Edge {
	return nil
}
