package api

import (
	"controle_financas/internal/core/config"
	"controle_financas/internal/core/transaction"
	"controle_financas/internal/features/auth"
	"controle_financas/internal/features/category"
	"controle_financas/internal/features/user"
)

type services struct {
	user.UserService
	auth.AuthService
	category.CategoryService
}

func NewServices(
	r *repositories,
	tx transaction.Manager,
	config config.Config,
) (*services, error) {
	userService := user.NewService(r.UserRepository, tx)
	authService, err := auth.NewService(userService, config)
	categoryService := category.NewService(r.CategoryRepository, tx)

	if err != nil {
		return nil, err
	}

	return &services{
		UserService:     userService,
		AuthService:     authService,
		CategoryService: categoryService,
	}, nil
}
