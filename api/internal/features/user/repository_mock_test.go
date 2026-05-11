package user

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type mockUserRepository struct {
	FindByIDFn    func(ctx context.Context, id uuid.UUID) (*Usuario, error)
	FindByEmailFn func(ctx context.Context, email string) (*Usuario, error)
	SaveFn        func(ctx context.Context, tx *sql.Tx, model *Usuario) error
	UpdateFn      func(ctx context.Context, tx *sql.Tx, model *Usuario) error
	DeleteByIdFn  func(ctx context.Context, tx *sql.Tx, id uuid.UUID) error
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*Usuario, error) {
	if m.FindByEmailFn != nil {
		return m.FindByEmailFn(ctx, email)
	}
	return nil, nil
}

func (m *mockUserRepository) Save(ctx context.Context, tx *sql.Tx, model *Usuario) error {
	if m.SaveFn != nil {
		return m.SaveFn(ctx, tx, model)
	}
	return nil
}

func (m *mockUserRepository) Update(ctx context.Context, tx *sql.Tx, model *Usuario) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, tx, model)
	}
	return nil
}

func (m *mockUserRepository) DeleteById(ctx context.Context, tx *sql.Tx, id uuid.UUID) error {
	if m.DeleteByIdFn != nil {
		return m.DeleteByIdFn(ctx, tx, id)
	}
	return nil
}
