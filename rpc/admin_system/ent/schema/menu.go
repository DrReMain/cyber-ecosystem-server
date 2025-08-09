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

type Menu struct {
	ent.Schema
}

func (Menu) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").
			Default("").
			Comment("Menu title"),
		field.String("icon").
			SchemaType(map[string]string{
				dialect.MySQL: "varchar(512)",
			}).
			Default("").
			Comment("Menu icon"),
		field.String("code").
			NotEmpty().
			Comment("Menu code"),
		field.String("code_path").
			SchemaType(map[string]string{
				dialect.MySQL: "varchar(512)",
			}).
			NotEmpty().
			Comment("Menu code path (code1_code2_code3)"),
		field.String("parent_id").
			Optional().
			Nillable().
			Comment("Parent MenuID"),
		field.String("menu_type").
			SchemaType(map[string]string{
				dialect.MySQL: "varchar(16)",
			}).
			NotEmpty().
			Comment("Menu type (divider/group/menu/page/button)"),
		field.String("properties").
			SchemaType(map[string]string{
				dialect.MySQL: "varchar(2048)",
			}).
			Default("{}").
			Comment("Menu properties (JSON字符串)"),
	}
}

func (Menu) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDStringMixin{},
		mixins.SortMixin{},
		mixins.StatusMixin{},
	}
}

func (Menu) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("roles", Role.Type).Ref("menus"),
		edge.To("children", Menu.Type),
		edge.From("parent", Menu.Type).Ref("children").Unique().Field("parent_id"),
		edge.To("resources", Resource.Type),
	}
}

func (Menu) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("code", "parent_id").Unique(),
		index.Fields("parent_id"),
		index.Fields("code_path").Unique(),
		index.Fields("menu_type"),
	}
}

func (Menu) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "admin_system_menus"},
	}
}
