package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"payment-link-mvp/internal/model"
	"time"

	"gorm.io/gorm"
)

type MiddlemanWalletService struct {
	db              *gorm.DB
	epusdtService   *EpusdtService
	tronService     *TronService
	callbackService *CallbackService
}

func NewMiddlemanWalletService(db *gorm.DB, epusdtService *EpusdtService, tronService *TronService, callbackService *CallbackService) *MiddlemanWalletService {
	return &MiddlemanWalletService{
		db:              db,
		epusdtService:   epusdtService,
		tronService:     tronService,
		callbackService: callbackService,
	}
}

// GetRandomMiddlemanWallet 获取随机的中间人钱包地址
func (s *MiddlemanWalletService) GetRandomMiddlemanWallet() (*model.MiddlemanWallet, error) {
	var wallets []model.MiddlemanWallet
	err := s.db.Where("status = ?", "active").Find(&wallets).Error
	if err != nil {
		return nil, fmt.Errorf("查询中间人钱包失败: %w", err)
	}

	if len(wallets) == 0 {
		return nil, errors.New("没有可用的中间人钱包")
	}

	// 随机选择一个钱包
	randomIndex := rand.Intn(len(wallets))
	return &wallets[randomIndex], nil
}

// AddMiddlemanWallet 添加中间人钱包
func (s *MiddlemanWalletService) AddMiddlemanWallet(name, walletAddress, privateKey string) error {
	// 检查钱包地址是否已存在
	var count int64
	err := s.db.Model(&model.MiddlemanWallet{}).Where("wallet_address = ?", walletAddress).Count(&count).Error
	if err != nil {
		return fmt.Errorf("检查钱包地址失败: %w", err)
	}

	if count > 0 {
		return errors.New("钱包地址已存在")
	}

	// 验证钱包地址格式
	if !s.isValidWalletAddress(walletAddress) {
		return errors.New("无效的钱包地址格式")
	}

	// 验证私钥格式（简单验证）
	if privateKey == "" {
		return errors.New("私钥不能为空")
	}

	// 创建中间人钱包记录
	middlemanWallet := &model.MiddlemanWallet{
		Name:          name,
		WalletAddress: walletAddress,
		PrivateKey:    privateKey, // 注意：实际生产环境应该加密存储
		Status:        "active",
		Balance:       0,
	}

	if err := s.db.Create(middlemanWallet).Error; err != nil {
		return fmt.Errorf("创建中间人钱包失败: %w", err)
	}

	// 将钱包地址添加到Epusdt系统
	err = s.epusdtService.AddWalletAddress(walletAddress)
	if err != nil {
		log.Printf("添加钱包地址到Epusdt失败: %v", err)
	}

	return nil
}

// ListMiddlemanWallets 获取中间人钱包列表
func (s *MiddlemanWalletService) ListMiddlemanWallets() ([]model.MiddlemanWallet, error) {
	var wallets []model.MiddlemanWallet
	err := s.db.Order("created_at DESC").Find(&wallets).Error
	return wallets, err
}

// UpdateMiddlemanWalletStatus 更新中间人钱包状态
func (s *MiddlemanWalletService) UpdateMiddlemanWalletStatus(id uint, status string) error {
	result := s.db.Model(&model.MiddlemanWallet{}).Where("id = ?", id).Update("status", status)
	if result.RowsAffected == 0 {
		return errors.New("中间人钱包不存在")
	}
	return result.Error
}

// DeleteMiddlemanWallet 删除中间人钱包
func (s *MiddlemanWalletService) DeleteMiddlemanWallet(id uint) error {
	result := s.db.Delete(&model.MiddlemanWallet{}, id)
	if result.RowsAffected == 0 {
		return errors.New("中间人钱包不存在")
	}
	return result.Error
}

// CreateTransferRecord 创建转账记录
func (s *MiddlemanWalletService) CreateTransferRecord(paymentOrderID uint, middlemanWalletID uint, toAddress string, amount float64) (*model.TransferRecord, error) {
	// 获取费率配置
	var feeConfig model.FeeConfig
	err := s.db.Where("is_active = ?", true).First(&feeConfig).Error
	if err != nil {
		return nil, fmt.Errorf("获取费率配置失败: %w", err)
	}

	// 计算固定手续费
	fixedFee := feeConfig.CalculateFee(amount)

	// 动态获取真实的网络手续费（2次转账）
	realNetworkFees := make([]float64, 2)
	totalNetworkFee := 0.0

	for i := 0; i < 2; i++ {
		networkFee, err := s.tronService.GetNetworkFee()
		if err != nil {
			log.Printf("获取网络手续费失败，使用默认值: %v", err)
			// 使用配置的默认值
			networkFee = (feeConfig.NetworkFeeMin + feeConfig.NetworkFeeMax) / 2
		}
		realNetworkFees[i] = networkFee
		totalNetworkFee += networkFee
	}

	// 计算总手续费和实际到账金额
	totalFee := fixedFee + totalNetworkFee
	netAmount := amount - totalFee

	// 将网络转账手续费转换为JSON字符串
	networkFeesJSON, err := json.Marshal(realNetworkFees)
	if err != nil {
		return nil, fmt.Errorf("序列化网络转账手续费失败: %w", err)
	}

	// 获取中间人钱包信息
	var middlemanWallet model.MiddlemanWallet
	err = s.db.First(&middlemanWallet, middlemanWalletID).Error
	if err != nil {
		return nil, fmt.Errorf("获取中间人钱包失败: %w", err)
	}

	// 创建转账记录
	transferRecord := &model.TransferRecord{
		PaymentOrderID:    paymentOrderID,
		MiddlemanWalletID: middlemanWalletID,
		FromAddress:       middlemanWallet.WalletAddress,
		ToAddress:         toAddress,
		Amount:            amount,
		FixedFeeAmount:    fixedFee,
		NetworkFeeAmount:  totalNetworkFee,
		TotalFeeAmount:    totalFee,
		NetAmount:         netAmount,
		NetworkFees:       string(networkFeesJSON),
		Status:            "pending",
	}

	if err := s.db.Create(transferRecord).Error; err != nil {
		return nil, fmt.Errorf("创建转账记录失败: %w", err)
	}

	return transferRecord, nil
}

// ProcessTransfer 处理转账
func (s *MiddlemanWalletService) ProcessTransfer(transferRecordID uint) error {
	var transferRecord model.TransferRecord
	err := s.db.Preload("MiddlemanWallet").First(&transferRecord, transferRecordID).Error
	if err != nil {
		return fmt.Errorf("获取转账记录失败: %w", err)
	}

	if transferRecord.Status != "pending" {
		return errors.New("转账记录状态不是pending")
	}

	// 使用TronService进行真实的USDT转账
	log.Printf("开始转账: 从 %s 到 %s, 金额: %.4f USDT",
		transferRecord.FromAddress,
		transferRecord.ToAddress,
		transferRecord.NetAmount)

	// 使用存储的私钥进行转账
	transferReq := TronTransferRequest{
		FromAddress: transferRecord.FromAddress,
		ToAddress:   transferRecord.ToAddress,
		Amount:      transferRecord.NetAmount,
		PrivateKey:  transferRecord.MiddlemanWallet.PrivateKey, // 从数据库获取私钥
	}

	// 调用TronService进行转账
	transferResp, err := s.tronService.TransferUSDT(transferReq)
	if err != nil {
		// 转账失败，更新状态
		updates := map[string]interface{}{
			"status":     "failed",
			"updated_at": time.Now(),
		}
		s.db.Model(&transferRecord).Updates(updates)
		return fmt.Errorf("转账失败: %w", err)
	}

	if !transferResp.Success {
		// 转账失败，更新状态
		updates := map[string]interface{}{
			"status":        "failed",
			"error_message": transferResp.Error,
			"updated_at":    time.Now(),
		}
		s.db.Model(&transferRecord).Updates(updates)
		
		// 发送转账失败回调通知
		callbackData := TransferCallbackData{
			TransferID:    fmt.Sprintf("%d", transferRecord.ID),
			FromAddress:   transferRecord.FromAddress,
			ToAddress:     transferRecord.ToAddress,
			Amount:        transferRecord.Amount,
			NetworkFee:    transferRecord.NetworkFeeAmount,
			TotalFee:      transferRecord.TotalFeeAmount,
			NetAmount:     transferRecord.NetAmount,
			ErrorMessage:  transferResp.Error,
		}

		if err := s.callbackService.NotifyTransferFail(callbackData); err != nil {
			log.Printf("发送转账失败回调失败: %v", err)
		}
		
		return fmt.Errorf("转账失败: %s", transferResp.Error)
	}

	// 转账成功，更新转账记录状态
	updates := map[string]interface{}{
		"status":           "success",
		"transaction_hash": transferResp.TxID,
		"updated_at":       time.Now(),
	}

	err = s.db.Model(&transferRecord).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("更新转账记录失败: %w", err)
	}

	// 更新中间人钱包余额（手续费收入）
	err = s.db.Model(&transferRecord.MiddlemanWallet).Update("balance",
		gorm.Expr("balance + ?", transferRecord.TotalFeeAmount)).Error
	if err != nil {
		return fmt.Errorf("更新中间人钱包余额失败: %w", err)
	}

	// 发送转账成功回调通知
	callbackData := TransferCallbackData{
		TransferID:      fmt.Sprintf("%d", transferRecord.ID),
		FromAddress:     transferRecord.FromAddress,
		ToAddress:       transferRecord.ToAddress,
		Amount:          transferRecord.Amount,
		NetworkFee:      transferRecord.NetworkFeeAmount,
		TotalFee:        transferRecord.TotalFeeAmount,
		NetAmount:       transferRecord.NetAmount,
		TransactionHash: transferResp.TxID,
	}

	if err := s.callbackService.NotifyTransferSuccess(callbackData); err != nil {
		log.Printf("发送转账成功回调失败: %v", err)
		// 不返回错误，因为转账本身是成功的
	}

	log.Printf("转账成功: 交易哈希 %s", transferResp.TxID)
	return nil
}

// GetTransferRecords 获取转账记录列表
func (s *MiddlemanWalletService) GetTransferRecords(page, pageSize int) ([]model.TransferRecord, int64, error) {
	var records []model.TransferRecord
	var total int64

	// 计算总数
	if err := s.db.Model(&model.TransferRecord{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询记录列表
	offset := (page - 1) * pageSize
	err := s.db.Preload("PaymentOrder").Preload("MiddlemanWallet").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&records).Error

	return records, total, err
}

// GetTransferRecord 获取单个转账记录
func (s *MiddlemanWalletService) GetTransferRecord(id uint) (*model.TransferRecord, error) {
	var record model.TransferRecord
	err := s.db.Preload("PaymentOrder").Preload("MiddlemanWallet").First(&record, id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// isValidWalletAddress 验证钱包地址格式
func (s *MiddlemanWalletService) isValidWalletAddress(address string) bool {
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
