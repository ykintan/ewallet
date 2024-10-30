package handler

import (
	"context"
	models "ewallet/user/entity"
	pb "ewallet/user/proto"
	services "ewallet/user/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	service *services.UserService
}

func NewUserHandler(svc *services.UserService) *UserHandler {
	return &UserHandler{service: svc}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := &models.User{
		Username: req.User.Username,
		Password: req.User.Password,
		Email:    req.User.Email,
	}

	createdUser, err := h.service.CreateUser(user.Username, user.Password, user.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &pb.CreateUserResponse{
		User: &pb.User{
			UserId:    uint32(createdUser.UserID),
			Username:  createdUser.Username,
			Password:  createdUser.Password,
			Email:     createdUser.Email,
			CreatedAt: timestamppb.New(createdUser.CreatedAt),
		},
	}, nil
}

func (h *UserHandler) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {
	user, err := h.service.GetUserByID(uint(req.UserId))
	if err != nil {
		if err.Error() == "user not found" {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user by ID: %v", err)
	}

	return &pb.GetUserByIDResponse{
		User: &pb.User{
			UserId:    uint32(user.UserID),
			Username:  user.Username,
			Password:  user.Password,
			Email:     user.Email,
			CreatedAt: timestamppb.New(user.CreatedAt),
		},
	}, nil
}

func (h *UserHandler) GetUserByUsername(ctx context.Context, req *pb.GetUserByUsernameRequest) (*pb.GetUserByUsernameResponse, error) {
	user, err := h.service.GetUserByUsername(req.Username)
	if err != nil {
		if err.Error() == "user not found" {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user by username: %v", err)
	}

	return &pb.GetUserByUsernameResponse{
		User: &pb.User{
			UserId:    uint32(user.UserID),
			Username:  user.Username,
			Password:  user.Password,
			Email:     user.Email,
			CreatedAt: timestamppb.New(user.CreatedAt),
		},
	}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	user := &models.User{
		UserID:   int32(req.User.UserId),
		Username: req.User.Username,
		Password: req.User.Password,
		Email:    req.User.Email,
	}

	err := h.service.UpdateUser(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	return &pb.UpdateUserResponse{
		User: &pb.User{
			UserId:    uint32(user.UserID),
			Username:  user.Username,
			Password:  user.Password,
			Email:     user.Email,
			CreatedAt: timestamppb.New(user.CreatedAt),
		},
	}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := h.service.DeleteUser(uint(req.UserId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	return &pb.DeleteUserResponse{
		Message: "User deleted successfully",
	}, nil
}
