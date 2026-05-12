package api

import (
	"controle_financas/internal/core/domain/errors"
	"controle_financas/internal/features/auth"
	"controle_financas/internal/features/category"
	"controle_financas/internal/features/user"
)

type handlers struct {
	services *services
	user.UserHandler
	auth.AuthHandler
	category.CategoyHandler
}

func NewHandlers(
	services *services,
	errHandler errors.ErrorHandler,
) *handlers {
	return &handlers{
		services:       services,
		UserHandler:    user.NewHandler(services.UserService, errHandler),
		AuthHandler:    auth.NewHandler(services.AuthService, errHandler),
		CategoyHandler: category.NewHandler(services.CategoryService, errHandler),
	}
}
