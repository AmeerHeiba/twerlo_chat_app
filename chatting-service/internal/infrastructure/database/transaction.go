package database

import (
	"context"

	"github.com/AmeerHeiba/chatting-service/internal/domain"
	"gorm.io/gorm"
)

type GormTransactionManager struct {
	db *gorm.DB
}

// NewTransactionManager creates a new TransactionManager that wraps a Gorm DB
// instance. The TransactionManager can be used to execute database operations
// within a transaction. If the operations are successful, the transaction is
// committed. If an error is encountered, the transaction is rolled back.
// should only be used for transaction where multiple operations are executed
func NewTransactionManager(db *gorm.DB) domain.TransactionManager {
	return &GormTransactionManager{db: db}
}

func (m *GormTransactionManager) WithTransaction(
	ctx context.Context,
	fn func(context.Context, *domain.Repositories) error,
) error {
	// Start transaction with context
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		repos := &domain.Repositories{
			//TODO
			// Users:             NewGormUserRepository(tx),
			// Messages:          NewGormMessageRepository(tx),
			// MessageRecipients: NewGormMessageRecipientRepository(tx),
		}
		return fn(ctx, repos) // Propagate context to the callback
	})
}
