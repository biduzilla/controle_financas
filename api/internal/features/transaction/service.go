package transaction

import (
	"context"
	"controle_financas/internal/core/cache"
	"controle_financas/internal/core/domain/apiError"
	"controle_financas/internal/core/filters"
	"controle_financas/internal/core/transaction"
	"controle_financas/internal/core/validator"
	"time"

	"github.com/google/uuid"
)

type TransactionService struct {
	repo       transactionRepository
	tx         transaction.Manager
	cache      cache.Cache
	keyBuilder cache.KeyBuilder
}

type transactionService interface {
	BalanceSummary(
		ctx context.Context,
		startDate, endDate *time.Time,
	) (BalanceSummary, error)

	FindById(
		ctx context.Context,
		id uuid.UUID,
	) (*Transaction, error)

	FindAll(
		ctx context.Context,
		categoryId *uuid.UUID,
		startDate, endDate *time.Time,
		f filters.Filters,
		search ...string,
	) ([]*Transaction, filters.Metadata, error)
	Insert(
		ctx context.Context,
		model *Transaction,
	) error

	Update(
		ctx context.Context,
		model *Transaction,
	) error

	DeleteById(
		ctx context.Context,
		id uuid.UUID,
	) error
}

func NewService(
	repo transactionRepository,
	tx transaction.Manager,
	cache cache.Cache,
	keyBuilder cache.KeyBuilder,
) *TransactionService {
	return &TransactionService{
		repo:       repo,
		tx:         tx,
		cache:      cache,
		keyBuilder: keyBuilder,
	}
}

func (s *TransactionService) BalanceSummary(
	ctx context.Context,
	startDate, endDate *time.Time,
) (BalanceSummary, error) {
	var endAdjusted *time.Time
	if endDate != nil {
		t := endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		endAdjusted = &t
	} else {
		endAdjusted = nil
	}
	return s.repo.GetBalanceSummary(ctx, startDate, endAdjusted)
}

func (s *TransactionService) FindById(
	ctx context.Context,
	id uuid.UUID,
) (*Transaction, error) {
	key := s.keyBuilder.BuildItemKey(id.String())

	return cache.FetchOrCache(ctx, s.cache, key, func() (*Transaction, error) {
		return s.repo.FindById(ctx, id)
	})
}

func (s *TransactionService) FindAll(
	ctx context.Context,
	categoryId *uuid.UUID,
	startDate, endDate *time.Time,
	f filters.Filters,
	search ...string,
) ([]*Transaction, filters.Metadata, error) {
	key := s.keyBuilder.BuildListKey(f.Page, f.Sort, f.PageSize, f.Limit(), f.Offset(), search)

	type cachePayload struct {
		Data     []*Transaction
		Metadata filters.Metadata
	}

	payload, err := cache.FetchOrCache(ctx, s.cache, key, func() (cachePayload, error) {
		models, metadata, err := s.repo.FindAll(ctx, categoryId, startDate, endDate, f, search...)
		return cachePayload{Data: models, Metadata: metadata}, err
	})

	if err != nil {
		return nil, filters.Metadata{}, err
	}

	return payload.Data, payload.Metadata, nil
}

func (s *TransactionService) Insert(
	ctx context.Context,
	model *Transaction,
) error {
	v := validator.New()
	if model.ValidateTransaction(v); !v.Valid() {
		return apiError.NewValidationError(v.Errors)
	}

	err := s.tx.RunInTx(ctx, func(ctx context.Context) error {
		return s.repo.Insert(ctx, model)
	})

	if err != nil {
		return err
	}

	go func() {
		_ = s.cache.DeleteByPrefix(context.Background(), s.keyBuilder.GetPrefix())
	}()

	return nil
}

func (s *TransactionService) Update(
	ctx context.Context,
	model *Transaction,
) error {
	v := validator.New()
	if model.ValidateTransaction(v); !v.Valid() {
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

func (s *TransactionService) DeleteById(
	ctx context.Context,
	id uuid.UUID,
) error {
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
