package auth

import (
	"controle_financas/internal/core/middleware"

	"github.com/go-chi/chi/v5"
)

type authRouter struct {
	handler AuthHandler
	m       middleware.Middleware
}

type AuthRouter interface {
	Routes(router chi.Router)
}

func NewRouter(
	handler AuthHandler,
	m middleware.Middleware,
) AuthRouter {
	return &authRouter{
		handler: handler,
		m:       m,
	}
}

func (r *authRouter) Routes(router chi.Router) {
	router.Route("/auth", func(router chi.Router) {
		router.Post("/", r.handler.Login)
		router.Post("/refresh-token", r.handler.RefreshToken)
	})
}
