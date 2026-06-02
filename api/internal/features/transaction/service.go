package transaction

import (
	"context"
	"controle_financas/internal/core/cache"
	"controle_financas/internal/core/domain/apiError"
	"controle_financas/internal/core/filters"
	"controle_financas/internal/core/transaction"
	"controle_financas/internal/core/validator"
	"database/sql"
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
		tx *sql.Tx,
		model *Transaction,
	) error

	Update(
		ctx context.Context,
		tx *sql.Tx,
		model *Transaction,
	) error

	DeleteById(
		ctx context.Context,
		tx *sql.Tx,
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
	tx *sql.Tx,
	model *Transaction,
) error {
	saveFn := func(tx *sql.Tx) error {
		v := validator.New()
		if model.ValidateTransaction(v); !v.Valid() {
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

func (s *TransactionService) Update(
	ctx context.Context,
	tx *sql.Tx,
	model *Transaction,
) error {
	updateFn := func(tx *sql.Tx) error {
		v := validator.New()
		if model.ValidateTransaction(v); !v.Valid() {
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

func (s *TransactionService) DeleteById(
	ctx context.Context,
	tx *sql.Tx,
	id uuid.UUID,
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
