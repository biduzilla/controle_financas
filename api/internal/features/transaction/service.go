package transaction

import (
	"context"
	"controle_financas/internal/core/domain/apiError"
	"controle_financas/internal/core/filters"
	"controle_financas/internal/core/transaction"
	"controle_financas/internal/core/validator"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type TransactionService struct {
	repo transactionRepository
	tx   transaction.Manager
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
) *TransactionService {
	return &TransactionService{
		repo: repo,
		tx:   tx,
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
	return s.repo.FindById(ctx, id)
}

func (s *TransactionService) FindAll(
	ctx context.Context,
	categoryId *uuid.UUID,
	startDate, endDate *time.Time,
	f filters.Filters,
	search ...string,
) ([]*Transaction, filters.Metadata, error) {
	return s.repo.FindAll(ctx, categoryId, startDate, endDate, f, search...)
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

	if tx != nil {
		return saveFn(tx)
	}
	return s.tx.RunInTx(ctx, saveFn)
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

	if tx != nil {
		return updateFn(tx)
	}
	return s.tx.RunInTx(ctx, updateFn)
}

func (s *TransactionService) DeleteById(
	ctx context.Context,
	tx *sql.Tx,
	id uuid.UUID,
) error {
	deleteFn := func(tx *sql.Tx) error {
		return s.repo.DeleteById(ctx, tx, id)
	}

	if tx != nil {
		return deleteFn(tx)
	}
	return s.tx.RunInTx(ctx, deleteFn)
}
