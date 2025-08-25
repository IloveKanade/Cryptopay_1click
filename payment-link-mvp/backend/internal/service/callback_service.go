package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"payment-link-mvp/internal/config"
)

type CallbackService struct {
	config *config.Config
	client *http.Client
}

type TransferCallbackData struct {
	TransferID    string  `json:"transfer_id"`
	FromAddress   string  `json:"from_address"`
	ToAddress     string  `json:"to_address"`
	Amount        float64 `json:"amount"`
	NetworkFee    float64 `json:"network_fee"`
	TotalFee      float64 `json:"total_fee"`
	NetAmount     float64 `json:"net_amount"`
	TransactionHash string `json:"transaction_hash"`
	Status        string  `json:"status"` // success, failed
	ErrorMessage  string  `json:"error_message,omitempty"`
	Timestamp     int64   `json:"timestamp"`
}

func NewCallbackService(cfg *config.Config) *CallbackService {
	return &CallbackService{
		config: cfg,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NotifyTransferSuccess 通知转账成功
func (s *CallbackService) NotifyTransferSuccess(data TransferCallbackData) error {
	if s.config.Callback.TransferSuccessURL == "" {
		fmt.Printf("转账成功回调地址未配置，跳过通知\n")
		return nil
	}

	data.Status = "success"
	data.Timestamp = time.Now().Unix()

	return s.sendCallback(s.config.Callback.TransferSuccessURL, data)
}

// NotifyTransferFail 通知转账失败
func (s *CallbackService) NotifyTransferFail(data TransferCallbackData) error {
	if s.config.Callback.TransferFailURL == "" {
		fmt.Printf("转账失败回调地址未配置，跳过通知\n")
		return nil
	}

	data.Status = "failed"
	data.Timestamp = time.Now().Unix()

	return s.sendCallback(s.config.Callback.TransferFailURL, data)
}

// sendCallback 发送回调请求
func (s *CallbackService) sendCallback(url string, data TransferCallbackData) error {
	// 序列化数据
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("序列化回调数据失败: %w", err)
	}

	// 生成签名
	signature := s.generateSignature(jsonData)

	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建回调请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Callback-Signature", signature)
	req.Header.Set("X-Callback-Timestamp", fmt.Sprintf("%d", data.Timestamp))

	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("发送回调请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Printf("回调通知发送成功: %s, 状态码: %d\n", url, resp.StatusCode)
		return nil
	} else {
		return fmt.Errorf("回调通知发送失败: %s, 状态码: %d", url, resp.StatusCode)
	}
}

// generateSignature 生成回调签名
func (s *CallbackService) generateSignature(data []byte) string {
	h := hmac.New(sha256.New, []byte(s.config.Callback.NotifySecret))
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// VerifyCallbackSignature 验证回调签名
func (s *CallbackService) VerifyCallbackSignature(data []byte, signature string) bool {
	expectedSignature := s.generateSignature(data)
	return signature == expectedSignature
}
