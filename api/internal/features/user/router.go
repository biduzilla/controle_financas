package user

import (
	"controle_financas/internal/core/middleware"

	"github.com/go-chi/chi/v5"
)

type userRouter struct {
	handler UserHandler
	m       middleware.Middleware
}

type UserRouter interface {
	Routes(router chi.Router)
}

func NewRouter(
	handler UserHandler,
	m middleware.Middleware,
) UserRouter {
	return &userRouter{
		handler: handler,
		m:       m,
	}
}

func (r *userRouter) Routes(router chi.Router) {
	router.Route("/user", func(router chi.Router) {
		router.Post("/", r.handler.Save)

		router.Group(func(router chi.Router) {
			router.Use(r.m.RequireActivatedUser)

			router.Get("/data", r.handler.FindAuthUserData)
			router.Get("/{id}", r.handler.FindByID)
			router.Put("/", r.handler.Update)
			router.Delete("/{id}", r.handler.Delete)
		})
	})
}
