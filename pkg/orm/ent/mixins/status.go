package mixins

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

const (
	StatusNormal uint8 = 1
	StatusBanned uint8 = 2
)

type StatusMixin struct {
	mixin.Schema
}

func (StatusMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Uint8("status").
			Default(1).
			Optional().
			Comment("Status | 1 正常 2 禁用"),
	}
}

func (StatusMixin) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("status"),
	}
}
