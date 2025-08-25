package model

import (
	"time"
	"gorm.io/gorm"
)

// PaymentOrder 支付订单模型
type PaymentOrder struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	PaymentLinkID   uint      `json:"payment_link_id" gorm:"not null"`
	TradeID         string    `json:"trade_id" gorm:"type:varchar(255);uniqueIndex;not null"` // Epusdt交易ID
	OrderID         string    `json:"order_id" gorm:"type:varchar(255);uniqueIndex;not null"` // 本地订单ID
	Amount          float64   `json:"amount" gorm:"not null"`               // 支付金额(CNY)
	ActualAmount    float64   `json:"actual_amount"`                        // 实际USDT金额
	Token           string    `json:"token" gorm:"type:varchar(255)"`       // 钱包地址
	Status          string    `json:"status" gorm:"type:varchar(50);not null;default:'pending'"` // 订单状态: pending, paid, expired, failed
	PaymentURL      string    `json:"payment_url" gorm:"type:text"`         // Epusdt收银台地址
	ExpirationTime  time.Time `json:"expiration_time"`                      // 过期时间
	BlockTransactionID string `json:"block_transaction_id" gorm:"type:varchar(255)"` // 区块链交易ID
	NotifyURL       string    `json:"notify_url" gorm:"type:text"`          // 回调地址
	RedirectURL     string    `json:"redirect_url" gorm:"type:text"`        // 跳转地址
	ClientIP        string    `json:"client_ip" gorm:"type:varchar(45)"`    // 客户端IP
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	
	// 关联关系
	PaymentLink     PaymentLink `json:"payment_link" gorm:"foreignKey:PaymentLinkID"`
}

// BeforeCreate 创建前钩子
func (po *PaymentOrder) BeforeCreate(tx *gorm.DB) error {
	if po.Status == "" {
		po.Status = "pending"
	}
	return nil
}

// IsExpired 检查订单是否过期
func (po *PaymentOrder) IsExpired() bool {
	return time.Now().After(po.ExpirationTime)
}

// IsPaid 检查订单是否已支付
func (po *PaymentOrder) IsPaid() bool {
	return po.Status == "paid"
}

// CanPay 检查订单是否可以支付
func (po *PaymentOrder) CanPay() bool {
	return po.Status == "pending" && !po.IsExpired()
}
