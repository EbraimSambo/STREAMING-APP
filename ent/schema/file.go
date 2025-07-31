package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// File holds the schema definition for the File entity.
type File struct {
	ent.Schema
}

// Mixin of the File.
func (File) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Fields of the File.
func (File) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique(), // Define ID as string
		field.Text("file_name").NotEmpty(),
		field.Bool("visibility").Default(false),
		field.String("status").Default("PENDING"),      // PENDING, PROCESSING, COMPLETED, FAILED
		field.Text("status_details").Optional(), // Store error messages
		field.JSON("metadata", map[string]interface{}{}).Optional(), // Store video metadata
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