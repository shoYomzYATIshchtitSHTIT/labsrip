package repository

import (
	"Backend-RIP/internal/app/ds"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// POST регистрация
func (r *UserRepository) RegisterUser(user *ds.Users) error {
	var existingUser ds.Users
	err := r.db.Where("login = ?", user.Login).First(&existingUser).Error
	if err == nil {
		return fmt.Errorf("user with login %s already exists", user.Login)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return r.db.Create(user).Error
}

// GET полей пользователя после аутентификации (для личного кабинета)
func (r *UserRepository) GetUserProfile(userID uint) (ds.Users, error) {
	var user ds.Users
	err := r.db.Select("user_id, login, is_moderator").Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		return ds.Users{}, err
	}
	return user, nil
}

// PUT пользователя (личный кабинет)
func (r *UserRepository) UpdateUserProfile(userID uint, updates map[string]interface{}) error {
	delete(updates, "user_id")
	delete(updates, "is_moderator")

	result := r.db.Model(&ds.Users{}).Where("user_id = ?", userID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", userID)
	}
	return nil
}

// POST аутентификация
func (r *UserRepository) AuthenticateUser(login, password string) (ds.Users, error) {
	var user ds.Users
	err := r.db.Where("login = ? AND password = ?", login, password).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.Users{}, fmt.Errorf("invalid credentials")
		}
		return ds.Users{}, err
	}

	user.Password = ""
	return user, nil
}

// POST деавторизация
func (r *UserRepository) LogoutUser(userID uint) error {
	return nil
}

// GetUserByID получает пользователя по ID (вспомогательный метод)
func (r *UserRepository) GetUserByID(userID uint) (ds.Users, error) {
	var user ds.Users
	err := r.db.Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		return ds.Users{}, err
	}
	return user, nil
}

// IsModerator проверяет, является ли пользователь модератором
func (r *UserRepository) IsModerator(userID uint) (bool, error) {
	var user ds.Users
	err := r.db.Select("is_moderator").Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		return false, err
	}
	return user.IsModerator, nil
}
