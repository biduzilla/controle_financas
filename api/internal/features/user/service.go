package user

import (
	"context"
	"controle_financas/internal/core/cache"
	"controle_financas/internal/core/domain/apiError"
	"controle_financas/internal/core/transaction"
	"controle_financas/internal/core/validator"
	"database/sql"

	"github.com/google/uuid"
)

type UserService struct {
	repo       userRepository
	tx         transaction.Manager
	cache      cache.Cache
	keyBuilder cache.KeyBuilder
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
	cache cache.Cache,
	keyBuilder cache.KeyBuilder,
) *UserService {
	return &UserService{
		repo:       repo,
		tx:         tx,
		cache:      cache,
		keyBuilder: keyBuilder,
	}
}

func (s *UserService) FindById(ctx context.Context, id uuid.UUID) (*Usuario, error) {
	key := s.keyBuilder.BuildItemKey(id.String())

	return cache.FetchOrCache(ctx, s.cache, key, func() (*Usuario, error) {
		return s.repo.FindByID(ctx, id)
	})
}

func (s *UserService) FindByEmail(ctx context.Context, email string) (*Usuario, error) {
	key := s.keyBuilder.BuildItemKey(email)

	return cache.FetchOrCache(ctx, s.cache, key, func() (*Usuario, error) {
		return s.repo.FindByEmail(ctx, email)
	})
}

func (s *UserService) Save(ctx context.Context, model *Usuario, tx *sql.Tx) error {
	saveFn := func(tx *sql.Tx) error {
		v := validator.New()
		if model.Validate(v); !v.Valid() {
			return apiError.NewValidationError(v.Errors)
		}

		return s.repo.Save(ctx, tx, model)
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

func (s *UserService) Update(ctx context.Context, model *Usuario, tx *sql.Tx) error {
	updateFn := func(tx *sql.Tx) error {
		v := validator.New()
		if model.Validate(v); !v.Valid() {
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

func (s *UserService) DeleteById(ctx context.Context, id uuid.UUID, tx *sql.Tx) error {
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
