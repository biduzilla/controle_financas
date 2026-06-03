package transaction

import (
	"context"
	"controle_financas/internal/core/contexts"
	"controle_financas/internal/core/domain/apiError"
	"controle_financas/internal/core/filters"
	"controle_financas/internal/core/jsonlog"
	"controle_financas/internal/core/repository"
	"controle_financas/internal/features/category"
	"controle_financas/internal/features/user"
	"controle_financas/pkg/sqlformat"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type TransactionRepository struct {
	db     *sql.DB
	logger jsonlog.Logger
	br     *repository.BaseRepository[Transaction]
}

type transactionRepository interface {
	GetBalanceSummary(
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

func NewRepository(
	db *sql.DB,
	logger jsonlog.Logger,
) *TransactionRepository {
	return &TransactionRepository{
		db:     db,
		logger: logger,
		br: repository.NewBaseRepository[Transaction](
			db,
			logger,
			"transactions",
			"t",
		),
	}
}

func (r TransactionRepository) GetBalanceSummary(
	ctx context.Context,
	startDate, endDate *time.Time,
) (BalanceSummary, error) {
	var summary BalanceSummary

	start := sql.NullTime{}
	if startDate != nil {
		start.Valid = true
		start.Time = *startDate
	}
	end := sql.NullTime{}
	if endDate != nil {
		end.Valid = true
		end.Time = *endDate
	}

	query := `
        SELECT
            COALESCE(SUM(CASE WHEN c.type = 1 THEN t.amount ELSE 0 END), 0) as total_income,
            COALESCE(SUM(CASE WHEN c.type = 2 THEN t.amount ELSE 0 END), 0) as total_expense
        FROM transactions t
        INNER JOIN categories c ON t.category_id = c.id
        WHERE t.user_id = :userId
          AND (:startDate::timestamptz IS NULL OR t.created_at >= :startDate::timestamptz)
          AND (:endDate::timestamptz IS NULL OR t.created_at <= :endDate::timestamptz)
    `

	params := map[string]any{
		"userId":    contexts.GetUser(ctx).GetID(),
		"startDate": start,
		"endDate":   end,
	}

	query, args := repository.NamedQuery(query, params)
	r.logger.PrintInfo(sqlformat.MinifySQL(query), nil)

	err := r.db.QueryRowContext(ctx, query, args).Scan(
		&summary.TotalIncome, &summary.TotalExpense,
	)

	return summary, err
}

func (r TransactionRepository) FindAll(
	ctx context.Context,
	categoryId *uuid.UUID,
	startDate, endDate *time.Time,
	f filters.Filters,
	search ...string,
) ([]*Transaction, filters.Metadata, error) {
	userAuth := contexts.GetUser(ctx)

	start := sql.NullTime{}
	if startDate != nil {
		start.Valid = true
		start.Time = *startDate
	}

	end := sql.NullTime{}
	if endDate != nil {
		end.Valid = true
		end.Time = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}

	query := `
	AND t.user_id = :userId
	AND (:categoryId = 0 OR t.category_id = :categoryId)
	AND (:startDate::timestamptz IS NULL OR t.created_at >= :startDate::timestamptz)
	AND (:endDate::timestamptz IS NULL OR t.created_at <= :endDate::timestamptz)
	`

	params := map[string]any{
		"userId":     userAuth.GetID(),
		"categoryId": categoryId,
		"startDate":  start,
		"endDate":    end,
	}

	joinType := "INNER"
	return r.br.FindWithFilters(
		ctx,
		f,
		repository.WithQueryExtraWhere(query, params),
		repository.WithJoin(
			category.Category{},
			"categories",
			"c",
			"t.category_id = c.id",
			&joinType,
		),
		repository.WithJoin(
			user.Usuario{},
			"usuarios",
			"u",
			"t.user_id = u.id",
			nil,
		),
	)
}

func (r TransactionRepository) FindById(
	ctx context.Context,
	id uuid.UUID,
) (*Transaction, error) {
	userAuth := contexts.GetUser(ctx)
	joinType := "INNER"

	return r.br.FindById(
		ctx,
		id,
		repository.WithQueryExtraWhere(`
			t.user_id = :userId
		`, map[string]any{
			"userId": userAuth.GetID(),
		}),
		repository.WithJoin(
			category.Category{},
			"categories",
			"c",
			"t.category_id = c.id",
			&joinType,
		),
		repository.WithJoin(
			user.Usuario{},
			"usuarios",
			"u",
			"t.user_id = u.id",
			nil,
		),
	)
}

func (r TransactionRepository) Insert(
	ctx context.Context,
	tx *sql.Tx,
	model *Transaction,
) error {
	userAuth := contexts.GetUser(ctx)
	if model.Category == nil {
		return apiError.ErrRecordNotFound
	}

	return r.br.Insert(
		ctx,
		tx,
		model,
		repository.WithExtraFields([]string{"user_id", "category_id"}, map[string]any{
			"user_id":     userAuth.GetID(),
			"category_id": model.Category.ID,
		}),
	)
}

func (r TransactionRepository) Update(
	ctx context.Context,
	tx *sql.Tx,
	model *Transaction,
) error {
	userAuth := contexts.GetUser(ctx)
	if model.Category == nil {
		return apiError.ErrRecordNotFound
	}

	return r.br.Update(
		ctx,
		tx,
		model,
		repository.WithExtraFields([]string{"user_id", "category_id"}, map[string]any{
			"user_id":     userAuth.GetID(),
			"category_id": model.Category.ID,
		}),
		repository.WithExtraWhere("user_id = :userId", map[string]any{
			"userId": userAuth.GetID(),
		}),
	)
}

func (r TransactionRepository) DeleteById(
	ctx context.Context,
	tx *sql.Tx,
	id uuid.UUID,
) error {
	userAuth := contexts.GetUser(ctx)

	return r.br.DeleteByQuery(
		ctx,
		tx,
		repository.WithQueryExtraWhere("id = :id and user_id = :userId",
			map[string]any{
				"id":     id,
				"userId": userAuth.GetID(),
			}),
	)
}
