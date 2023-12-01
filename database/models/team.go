package models

import (
	"gorm.io/gorm"
	"time"
)

type Team struct {
	ID                 uint64          `gorm:"primaryKey" json:"id"`
	OwnerID            *uint64         `gorm:"column:owner_id" json:"owner_id"`
	Name               string          `gorm:"not null" json:"name" validate:"required,max=255"`
	Logo               string          `gorm:"column:logo" json:"logo"`
	Address            string          `gorm:"column:address" json:"address"`
	CallbackURL        string          `gorm:"column:callback_url" json:"callback_url"`
	ReturnURL          string          `gorm:"column:return_url" json:"return_url"`
	BankID             string          `gorm:"column:bank_id" json:"bank_id"`
	Status             string          `gorm:"type:enum('unverified','verified','pending','rejected');default:'unverified'" json:"status"`
	StatusDisbursement string          `gorm:"type:enum('unverified','verified','pending','rejected');default:'unverified'" json:"status_disbursement"`
	Rate               string          `gorm:"type:enum('REGULER','CUSTOM');default:'REGULER'" json:"rate"`
	UUID               string          `gorm:"column:uuid" json:"uuid"`
	Secret             string          `gorm:"column:secret" json:"secret"`
	SSLVerification    bool            `gorm:"column:ssl_verification" json:"ssl_verification"`
	OtherBank          string          `gorm:"column:other_bank" json:"other_bank"`
	SettleTime         string          `gorm:"type:enum('DEFAULT','CUSTOM');default:'DEFAULT'" json:"settle_time"`
	FeeChargedTo       string          `gorm:"type:enum('MERCHANT','USER');default:'MERCHANT'" json:"fee_charged_to"`
	CampaignName       string          `gorm:"column:campaign_name" json:"campaign_name"`
	CreatedAt          time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedAt          time.Time       `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt          *gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}
