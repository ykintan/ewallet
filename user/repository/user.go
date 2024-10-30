package repositories

import (
	"context"
	"errors"
	models "ewallet/user/entity"
	services "ewallet/user/service"

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

type userRepository struct {
	db GormDBIface
}

func NewUserRepository(db GormDBIface) services.IUserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) (models.User, error) {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return models.User{}, err
	}
	return *user, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, userID uint) (models.User, error) {
	var user models.User

	if err := r.db.WithContext(ctx).First(&user, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}
	return user, nil
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	var user models.User

	if err := r.db.WithContext(ctx).First(&user, "username = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}
	return user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("user_id = ?", user.UserID).Updates(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) DeleteUser(ctx context.Context, userID uint) error {
	if err := r.db.WithContext(ctx).Delete(&models.User{}, "user_id = ?", userID).Error; err != nil {
		return err
	}
	return nil
}
