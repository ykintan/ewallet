package repository

import (
	"context"
	"errors"
	"ewallet/wallet/entity"
	"ewallet/wallet/service"

	"gorm.io/gorm"
)

type GormDBIface interface {
	WithContext(ctx context.Context) *gorm.DB
	Create(value interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
	Delete(value interface{}, conds ...interface{}) *gorm.DB
	Find(dest interface{}, conds ...interface{}) *gorm.DB
}

type transactionRepository struct {
	db GormDBIface
}

func NewTransactionRepository(db GormDBIface) service.ITransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateWallet(ctx context.Context, wallet *entity.Wallet) (entity.Wallet, error) {
	if err := r.db.WithContext(ctx).Create(wallet).Error; err != nil {
		return entity.Wallet{}, err
	}
	return *wallet, nil
}

func (r *transactionRepository) CreateTransaction(ctx context.Context, transaction *entity.Transaction) (entity.Transaction, error) {
	if err := r.db.WithContext(ctx).Create(transaction).Error; err != nil {
		return entity.Transaction{}, err
	}
	return *transaction, nil
}

// GetWalletByID retrieves a wallet by its ID from the repository
func (r *transactionRepository) GetWalletByID(ctx context.Context, walletID int) (entity.Wallet, error) {
	var wallet entity.Wallet

	if err := r.db.WithContext(ctx).First(&wallet, "Wallet_id = ?", walletID).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Wallet{}, err
	}
	return wallet, nil
}

// GetWalletByUserID retrieves wallets by user ID from the repository
func (r *transactionRepository) GetWalletByUserID(ctx context.Context, userID int) (entity.Wallet, error) {
	var wallets entity.Wallet

	if err := r.db.WithContext(ctx).Find(&wallets, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Wallet{}, errors.New("wallets not found")
		}
		return entity.Wallet{}, err
	}
	return wallets, nil
}

// GetTransaction retrieves a transaction by its ID from the repository
func (r *transactionRepository) GetTransaction(ctx context.Context, id int32) (entity.Transaction, error) {
	var transaction entity.Transaction

	if err := r.db.WithContext(ctx).First(&transaction, "transaction_id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Transaction{}, errors.New("transaction not found")
		}
		return entity.Transaction{}, err
	}
	return transaction, nil
}

// UpdateWallet updates a wallet in the repository
func (r *transactionRepository) UpdateWallet(ctx context.Context, wallet *entity.Wallet) error {
	if err := r.db.WithContext(ctx).Model(&entity.Wallet{}).Where("Wallet_id = ?", wallet.Walletid).Updates(wallet).Error; err != nil {
		return err
	}
	return nil
}

// GetTransactionByUserID retrieves transactions by user ID
func (r *transactionRepository) GetTransactionByUserID(ctx context.Context, userID int) ([]entity.Transaction, error) {
	var transactions []entity.Transaction

	if err := r.db.WithContext(ctx).
		// Debug().
		Joins("JOIN wallets ON transactions.wallet_id = wallets.wallet_id").
		Where("wallets.user_id = ?", userID).
		Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}
