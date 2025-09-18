package handler

import (
	"log"
	"net/http"
	"payment-link-mvp/internal/model"
	"payment-link-mvp/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentLinkHandler struct {
	paymentLinkService *service.PaymentLinkService
}

func NewPaymentLinkHandler(paymentLinkService *service.PaymentLinkService) *PaymentLinkHandler {
	return &PaymentLinkHandler{paymentLinkService: paymentLinkService}
}

func (h *PaymentLinkHandler) Create(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	var req service.CreatePaymentLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	userModel := user.(model.User)
	paymentLink, err := h.paymentLinkService.Create(userModel.ID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "创建成功",
		"payment_link": paymentLink,
	})
}

func (h *PaymentLinkHandler) List(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	userModel := user.(model.User)
	paymentLinks, err := h.paymentLinkService.List(userModel.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payment_links": paymentLinks,
	})
}

func (h *PaymentLinkHandler) Get(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	userModel := user.(model.User)
	paymentLink, err := h.paymentLinkService.Get(userModel.ID, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "支付链接不存在"})
		return
	}

	c.JSON(http.StatusOK, paymentLink)
}

func (h *PaymentLinkHandler) Delete(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	userModel := user.(model.User)
	if err := h.paymentLinkService.Delete(userModel.ID, uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func (h *PaymentLinkHandler) GetPublicPayment(c *gin.Context) {
	linkID := c.Param("id")
	if linkID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的链接ID"})
		return
	}

	// 添加调试日志
	log.Printf("查找支付链接，LinkID: %s", linkID)

	paymentLink, err := h.paymentLinkService.GetByLinkID(linkID)
	if err != nil {
		log.Printf("查找支付链接失败: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "支付链接不存在或已失效"})
		return
	}

	log.Printf("找到支付链接: ID=%d, Title=%s, Status=%s", paymentLink.ID, paymentLink.Title, paymentLink.Status)

	c.JSON(http.StatusOK, gin.H{
		"payment_link": paymentLink,
	})
}
