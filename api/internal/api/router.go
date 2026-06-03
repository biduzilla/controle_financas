package api

import (
	"controle_financas/internal/core/domain/apiError"
	"controle_financas/internal/core/middleware"
	"controle_financas/internal/features/auth"
	"controle_financas/internal/features/category"
	"controle_financas/internal/features/transaction"
	"controle_financas/internal/features/user"
	"database/sql"
	"expvar"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	errHandler  apiError.ErrorHandler
	m           middleware.Middleware
	user        *user.UserRouter
	auth        *auth.AuthRouter
	category    *category.CategoryRouter
	transaction *transaction.TransactionRouter
}

func NewRouter(
	handlers *handlers,
	errHandler apiError.ErrorHandler,
	m middleware.Middleware,
) *Router {
	return &Router{
		m:           m,
		errHandler:  errHandler,
		user:        user.NewRouter(handlers.UserHandler, m),
		auth:        auth.NewRouter(handlers.AuthHandler, m),
		category:    category.NewRouter(handlers.CategoyHandler, m),
		transaction: transaction.NewRouter(handlers.TransactionHandler, m),
	}
}

func (router *Router) RegisterRoutes(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()
	r.Use(router.m.RecoverPanic)
	r.Use(router.m.TimeoutMiddleWare)
	r.Use(router.m.Metrics)
	r.Use(router.m.RateLimit)
	r.Use(router.m.EnableCORS)
	r.Use(router.m.Authenticate)
	r.Use(router.m.Logging)

	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		router.errHandler.NotFoundResponse(w, req)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, req *http.Request) {
		router.errHandler.MethodNotAllowedResponse(w, req)
	})

	r.Route("/v1", func(r chi.Router) {
		r.Mount("/debug/vars", expvar.Handler())
		r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
			middleware.MetricsHandler(db).ServeHTTP(w, r)
		})
		router.user.Routes(r)
		router.auth.Routes(r)
		router.category.Routes(r)
		router.transaction.Routes(r)
	})

	return r
}
