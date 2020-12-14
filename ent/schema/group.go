package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
	"regexp"
)

// Group holds the schema definition for the Group entity.
type Group struct {
	ent.Schema
}

// Edges of the Group.
func (Group) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("users", User.Type),
	}
}

// Fields of the Car.
func (Car) Fields() []ent.Field {
	return []ent.Field{
		field.String("model"),
		field.Time("registered_at"),
	}
}

// Fields of the Group.
func (Group) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			// Regexp validation for group name.
			Match(regexp.MustCompile("[a-zA-Z_]+$")),
	}
}
