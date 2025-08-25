package handler

import (
	"fmt"
	"net/http"
	"payment-link-mvp/internal/model"
	"payment-link-mvp/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
}

func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

// CreatePayment 创建支付订单
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req service.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 获取客户端IP
	if req.ClientIP == "" {
		req.ClientIP = c.ClientIP()
	}

	result, err := h.paymentService.CreatePayment(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "支付订单创建成功",
		"data":    result,
	})
}

// GetPaymentOrder 获取支付订单详情
func (h *PaymentHandler) GetPaymentOrder(c *gin.Context) {
	orderID := c.Param("orderId")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "订单ID不能为空"})
		return
	}

	paymentOrder, err := h.paymentService.GetPaymentOrder(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取订单成功",
		"data":    paymentOrder,
	})
}

// QueryOrderStatus 查询订单状态
func (h *PaymentHandler) QueryOrderStatus(c *gin.Context) {
	tradeID := c.Param("tradeId")
	if tradeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "交易ID不能为空"})
		return
	}

	paymentOrder, err := h.paymentService.QueryOrderStatus(tradeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "查询订单状态成功",
		"data":    paymentOrder,
	})
}

// GetUserPaymentOrders 获取用户的支付订单列表
func (h *PaymentHandler) GetUserPaymentOrders(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	userModel := user.(model.User)

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	orders, total, err := h.paymentService.GetUserPaymentOrders(userModel.ID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取订单列表成功",
		"data": gin.H{
			"orders":     orders,
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// HandleNotify 处理Epusdt支付回调
func (h *PaymentHandler) HandleNotify(c *gin.Context) {
	// 读取请求体
	var notifyData map[string]interface{}
	if err := c.ShouldBindJSON(&notifyData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 验证签名
	if !h.paymentService.VerifyNotifySignature(notifyData) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "签名验证失败"})
		return
	}

	// 提取回调数据
	tradeID, ok := notifyData["trade_id"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少trade_id"})
		return
	}

	status, ok := notifyData["status"].(float64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少status"})
		return
	}

	blockTransactionID, _ := notifyData["block_transaction_id"].(string)

	// 根据状态更新订单
	var orderStatus string
	switch int(status) {
	case 1: // 等待支付
		orderStatus = "pending"
	case 2: // 支付成功
		orderStatus = "paid"
	case 3: // 已过期
		orderStatus = "expired"
	default:
		orderStatus = "pending"
	}

	// 更新订单状态
	if err := h.paymentService.UpdatePaymentOrderStatus(tradeID, orderStatus, blockTransactionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新订单状态失败"})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// GetPaymentPage 获取支付页面数据
func (h *PaymentHandler) GetPaymentPage(c *gin.Context) {
	linkID := c.Param("linkId")
	if linkID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "链接ID不能为空"})
		return
	}

	// 这里需要调用支付链接服务获取链接信息
	// 暂时返回简单的响应
	c.JSON(http.StatusOK, gin.H{
		"message": "获取支付页面成功",
		"data": gin.H{
			"link_id": linkID,
			"status":  "ready",
		},
	})
}

// RedirectToEpusdtPayment 直接重定向到Epusdt付款页面
func (h *PaymentHandler) RedirectToEpusdtPayment(c *gin.Context) {
	// 从查询参数获取必要信息
	amount := c.Query("amount")
	orderID := c.Query("order_id")
	title := c.Query("title") // 支付链接标题

	if amount == "" || orderID == "" || title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少必要参数: amount, order_id, title"})
		return
	}

	// 创建Epusdt订单（使用随机中间人钱包地址）
	response, err := h.paymentService.CreateDirectEpusdtOrderWithRandomWallet(orderID, amount, title, c.ClientIP())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("创建支付订单失败: %v", err)})
		return
	}

	// 直接重定向到Epusdt的付款页面
	if response.Data.PaymentURL != "" {
		c.Redirect(http.StatusFound, response.Data.PaymentURL)
		return
	}

	// 如果重定向失败，返回错误
	c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取付款页面URL"})
}

// GetEpusdtPaymentURL 获取Epusdt付款页面URL（不重定向，返回URL）
func (h *PaymentHandler) GetEpusdtPaymentURL(c *gin.Context) {
	// 从查询参数获取必要信息
	amount := c.Query("amount")
	orderID := c.Query("order_id")
	title := c.Query("title") // 支付链接标题

	if amount == "" || orderID == "" || title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少必要参数: amount, order_id, title"})
		return
	}

	// 创建Epusdt订单（使用随机中间人钱包地址）
	response, err := h.paymentService.CreateDirectEpusdtOrderWithRandomWallet(orderID, amount, title, c.ClientIP())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("创建支付订单失败: %v", err)})
		return
	}

	// 返回付款页面URL
	c.JSON(http.StatusOK, gin.H{
		"payment_url": response.Data.PaymentURL,
		"trade_id":    response.Data.TradeID,
		"order_id":    response.Data.OrderID,
		"amount":      response.Data.Amount,
		"expires_at":  response.Data.ExpirationTime,
	})
}
