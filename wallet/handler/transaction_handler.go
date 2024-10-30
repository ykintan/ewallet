package handler

import (
	"context"
	"ewallet/wallet/entity"
	pb "ewallet/wallet/proto"
	"ewallet/wallet/service"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TransactionHandler implements the gRPC service defined in the protobuf file
type TransactionHandler struct {
	pb.UnimplementedTransactionServiceServer
	service service.ITransactionService
}

// NewTransactionHandler creates a new instance of TransactionHandler
func NewTransactionHandler(svc service.ITransactionService) *TransactionHandler {
	return &TransactionHandler{service: svc}
}

// CreateTransaction handles the gRPC request to create a transaction
func (h *TransactionHandler) CreateTransaction(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.CreateTransactionResponse, error) {
	transaction := &entity.Transaction{
		WalletID:        int(req.Transaction.WalletId),
		Amount:          float64(req.Transaction.Amount),
		TransactionType: req.Transaction.TransactionType,
		CreatedAt:       req.Transaction.CreatedAt.AsTime(),
	}
	createdTransaction, err := h.service.CreateTransaction(ctx, transaction)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create transaction: %v", err)
	}
	return &pb.CreateTransactionResponse{
		Transaction: &pb.Transaction{
			TransactionId:   uint32(createdTransaction.TransactionID),
			WalletId:        int32(createdTransaction.WalletID),
			Amount:          float32(createdTransaction.Amount),
			TransactionType: createdTransaction.TransactionType,
			CreatedAt:       timestamppb.New(createdTransaction.CreatedAt),
		},
	}, nil
}

// GetTransaction handles the gRPC request to get a transaction by ID
func (h *TransactionHandler) GetTransaction(ctx context.Context, req *pb.GetTransactionRequest) (*pb.GetTransactionResponse, error) {
	transactionID, err := strconv.Atoi(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid transaction ID: %v", err)
	}

	transaction, err := h.service.GetTransaction(ctx, int32(transactionID))
	if err != nil {
		if err.Error() == "transaction not found" {
			return nil, status.Errorf(codes.NotFound, "transaction not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get transaction: %v", err)
	}
	return &pb.GetTransactionResponse{
		Transaction: &pb.Transaction{
			TransactionId:   uint32(transaction.TransactionID),
			WalletId:        int32(transaction.WalletID),
			Amount:          float32(transaction.Amount),
			TransactionType: transaction.TransactionType,
			CreatedAt:       timestamppb.New(transaction.CreatedAt),
			Walletidsource:  int32(transaction.WalletIDSource),
		},
	}, nil
}

// CreateWallet handles the gRPC request to create a wallet
func (h *TransactionHandler) CreateWallet(ctx context.Context, req *pb.CreateWalletRequest) (*pb.CreateWalletResponse, error) {
	wallet := &entity.Wallet{
		UserID:    uint(req.Wallet.UserId),
		Balance:   float64(req.Wallet.Balance),
		CreatedAt: req.Wallet.CreatedAt.AsTime(),
		UpdatedAt: req.Wallet.UpdatedAt.AsTime(),
	}
	createdWallet, err := h.service.CreateWallet(ctx, wallet)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create wallet: %v", err)
	}
	return &pb.CreateWalletResponse{
		Wallet: &pb.Wallet{
			UserId:    uint32(createdWallet.UserID),
			Balance:   float32(createdWallet.Balance),
			CreatedAt: timestamppb.New(createdWallet.CreatedAt),
			UpdatedAt: timestamppb.New(createdWallet.UpdatedAt),
		},
	}, nil
}

// TransferWallet handles the gRPC request to transfer funds between wallets
func (h *TransactionHandler) TransferWallet(ctx context.Context, req *pb.TransferWalletRequest) (*pb.TransferWalletResponse, error) {
	fromWalletID := int(req.FromWalletId)
	toWalletID := int(req.ToWalletId)
	amount := req.Amount

	err := h.service.TransferWallet(ctx, fromWalletID, toWalletID, float64(amount))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to transfer wallet: %v", err)
	}

	return &pb.TransferWalletResponse{
		Message: "Transfer successful",
	}, nil
}

// TopUp handles the gRPC request to top up a wallet
func (h *TransactionHandler) TopUp(ctx context.Context, req *pb.TopUpRequest) (*pb.TopUpResponse, error) {
	walletID := int(req.WalletId)
	amount := req.Amount

	err := h.service.TopUp(ctx, walletID, float64(amount))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to top up wallet: %v", err)
	}

	// Fetch the updated wallet to return the latest transaction details
	updatedWallet, err := h.service.GetWalletByID(ctx, walletID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get updated wallet: %v", err)
	}

	// Create the transaction record for the top-up
	transaction := &pb.Transaction{
		WalletId:        int32(updatedWallet.Walletid),
		Amount:          float32(amount),
		TransactionType: "in",
		CreatedAt:       timestamppb.Now(),
	}

	return &pb.TopUpResponse{
		Transaction: transaction,
	}, nil
}

// Payment handles the gRPC request to make a payment from a wallet
func (h *TransactionHandler) Payment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	walletID := int(req.WalletId)
	amount := req.Amount

	err := h.service.Payment(ctx, walletID, float64(amount))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to make payment: %v", err)
	}

	// Fetch the updated wallet to return the latest transaction details
	updatedWallet, err := h.service.GetWalletByID(ctx, walletID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get updated wallet: %v", err)
	}

	// Create the transaction record for the payment
	transaction := &pb.Transaction{
		WalletId:        int32(updatedWallet.Walletid),
		Amount:          float32(amount),
		TransactionType: "out",
		CreatedAt:       timestamppb.Now(),
	}

	return &pb.PaymentResponse{
		Transaction: transaction,
	}, nil
}

// GetWalletByUserID handles the gRPC request to get a wallet by user ID
func (h *TransactionHandler) GetWalletByUserID(ctx context.Context, req *pb.GetWalletByUserIDRequest) (*pb.GetWalletByUserIDResponse, error) {
	userID := int(req.UserId)

	wallet, err := h.service.GetWalletByUserID(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get wallet by user ID: %v", err)
	}

	// Convert the wallet to the protobuf format
	pbWallet := &pb.Wallet{
		Id:        int32(wallet.Walletid),
		UserId:    uint32(wallet.UserID),
		Balance:   float32(wallet.Balance),
		CreatedAt: timestamppb.New(wallet.CreatedAt),
		UpdatedAt: timestamppb.New(wallet.UpdatedAt),
	}

	return &pb.GetWalletByUserIDResponse{
		Wallets: pbWallet,
	}, nil
}

// GetTransactionByUserID handles the gRPC request to get transactions by user ID
func (h *TransactionHandler) GetTransactionByUserID(ctx context.Context, req *pb.GetTransactionByUserIDRequest) (*pb.GetTransactionByUserIDResponse, error) {
	userID := int(req.UserId)

	transactions, err := h.service.GetTransactionByUserID(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get transactions by user ID: %v", err)
	}

	// Convert the transactions to the protobuf format
	var pbTransactions []*pb.Transaction
	for _, transaction := range transactions {
		pbTransactions = append(pbTransactions, &pb.Transaction{
			TransactionId:   uint32(transaction.TransactionID),
			WalletId:        int32(transaction.WalletID),
			Amount:          float32(transaction.Amount),
			TransactionType: transaction.TransactionType,
			CreatedAt:       timestamppb.New(transaction.CreatedAt),
			Walletidsource:  int32(transaction.WalletIDSource),
		})
	}

	return &pb.GetTransactionByUserIDResponse{
		Transactions: pbTransactions,
	}, nil
}

func (h *TransactionHandler) GetWalletByID(ctx context.Context, req *pb.GetWalletByIdrequest) (*pb.GetWalletByIdrespon, error) {
	id := int(req.Id)

	wallet, err := h.service.GetWalletByID(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get wallet by user ID: %v", err)
	}

	// Convert the wallet to the protobuf format
	pbWallet := &pb.Wallet{
		Id:        int32(wallet.Walletid),
		UserId:    uint32(wallet.UserID),
		Balance:   float32(wallet.Balance),
		CreatedAt: timestamppb.New(wallet.CreatedAt),
		UpdatedAt: timestamppb.New(wallet.UpdatedAt),
	}

	return &pb.GetWalletByIdrespon{
		Wallet: pbWallet,
	}, nil
}
