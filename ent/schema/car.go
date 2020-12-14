package schema

import (
	_ "log"

	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
)

// Car holds the schema definition for the Car entity.
type Car struct {
	ent.Schema
}

// Edges of the Car.
func (Car) Edges() []ent.Edge {
	return []ent.Edge{
		// create an inverse-edge called "owner" of type "User"
		// and reference it to the "cars" edge(in User schema)
		// explicitly user the `Ref` method.
		edge.From("Owner", User.Type).
			Ref("cars").
			Unique(),
	}
}
