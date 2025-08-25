package service

import (
	"errors"
	"payment-link-mvp/internal/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentLinkService struct {
	db            *gorm.DB
	epusdtService *EpusdtService
}

func NewPaymentLinkService(db *gorm.DB, epusdtService *EpusdtService) *PaymentLinkService {
	return &PaymentLinkService{
		db:            db,
		epusdtService: epusdtService,
	}
}

type CreatePaymentLinkRequest struct {
	Title         string  `json:"title" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	Currency      string  `json:"currency"`
	WalletAddress string  `json:"wallet_address"`
	Description   string  `json:"description"`
	CallbackURL   string  `json:"callback_url"` // 添加回调地址字段
}

func (s *PaymentLinkService) Create(userID uint, req CreatePaymentLinkRequest) (*model.PaymentLink, error) {
	// 设置默认值
	if req.Currency == "" {
		req.Currency = "USDT"
	}

	// 如果没有提供钱包地址，从epusdt获取
	if req.WalletAddress == "" {
		walletAddress, err := s.epusdtService.GetWalletAddress()
		if err != nil {
			// 如果获取失败，使用默认地址
			req.WalletAddress = "default_wallet_address"
		} else {
			req.WalletAddress = walletAddress
		}
	}

	// 自动将钱包地址添加到Epusdt系统
	if req.WalletAddress != "" && req.WalletAddress != "default_wallet_address" {
		err := s.epusdtService.AddWalletAddress(req.WalletAddress)
		if err != nil {
			// 如果添加失败，记录错误但不阻止创建支付链接
			// 因为可能是钱包地址已存在
			// log.Printf("添加钱包地址到Epusdt失败: %v", err)
		}
	}

	// 生成唯一的LinkID
	linkID := uuid.New().String()

	paymentLink := &model.PaymentLink{
		UserID:        userID,
		Title:         req.Title,
		Amount:        req.Amount,
		Currency:      req.Currency,
		WalletAddress: req.WalletAddress,
		LinkID:        linkID,
		Description:   req.Description,
		CallbackURL:   req.CallbackURL, // 添加回调地址
		Status:        "active",
	}

	if err := s.db.Create(paymentLink).Error; err != nil {
		return nil, err
	}

	return paymentLink, nil
}

func (s *PaymentLinkService) List(userID uint) ([]model.PaymentLink, error) {
	var paymentLinks []model.PaymentLink
	err := s.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&paymentLinks).Error
	return paymentLinks, err
}

func (s *PaymentLinkService) Get(userID uint, id uint) (*model.PaymentLink, error) {
	var paymentLink model.PaymentLink
	err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&paymentLink).Error
	if err != nil {
		return nil, err
	}
	return &paymentLink, nil
}

func (s *PaymentLinkService) GetByLinkID(linkID string) (*model.PaymentLink, error) {
	var paymentLink model.PaymentLink
	err := s.db.Where("link_id = ? AND status = ?", linkID, "active").First(&paymentLink).Error
	if err != nil {
		return nil, err
	}
	return &paymentLink, nil
}

func (s *PaymentLinkService) Delete(userID uint, id uint) error {
	result := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.PaymentLink{})
	if result.RowsAffected == 0 {
		return errors.New("支付链接不存在或无权限删除")
	}
	return result.Error
}

// GetPaymentLinkWithStats 获取支付链接及其统计信息
func (s *PaymentLinkService) GetPaymentLinkWithStats(userID uint, id uint) (*model.PaymentLink, map[string]interface{}, error) {
	paymentLink, err := s.Get(userID, id)
	if err != nil {
		return nil, nil, err
	}

	// 获取统计信息
	var stats map[string]interface{}

	// 总订单数
	var totalOrders int64
	s.db.Model(&model.PaymentOrder{}).Where("payment_link_id = ?", id).Count(&totalOrders)

	// 已支付订单数
	var paidOrders int64
	s.db.Model(&model.PaymentOrder{}).Where("payment_link_id = ? AND status = ?", id, "paid").Count(&paidOrders)

	// 总支付金额
	var totalAmount float64
	s.db.Model(&model.PaymentOrder{}).Where("payment_link_id = ? AND status = ?", id, "paid").Select("COALESCE(SUM(amount), 0)").Scan(&totalAmount)

	// 今日订单数
	today := time.Now().Format("2006-01-02")
	var todayOrders int64
	s.db.Model(&model.PaymentOrder{}).Where("payment_link_id = ? AND DATE(created_at) = ?", id, today).Count(&todayOrders)

	stats = map[string]interface{}{
		"total_orders": totalOrders,
		"paid_orders":  paidOrders,
		"total_amount": totalAmount,
		"today_orders": todayOrders,
		"success_rate": float64(0),
	}

	if totalOrders > 0 {
		stats["success_rate"] = float64(paidOrders) / float64(totalOrders) * 100
	}

	return paymentLink, stats, nil
}

// UpdateWalletAddress 更新钱包地址
func (s *PaymentLinkService) UpdateWalletAddress(userID uint, id uint, walletAddress string) error {
	result := s.db.Model(&model.PaymentLink{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("wallet_address", walletAddress)

	if result.RowsAffected == 0 {
		return errors.New("支付链接不存在或无权限更新")
	}
	return result.Error
}

// GetDefaultWalletAddress 获取默认钱包地址
func (s *PaymentLinkService) GetDefaultWalletAddress() (string, error) {
	return s.epusdtService.GetWalletAddress()
}
