package services

import (
	"errors"
	"service_system/models"
	"service_system/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
}

func (a *AuthService) Register(username, password, email string) (*models.User, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: username,
		Password: string(hashedPassword),
		Email:    email,
	}

	if err := a.DB.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (a *AuthService) Login(username, password string) (string, error) {
	var user models.User
	if err := a.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Cek session aktif (belum expired & belum dihapus)
	var existingSession models.UserSession
	now := time.Now()
	if err := a.DB.
		Select("id").
		Where("user_id = ? AND is_deleted = ? AND created_at > ?", user.ID, false, now.Add(-24*time.Hour)).
		First(&existingSession).Error; err == nil {
		return "", errors.New("user already logged in, please logout first")
	}

	token, err := utils.CreateJWT(user.Username)
	if err != nil {
		return "", err
	}

	session := &models.UserSession{
		UserID:    user.ID,
		Token:     token,
		CreatedAt: now,
		UpdatedAt: now,
		IsDeleted: false,
	}
	if err := a.DB.Create(session).Error; err != nil {
		return "", err
	}

	return token, nil
}

func (a *AuthService) Logout(token string) error {
	var session models.UserSession
	if err := a.DB.Where("token = ? AND is_deleted = ?", token, false).First(&session).Error; err != nil {
		return errors.New("session not found")
	}
	session.IsDeleted = true
	now := time.Now()
	session.DeletedAt = &now
	return a.DB.Save(&session).Error
}
