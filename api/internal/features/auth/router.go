package auth

import (
	"controle_financas/internal/core/middleware"

	"github.com/go-chi/chi/v5"
)

type AuthRouter struct {
	handler authHandler
	m       middleware.Middleware
}

type authRouter interface {
	Routes(router chi.Router)
}

func NewRouter(
	handler authHandler,
	m middleware.Middleware,
) *AuthRouter {
	return &AuthRouter{
		handler: handler,
		m:       m,
	}
}

func (r *AuthRouter) Routes(router chi.Router) {
	router.Route("/auth", func(router chi.Router) {
		router.Post("/", r.handler.Login)
		router.Post("/refresh-token", r.handler.RefreshToken)
	})
}
