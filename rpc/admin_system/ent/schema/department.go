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

type Department struct {
	ent.Schema
}

func (Department) Fields() []ent.Field {
	return []ent.Field{
		field.String("department_name").
			NotEmpty().
			Comment("Department name"),
		field.String("remark").
			Default("").
			Comment("Department remark"),
		field.String("parent_id").
			Optional().
			Nillable().
			Comment("Parent ID"),
		field.String("id_path").
			SchemaType(map[string]string{
				dialect.MySQL: "varchar(512)",
			}).
			NotEmpty().
			Comment("Id path"),
	}
}

func (Department) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDStringMixin{},
		mixins.SortMixin{},
	}
}

func (Department) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("children", Department.Type),
		edge.From("parent", Department.Type).Ref("children").Unique().Field("parent_id"),
		edge.From("users", User.Type).Ref("department"),
	}
}

func (Department) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("department_name", "parent_id").Unique(),
		index.Fields("parent_id"),
		index.Fields("id_path").Unique(),
	}
}

func (Department) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "admin_system_departments"},
	}
}
