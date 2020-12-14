package schema

import "github.com/facebook/ent"

// Owner holds the schema definition for the Owner entity.
type Owner struct {
	ent.Schema
}

// Fields of the Owner.
func (Owner) Fields() []ent.Field {
	return nil
}

// Edges of the Owner.
func (Owner) Edges() []ent.Edge {
	return nil
}
