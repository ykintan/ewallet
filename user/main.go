package main

import (
	"ewallet/user/handler"
	repositories "ewallet/user/repository"
	services "ewallet/user/service"
	"log"
	"net"

	pb "ewallet/user/proto"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// DSN (Data Source Name) for PostgreSQL connection
	dsn := "postgresql://postgres:admin@localhost:5432/tugas_user"

	// Open connection to PostgreSQL using GORM
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Setup repository, service, and handler for User
	userRepo := repositories.NewUserRepository(gormDB)
	userService := services.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Initialize gRPC server
	grpcServer := grpc.NewServer()
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Register the User service with the gRPC server
	pb.RegisterUserServiceServer(grpcServer, userHandler)
	log.Println("Server is running on port :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
