package goal

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

type Repository struct {
	db     *sql.DB
	logger jsonlog.Logger
	br     repository.Orm[Goal]
}

type goalRepository interface {
	FindAllByUserId(
		ctx context.Context,
		name string,
		f filters.Filters,
	) ([]*Goal, filters.Metadata, error)

	FindById(
		ctx context.Context,
		id uuid.UUID,
	) (*Goal, error)

	Insert(
		ctx context.Context,
		model *Goal,
	) error

	Update(
		ctx context.Context,
		model *Goal,
	) error
	DeleteById(
		ctx context.Context,
		id uuid.UUID,
	) error
}

func (r *Repository) parseConstraintError(err error) error {
	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Constraint {
		case "unique_user_goal_name":
			return apiError.ValidationAlreadyExists("nome")
		}
	}
	return err
}

func NewRepository(
	db *sql.DB,
	logger jsonlog.Logger,
) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
		br:     repository.NewBaseRepository[Goal](db, logger, "goals", "g"),
	}
}

func (r *Repository) baseQueryOption(ctx context.Context, query string) []repository.QueryOption {
	userAuth := contexts.GetUser(ctx)

	return []repository.QueryOption{
		repository.WithQueryExtraWhere(query,
			map[string]any{
				"userId": userAuth.GetID(),
			}),
		repository.WithJoin(
			user.Usuario{},
			"usuarios",
			"u",
			"c.user_id = u.id",
			nil,
		),
	}
}

func (r *Repository) FindAllByUserId(
	ctx context.Context,
	name string,
	f filters.Filters,
) ([]*Goal, filters.Metadata, error) {

	query := fmt.Sprintf(`
		%s
		and user_id = :userId
	`, repository.BuildFilterQuery("c", name))

	return r.br.FindWithFilters(
		ctx,
		f,
		r.baseQueryOption(ctx, query)...,
	)
}

func (r *Repository) FindById(
	ctx context.Context,
	id uuid.UUID,
) (*Goal, error) {
	return r.br.FindById(
		ctx,
		id,
		r.baseQueryOption(ctx, "user_id = :userId")...,
	)
}

func (r *Repository) Insert(
	ctx context.Context,
	model *Goal,
) error {
	userAuth := contexts.GetUser(ctx)

	err := r.br.Insert(
		ctx,
		model,
		repository.WithExtraFields([]string{"user_id"}, map[string]any{
			"user_id": userAuth.GetID(),
		}),
	)

	if err != nil {
		return r.parseConstraintError(err)
	}

	return nil
}

func (r *Repository) Update(
	ctx context.Context,
	model *Goal,
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

func (r *Repository) DeleteById(
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
