package services

import (
	"service_system/models"
	"strconv"
	"time"

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
