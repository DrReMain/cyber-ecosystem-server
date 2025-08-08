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
			Comment("Menu title | 菜单标题"),
		field.String("icon").
			SchemaType(map[string]string{
				dialect.MySQL: "varchar(512)",
			}).
			Default("").
			Comment("Menu icon | 菜单图标"),
		field.String("code").
			NotEmpty().
			Comment("Menu code | 菜单CODE"),
		field.String("code_path").
			SchemaType(map[string]string{
				dialect.MySQL: "varchar(512)",
			}).
			Comment("Menu code path | 菜单CODE路径 (code1.code2.code3)"),
		field.String("parent_id").
			Default("").
			Comment("Parent MenuID | 父级菜单ID"),
		field.String("menu_type").
			SchemaType(map[string]string{
				dialect.MySQL: "varchar(16)",
			}).
			NotEmpty().
			Comment("Menu type | 菜单类型 (page/button)"),
		field.String("menu_path").
			Default("").
			Comment("Menu path | 菜单路径"),
		field.String("properties").
			SchemaType(map[string]string{
				dialect.MySQL: "varchar(2048)",
			}).
			Default("{}").
			Comment("Menu properties | 菜单属性 (JSON字符串)"),
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
		edge.To("children", Menu.Type).From("parent").Required().Unique().Field("parent_id"),
		edge.To("resources", Resource.Type),
	}
}

func (Menu) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("code").Unique(),
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
