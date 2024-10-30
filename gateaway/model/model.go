// model/model.go
package model

import (
	pb "ewallet/gateaway/proto"
	"time"
)

type TransferWalletRequest struct {
	UserIDFrom int32   `json:"user_idfrom"`
	UserIDTo   int32   `json:"user_idto"`
	Amount     float32 `json:"amount"`
}

type UserAndWalletResponse struct {
	User   *pb.User   `json:"user"`
	Wallet *pb.Wallet `json:"wallet"`
}

type TopUpRequest struct {
	UserID int32   `json:"user_id"`
	Amount float32 `json:"amount"`
}

type UserBalance struct {
	UserID    uint32    `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	Balance   float32   `json:"balance"`
}

type Transaction struct {
	TransactionID  uint32    `json:"transaction_id"`
	UserID         uint32    `json:"user_id"`
	Username       string    `json:"username"`
	Amount         float32   `json:"amount"`
	Description    string    `json:"description"`
	Type           string    `json:"type"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	SourceUserID   uint32    `json:"SourceUserID"`
	SourceUserName string    `json:"SourceUserName"`
}

type TransactionsResponse struct {
	Transactions []Transaction `json:"transactions"`
}
