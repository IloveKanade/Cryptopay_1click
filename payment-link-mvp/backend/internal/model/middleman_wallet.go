package model

import (
	"time"

	"gorm.io/gorm"
)

// MiddlemanWallet 中间人钱包模型
type MiddlemanWallet struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Name          string         `json:"name" gorm:"type:varchar(255)"`                                // 钱包名称
	WalletAddress string         `json:"wallet_address" gorm:"type:varchar(255);uniqueIndex;not null"` // 中间人钱包地址
	PrivateKey    string         `json:"-" gorm:"type:text"`                                           // 私钥（加密存储，不返回给前端）
	Status        string         `json:"status" gorm:"type:varchar(20);default:'active'"`              // active, inactive
	Balance       float64        `json:"balance" gorm:"default:0"`                                     // 当前余额
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

// TransferRecord 转账记录模型
type TransferRecord struct {
	ID                uint            `json:"id" gorm:"primaryKey"`
	PaymentOrderID    uint            `json:"payment_order_id" gorm:"not null"`                     // 关联的支付订单ID
	PaymentOrder      PaymentOrder    `json:"payment_order" gorm:"foreignKey:PaymentOrderID"`       // 关联的支付订单
	MiddlemanWalletID uint            `json:"middleman_wallet_id" gorm:"not null"`                  // 中间人钱包ID
	MiddlemanWallet   MiddlemanWallet `json:"middleman_wallet" gorm:"foreignKey:MiddlemanWalletID"` // 中间人钱包
	FromAddress       string          `json:"from_address" gorm:"type:varchar(255);not null"`       // 转出地址（中间人钱包）
	ToAddress         string          `json:"to_address" gorm:"type:varchar(255);not null"`         // 转入地址（用户钱包）
	Amount            float64         `json:"amount" gorm:"not null"`                               // 转账金额
	FixedFeeAmount    float64         `json:"fixed_fee_amount" gorm:"default:0"`                    // 固定手续费金额
	NetworkFeeAmount  float64         `json:"network_fee_amount" gorm:"default:0"`                  // 网络转账手续费金额
	TotalFeeAmount    float64         `json:"total_fee_amount" gorm:"default:0"`                    // 总手续费金额
	NetAmount         float64         `json:"net_amount" gorm:"not null"`                           // 实际到账金额
	NetworkFees       string          `json:"network_fees" gorm:"type:text"`                        // 网络转账手续费详情（JSON格式）
	Status            string          `json:"status" gorm:"type:varchar(20);default:'pending'"`     // pending, success, failed
	TransactionHash   string          `json:"transaction_hash" gorm:"type:varchar(255)"`            // 区块链交易哈希
	ErrorMessage      string          `json:"error_message" gorm:"type:text"`                       // 错误信息
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	DeletedAt         gorm.DeletedAt  `json:"-" gorm:"index"`
}

// FeeConfig 费率配置模型
type FeeConfig struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	FeeRate       float64        `json:"fee_rate" gorm:"not null;default:0.01"`         // 固定费率，默认1%
	MinFee        float64        `json:"min_fee" gorm:"not null;default:0.1"`           // 最小手续费，默认0.1 USDT
	MaxFee        float64        `json:"max_fee" gorm:"not null;default:10"`            // 最大手续费，默认10 USDT
	NetworkFeeMin float64        `json:"network_fee_min" gorm:"not null;default:0.001"` // 网络转账手续费最小值，默认0.001 USDT
	NetworkFeeMax float64        `json:"network_fee_max" gorm:"not null;default:0.01"`  // 网络转账手续费最大值，默认0.01 USDT
	IsActive      bool           `json:"is_active" gorm:"default:true"`                 // 是否启用
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

// BeforeCreate 创建前钩子
func (tr *TransferRecord) BeforeCreate(tx *gorm.DB) error {
	if tr.Status == "" {
		tr.Status = "pending"
	}
	return nil
}

// CalculateFee 计算手续费
func (fc *FeeConfig) CalculateFee(amount float64) float64 {
	// 计算固定比例手续费
	fixedFee := amount * fc.FeeRate

	if fixedFee < fc.MinFee {
		fixedFee = fc.MinFee
	}

	if fixedFee > fc.MaxFee {
		fixedFee = fc.MaxFee
	}

	return fixedFee
}

// CalculateNetworkFees 计算网络转账手续费
// 注意：这个方法现在返回固定值，实际网络手续费应该在转账时动态获取
func (fc *FeeConfig) CalculateNetworkFees() []float64 {
	var networkFees []float64

	// 固定2次网络转账：从中间人钱包到用户钱包
	// 这里返回配置的默认值，实际转账时会动态获取真实网络手续费
	for i := 0; i < 2; i++ {
		// 使用配置的默认值作为预估
		estimatedFee := (fc.NetworkFeeMin + fc.NetworkFeeMax) / 2
		networkFees = append(networkFees, estimatedFee)
	}

	return networkFees
}

// CalculateTotalFees 计算总手续费（固定手续费 + 网络转账手续费）
func (fc *FeeConfig) CalculateTotalFees(amount float64) (float64, []float64) {
	fixedFee := fc.CalculateFee(amount)
	networkFees := fc.CalculateNetworkFees()

	totalNetworkFee := 0.0
	for _, fee := range networkFees {
		totalNetworkFee += fee
	}

	return fixedFee + totalNetworkFee, networkFees
}
