package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/orm/ent/mixins"
)

type Position struct {
	ent.Schema
}

func (Position) Fields() []ent.Field {
	return []ent.Field{
		field.String("position_name").
			Comment("Position name | 职位名称"),
		field.String("code").
			Comment("Position code | 职位编号"),
		field.String("remark").
			Default("").
			Comment("Position remark | 职位备注"),
	}
}

func (Position) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDStringMixin{},
		mixins.SortMixin{},
	}
}

func (Position) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("users", User.Type).Ref("positions"),
	}
}

func (Position) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("position_name").Unique(),
		index.Fields("code").Unique(),
	}
}

func (Position) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "admin_system_positions"},
	}
}
