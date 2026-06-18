package contexts

import (
	"context"
	"controle_financas/internal/core/security"
	"database/sql"
)

type contextKey string

const userContextKey = contextKey("user")
const txContextKey = contextKey("tx")
const requestIDKey = contextKey("request_id")

func SetRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}

func SetUser(ctx context.Context, user security.UserDetails) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func GetUser(ctx context.Context) security.UserDetails {
	user, ok := ctx.Value(userContextKey).(security.UserDetails)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}

func SetTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txContextKey, tx)
}

func GetTx(ctx context.Context) *sql.Tx {
	tx, _ := ctx.Value(txContextKey).(*sql.Tx)
	return tx
}
