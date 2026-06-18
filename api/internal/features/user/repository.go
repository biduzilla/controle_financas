package user

import (
	"context"
	"controle_financas/internal/core/domain/apiError"
	"controle_financas/internal/core/jsonlog"
	"controle_financas/internal/core/repository"
	"database/sql"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type UserRepository struct {
	db     *sql.DB
	logger jsonlog.Logger
	br     *repository.BaseRepository[Usuario]
}

type userRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*Usuario, error)
	FindByEmail(ctx context.Context, email string) (*Usuario, error)
	Save(ctx context.Context, model *Usuario) error
	Update(ctx context.Context, model *Usuario) error
	DeleteById(ctx context.Context, id uuid.UUID) error
}

func NewRepository(
	db *sql.DB,
	logger jsonlog.Logger,
) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
		br:     repository.NewBaseRepository[Usuario](db, logger, "usuarios", "u"),
	}
}

func parseUserConstraintError(err error) error {
	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Constraint {
		case "user_email_key":
			return apiError.ValidationAlreadyExists("email")
		case "user_telefone_key":
			return apiError.ValidationAlreadyExists("telefone")
		}
	}
	return err
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*Usuario, error) {
	return r.br.FindById(ctx, id)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*Usuario, error) {
	return r.br.FindOne(ctx, repository.WithQueryExtraWhere("u.email = :email", map[string]any{
		"email": email,
	}))
}

func (r *UserRepository) Save(ctx context.Context, model *Usuario) error {
	err := r.br.Insert(ctx, model)
	if err != nil {
		return parseUserConstraintError(err)
	}

	return nil
}

func (r *UserRepository) Update(ctx context.Context, model *Usuario) error {
	err := r.br.Update(ctx, model)
	if err != nil {
		return parseUserConstraintError(err)
	}

	return nil
}

func (r *UserRepository) DeleteById(ctx context.Context, id uuid.UUID) error {
	return r.br.DeleteByQuery(
		ctx,
		repository.WithQueryExtraWhere("u.id = :id", map[string]any{
			"id": id,
		}))
}
