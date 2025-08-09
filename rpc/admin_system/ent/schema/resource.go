package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/orm/ent/mixins"
)

type Resource struct {
	ent.Schema
}

func (Resource) Fields() []ent.Field {
	return []ent.Field{
		field.String("menu_id").
			NotEmpty().
			Immutable().
			Comment("Menu ID"),
		field.String("method").
			SchemaType(map[string]string{
				dialect.MySQL: "varchar(32)",
			}).
			NotEmpty().
			Comment("HTTP method"),
		field.String("path").
			SchemaType(map[string]string{
				dialect.MySQL: "varchar(512)",
			}).
			NotEmpty().
			Comment("HTTP path"),
	}
}

func (Resource) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDStringMixin{},
	}
}

func (Resource) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("menu", Menu.Type).Ref("resources").Required().Unique().Immutable().Field("menu_id"),
	}
}

func (Resource) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("menu_id", "method", "path").Unique(),
		index.Fields("menu_id"),
		index.Fields("method", "path"),
	}
}

func (Resource) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "admin_system_resources"},
	}
}
