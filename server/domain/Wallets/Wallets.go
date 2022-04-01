package Wallets

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type WalletRepositoryDb struct {
	Client *sqlx.DB
}

func NewWalletRepositoryDb(client *sqlx.DB) WalletRepositoryDb {
	return WalletRepositoryDb{
		Client: client,
	}
}

func (repo *WalletRepositoryDb) CreateWallet(w WalletModel) (*WalletModel, error) {
	// Inserting a new wallet
	_, err := repo.Client.NamedExec("INSERT INTO wallets (currency, money, user_id) VALUES (:currency, :money, :user_id)", w)
	if err != nil {
		return nil, err
	}

	return &w, nil
}

func (repo *WalletRepositoryDb) CreateTransaction(t TransactionModel) (*TransactionModel, error) {
	// Inserting a new transaction
	_, err := repo.Client.NamedExec("INSERT INTO transactions (sender, recipient, money) VALUES (:sender, :recipient, :money)", t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (repo *WalletRepositoryDb) SelectUserWallets(id uuid.UUID) (*[]WalletModel, error) {
	// Creating an array with wallets
	w := []WalletModel{}
	err := repo.Client.Select(&w, "SELECT * FROM wallets WHERE user_id = $1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &w, nil
}

func (repo *WalletRepositoryDb) SelectWalletByUUID(id uuid.UUID) (*WalletModel, error) {
	// Creating wallet model
	w := WalletModel{}

	// Filling the wallet model
	err := repo.Client.Get(&w, "SELECT * FROM wallets WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &w, nil
}

func (repo *WalletRepositoryDb) UpdateWalletByUUID(w WalletModel) (*WalletModel, error) {
	// Updating the wallet by ID
	_, err := repo.Client.NamedExec("UPDATE wallets SET currency = :currency, money = :money, user_id = :user_id WHERE id = :id", w)
	if err != nil {
		return nil, err
	}

	return &w, nil
}

func (repo *WalletRepositoryDb) DeleteWalletByUUID(id uuid.UUID) error {
	//Deleting the wallet by ID
	_, err := repo.Client.Exec("DELETE FROM wallets WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
