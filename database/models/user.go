package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID              uint64          `gorm:"primaryKey" json:"id"`
	Name            string          `gorm:"not null" json:"name" validate:"required,max=255"`
	Email           string          `gorm:"not null;unique" json:"email" validate:"required,email,max=255"`
	EmailVerifiedAt *time.Time      `gorm:"column:email_verified_at" json:"email_verified_at"`
	Password        string          `gorm:"not null" json:"-"`
	Phone           string          `gorm:"column:phone" json:"phone"`
	TeamID          *uint64         `gorm:"column:team_id" json:"team_id"`
	CreatedAt       time.Time       `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time       `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt       *gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
	Team            Team            `gorm:"foreignKey:OwnerID" json:"team,omitempty"`
}

//// Validate menggunakan validator untuk memvalidasi struktur User
//func (u *User) Validate() error {
//	validate := validator.New()
//	return validate.Struct(u)
//}
