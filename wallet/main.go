package main

import (
	grpcHandler "ewallet/wallet/handler"
	"ewallet/wallet/repository"
	"ewallet/wallet/service"
	"log"
	"net"

	pb "ewallet/wallet/proto"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// DSN (Data Source Name) for PostgreSQL connection
	dsn := "postgresql://postgres:P4ssw0rd@localhost:5432/go_ewallet"

	// Open connection to PostgreSQL using GORM
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Setup service and handler
	transactionRepo := repository.NewTransactionRepository(gormDB)
	transactionService := service.NewTransactionService(transactionRepo)
	transactionHandler := grpcHandler.NewTransactionHandler(transactionService)

	// Initialize gRPC server
	grpcServer := grpc.NewServer()
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pb.RegisterTransactionServiceServer(grpcServer, transactionHandler)
	log.Println("Server is running on port :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
