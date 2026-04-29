package user

import (
	"context"
	e "controle_financas/internal/core/domain/errors"
	"controle_financas/internal/core/jsonlog"
	"controle_financas/internal/core/repository"
	"database/sql"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type userRepository struct {
	db     *sql.DB
	logger jsonlog.Logger
	repository.BaseRepository[Usuario]
}

type UserRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*Usuario, error)
	FindByEmail(ctx context.Context, email string) (*Usuario, error)
	Save(ctx context.Context, tx *sql.Tx, model *Usuario) error
	Update(ctx context.Context, tx *sql.Tx, model *Usuario) error
	DeleteById(ctx context.Context, tx *sql.Tx, id uuid.UUID) error
}

func NewRepository(
	db *sql.DB,
	logger jsonlog.Logger,
) UserRepository {
	return &userRepository{
		db:             db,
		logger:         logger,
		BaseRepository: repository.NewBaseRepository[Usuario](db, logger, "usuarios", "u"),
	}
}

func parseUserConstraintError(err error) error {
	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Constraint {
		case "user_email_key":
			return e.ValidationAlreadyExists("email")
		case "user_telefone_key":
			return e.ValidationAlreadyExists("telefone")
		}
	}
	return err
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*Usuario, error) {
	return r.BaseRepository.FindById(ctx, id)
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*Usuario, error) {
	return r.BaseRepository.FindOne(ctx, repository.WithQueryExtraWhere("u.email = :email", map[string]any{
		"email": email,
	}))
}

func (r *userRepository) Save(ctx context.Context, tx *sql.Tx, model *Usuario) error {
	err := r.BaseRepository.Insert(ctx, tx, model)
	if err != nil {
		return parseUserConstraintError(err)
	}

	return nil
}

func (r *userRepository) Update(ctx context.Context, tx *sql.Tx, model *Usuario) error {
	err := r.BaseRepository.Update(ctx, tx, model)
	if err != nil {
		return parseUserConstraintError(err)
	}

	return nil
}

func (r *userRepository) DeleteById(ctx context.Context, tx *sql.Tx, id uuid.UUID) error {
	return r.BaseRepository.DeleteByQuery(
		ctx,
		tx,
		repository.WithQueryExtraWhere("u.id = :id", map[string]any{
			"id": id,
		}))
}
