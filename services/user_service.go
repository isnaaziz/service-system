package services

import (
	"errors"
	"service_system/models"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func (u *UserService) GetActiveUsers() ([]models.User, error) {
	var users []models.User
	if err := u.DB.Where("is_deleted = ?", false).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserService) SoftDeleteUser(userID string) error {
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return err
	}
	var user models.User
	if err := u.DB.Where("id = ? AND is_deleted = ?", uint(id), false).First(&user).Error; err != nil {
		return err
	}

	now := time.Now()
	user.IsDeleted = true
	user.DeletedAt = &now
	if err := u.DB.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func (u *UserService) ChangePassword(username, oldPassword, newPassword string) error {
	var user models.User
	if err := u.DB.Where("username = ? AND is_deleted = ?", username, false).First(&user).Error; err != nil {
		return err
	}

	// Cek password lama
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("old password is incorrect")
	}

	// Hash password baru
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashed)
	return u.DB.Save(&user).Error
}
