package category

import (
	"context"
	"controle_financas/internal/core/contexts"
	e "controle_financas/internal/core/domain/errors"
	"controle_financas/internal/core/filters"
	"controle_financas/internal/core/jsonlog"
	"controle_financas/internal/core/repository"
	"controle_financas/internal/features/user"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type categoryRepository struct {
	db             *sql.DB
	logger         jsonlog.Logger
	baseRepository repository.BaseRepository[Category]
}

type CategoryRepository interface {
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
		tx *sql.Tx,
		model *Category,
	) error

	Update(
		ctx context.Context,
		tx *sql.Tx,
		model *Category,
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
) CategoryRepository {
	return &categoryRepository{
		db:             db,
		logger:         logger,
		baseRepository: repository.NewBaseRepository[Category](db, logger, "categories", "c"),
	}
}

func (r *categoryRepository) parseConstraintError(err error) error {
	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Constraint {
		case "categories_nome_key":
			return e.ValidationAlreadyExists("nome")
		}
	}
	return err
}

func (r *categoryRepository) FindById(
	ctx context.Context,
	id uuid.UUID,
) (*Category, error) {
	userLogado := contexts.GetUser(ctx)

	return r.baseRepository.FindOne(
		ctx,
		repository.WithQueryExtraWhere(`
		id = :id
		and user_id = :userId
	`, map[string]any{
			"id":     id,
			"userId": userLogado.GetID(),
		}),
		repository.WithJoin(
			user.Usuario{}, "users", "u", "c.user_id = u.id",
		),
	)
}

func (r *categoryRepository) FindAll(
	ctx context.Context,
	f filters.Filters,
	search ...string,
) ([]*Category, filters.Metadata, error) {
	userLogado := contexts.GetUser(ctx)
	query := fmt.Sprintf(`
		%s
		and user_id = :userId
	`, repository.BuildFilterQuery("c", search...))

	return r.baseRepository.FindWithFilters(
		ctx,
		f,
		repository.WithQueryExtraWhere(query,
			map[string]any{
				"userId": userLogado.GetID(),
			}),
		repository.WithJoin(
			user.Usuario{}, "users", "u", "c.user_id = u.id",
		),
	)
}

func (r *categoryRepository) Insert(
	ctx context.Context,
	tx *sql.Tx,
	model *Category,
) error {
	err := r.baseRepository.Insert(
		ctx,
		tx,
		model,
		repository.WithExtraFields([]string{"user_id"}, map[string]any{
			"tutor_id": model.User.ID,
		}),
	)

	if err != nil {
		return r.parseConstraintError(err)
	}

	return nil
}

func (r *categoryRepository) Update(
	ctx context.Context,
	tx *sql.Tx,
	model *Category,
) error {
	userAuth := contexts.GetUser(ctx)
	err := r.baseRepository.Update(
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

func (r *categoryRepository) DeleteById(
	ctx context.Context,
	tx *sql.Tx,
	id uuid.UUID,
) error {
	userAuth := contexts.GetUser(ctx)

	return r.baseRepository.DeleteByQuery(
		ctx,
		tx,
		repository.WithQueryExtraWhere("id = :id and user_id = :userId",
			map[string]any{
				"id":     id,
				"userId": userAuth.GetID(),
			}),
	)
}
