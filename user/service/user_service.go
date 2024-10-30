package services

import (
	"context"
	models "ewallet/user/entity"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (models.User, error)
	GetUserByID(ctx context.Context, id uint) (models.User, error)
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id uint) error
}

type UserService struct {
	userRepository IUserRepository
}

func NewUserService(userRepository IUserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) CreateUser(username, password, email string) (*models.User, error) {
	user := &models.User{
		Username: username,
		Password: password,
		Email:    email,
	}

	createdUser, err := s.userRepository.CreateUser(context.Background(), user)
	if err != nil {
		return nil, err
	}
	return &createdUser, nil
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	user, err := s.userRepository.GetUserByID(context.Background(), id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	user, err := s.userRepository.GetUserByUsername(context.Background(), username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateUser(user *models.User) error {
	return s.userRepository.UpdateUser(context.Background(), user)
}

func (s *UserService) DeleteUser(id uint) error {
	return s.userRepository.DeleteUser(context.Background(), id)
}
