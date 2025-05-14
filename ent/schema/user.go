package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// User — модель пользователя.
type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique(),
		field.String("fio").NotEmpty().Comment("ФИО"),
		field.String("email").
			NotEmpty().
			Unique().
			Comment("уникальный e-mail"),
		field.String("password_hash").
			Sensitive().
			Comment("bcrypt-хеш пароля"),
		field.Enum("role").Values("admin", "master", "client").Default("client"),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

func (User) Edges() []ent.Edge {
	return nil
}
