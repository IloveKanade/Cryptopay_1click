package service

import (
	"errors"
	"fmt"
	"payment-link-mvp/internal/model"

	"gorm.io/gorm"
)

type FeeConfigService struct {
	db *gorm.DB
}

func NewFeeConfigService(db *gorm.DB) *FeeConfigService {
	return &FeeConfigService{
		db: db,
	}
}

// GetFeeConfig 获取当前费率配置
func (s *FeeConfigService) GetFeeConfig() (*model.FeeConfig, error) {
	var feeConfig model.FeeConfig
	err := s.db.Where("is_active = ?", true).First(&feeConfig).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("没有找到有效的费率配置")
		}
		return nil, fmt.Errorf("获取费率配置失败: %w", err)
	}
	return &feeConfig, nil
}

// UpdateFeeConfig 更新费率配置
func (s *FeeConfigService) UpdateFeeConfig(feeRate, minFee, maxFee, networkFeeMin, networkFeeMax float64) error {
	// 验证参数
	if feeRate < 0 || feeRate > 1 {
		return errors.New("费率必须在0-1之间")
	}
	if minFee < 0 {
		return errors.New("最小手续费不能为负数")
	}
	if maxFee < minFee {
		return errors.New("最大手续费不能小于最小手续费")
	}
	if networkFeeMin < 0 {
		return errors.New("网络手续费最小值不能为负数")
	}
	if networkFeeMax < networkFeeMin {
		return errors.New("网络手续费最大值不能小于最小值")
	}

	// 获取当前配置
	var feeConfig model.FeeConfig
	err := s.db.Where("is_active = ?", true).First(&feeConfig).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新的费率配置
			newFeeConfig := model.FeeConfig{
				FeeRate:       feeRate,
				MinFee:        minFee,
				MaxFee:        maxFee,
				NetworkFeeMin: networkFeeMin,
				NetworkFeeMax: networkFeeMax,
				IsActive:      true,
			}
			return s.db.Create(&newFeeConfig).Error
		}
		return fmt.Errorf("获取费率配置失败: %w", err)
	}

	// 更新现有配置
	updates := map[string]interface{}{
		"fee_rate":        feeRate,
		"min_fee":         minFee,
		"max_fee":         maxFee,
		"network_fee_min": networkFeeMin,
		"network_fee_max": networkFeeMax,
	}

	return s.db.Model(&feeConfig).Updates(updates).Error
}

// CreateFeeConfig 创建新的费率配置
func (s *FeeConfigService) CreateFeeConfig(feeRate, minFee, maxFee, networkFeeMin, networkFeeMax float64) error {
	// 验证参数
	if feeRate < 0 || feeRate > 1 {
		return errors.New("费率必须在0-1之间")
	}
	if minFee < 0 {
		return errors.New("最小手续费不能为负数")
	}
	if maxFee < minFee {
		return errors.New("最大手续费不能小于最小手续费")
	}
	if networkFeeMin < 0 {
		return errors.New("网络手续费最小值不能为负数")
	}
	if networkFeeMax < networkFeeMin {
		return errors.New("网络手续费最大值不能小于最小值")
	}

	// 禁用所有现有配置
	err := s.db.Model(&model.FeeConfig{}).Update("is_active", false).Error
	if err != nil {
		return fmt.Errorf("禁用现有配置失败: %w", err)
	}

	// 创建新配置
	newFeeConfig := model.FeeConfig{
		FeeRate:       feeRate,
		MinFee:        minFee,
		MaxFee:        maxFee,
		NetworkFeeMin: networkFeeMin,
		NetworkFeeMax: networkFeeMax,
		IsActive:      true,
	}

	return s.db.Create(&newFeeConfig).Error
}

// ListFeeConfigs 获取所有费率配置历史
func (s *FeeConfigService) ListFeeConfigs() ([]model.FeeConfig, error) {
	var configs []model.FeeConfig
	err := s.db.Order("created_at DESC").Find(&configs).Error
	return configs, err
}

// CalculateFee 计算手续费
func (s *FeeConfigService) CalculateFee(amount float64) (float64, error) {
	feeConfig, err := s.GetFeeConfig()
	if err != nil {
		return 0, err
	}

	return feeConfig.CalculateFee(amount), nil
}

// CalculateTotalFees 计算总手续费（固定手续费 + 网络转账手续费）
func (s *FeeConfigService) CalculateTotalFees(amount float64) (float64, []float64, error) {
	feeConfig, err := s.GetFeeConfig()
	if err != nil {
		return 0, nil, err
	}

	totalFee, networkFees := feeConfig.CalculateTotalFees(amount)
	return totalFee, networkFees, nil
}

// CalculateNetworkFees 计算网络转账手续费
func (s *FeeConfigService) CalculateNetworkFees() ([]float64, error) {
	feeConfig, err := s.GetFeeConfig()
	if err != nil {
		return nil, err
	}

	return feeConfig.CalculateNetworkFees(), nil
}

// GetFeeConfigHistory 获取费率配置历史
func (s *FeeConfigService) GetFeeConfigHistory(page, pageSize int) ([]model.FeeConfig, int64, error) {
	var configs []model.FeeConfig
	var total int64

	// 计算总数
	if err := s.db.Model(&model.FeeConfig{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询配置列表
	offset := (page - 1) * pageSize
	err := s.db.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&configs).Error

	return configs, total, err
}
