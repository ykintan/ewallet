package service

import (
	"context"
	"ewallet/wallet/entity"
	"fmt"
)

// ITransactionService defines the interface for transaction services
type ITransactionService interface {
	CreateTransaction(ctx context.Context, transaction *entity.Transaction) (entity.Transaction, error)
	GetTransaction(ctx context.Context, id int32) (entity.Transaction, error)
	CreateWallet(ctx context.Context, wallet *entity.Wallet) (entity.Wallet, error)
	TransferWallet(ctx context.Context, fromWalletID, toWalletID int, amount float64) error
	TopUp(ctx context.Context, walletID int, amount float64) error
	Payment(ctx context.Context, walletID int, amount float64) error
	GetWalletByID(ctx context.Context, walletID int) (entity.Wallet, error)
	GetWalletByUserID(ctx context.Context, userID int) (entity.Wallet, error)
	GetTransactionByUserID(ctx context.Context, userID int) ([]entity.Transaction, error)
}

// ITransactionRepository defines the interface for transaction repositories
type ITransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction *entity.Transaction) (entity.Transaction, error)
	GetTransaction(ctx context.Context, id int32) (entity.Transaction, error)
	CreateWallet(ctx context.Context, wallet *entity.Wallet) (entity.Wallet, error)
	GetWalletByID(ctx context.Context, walletID int) (entity.Wallet, error)
	UpdateWallet(ctx context.Context, wallet *entity.Wallet) error
	GetWalletByUserID(ctx context.Context, userID int) (entity.Wallet, error)
	GetTransactionByUserID(ctx context.Context, userID int) ([]entity.Transaction, error)
}

// transactionService is the implementation of ITransactionService that uses ITransactionRepository
type transactionService struct {
	transactionRepo ITransactionRepository
}

// NewTransactionService creates a new instance of transactionService
func NewTransactionService(repo ITransactionRepository) ITransactionService {
	return &transactionService{transactionRepo: repo}
}

// CreateTransaction creates a new transaction
func (s *transactionService) CreateTransaction(ctx context.Context, transaction *entity.Transaction) (entity.Transaction, error) {
	createdTransaction, err := s.transactionRepo.CreateTransaction(ctx, transaction)
	if err != nil {
		return entity.Transaction{}, fmt.Errorf("failed to create transaction: %v", err)
	}
	return createdTransaction, nil
}

// CreateWallet creates a new wallet
func (s *transactionService) CreateWallet(ctx context.Context, wallet *entity.Wallet) (entity.Wallet, error) {
	createdWallet, err := s.transactionRepo.CreateWallet(ctx, wallet)
	if err != nil {
		return entity.Wallet{}, fmt.Errorf("failed to create wallet: %v", err)
	}
	return createdWallet, nil
}

// GetTransaction retrieves a transaction by its ID
func (s *transactionService) GetTransaction(ctx context.Context, id int32) (entity.Transaction, error) {
	transaction, err := s.transactionRepo.GetTransaction(ctx, id)
	if err != nil {
		return entity.Transaction{}, fmt.Errorf("failed to get transaction: %v", err)
	}
	return transaction, nil
}

// TransferWallet transfers funds from one wallet to another
func (s *transactionService) TransferWallet(ctx context.Context, fromWalletID, toWalletID int, amount float64) error {
	fromWallet, err := s.transactionRepo.GetWalletByID(ctx, fromWalletID)
	if err != nil {
		return fmt.Errorf("failed to retrieve source wallet: %v", err)
	}

	toWallet, err := s.transactionRepo.GetWalletByID(ctx, toWalletID)
	if err != nil {
		return fmt.Errorf("failed to retrieve destination wallet: %v", err)
	}

	if fromWallet.Balance < amount {
		return fmt.Errorf("insufficient funds in source wallet")
	}

	fromWallet.Balance -= amount
	toWallet.Balance += amount

	err = s.transactionRepo.UpdateWallet(ctx, &fromWallet)
	if err != nil {
		return fmt.Errorf("failed to update source wallet: %v", err)
	}

	err = s.transactionRepo.UpdateWallet(ctx, &toWallet)
	if err != nil {
		return fmt.Errorf("failed to update destination wallet: %v", err)
	}

	transactionOut := &entity.Transaction{
		WalletID:        fromWalletID,
		WalletIDSource:  toWalletID,
		Amount:          amount,
		TransactionType: "out",
	}
	_, err = s.transactionRepo.CreateTransaction(ctx, transactionOut)
	if err != nil {
		return fmt.Errorf("failed to create transaction record for source wallet: %v", err)
	}

	transactionIn := &entity.Transaction{
		WalletID:        toWalletID,
		WalletIDSource:  fromWalletID,
		Amount:          amount,
		TransactionType: "in",
	}
	_, err = s.transactionRepo.CreateTransaction(ctx, transactionIn)
	if err != nil {
		return fmt.Errorf("failed to create transaction record for destination wallet: %v", err)
	}

	return nil
}

// TopUp adds funds to a wallet and creates an "in" transaction
func (s *transactionService) TopUp(ctx context.Context, walletID int, amount float64) error {
	wallet, err := s.transactionRepo.GetWalletByID(ctx, walletID)
	if err != nil {
		return fmt.Errorf("failed to retrieve wallet: %v", err)
	}

	wallet.Balance += amount

	err = s.transactionRepo.UpdateWallet(ctx, &wallet)
	if err != nil {
		return fmt.Errorf("failed to update wallet: %v", err)
	}

	transaction := &entity.Transaction{
		WalletID:        walletID,
		WalletIDSource:  0,
		Amount:          amount,
		TransactionType: "in",
	}
	_, err = s.transactionRepo.CreateTransaction(ctx, transaction)
	if err != nil {
		return fmt.Errorf("failed to create transaction record for top-up: %v", err)
	}

	return nil
}

// Payment deducts funds from a wallet and creates an "out" transaction
func (s *transactionService) Payment(ctx context.Context, walletID int, amount float64) error {
	wallet, err := s.transactionRepo.GetWalletByID(ctx, walletID)
	if err != nil {
		return fmt.Errorf("failed to retrieve wallet: %v", err)
	}

	if wallet.Balance < amount {
		return fmt.Errorf("insufficient funds in wallet")
	}

	wallet.Balance -= amount

	err = s.transactionRepo.UpdateWallet(ctx, &wallet)
	if err != nil {
		return fmt.Errorf("failed to update wallet: %v", err)
	}

	transaction := &entity.Transaction{
		WalletID:        walletID,
		WalletIDSource:  0,
		Amount:          amount,
		TransactionType: "out",
	}
	_, err = s.transactionRepo.CreateTransaction(ctx, transaction)
	if err != nil {
		return fmt.Errorf("failed to create transaction record for payment: %v", err)
	}

	return nil
}

// GetWalletByID retrieves a wallet by its ID
func (s *transactionService) GetWalletByID(ctx context.Context, walletID int) (entity.Wallet, error) {
	wallet, err := s.transactionRepo.GetWalletByID(ctx, walletID)
	if err != nil {
		return entity.Wallet{}, fmt.Errorf("failed to get wallet: %v", err)
	}
	return wallet, nil
}

// GetWalletByUserID retrieves wallets by user ID
func (s *transactionService) GetWalletByUserID(ctx context.Context, userID int) (entity.Wallet, error) {
	wallets, err := s.transactionRepo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return entity.Wallet{}, fmt.Errorf("failed to get wallets: %v", err)
	}
	return wallets, nil
}

// GetTransactionByUserID retrieves transactions by user ID
func (s *transactionService) GetTransactionByUserID(ctx context.Context, userID int) ([]entity.Transaction, error) {
	transactions, err := s.transactionRepo.GetTransactionByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %v", err)
	}
	// log.Printf("transactions: %+v", transactions)
	return transactions, nil

}
