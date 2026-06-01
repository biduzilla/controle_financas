package user

import (
	"context"
	"controle_financas/internal/core/domain/apiError"
	"controle_financas/internal/core/transaction"
	"controle_financas/internal/core/validator"
	"database/sql"

	"github.com/google/uuid"
)

type UserService struct {
	repo userRepository
	tx   transaction.Manager
}

type userService interface {
	FindById(ctx context.Context, id uuid.UUID) (*Usuario, error)
	FindByEmail(ctx context.Context, email string) (*Usuario, error)
	Save(ctx context.Context, model *Usuario, tx *sql.Tx) error
	Update(ctx context.Context, model *Usuario, tx *sql.Tx) error
	DeleteById(ctx context.Context, id uuid.UUID, tx *sql.Tx) error
}

func NewService(
	repo userRepository,
	tx transaction.Manager,
) *UserService {
	return &UserService{
		repo: repo,
		tx:   tx,
	}
}

func (s *UserService) FindById(ctx context.Context, id uuid.UUID) (*Usuario, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *UserService) FindByEmail(ctx context.Context, email string) (*Usuario, error) {
	return s.repo.FindByEmail(ctx, email)
}

func (s *UserService) Save(ctx context.Context, model *Usuario, tx *sql.Tx) error {
	fn := func(tx *sql.Tx) error {
		v := validator.New()
		if model.Validate(v); !v.Valid() {
			return apiError.NewValidationError(v.Errors)
		}

		return s.repo.Save(ctx, tx, model)
	}

	if tx != nil {
		return fn(tx)
	}

	return s.tx.RunInTx(ctx, fn)
}

func (s *UserService) Update(ctx context.Context, model *Usuario, tx *sql.Tx) error {
	fn := func(tx *sql.Tx) error {
		v := validator.New()
		if model.Validate(v); !v.Valid() {
			return apiError.NewValidationError(v.Errors)
		}

		return s.repo.Update(ctx, tx, model)
	}

	if tx != nil {
		return fn(tx)
	}

	return s.tx.RunInTx(ctx, fn)
}

func (s *UserService) DeleteById(ctx context.Context, id uuid.UUID, tx *sql.Tx) error {
	fn := func(tx *sql.Tx) error {
		return s.repo.DeleteById(ctx, tx, id)
	}

	if tx != nil {
		return fn(tx)
	}

	return s.tx.RunInTx(ctx, fn)
}
