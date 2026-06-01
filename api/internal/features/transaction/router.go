package transaction

import (
	"controle_financas/internal/core/middleware"

	"github.com/go-chi/chi/v5"
)

type TransactionRouter struct {
	handler transactionHandler
	m       middleware.Middleware
}

type transactionRouter interface {
	Routes(router chi.Router)
}

func NewRouter(
	handler transactionHandler,
	m middleware.Middleware,
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
