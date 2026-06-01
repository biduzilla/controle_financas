package category

import (
	"controle_financas/internal/core/middleware"

	"github.com/go-chi/chi/v5"
)

type CategoryRouter struct {
	handler categoyHandler
	m       middleware.Middleware
}

type categoryRouter interface {
	Routes(router chi.Router)
}

func NewRouter(
	handler categoyHandler,
	m middleware.Middleware,
) *CategoryRouter {
	return &CategoryRouter{
		handler: handler,
		m:       m,
	}
}

func (r *CategoryRouter) Routes(router chi.Router) {
	router.Route("/category", func(router chi.Router) {
		router.Group(func(router chi.Router) {
			router.Use(r.m.RequireActivatedUser)

			router.Get("/{id}", r.handler.FindByID)
			router.Get("/", r.handler.FindAll)
			router.Post("/", r.handler.Save)
			router.Put("/", r.handler.Update)
			router.Delete("/{id}", r.handler.DeleteByID)
		})
	})
}
