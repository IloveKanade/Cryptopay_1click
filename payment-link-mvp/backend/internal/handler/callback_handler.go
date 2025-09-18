package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"payment-link-mvp/internal/config"
	"payment-link-mvp/internal/service"
)

type CallbackHandler struct {
	callbackService *service.CallbackService
	config          *config.Config
}

func NewCallbackHandler(callbackService *service.CallbackService, cfg *config.Config) *CallbackHandler {
	return &CallbackHandler{
		callbackService: callbackService,
		config:          cfg,
	}
}

// GetCallbackConfig 获取回调配置
func (h *CallbackHandler) GetCallbackConfig(c *gin.Context) {
	config := gin.H{
		"transfer_success_url": h.config.Callback.TransferSuccessURL,
		"transfer_fail_url":    h.config.Callback.TransferFailURL,
		"notify_secret":        h.config.Callback.NotifySecret,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    config,
	})
}

// UpdateCallbackConfig 更新回调配置
func (h *CallbackHandler) UpdateCallbackConfig(c *gin.Context) {
	var req struct {
		TransferSuccessURL string `json:"transfer_success_url"`
		TransferFailURL    string `json:"transfer_fail_url"`
		NotifySecret       string `json:"notify_secret"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "请求参数错误: " + err.Error(),
		})
		return
	}

	// 更新配置（这里只是示例，实际应该持久化到数据库或配置文件）
	h.config.Callback.TransferSuccessURL = req.TransferSuccessURL
	h.config.Callback.TransferFailURL = req.TransferFailURL
	if req.NotifySecret != "" {
		h.config.Callback.NotifySecret = req.NotifySecret
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "回调配置更新成功",
	})
}

// TestCallback 测试回调
func (h *CallbackHandler) TestCallback(c *gin.Context) {
	var req struct {
		Type string `json:"type" binding:"required"` // success, fail
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "请求参数错误: " + err.Error(),
		})
		return
	}

	// 创建测试数据
	testData := service.TransferCallbackData{
		TransferID:    "test_123",
		FromAddress:   "TQn9Y2khDD95J42FQtQTdwVVRj2qE8qXK1",
		ToAddress:     "TQn9Y2khDD95J42FQtQTdwVVRj2qE8qXK2",
		Amount:        100.0,
		NetworkFee:    0.001,
		TotalFee:      2.001,
		NetAmount:     97.999,
		TransactionHash: "0x1234567890abcdef",
	}

	var err error
	if req.Type == "success" {
		err = h.callbackService.NotifyTransferSuccess(testData)
	} else if req.Type == "fail" {
		testData.ErrorMessage = "测试失败消息"
		err = h.callbackService.NotifyTransferFail(testData)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的测试类型，支持: success, fail",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "测试回调失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "测试回调发送成功",
	})
}
