package api

import (
	"controle_financas/internal/core/cache"
	"controle_financas/internal/core/config"
	"controle_financas/internal/core/jsonlog"
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
	logger jsonlog.Logger,
) (*services, error) {
	cacheClient, err := cache.NewRedisCache(config.Cache.Addr, config.Cache.Password, config.Cache.Db)
	if err != nil {
		return nil, err
	}

	logger.PrintInfo("reddis connection pool established", nil)

	userService := user.NewService(r.UserRepository, tx, cacheClient, cache.NewKeyBuilder("user"))
	authService, err := auth.NewService(userService, config)
	categoryService := category.NewService(r.CategoryRepository, tx, cacheClient, cache.NewKeyBuilder("category"))
	transactionService := t.NewService(r.TransactionRepository, tx, cacheClient, cache.NewKeyBuilder("transaction"))

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
