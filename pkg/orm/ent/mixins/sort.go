package mixins

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

type SortMixin struct {
	mixin.Schema
}

func (SortMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("sort").Default(0),
	}
}

func (SortMixin) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("sort"),
	}
}
