package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type EpusdtService struct {
	baseURL    string
	authToken  string
	httpClient *http.Client
}

type EpusdtOrderRequest struct {
	Amount      string `json:"amount"`
	OrderID     string `json:"order_id"`
	NotifyURL   string `json:"notify_url"`
	RedirectURL string `json:"redirect_url"`
	OrderName   string `json:"order_name"`
	ClientIP    string `json:"client_ip"`
	Signature   string `json:"signature"`
	Timestamp   int64  `json:"timestamp"`
}

type EpusdtOrderResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		TradeID      string `json:"trade_id"`
		OrderID      string `json:"order_id"`
		Amount       string `json:"amount"`
		ActualAmount string `json:"actual_amount"`
		Token        string `json:"token"`
		BlockNumber  string `json:"block_number"`
		Signature    string `json:"signature"`
		Status       int    `json:"status"`
		URL          string `json:"url"`
	} `json:"data"`
}

type EpusdtOrderQueryResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		TradeID      string `json:"trade_id"`
		OrderID      string `json:"order_id"`
		Amount       string `json:"amount"`
		ActualAmount string `json:"actual_amount"`
		Token        string `json:"token"`
		BlockNumber  string `json:"block_number"`
		Status       int    `json:"status"`
		NotifyURL    string `json:"notify_url"`
		RedirectURL  string `json:"redirect_url"`
		OrderName    string `json:"order_name"`
		ClientIP     string `json:"client_ip"`
		CreatedAt    string `json:"created_at"`
		UpdatedAt    string `json:"updated_at"`
	} `json:"data"`
}

func NewEpusdtService(baseURL, authToken string) *EpusdtService {
	return &EpusdtService{
		baseURL:   baseURL,
		authToken: authToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CreateOrder 创建Epusdt支付订单
func (s *EpusdtService) CreateOrder(orderID, amount, orderName, notifyURL, redirectURL, clientIP string) (*EpusdtOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v1/order/create-transaction", s.baseURL)

	// 生成时间戳
	timestamp := time.Now().Unix()

	// 构建请求数据
	requestData := EpusdtOrderRequest{
		Amount:      amount,
		OrderID:     orderID,
		NotifyURL:   notifyURL,
		RedirectURL: redirectURL,
		OrderName:   orderName,
		ClientIP:    clientIP,
		Timestamp:   timestamp,
	}

	// 生成签名（这里需要根据Epusdt的签名算法实现）
	// signature := s.generateSignature(requestData)
	// requestData.Signature = signature

	// 序列化请求数据
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("序列化请求数据失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", s.authToken)

	// 发送请求
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var response EpusdtOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if response.Code != 200 {
		return nil, fmt.Errorf("Epusdt API错误: %s", response.Message)
	}

	return &response, nil
}

// QueryOrder 查询Epusdt订单状态
func (s *EpusdtService) QueryOrder(tradeID string) (*EpusdtOrderQueryResponse, error) {
	url := fmt.Sprintf("%s/api/v1/order/query-transaction", s.baseURL)

	// 构建请求数据
	requestData := map[string]string{
		"trade_id": tradeID,
	}

	// 序列化请求数据
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("序列化请求数据失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", s.authToken)

	// 发送请求
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var response EpusdtOrderQueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if response.Code != 200 {
		return nil, fmt.Errorf("Epusdt API错误: %s", response.Message)
	}

	return &response, nil
}

// CloseExpiredOrder 关闭过期订单
func (s *EpusdtService) CloseExpiredOrder(tradeID string) error {
	url := fmt.Sprintf("%s/api/v1/order/close-transaction", s.baseURL)

	// 构建请求数据
	requestData := map[string]string{
		"trade_id": tradeID,
	}

	// 序列化请求数据
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("序列化请求数据失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", s.authToken)

	// 发送请求
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("发送HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var response struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	if response.Code != 200 {
		return fmt.Errorf("Epusdt API错误: %s", response.Message)
	}

	return nil
}

// GetWalletAddress 获取钱包地址
func (s *EpusdtService) GetWalletAddress() (string, error) {
	url := fmt.Sprintf("%s/api/v1/wallet/address", s.baseURL)

	// 创建HTTP请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Authorization", s.authToken)

	// 发送请求
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var response struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Address string `json:"address"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if response.Code != 200 {
		return "", fmt.Errorf("Epusdt API错误: %s", response.Message)
	}

	return response.Data.Address, nil
}
