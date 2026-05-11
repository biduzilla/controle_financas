package transaction

import (
	"context"
	"database/sql"
)

type MockTransactionManager struct {
	RunInTxFn func(ctx context.Context, fn func(*sql.Tx) error) error
}

func (m *MockTransactionManager) RunInTx(ctx context.Context, fn func(*sql.Tx) error) error {
	if m.RunInTxFn != nil {
		return m.RunInTxFn(ctx, fn)
	}
	return fn(nil)
}
