package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentLink struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	UserID        uint           `json:"user_id" gorm:"not null"`
	User          User           `json:"user" gorm:"foreignKey:UserID"`
	Title         string         `json:"title" gorm:"not null"`
	Amount        float64        `json:"amount" gorm:"not null"`
	Currency      string         `json:"currency" gorm:"default:'CNY'"`
	WalletAddress string         `json:"wallet_address" gorm:"not null"`
	LinkID        string         `json:"link_id" gorm:"uniqueIndex:idx_payment_links_link_id,length:191;not null"`
	Description   string         `json:"description"`
	CallbackURL   string         `json:"callback_url"` // 添加回调地址字段
	Status        string         `json:"status" gorm:"default:'active'"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (pl *PaymentLink) BeforeCreate(tx *gorm.DB) error {
	if pl.LinkID == "" {
		pl.LinkID = uuid.New().String()
	}
	return nil
}
