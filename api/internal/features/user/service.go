package user

import (
	"context"
	e "controle_financas/internal/core/domain/errors"
	"controle_financas/internal/core/transaction"
	"controle_financas/internal/core/validator"
	"database/sql"

	"github.com/google/uuid"
)

type userService struct {
	repo UserRepository
	tx   transaction.Manager
}

type UserService interface {
	FindById(ctx context.Context, id uuid.UUID) (*Usuario, error)
	FindByEmail(ctx context.Context, email string) (*Usuario, error)
	Save(ctx context.Context, model *Usuario, tx *sql.Tx) error
	Update(ctx context.Context, model *Usuario, tx *sql.Tx) error
	DeleteById(ctx context.Context, id uuid.UUID, tx *sql.Tx) error
}

func NewService(
	repo UserRepository,
	tx transaction.Manager,
) UserService {
	return &userService{
		repo: repo,
		tx:   tx,
	}
}

func (s *userService) FindById(ctx context.Context, id uuid.UUID) (*Usuario, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *userService) FindByEmail(ctx context.Context, email string) (*Usuario, error) {
	return s.repo.FindByEmail(ctx, email)
}

func (s *userService) Save(ctx context.Context, model *Usuario, tx *sql.Tx) error {
	fn := func(tx *sql.Tx) error {
		v := validator.New()
		if model.Validate(v); !v.Valid() {
			return e.NewValidationError(v.Errors)
		}

		return s.repo.Save(ctx, tx, model)
	}

	if tx != nil {
		return fn(tx)
	}

	return s.tx.RunInTx(ctx, fn)
}

func (s *userService) Update(ctx context.Context, model *Usuario, tx *sql.Tx) error {
	fn := func(tx *sql.Tx) error {
		v := validator.New()
		if model.Validate(v); !v.Valid() {
			return e.NewValidationError(v.Errors)
		}

		return s.repo.Update(ctx, tx, model)
	}

	if tx != nil {
		return fn(tx)
	}

	return s.tx.RunInTx(ctx, fn)
}

func (s *userService) DeleteById(ctx context.Context, id uuid.UUID, tx *sql.Tx) error {
	fn := func(tx *sql.Tx) error {
		return s.repo.DeleteById(ctx, tx, id)
	}

	if tx != nil {
		return fn(tx)
	}

	return s.tx.RunInTx(ctx, fn)
}
