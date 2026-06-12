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

func (r *Repository) FindAllByUserId(
	ctx context.Context,
	name string,
	f filters.Filters,
) ([]*Goal, filters.Metadata, error) {
	userAuth := contexts.GetUser(ctx)

	query := fmt.Sprintf(`
		%s
		and user_id = :userId
	`, repository.BuildFilterQuery("c", name))

	return r.br.FindWithFilters(
		ctx,
		f,
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
	)
}

func (r *Repository) FindById(
	ctx context.Context,
	id uuid.UUID,
) (*Goal, error) {
	userAuth := contexts.GetUser(ctx)

	return r.br.FindById(
		ctx,
		id,
		repository.WithJoin(
			user.Usuario{},
			"usuarios",
			"u",
			"c.user_id = u.id",
			nil,
		),
		repository.WithQueryExtraWhere("user_id = :userId", map[string]any{
			"userId": userAuth.GetID(),
		}),
	)
}

func (r *Repository) Insert(
	ctx context.Context,
	tx *sql.Tx,
	model *Goal,
) error {
	userAuth := contexts.GetUser(ctx)

	err := r.br.Insert(
		ctx,
		tx,
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
	tx *sql.Tx,
	model *Goal,
) error {
	userAuth := contexts.GetUser(ctx)
	err := r.br.Update(
		ctx,
		tx,
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
