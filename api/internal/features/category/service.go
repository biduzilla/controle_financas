package category

import (
	"context"
	"controle_financas/internal/core/cache"
	"controle_financas/internal/core/domain/apiError"
	"controle_financas/internal/core/filters"
	"controle_financas/internal/core/transaction"
	"controle_financas/internal/core/validator"
	"database/sql"

	"github.com/google/uuid"
)

type CategoryService struct {
	repo       categoryRepository
	tx         transaction.Manager
	cache      cache.Cache
	keyBuilder cache.KeyBuilder
}

type categoryService interface {
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
	repo categoryRepository,
	tx transaction.Manager,
	cache cache.Cache,
	keyBuilder cache.KeyBuilder,
) *CategoryService {
	return &CategoryService{
		repo:       repo,
		tx:         tx,
		cache:      cache,
		keyBuilder: keyBuilder,
	}
}

func (s *CategoryService) FindById(
	ctx context.Context,
	id uuid.UUID,
) (*Category, error) {
	key := s.keyBuilder.BuildItemKey(id.String())

	return cache.FetchOrCache(ctx, s.cache, key, func() (*Category, error) {
		return s.repo.FindById(ctx, id)
	})
}

func (s *CategoryService) FindAll(
	ctx context.Context,
	f filters.Filters,
	search ...string,
) ([]*Category, filters.Metadata, error) {
	key := s.keyBuilder.BuildListKey(f.Page, f.Sort, f.PageSize, f.Limit(), f.Offset(), search)

	type cachePayload struct {
		Data     []*Category
		Metadata filters.Metadata
	}

	payload, err := cache.FetchOrCache(ctx, s.cache, key, func() (cachePayload, error) {
		models, metadata, err := s.repo.FindAll(ctx, f, search...)
		return cachePayload{Data: models, Metadata: metadata}, err
	})

	if err != nil {
		return nil, filters.Metadata{}, err
	}

	return payload.Data, payload.Metadata, nil
}

func (s *CategoryService) Insert(
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

	var err error
	if tx != nil {
		err = saveFn(tx)
	} else {
		err = s.tx.RunInTx(ctx, saveFn)
	}

	if err != nil {
		return err
	}

	go func() {
		_ = s.cache.DeleteByPrefix(context.Background(), s.keyBuilder.GetPrefix())
	}()

	return nil
}

func (s *CategoryService) Update(
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

	var err error
	if tx != nil {
		err = updateFn(tx)
	} else {
		err = s.tx.RunInTx(ctx, updateFn)
	}

	if err != nil {
		return err
	}

	go func() {
		_ = s.cache.DeleteByPrefix(context.Background(), s.keyBuilder.GetPrefix())
	}()
	return nil
}

func (s *CategoryService) DeleteById(
	ctx context.Context,
	id uuid.UUID,
	tx *sql.Tx,
) error {
	deleteFn := func(tx *sql.Tx) error {
		return s.repo.DeleteById(ctx, tx, id)
	}

	var err error
	if tx != nil {
		err = deleteFn(tx)
	} else {
		err = s.tx.RunInTx(ctx, deleteFn)
	}

	if err != nil {
		return err
	}

	go func() {
		_ = s.cache.DeleteByPrefix(context.Background(), s.keyBuilder.GetPrefix())
	}()

	return nil
}
