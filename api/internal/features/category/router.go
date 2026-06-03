package category

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type authMiddleware interface {
	RequireActivatedUser(next http.Handler) http.Handler
}

type CategoryRouter struct {
	handler categoyHandler
	m       authMiddleware
}

type categoryRouter interface {
	Routes(router chi.Router)
}

func NewRouter(
	handler categoyHandler,
	m authMiddleware,
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
