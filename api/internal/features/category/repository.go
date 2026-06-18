package category

import (
	"context"
	"controle_financas/internal/core/contexts"
	"controle_financas/internal/core/domain/apiError"
	"controle_financas/internal/core/filters"
	"controle_financas/internal/core/jsonlog"
	"controle_financas/internal/core/repository"
	"controle_financas/internal/features/user"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type CategoryRepository struct {
	db     *sql.DB
	logger jsonlog.Logger
	br     *repository.BaseRepository[Category]
}

type categoryRepository interface {
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
	) error

	Update(
		ctx context.Context,
		model *Category,
	) error

	DeleteById(
		ctx context.Context,
		id uuid.UUID,
	) error
}

func NewRepository(
	db *sql.DB,
	logger jsonlog.Logger,
) *CategoryRepository {
	return &CategoryRepository{
		db:     db,
		logger: logger,
		br:     repository.NewBaseRepository[Category](db, logger, "categories", "c"),
	}
}

func (r *CategoryRepository) parseConstraintError(err error) error {
	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Constraint {
		case "categories_name_user_id_unique":
			return apiError.ValidationAlreadyExists("nome")
		}
	}
	return err
}

func (r *CategoryRepository) FindById(
	ctx context.Context,
	id uuid.UUID,
) (*Category, error) {
	userLogado := contexts.GetUser(ctx)

	return r.br.FindOne(
		ctx,
		repository.WithQueryExtraWhere(`
		c.id = :id
		and c.user_id = :userId
	`, map[string]any{
			"id":     id,
			"userId": userLogado.GetID(),
		}),
		repository.WithJoin(
			user.Usuario{},
			"usuarios",
			"u",
			"c.user_id = u.id",
			nil,
		),
	)
}

func (r *CategoryRepository) FindAll(
	ctx context.Context,
	f filters.Filters,
	search ...string,
) ([]*Category, filters.Metadata, error) {
	userLogado := contexts.GetUser(ctx)
	query := fmt.Sprintf(`
		%s
		and user_id = :userId
	`, repository.BuildFilterQuery("c", search...))

	return r.br.FindWithFilters(
		ctx,
		f,
		repository.WithQueryExtraWhere(query,
			map[string]any{
				"userId": userLogado.GetID(),
			}),
		repository.WithJoin(
			user.Usuario{},
			"usuarios",
			"u",
			"c.user_id = u.id",
			nil,
		),
	)
}

func (r *CategoryRepository) Insert(
	ctx context.Context,
	model *Category,
) error {
	userLogado := contexts.GetUser(ctx)
	err := r.br.Insert(
		ctx,
		model,
		repository.WithExtraFields([]string{"user_id"}, map[string]any{
			"user_id": userLogado.GetID(),
		}),
	)

	if err != nil {
		return r.parseConstraintError(err)
	}

	return nil
}

func (r *CategoryRepository) Update(
	ctx context.Context,
	model *Category,
) error {
	userAuth := contexts.GetUser(ctx)
	err := r.br.Update(
		ctx,
		model,
		repository.WithExtraFields([]string{"user_id"}, map[string]any{
			"user_id": model.User.ID,
		}),
		repository.WithExtraWhere("user_id = :userId", map[string]any{
			"userId": userAuth.GetID(),
		}),
	)

	if err != nil {
		return r.parseConstraintError(err)
	}

	return nil
}

func (r *CategoryRepository) DeleteById(
	ctx context.Context,
	id uuid.UUID,
) error {
	userAuth := contexts.GetUser(ctx)

	return r.br.DeleteByQuery(
		ctx,
		repository.WithQueryExtraWhere("id = :id and user_id = :userId",
			map[string]any{
				"id":     id,
				"userId": userAuth.GetID(),
			}),
	)
}
