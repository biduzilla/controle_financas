package api

import (
	"controle_financas/internal/core/config"
	"controle_financas/internal/core/transaction"
	"controle_financas/internal/features/auth"
	"controle_financas/internal/features/category"
	t "controle_financas/internal/features/transaction"
	"controle_financas/internal/features/user"
)

type services struct {
	*user.UserService
	*auth.AuthService
	*category.CategoryService
	*t.TransactionService
}

func NewServices(
	r *repositories,
	tx transaction.Manager,
	config config.Config,
) (*services, error) {
	userService := user.NewService(r.UserRepository, tx)
	authService, err := auth.NewService(userService, config)
	categoryService := category.NewService(r.CategoryRepository, tx)
	transactionService := t.NewService(r.TransactionRepository, tx)

	if err != nil {
		return nil, err
	}

	return &services{
		UserService:        userService,
		AuthService:        authService,
		CategoryService:    categoryService,
		TransactionService: transactionService,
	}, nil
}
