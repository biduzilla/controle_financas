package transaction

import (
	"controle_financas/internal/core/domain/models"
	"controle_financas/internal/core/validator"
	"controle_financas/internal/features/category"
	"controle_financas/internal/features/user"
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	models.BaseModel
	ID          uuid.UUID          `db:"id" repo:"pk,auto"`
	User        *user.Usuario      `db:"-" repo:"-"`
	Category    *category.Category `db:"-" repo:"-"`
	Description string             `db:"description" repo:"insert,update"`
	Amount      float64            `db:"amount" repo:"insert,update"`
}

type TransactionDTO struct {
	ID          *uuid.UUID            `json:"transaction_id"`
	Version     *int                  `json:"version"`
	User        *user.UsuarioDTO      `json:"user"`
	Category    *category.CategoryDTO `json:"category"`
	Description *string               `json:"description"`
	Amount      *float64              `json:"amount"`
	CreatedAt   *time.Time            `json:"created_at"`
}

func (t *Transaction) ToDTO() *TransactionDTO {
	dto := TransactionDTO{}

	dto.ID = &t.ID
	dto.CreatedAt = &t.CreatedAt
	dto.Version = &t.Version
	dto.Description = &t.Description
	dto.Amount = &t.Amount

	if t.User != nil {
		dto.User = t.User.ToDTO()
	}

	if t.Category != nil {
		dto.Category = t.Category.ToDTO()
	}

	return &dto
}

func (t *TransactionDTO) ToModel() (*Transaction, error) {
	transaction := &Transaction{}

	if t.ID != nil {
		transaction.ID = *t.ID
	}
	if t.Version != nil {
		transaction.Version = *t.Version
	}
	if t.User != nil {
		user, err := t.User.ToModel()
		if err != nil {
			return nil, err
		}
		transaction.User = user
	}
	if t.Category != nil {
		c, err := t.Category.ToModel()
		if err != nil {
			return nil, err
		}
		transaction.Category = c
	}
	if t.Description != nil {
		transaction.Description = *t.Description
	}
	if t.Amount != nil {
		transaction.Amount = *t.Amount
	}

	return transaction, nil
}

func (t *Transaction) ValidateTransaction(v *validator.Validator) {
	v.Check(t.User != nil, "user", "must be provided")
	v.Check(t.Category != nil, "category", "must be provided")
	v.Check(t.Description != "", "description", "must be provided")
	v.Check(len(t.Description) <= 500, "description", "must not be more than 500 bytes long")
	v.Check(t.Amount > 0, "amount", "must be positive")
	v.Check(t.Amount != 0, "amount", "must be provided")
}
