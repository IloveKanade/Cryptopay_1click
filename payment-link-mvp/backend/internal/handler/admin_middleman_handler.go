package handler

import (
	"net/http"
	"payment-link-mvp/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminMiddlemanHandler struct {
	middlemanWalletService *service.MiddlemanWalletService
	feeConfigService       *service.FeeConfigService
}

func NewAdminMiddlemanHandler(middlemanWalletService *service.MiddlemanWalletService, feeConfigService *service.FeeConfigService) *AdminMiddlemanHandler {
	return &AdminMiddlemanHandler{
		middlemanWalletService: middlemanWalletService,
		feeConfigService:       feeConfigService,
	}
}

// 中间人钱包相关接口

// ListMiddlemanWallets 获取中间人钱包列表
func (h *AdminMiddlemanHandler) ListMiddlemanWallets(c *gin.Context) {
	wallets, err := h.middlemanWalletService.ListMiddlemanWallets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": wallets,
	})
}

// AddMiddlemanWallet 添加中间人钱包
func (h *AdminMiddlemanHandler) AddMiddlemanWallet(c *gin.Context) {
	var req struct {
		Name          string `json:"name"`
		WalletAddress string `json:"wallet_address" binding:"required"`
		PrivateKey    string `json:"private_key" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果没有提供名称，使用钱包地址的前8位作为默认名称
	if req.Name == "" {
		if len(req.WalletAddress) >= 8 {
			req.Name = req.WalletAddress[:8] + "..."
		} else {
			req.Name = req.WalletAddress
		}
	}

	err := h.middlemanWalletService.AddMiddlemanWallet(req.Name, req.WalletAddress, req.PrivateKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "中间人钱包添加成功",
	})
}

// UpdateMiddlemanWalletStatus 更新中间人钱包状态
func (h *AdminMiddlemanHandler) UpdateMiddlemanWalletStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.middlemanWalletService.UpdateMiddlemanWalletStatus(uint(id), req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "中间人钱包状态更新成功",
	})
}

// DeleteMiddlemanWallet 删除中间人钱包
func (h *AdminMiddlemanHandler) DeleteMiddlemanWallet(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	err = h.middlemanWalletService.DeleteMiddlemanWallet(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "中间人钱包删除成功",
	})
}

// 转账记录相关接口

// ListTransferRecords 获取转账记录列表
func (h *AdminMiddlemanHandler) ListTransferRecords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	records, total, err := h.middlemanWalletService.GetTransferRecords(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      records,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetTransferRecord 获取单个转账记录
func (h *AdminMiddlemanHandler) GetTransferRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	record, err := h.middlemanWalletService.GetTransferRecord(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": record,
	})
}

// ProcessTransfer 手动处理转账
func (h *AdminMiddlemanHandler) ProcessTransfer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	err = h.middlemanWalletService.ProcessTransfer(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "转账处理成功",
	})
}

// 费率配置相关接口

// GetFeeConfig 获取当前费率配置
func (h *AdminMiddlemanHandler) GetFeeConfig(c *gin.Context) {
	config, err := h.feeConfigService.GetFeeConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": config,
	})
}

// UpdateFeeConfig 更新费率配置
func (h *AdminMiddlemanHandler) UpdateFeeConfig(c *gin.Context) {
	var req struct {
		FeeRate       float64 `json:"fee_rate" binding:"required"`
		MinFee        float64 `json:"min_fee" binding:"required"`
		MaxFee        float64 `json:"max_fee" binding:"required"`
		NetworkFeeMin float64 `json:"network_fee_min" binding:"required"`
		NetworkFeeMax float64 `json:"network_fee_max" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.feeConfigService.UpdateFeeConfig(req.FeeRate, req.MinFee, req.MaxFee, req.NetworkFeeMin, req.NetworkFeeMax)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "费率配置更新成功",
	})
}

// CreateFeeConfig 创建新的费率配置
func (h *AdminMiddlemanHandler) CreateFeeConfig(c *gin.Context) {
	var req struct {
		FeeRate       float64 `json:"fee_rate" binding:"required"`
		MinFee        float64 `json:"min_fee" binding:"required"`
		MaxFee        float64 `json:"max_fee" binding:"required"`
		NetworkFeeMin float64 `json:"network_fee_min" binding:"required"`
		NetworkFeeMax float64 `json:"network_fee_max" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.feeConfigService.CreateFeeConfig(req.FeeRate, req.MinFee, req.MaxFee, req.NetworkFeeMin, req.NetworkFeeMax)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "费率配置创建成功",
	})
}

// GetFeeConfigHistory 获取费率配置历史
func (h *AdminMiddlemanHandler) GetFeeConfigHistory(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	configs, total, err := h.feeConfigService.GetFeeConfigHistory(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      configs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// CalculateFee 计算手续费
func (h *AdminMiddlemanHandler) CalculateFee(c *gin.Context) {
	amountStr := c.Query("amount")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的金额"})
		return
	}

	totalFee, networkFees, err := h.feeConfigService.CalculateTotalFees(amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fixedFee, err := h.feeConfigService.CalculateFee(amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"amount":       amount,
		"fixed_fee":    fixedFee,
		"network_fees": networkFees,
		"total_fee":    totalFee,
		"net_amount":   amount - totalFee,
	})
}
