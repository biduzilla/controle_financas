package category

import (
	"controle_financas/internal/core/domain/models"
	"controle_financas/internal/core/validator"
	"controle_financas/internal/features/user"
	"time"
)

type TypeCategoria int

const (
	RECEITA TypeCategoria = iota + 1
	DESPESA
)

func (t TypeCategoria) String() string {
	switch t {
	case RECEITA:
		return "RECEITA"
	case DESPESA:
		return "DESPESA"
	default:
		return "Unknown"
	}
}

func TypeCategoriaFromString(s string) TypeCategoria {
	switch s {
	case "RECEITA":
		return RECEITA
	case "DESPESA":
		return DESPESA
	default:
		return 0
	}
}

type Category struct {
	models.BaseModel
	ID    int64
	Name  string
	Type  TypeCategoria
	Color string
	User  *user.Usuario
}

type CategoryDTO struct {
	ID        *int64           `json:"category_id"`
	CreatedAt *time.Time       `json:"-"`
	Name      *string          `json:"name"`
	Type      *string          `json:"type"`
	Color     *string          `json:"color"`
	User      *user.UsuarioDTO `json:"user"`
	Version   *int             `json:"version"`
}

func (m *CategoryDTO) ToModel() (*Category, error) {
	category := &Category{}
	if m.ID != nil {
		category.ID = *m.ID
	}

	if m.CreatedAt != nil {
		category.CreatedAt = *m.CreatedAt
	}

	if m.Name != nil {
		category.Name = *m.Name
	}

	if m.Type != nil {
		category.Type = TypeCategoriaFromString(*m.Type)
	}

	if m.Color != nil {
		category.Color = *m.Color
	}

	if m.User != nil {
		user, err := m.User.ToModel()
		if err != nil {
			return nil, err
		}

		category.User = user
	}

	if m.Version != nil {
		category.Version = *m.Version
	}

	return category, nil
}

func (m *Category) ToDTO() *CategoryDTO {
	category := &CategoryDTO{}

	category.ID = &m.ID
	category.CreatedAt = &m.CreatedAt
	category.Name = &m.Name
	typeStr := m.Type.String()
	category.Type = &typeStr
	category.Color = &m.Color
	category.User = m.User.ToDTO()
	category.Version = &m.Version

	return category
}

func (c *Category) ValidateCategory(v *validator.Validator) {
	v.Check(c.Name != "", "name", "must be provided")
	v.Check(len(c.Name) <= 500, "name", "must not be more than 500 bytes long")
	v.Check(c.Type.String() != "", "type", "must be provided")
	v.Check(c.Type.String() != "Unknown", "type", "invalid type")
	v.Check(c.Color != "", "color", "must be provided")
}
