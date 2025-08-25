package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"payment-link-mvp/internal/model"
	"time"

	"gorm.io/gorm"
)

type PaymentService struct {
	db                     *gorm.DB
	epusdtService          *EpusdtService
	middlemanWalletService *MiddlemanWalletService
	config                 *Config
}

type Config struct {
	BaseURL     string
	FrontendURL string
}

type CreatePaymentRequest struct {
	PaymentLinkID uint    `json:"payment_link_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	ClientIP      string  `json:"client_ip"`
}

type PaymentResult struct {
	OrderID      string    `json:"order_id"`
	TradeID      string    `json:"trade_id"`
	Amount       float64   `json:"amount"`
	ActualAmount float64   `json:"actual_amount"`
	PaymentURL   string    `json:"payment_url"`
	ExpiresAt    time.Time `json:"expires_at"`
	Status       string    `json:"status"`
}

func NewPaymentService(db *gorm.DB, epusdtService *EpusdtService, middlemanWalletService *MiddlemanWalletService, config *Config) *PaymentService {
	return &PaymentService{
		db:                     db,
		epusdtService:          epusdtService,
		middlemanWalletService: middlemanWalletService,
		config:                 config,
	}
}

// CreateDirectEpusdtOrder 直接创建Epusdt订单（用于重定向）
func (s *PaymentService) CreateDirectEpusdtOrder(orderID, amount, walletAddress, clientIP string) (*EpusdtOrderResponse, error) {
	// 构建回调URL
	notifyURL := fmt.Sprintf("%s/api/payment/notify", s.config.BaseURL)
	redirectURL := fmt.Sprintf("%s/payment/success?order_id=%s", s.config.FrontendURL, orderID)

	// 创建Epusdt订单
	response, err := s.epusdtService.CreateOrder(orderID, amount, "支付订单", notifyURL, redirectURL, clientIP, walletAddress)
	if err != nil {
		return nil, fmt.Errorf("创建Epusdt订单失败: %w", err)
	}

	return response, nil
}

// CreateDirectEpusdtOrderWithRandomWallet 使用随机中间人钱包地址创建Epusdt订单
func (s *PaymentService) CreateDirectEpusdtOrderWithRandomWallet(orderID, amount, title, clientIP string) (*EpusdtOrderResponse, error) {
	// 1. 获取随机的中间人钱包地址
	middlemanWallet, err := s.middlemanWalletService.GetRandomMiddlemanWallet()
	if err != nil {
		return nil, fmt.Errorf("获取中间人钱包失败: %w", err)
	}

	// 2. 确保中间人钱包地址在Epusdt系统中
	err = s.epusdtService.AddWalletAddress(middlemanWallet.WalletAddress)
	if err != nil {
		log.Printf("添加中间人钱包地址到Epusdt失败: %v", err)
	}

	// 3. 构建回调URL
	notifyURL := fmt.Sprintf("%s/api/payment/notify", s.config.BaseURL)
	redirectURL := fmt.Sprintf("%s/payment/success?order_id=%s", s.config.FrontendURL, orderID)

	// 4. 创建Epusdt订单，使用随机中间人钱包地址
	response, err := s.epusdtService.CreateOrder(orderID, amount, title, notifyURL, redirectURL, clientIP, middlemanWallet.WalletAddress)
	if err != nil {
		return nil, fmt.Errorf("创建Epusdt订单失败: %w", err)
	}

	return response, nil
}

// CreatePayment 创建支付订单
func (s *PaymentService) CreatePayment(req CreatePaymentRequest) (*PaymentResult, error) {
	// 1. 验证支付链接
	var paymentLink model.PaymentLink
	if err := s.db.First(&paymentLink, req.PaymentLinkID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("支付链接不存在")
		}
		return nil, fmt.Errorf("查询支付链接失败: %w", err)
	}

	if paymentLink.Status != "active" {
		return nil, errors.New("支付链接已失效")
	}

	// 2. 获取随机的中间人钱包地址
	middlemanWallet, err := s.middlemanWalletService.GetRandomMiddlemanWallet()
	if err != nil {
		return nil, fmt.Errorf("获取中间人钱包失败: %w", err)
	}

	// 3. 确保中间人钱包地址在Epusdt系统中
	err = s.epusdtService.AddWalletAddress(middlemanWallet.WalletAddress)
	if err != nil {
		log.Printf("添加中间人钱包地址到Epusdt失败: %v", err)
	}

	// 4. 生成订单ID
	orderID := fmt.Sprintf("ORDER_%d_%d", paymentLink.ID, time.Now().Unix())

	// 5. 构建回调地址
	notifyURL := fmt.Sprintf("%s/api/payment/notify", s.config.BaseURL)
	redirectURL := fmt.Sprintf("%s/payment/success/%s", s.config.FrontendURL, orderID)

	// 6. 调用Epusdt创建订单，使用中间人钱包地址
	epusdtResp, err := s.epusdtService.CreateOrder(
		orderID,
		fmt.Sprintf("%.2f", req.Amount),
		paymentLink.Title,
		notifyURL,
		redirectURL,
		req.ClientIP,
		middlemanWallet.WalletAddress, // 使用中间人钱包地址
	)
	if err != nil {
		return nil, fmt.Errorf("创建Epusdt订单失败: %w", err)
	}

	// 7. 创建本地支付订单记录
	paymentOrder := &model.PaymentOrder{
		PaymentLinkID:  req.PaymentLinkID,
		TradeID:        epusdtResp.Data.TradeID,
		OrderID:        orderID,
		Amount:         req.Amount,
		ActualAmount:   epusdtResp.Data.ActualAmount,
		Token:          epusdtResp.Data.Token,
		Status:         "pending",
		PaymentURL:     epusdtResp.Data.PaymentURL,
		ExpirationTime: time.Unix(epusdtResp.Data.ExpirationTime, 0),
		NotifyURL:      notifyURL,
		RedirectURL:    redirectURL,
		ClientIP:       req.ClientIP,
	}

	if err := s.db.Create(paymentOrder).Error; err != nil {
		return nil, fmt.Errorf("保存支付订单失败: %w", err)
	}

	// 8. 返回支付结果
	result := &PaymentResult{
		OrderID:      orderID,
		TradeID:      epusdtResp.Data.TradeID,
		Amount:       req.Amount,
		ActualAmount: epusdtResp.Data.ActualAmount,
		PaymentURL:   epusdtResp.Data.PaymentURL,
		ExpiresAt:    time.Unix(epusdtResp.Data.ExpirationTime, 0),
		Status:       "pending",
	}

	return result, nil
}

// GetPaymentOrder 获取支付订单
func (s *PaymentService) GetPaymentOrder(orderID string) (*model.PaymentOrder, error) {
	var paymentOrder model.PaymentOrder
	err := s.db.Preload("PaymentLink").Where("order_id = ?", orderID).First(&paymentOrder).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("支付订单不存在")
		}
		return nil, fmt.Errorf("查询支付订单失败: %w", err)
	}
	return &paymentOrder, nil
}

// GetPaymentOrderByTradeID 通过TradeID获取支付订单
func (s *PaymentService) GetPaymentOrderByTradeID(tradeID string) (*model.PaymentOrder, error) {
	var paymentOrder model.PaymentOrder
	err := s.db.Preload("PaymentLink").Where("trade_id = ?", tradeID).First(&paymentOrder).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("支付订单不存在")
		}
		return nil, fmt.Errorf("查询支付订单失败: %w", err)
	}
	return &paymentOrder, nil
}

// UpdatePaymentOrderStatus 更新支付订单状态
func (s *PaymentService) UpdatePaymentOrderStatus(tradeID, status, blockTransactionID string) error {
	updates := map[string]interface{}{
		"status": status,
	}

	if blockTransactionID != "" {
		updates["block_transaction_id"] = blockTransactionID
	}

	if status == "paid" {
		updates["updated_at"] = time.Now()
	}

	err := s.db.Model(&model.PaymentOrder{}).Where("trade_id = ?", tradeID).Updates(updates).Error
	if err != nil {
		return err
	}

	// 如果支付成功，触发转账流程
	if status == "paid" {
		go s.handlePaymentSuccess(tradeID)
	}

	return nil
}

// handlePaymentSuccess 处理支付成功后的转账
func (s *PaymentService) handlePaymentSuccess(tradeID string) {
	// 获取支付订单信息
	paymentOrder, err := s.GetPaymentOrderByTradeID(tradeID)
	if err != nil {
		log.Printf("获取支付订单失败: %v", err)
		return
	}

	// 获取支付链接信息
	var paymentLink model.PaymentLink
	if err := s.db.First(&paymentLink, paymentOrder.PaymentLinkID).Error; err != nil {
		log.Printf("获取支付链接失败: %v", err)
		return
	}

	// 获取中间人钱包信息
	var middlemanWallet model.MiddlemanWallet
	if err := s.db.Where("wallet_address = ?", paymentOrder.Token).First(&middlemanWallet).Error; err != nil {
		log.Printf("获取中间人钱包失败: %v", err)
		return
	}

	// 创建转账记录
	transferRecord, err := s.middlemanWalletService.CreateTransferRecord(
		paymentOrder.ID,
		middlemanWallet.ID,
		paymentLink.WalletAddress, // 用户登记的钱包地址
		paymentOrder.ActualAmount, // 实际USDT金额
	)
	if err != nil {
		log.Printf("创建转账记录失败: %v", err)
		return
	}

	// 执行转账
	err = s.middlemanWalletService.ProcessTransfer(transferRecord.ID)
	if err != nil {
		log.Printf("执行转账失败: %v", err)
		return
	}

	// 如果支付链接配置了回调地址，发送回调通知
	if paymentLink.CallbackURL != "" {
		go s.sendPaymentCallback(paymentLink, paymentOrder, transferRecord)
	}

	log.Printf("支付成功，转账完成: 订单ID %s, 转账记录ID %d", paymentOrder.OrderID, transferRecord.ID)
}

// sendPaymentCallback 发送支付回调通知
func (s *PaymentService) sendPaymentCallback(paymentLink model.PaymentLink, paymentOrder *model.PaymentOrder, transferRecord *model.TransferRecord) {
	// 构建回调数据
	callbackData := map[string]interface{}{
		"payment_link_id":  paymentLink.ID,
		"order_id":         paymentOrder.OrderID,
		"trade_id":         paymentOrder.TradeID,
		"amount":           paymentOrder.Amount,
		"actual_amount":    paymentOrder.ActualAmount,
		"status":           "success",
		"transfer_id":      transferRecord.ID,
		"from_address":     transferRecord.FromAddress,
		"to_address":       transferRecord.ToAddress,
		"transfer_amount":  transferRecord.NetAmount,
		"transaction_hash": transferRecord.TransactionHash,
		"timestamp":        time.Now().Unix(),
	}

	// 发送HTTP POST请求到回调地址
	client := &http.Client{Timeout: 30 * time.Second}

	jsonData, err := json.Marshal(callbackData)
	if err != nil {
		log.Printf("序列化回调数据失败: %v", err)
		return
	}

	req, err := http.NewRequest("POST", paymentLink.CallbackURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("创建回调请求失败: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Payment-System/1.0")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("发送回调请求失败: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("支付回调发送成功: %s, 状态码: %d", paymentLink.CallbackURL, resp.StatusCode)
	} else {
		log.Printf("支付回调发送失败: %s, 状态码: %d", paymentLink.CallbackURL, resp.StatusCode)
	}
}

// QueryOrderStatus 查询订单状态
func (s *PaymentService) QueryOrderStatus(tradeID string) (*model.PaymentOrder, error) {
	// 1. 从本地数据库查询
	paymentOrder, err := s.GetPaymentOrderByTradeID(tradeID)
	if err != nil {
		return nil, err
	}

	// 2. 如果订单状态不是pending，直接返回
	if paymentOrder.Status != "pending" {
		return paymentOrder, nil
	}

	// 3. 检查是否过期
	if paymentOrder.IsExpired() {
		// 更新状态为过期
		if err := s.UpdatePaymentOrderStatus(tradeID, "expired", ""); err != nil {
			return nil, fmt.Errorf("更新订单状态失败: %w", err)
		}
		paymentOrder.Status = "expired"
		return paymentOrder, nil
	}

	// 4. 调用Epusdt API查询最新状态
	epusdtResp, err := s.epusdtService.QueryOrder(tradeID)
	if err != nil {
		// API调用失败，返回本地状态
		return paymentOrder, nil
	}

	// 5. 根据Epusdt状态更新本地状态
	var newStatus string
	switch epusdtResp.Data.Status {
	case 1: // 等待支付
		newStatus = "pending"
	case 2: // 支付成功
		newStatus = "paid"
	case 3: // 已过期
		newStatus = "expired"
	default:
		newStatus = "pending"
	}

	// 6. 如果状态有变化，更新本地数据库
	if newStatus != paymentOrder.Status {
		if err := s.UpdatePaymentOrderStatus(tradeID, newStatus, ""); err != nil {
			return nil, fmt.Errorf("更新订单状态失败: %w", err)
		}
		paymentOrder.Status = newStatus
	}

	return paymentOrder, nil
}

// GetUserPaymentOrders 获取用户的支付订单列表
func (s *PaymentService) GetUserPaymentOrders(userID uint, page, pageSize int) ([]model.PaymentOrder, int64, error) {
	var orders []model.PaymentOrder
	var total int64

	// 计算总数
	if err := s.db.Model(&model.PaymentOrder{}).
		Joins("JOIN payment_links ON payment_orders.payment_link_id = payment_links.id").
		Where("payment_links.user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询订单列表
	offset := (page - 1) * pageSize
	err := s.db.Preload("PaymentLink").
		Joins("JOIN payment_links ON payment_orders.payment_link_id = payment_links.id").
		Where("payment_links.user_id = ?", userID).
		Order("payment_orders.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&orders).Error

	return orders, total, err
}

// CleanExpiredOrders 清理过期订单
func (s *PaymentService) CleanExpiredOrders() error {
	return s.db.Model(&model.PaymentOrder{}).
		Where("status = ? AND expiration_time < ?", "pending", time.Now()).
		Update("status", "expired").Error
}

// isValidWalletAddress 验证钱包地址格式
func (s *PaymentService) isValidWalletAddress(address string) bool {
	// TRC20地址格式验证：34个字符，以T开头
	if len(address) != 34 || address[0] != 'T' {
		return false
	}

	// 检查是否只包含字母和数字（允许大小写字母）
	for _, char := range address {
		if !((char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9')) {
			return false
		}
	}

	return true
}

// VerifyNotifySignature 验证回调签名
func (s *PaymentService) VerifyNotifySignature(data map[string]interface{}) bool {
	return s.epusdtService.VerifyNotifySignature(data)
}
