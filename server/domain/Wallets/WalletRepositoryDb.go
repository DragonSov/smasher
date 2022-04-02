package Wallets

import (
	"github.com/google/uuid"
	"time"
)

type WalletModel struct {
	ID        uuid.UUID `db:"id"`
	Currency  string    `db:"currency"`
	Money     int       `db:"money"`
	UserID    uuid.UUID `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}

type TransactionModel struct {
	ID        uuid.UUID `db:"id"`
	Sender    uuid.UUID `db:"sender"`
	Recipient uuid.UUID `db:"recipient"`
	Money     int       `db:"money"`
	CreatedAt time.Time `db:"created_at"`
}

type WalletRepository interface {
	CreateWallet(WalletModel) (*WalletModel, error)
	CreateTransaction(TransactionModel) (*TransactionModel, error)
	SelectUserWallets(uuid.UUID) (*[]WalletModel, error)
	SelectWalletTransactions(uuid.UUID) (*[]TransactionModel, error)
	SelectWalletByUUID(uuid.UUID) (*WalletModel, error)
	UpdateWalletByUUID(WalletModel) (*WalletModel, error)
	DeleteWalletByUUID(uuid.UUID) error
}
