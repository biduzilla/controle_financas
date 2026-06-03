package api

import (
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

type mw interface {
	Metrics(next http.Handler) http.Handler
	EnableCORS(next http.Handler) http.Handler
	RequireAuthenticatedUser(next http.Handler) http.Handler
	RequireActivatedUser(next http.Handler) http.Handler
	Authenticate(next http.Handler) http.Handler
	RateLimit(next http.Handler) http.Handler
	RecoverPanic(next http.Handler) http.Handler
	Logging(next http.Handler) http.Handler
	TimeoutMiddleWare(next http.Handler) http.Handler
	RequestID(next http.Handler) http.Handler
}

type errorHandler interface {
	NotFoundResponse(w http.ResponseWriter, r *http.Request)
	MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request)
}

type Router struct {
	errHandler  errorHandler
	m           mw
	user        *user.UserRouter
	auth        *auth.AuthRouter
	category    *category.CategoryRouter
	transaction *transaction.TransactionRouter
}

func NewRouter(
	handlers *handlers,
	errHandler errorHandler,
	m mw,
) *Router {
	return &Router{
		m:           m,
		errHandler:  errHandler,
		user:        user.NewRouter(handlers.UserHandler, m),
		auth:        auth.NewRouter(handlers.AuthHandler),
		category:    category.NewRouter(handlers.CategoyHandler, m),
		transaction: transaction.NewRouter(handlers.TransactionHandler, m),
	}
}

func (router *Router) RegisterRoutes(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Use(router.m.RecoverPanic)
	r.Use(router.m.TimeoutMiddleWare)
	r.Use(router.m.RequestID)
	r.Use(router.m.Metrics)
	r.Use(router.m.Logging)

	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		router.errHandler.NotFoundResponse(w, req)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, req *http.Request) {
		router.errHandler.MethodNotAllowedResponse(w, req)
	})

	r.Get("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Handle("/metrics", middleware.MetricsHandler(db))
	r.Mount("/debug/vars", expvar.Handler())

	r.Route("/v1", func(r chi.Router) {
		r.Use(router.m.RateLimit)
		r.Use(router.m.EnableCORS)
		r.Use(router.m.Authenticate)

		router.user.Routes(r)
		router.auth.Routes(r)
		router.category.Routes(r)
		router.transaction.Routes(r)
	})

	return r
}
