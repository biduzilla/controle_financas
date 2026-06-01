package api

import (
	"controle_financas/internal/core/jsonlog"
	"controle_financas/internal/features/category"
	"controle_financas/internal/features/transaction"
	"controle_financas/internal/features/user"
	"database/sql"
)

type repositories struct {
	*user.UserRepository
	*category.CategoryRepository
	*transaction.TransactionRepository
}

func NewRepositories(
	db *sql.DB,
	logger jsonlog.Logger,
) *repositories {
	return &repositories{
		UserRepository:        user.NewRepository(db, logger),
		CategoryRepository:    category.NewRepository(db, logger),
		TransactionRepository: transaction.NewRepository(db, logger),
	}
}
