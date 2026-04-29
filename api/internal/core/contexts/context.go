package contexts

import (
	"context"
	"controle_financas/internal/core/security"
	"net/http"
)

type contextKey string

const userContextKey = contextKey("user")

func SetUser(r *http.Request, user security.UserDetails) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func GetUser(ctx context.Context) security.UserDetails {
	user, ok := ctx.Value(userContextKey).(security.UserDetails)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}
