package services

import (
	"github.com/DragonSov/smasher/server/domain/Wallets"
	"github.com/google/uuid"
)

type DefaultWalletService struct {
	repo Wallets.WalletRepositoryDb
}

type WalletService interface {
	CreateWallet(Wallets.WalletModel) (*Wallets.WalletModel, error)
	CreateTransaction(Wallets.TransactionModel) (*Wallets.TransactionModel, error)
	SelectUserWallets(uuid.UUID) (*[]Wallets.WalletModel, error)
	SelectWalletByUUID(uuid.UUID) (*Wallets.WalletModel, error)
	SelectWalletTransactions(uuid.UUID) (*[]Wallets.TransactionModel, error)
	UpdateWalletByUUID(Wallets.WalletModel) (*Wallets.WalletModel, error)
	DeleteWalletByUUID(uuid.UUID) error
}

func NewWalletService(repository Wallets.WalletRepositoryDb) DefaultWalletService {
	return DefaultWalletService{repo: repository}
}

func (db DefaultWalletService) CreateWallet(w Wallets.WalletModel) (*Wallets.WalletModel, error) {
	// Creating wallet
	wallet, err := db.repo.CreateWallet(w)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (db DefaultWalletService) CreateTransaction(t Wallets.TransactionModel) (*Wallets.TransactionModel, error) {
	transaction, err := db.repo.CreateTransaction(t)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (db DefaultWalletService) SelectUserWallets(id uuid.UUID) (*[]Wallets.WalletModel, error) {
	// Selecting wallets
	wallets, err := db.repo.SelectUserWallets(id)
	if err != nil {
		return nil, err
	}

	return wallets, nil
}

func (db DefaultWalletService) SelectWalletByUUID(id uuid.UUID) (*Wallets.WalletModel, error) {
	// Selecting wallet by ID
	wallet, err := db.repo.SelectWalletByUUID(id)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (db DefaultWalletService) SelectWalletTransactions(id uuid.UUID) (*[]Wallets.TransactionModel, error) {
	// Selecting transactions
	transactions, err := db.repo.SelectWalletTransactions(id)
	if err != nil {
		return nil, err
	}

	return transactions, err
}

func (db DefaultWalletService) UpdateWalletByUUID(w Wallets.WalletModel) (*Wallets.WalletModel, error) {
	// Updating wallet by ID
	wallet, err := db.repo.UpdateWalletByUUID(w)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (db DefaultWalletService) DeleteWalletByUUID(id uuid.UUID) error {
	// Deleting wallet by ID
	err := db.repo.DeleteWalletByUUID(id)
	if err != nil {
		return err
	}

	return nil
}
