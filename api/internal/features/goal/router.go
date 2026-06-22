package goal

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type authMiddleware interface {
	RequireActivatedUser(next http.Handler) http.Handler
}

type Router struct {
	handler goalHandler
	m       authMiddleware
}

type goalRouter interface {
	Routes(router chi.Router)
}

func NewRouter(
	handler goalHandler,
	m authMiddleware,
) *Router {
	return &Router{
		handler: handler,
		m:       m,
	}
}

func (r *Router) Routes(router chi.Router) {
	router.Route("/goals", func(router chi.Router) {
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
