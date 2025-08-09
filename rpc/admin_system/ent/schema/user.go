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
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/schema/local_mixins"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("password").
			NotEmpty().
			Comment("User password"),
		field.String("email").
			NotEmpty().
			Comment("User email"),
		field.String("name").
			Default("").
			Comment("User name"),
		field.String("nickname").
			Default("").
			Comment("User nickname"),
		field.String("phone").
			SchemaType(map[string]string{
				dialect.MySQL: "varchar(32)",
			}).
			Default("").
			Comment("User phone"),
		field.String("avatar").
			SchemaType(map[string]string{
				dialect.MySQL: "varchar(1024)",
			}).
			Default("").
			Comment("User avatar"),
		field.String("remark").
			Default("").
			Comment("User remark"),
		field.String("department_id").
			Optional().
			Nillable().
			Comment("User DepartmentID"),
	}
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDStringMixin{},
		local_mixins.SoftDeleteMixin{},
		mixins.StatusMixin{},
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("department", Department.Type).Unique().Field("department_id"),
		edge.To("positions", Position.Type).StorageKey(edge.Table("admin_system_user_positions")),
		edge.To("roles", Role.Type).StorageKey(edge.Table("admin_system_user_roles")),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("email").Unique(),
		index.Fields("name"),
		index.Fields("nickname"),
		index.Fields("phone"),
		index.Fields("department_id"),
	}
}

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "admin_system_users"},
	}
}
