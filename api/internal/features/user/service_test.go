package user

import (
	"testing"

	"github.com/google/uuid"
)

func newTestUsuario(email, telefone string) *Usuario {
	return &Usuario{
		ID:       uuid.New(),
		Email:    email,
		Telefone: telefone,
		Nome:     "Test User",
	}
}

func TestUserService_Save(t *testing.T) {
	t.Parallel()

	// tests := []struct {
	// 	name          string
	// 	model         *Usuario
	// 	mockRepoFn    func() *mockUserRepository
	// 	mockTxFn      func() *transaction.MockTransactionManager
	// 	tx            *sql.Tx
	// 	expectedError error
	// }{}
}
