package models

import (
	"time"
)

type UserSession struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"index" json:"user_id"`
	Token     string     `gorm:"uniqueIndex" json:"token"`
	CreatedAt time.Time  `gorm:"index" json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	IsDeleted bool       `gorm:"index;default:false" json:"is_deleted"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
}

func (UserSession) TableName() string {
	return "user_sessions"
}
