package transaction

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type authMiddleware interface {
	RequireActivatedUser(next http.Handler) http.Handler
}

type TransactionRouter struct {
	handler transactionHandler
	m       authMiddleware
}

type transactionRouter interface {
	Routes(router chi.Router)
}

func NewRouter(
	handler transactionHandler,
	m authMiddleware,
) *TransactionRouter {
	return &TransactionRouter{
		handler: handler,
		m:       m,
	}
}

func (r *TransactionRouter) Routes(router chi.Router) {
	router.Route("/transactions", func(router chi.Router) {
		router.Group(func(router chi.Router) {
			router.Use(r.m.RequireActivatedUser)

			router.Get("/", r.handler.FindAll)
			router.Get("/{id}", r.handler.FindByID)
			router.Post("/", r.handler.Save)
			router.Put("/", r.handler.Update)
			router.Delete("/{id}", r.handler.DeleteByID)

			router.Get("/summary", r.handler.BalanceSummary)
		})
	})
}
