package user

import (
	"context"
	"controle_financas/internal/core/cache"
	"controle_financas/internal/core/domain/apiError"
	"controle_financas/internal/core/transaction"
	"controle_financas/internal/core/validator"

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
	Save(ctx context.Context, model *Usuario) error
	Update(ctx context.Context, model *Usuario) error
	DeleteById(ctx context.Context, id uuid.UUID) error
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

func (s *UserService) Save(ctx context.Context, model *Usuario) error {
	v := validator.New()
	if model.Validate(v); !v.Valid() {
		return apiError.NewValidationError(v.Errors)
	}

	err := s.tx.RunInTx(ctx, func(ctx context.Context) error {
		return s.repo.Save(ctx, model)
	})

	if err != nil {
		return err
	}

	go func() {
		_ = s.cache.DeleteByPrefix(context.Background(), s.keyBuilder.GetPrefix())
	}()

	return nil
}

func (s *UserService) Update(ctx context.Context, model *Usuario) error {
	v := validator.New()
	if model.Validate(v); !v.Valid() {
		return apiError.NewValidationError(v.Errors)
	}

	err := s.tx.RunInTx(ctx, func(ctx context.Context) error {
		return s.repo.Update(ctx, model)
	})

	if err != nil {
		return err
	}

	go func() {
		_ = s.cache.DeleteByPrefix(context.Background(), s.keyBuilder.GetPrefix())
	}()
	return nil
}

func (s *UserService) DeleteById(ctx context.Context, id uuid.UUID) error {
	err := s.tx.RunInTx(ctx, func(ctx context.Context) error {
		return s.repo.DeleteById(ctx, id)
	})

	if err != nil {
		return err
	}

	go func() {
		_ = s.cache.DeleteByPrefix(context.Background(), s.keyBuilder.GetPrefix())
	}()

	return nil
}
