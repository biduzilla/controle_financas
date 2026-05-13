package contexts

import (
	"context"
	"controle_financas/internal/core/security"
	"database/sql"
	"net/http"
)

type contextKey string

const userContextKey = contextKey("user")
const txContextKey = contextKey("tx")

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

func SetTx(r *http.Request, tx *sql.Tx) *http.Request {
	ctx := context.WithValue(r.Context(), txContextKey, tx)
	return r.WithContext(ctx)
}

func GetTx(ctx context.Context) *sql.Tx {
	tx, ok := ctx.Value(txContextKey).(*sql.Tx)
	if !ok {
		panic("missing tx value in request context")
	}
	return tx
}
