package category

import (
	"context"
	"controle_financas/internal/core/domain/apiError"
	"controle_financas/internal/core/filters"
	"controle_financas/internal/core/transaction"
	"controle_financas/internal/core/validator"
	"database/sql"

	"github.com/google/uuid"
)

type categoryService struct {
	repo CategoryRepository
	tx   transaction.Manager
}

type CategoryService interface {
	FindById(
		ctx context.Context,
		id uuid.UUID,
	) (*Category, error)

	FindAll(
		ctx context.Context,
		f filters.Filters,
		search ...string,
	) ([]*Category, filters.Metadata, error)
	Insert(
		ctx context.Context,
		model *Category,
		tx *sql.Tx,
	) error

	Update(
		ctx context.Context,
		model *Category,
		tx *sql.Tx,
	) error

	DeleteById(
		ctx context.Context,
		id uuid.UUID,
		tx *sql.Tx,
	) error
}

func NewService(
	repo CategoryRepository,
	tx transaction.Manager,
) CategoryService {
	return &categoryService{
		repo: repo,
		tx:   tx,
	}
}

func (s *categoryService) FindById(
	ctx context.Context,
	id uuid.UUID,
) (*Category, error) {
	return s.repo.FindById(ctx, id)
}

func (s *categoryService) FindAll(
	ctx context.Context,
	f filters.Filters,
	search ...string,
) ([]*Category, filters.Metadata, error) {
	return s.repo.FindAll(ctx, f, search...)
}

func (s *categoryService) Insert(
	ctx context.Context,
	model *Category,
	tx *sql.Tx,
) error {
	saveFn := func(tx *sql.Tx) error {
		v := validator.New()
		if model.ValidateCategory(v); !v.Valid() {
			return apiError.NewValidationError(v.Errors)
		}

		return s.repo.Insert(ctx, tx, model)
	}

	if tx != nil {
		return saveFn(tx)
	}
	return s.tx.RunInTx(ctx, saveFn)
}

func (s *categoryService) Update(
	ctx context.Context,
	model *Category,
	tx *sql.Tx,
) error {
	updateFn := func(tx *sql.Tx) error {
		v := validator.New()
		if model.ValidateCategory(v); !v.Valid() {
			return apiError.NewValidationError(v.Errors)
		}

		return s.repo.Update(ctx, tx, model)
	}

	if tx != nil {
		return updateFn(tx)
	}
	return s.tx.RunInTx(ctx, updateFn)
}

func (s *categoryService) DeleteById(
	ctx context.Context,
	id uuid.UUID,
	tx *sql.Tx,
) error {
	deleteFn := func(tx *sql.Tx) error {
		return s.repo.DeleteById(ctx, tx, id)
	}

	if tx != nil {
		return deleteFn(tx)
	}
	return s.tx.RunInTx(ctx, deleteFn)
}
