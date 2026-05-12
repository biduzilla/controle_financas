package api

import (
	"controle_financas/internal/core/domain/errors"
	"controle_financas/internal/core/middleware"
	"controle_financas/internal/features/auth"
	"controle_financas/internal/features/category"
	"controle_financas/internal/features/user"
	"expvar"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type router struct {
	errHandler errors.ErrorHandler
	m          middleware.Middleware
	user       user.UserRouter
	auth       auth.AuthRouter
	category   category.CategoryRouter
}

type Router interface {
	RegisterRoutes() *chi.Mux
}

func NewRouter(
	handlers *handlers,
	errHandler errors.ErrorHandler,
	m middleware.Middleware,
) Router {
	return &router{
		m:          m,
		errHandler: errHandler,
		user:       user.NewRouter(handlers.UserHandler, m),
		auth:       auth.NewRouter(handlers.AuthHandler, m),
		category:   category.NewRouter(handlers.CategoyHandler, m),
	}
}

func (router *router) RegisterRoutes() *chi.Mux {
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
		router.user.Routes(r)
		router.auth.Routes(r)
		router.category.Routes(r)
	})

	return r
}
