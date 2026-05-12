package category

import (
	"controle_financas/internal/core/middleware"

	"github.com/go-chi/chi/v5"
)

type categoryRouter struct {
	handler CategoyHandler
	m       middleware.Middleware
}

type CategoryRouter interface {
	Routes(router chi.Router)
}

func NewRouter(
	handler CategoyHandler,
	m middleware.Middleware,
) CategoryRouter {
	return &categoryRouter{
		handler: handler,
		m:       m,
	}
}

func (r *categoryRouter) Routes(router chi.Router) {
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
