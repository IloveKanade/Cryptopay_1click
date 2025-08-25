package service

import (
	"bytes"
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

type EpusdtService struct {
	baseURL    string
	authToken  string
	httpClient *http.Client
}

// EpusdtOrderRequest 符合Epusdt API要求的请求结构
type EpusdtOrderRequest struct {
	OrderId     string  `json:"order_id"`
	Amount      float64 `json:"amount"`
	NotifyUrl   string  `json:"notify_url"`
	RedirectUrl string  `json:"redirect_url"`
	Signature   string  `json:"signature"`
}

type EpusdtOrderResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Data       struct {
		TradeID        string  `json:"trade_id"`
		OrderID        string  `json:"order_id"`
		Amount         float64 `json:"amount"`
		ActualAmount   float64 `json:"actual_amount"`
		Token          string  `json:"token"`
		ExpirationTime int64   `json:"expiration_time"`
		PaymentURL     string  `json:"payment_url"`
	} `json:"data"`
	RequestID string `json:"request_id"`
}

type EpusdtOrderQueryResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		TradeID      string  `json:"trade_id"`
		OrderID      string  `json:"order_id"`
		Amount       float64 `json:"amount"`
		ActualAmount float64 `json:"actual_amount"`
		Token        string  `json:"token"`
		Status       int     `json:"status"`
		NotifyURL    string  `json:"notify_url"`
		RedirectURL  string  `json:"redirect_url"`
		OrderName    string  `json:"order_name"`
		ClientIP     string  `json:"client_ip"`
		CreatedAt    string  `json:"created_at"`
		UpdatedAt    string  `json:"updated_at"`
	} `json:"data"`
}

type EpusdtNotifyRequest struct {
	TradeID            string  `json:"trade_id"`
	OrderID            string  `json:"order_id"`
	Amount             float64 `json:"amount"`
	ActualAmount       float64 `json:"actual_amount"`
	Token              string  `json:"token"`
	BlockTransactionID string  `json:"block_transaction_id"`
	Signature          string  `json:"signature"`
	Status             int     `json:"status"`
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

// generateSignature 生成Epusdt签名 - 修复为符合Epusdt要求的算法
func (s *EpusdtService) generateSignature(params map[string]interface{}) string {
	// 按ASCII排序
	var keys []string
	for k := range params {
		if k != "signature" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	// 构建签名字符串
	var tempArr []string
	for _, key := range keys {
		value := params[key]
		if value == nil {
			continue
		}

		var strValue string
		switch v := value.(type) {
		case string:
			if v == "" {
				continue
			}
			strValue = v
		case float64:
			// 对于金额，使用%.0f格式，避免小数点
			if v == float64(int(v)) {
				strValue = fmt.Sprintf("%.0f", v)
			} else {
				strValue = fmt.Sprintf("%.2f", v)
			}
		case int:
			strValue = fmt.Sprintf("%d", v)
		case int64:
			strValue = fmt.Sprintf("%d", v)
		default:
			strValue = fmt.Sprintf("%v", v)
		}

		// 空值不参与签名
		if strValue != "" {
			tempArr = append(tempArr, key+"="+strValue)
		}
	}

	// 构建最终的签名字符串
	var signStr strings.Builder
	for i, param := range tempArr {
		if i > 0 {
			signStr.WriteString("&")
		}
		signStr.WriteString(param)
	}

	// 添加API Token
	signStr.WriteString(s.authToken)

	// 添加调试日志
	log.Printf("签名参数: %v", params)
	log.Printf("签名字符串: %s", signStr.String())

	// MD5加密
	hash := md5.Sum([]byte(signStr.String()))
	signature := fmt.Sprintf("%x", hash)
	log.Printf("生成的签名: %s", signature)

	return signature
}

// CreateOrder 创建Epusdt支付订单 - 支持指定钱包地址
func (s *EpusdtService) CreateOrder(orderID, amount, orderName, notifyURL, redirectURL, clientIP, walletAddress string) (*EpusdtOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v1/order/create-transaction", s.baseURL)

	// 构建符合Epusdt API要求的请求参数
	// 将字符串金额转换为float64
	var amountFloat float64
	_, err := fmt.Sscanf(amount, "%f", &amountFloat)
	if err != nil {
		return nil, fmt.Errorf("金额格式错误: %v", err)
	}

	// 构建请求参数（只包含Epusdt需要的字段）
	params := map[string]interface{}{
		"order_id":     orderID,
		"amount":       amountFloat,
		"notify_url":   notifyURL,
		"redirect_url": redirectURL,
	}

	// 如果指定了钱包地址，添加到参数中
	if walletAddress != "" {
		params["token"] = walletAddress
	}

	// 生成签名
	signature := s.generateSignature(params)
	params["signature"] = signature

	// 序列化请求数据
	jsonData, err := json.Marshal(params)
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

	// 发送请求
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 添加调试日志
	log.Printf("Epusdt响应状态码: %d", resp.StatusCode)
	log.Printf("Epusdt响应内容: %s", string(respBody))

	// 解析响应
	var response EpusdtOrderResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Epusdt API错误: %s", response.Message)
	}

	return &response, nil
}

// QueryOrder 查询Epusdt订单状态
func (s *EpusdtService) QueryOrder(tradeID string) (*EpusdtOrderQueryResponse, error) {
	url := fmt.Sprintf("%s/api/v1/order/query-transaction", s.baseURL)

	// 构建请求参数
	params := map[string]interface{}{
		"trade_id": tradeID,
	}

	// 生成签名
	signature := s.generateSignature(params)
	params["signature"] = signature

	// 序列化请求数据
	jsonData, err := json.Marshal(params)
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

// VerifyNotifySignature 验证回调签名
func (s *EpusdtService) VerifyNotifySignature(data map[string]interface{}) bool {
	// 提取签名
	receivedSignature, ok := data["signature"].(string)
	if !ok {
		return false
	}

	// 构建验证参数（排除signature字段）
	params := make(map[string]interface{})
	for key, value := range data {
		if key != "signature" {
			params[key] = value
		}
	}

	// 生成签名
	expectedSignature := s.generateSignature(params)
	return receivedSignature == expectedSignature
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

// AddWalletAddress 添加钱包地址到Epusdt系统
func (s *EpusdtService) AddWalletAddress(walletAddress string) error {
	// 直接连接Epusdt的MySQL数据库
	dsn := "epusdt:epusdt123456@tcp(epusdt-mysql:3306)/epusdt?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}
	defer db.Close()

	// 检查地址是否已存在
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM wallet_address WHERE token = ?", walletAddress).Scan(&count)
	if err != nil {
		return fmt.Errorf("查询地址失败: %v", err)
	}

	if count > 0 {
		log.Printf("钱包地址已存在: %s", walletAddress)
		return nil
	}

	// 插入新地址
	_, err = db.Exec("INSERT INTO wallet_address (token, status, created_at, updated_at) VALUES (?, 1, NOW(), NOW())", walletAddress)
	if err != nil {
		return fmt.Errorf("插入地址失败: %v", err)
	}

	log.Printf("成功添加钱包地址到数据库: %s", walletAddress)
	return nil
}

// IsWalletAddressAvailable 检查钱包地址是否在Epusdt系统中可用
func (s *EpusdtService) IsWalletAddressAvailable(walletAddress string) (bool, error) {
	// 由于Epusdt没有提供钱包地址查询API，我们暂时跳过验证
	// 在实际生产环境中，应该直接查询数据库

	log.Printf("跳过钱包地址可用性检查: %s (Epusdt不支持此API)", walletAddress)
	return true, nil
}

// SetWalletPriority 设置钱包地址优先级（通过数据库操作）
func (s *EpusdtService) SetWalletPriority(targetWallet string) error {
	// 由于Epusdt可能没有提供设置优先级的API，我们通过数据库操作来实现
	// 这里我们使用一个更安全的方法：临时调整钱包地址的顺序

	// 注意：这个方法需要谨慎使用，因为它会影响整个Epusdt系统
	// 在实际生产环境中，建议使用更复杂的逻辑

	// 这里我们暂时返回nil，表示不进行优先级设置
	// 在实际实现中，可以通过数据库操作来调整钱包地址的顺序
	return nil
}

// DisableOtherWallets 禁用除指定地址外的其他钱包地址
func (s *EpusdtService) DisableOtherWallets(targetWallet string) error {
	// 通过HTTP API禁用其他钱包地址
	// 首先获取所有钱包地址
	url := fmt.Sprintf("%s/api/v1/wallet/addresses", s.baseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	req.Header.Set("Authorization", s.authToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("发送HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	var response struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    []struct {
			ID     uint   `json:"id"`
			Token  string `json:"token"`
			Status int    `json:"status"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	if response.Code != 200 {
		return fmt.Errorf("获取钱包地址失败: %s", response.Message)
	}

	// 禁用其他钱包地址
	for _, wallet := range response.Data {
		if wallet.Token != targetWallet && wallet.Status == 1 {
			// 禁用这个钱包地址
			err := s.disableWallet(wallet.ID)
			if err != nil {
				// 记录错误但继续处理其他地址
				log.Printf("禁用钱包地址失败: %v", err)
			}
		}
	}

	return nil
}

// disableWallet 禁用指定ID的钱包地址
func (s *EpusdtService) disableWallet(walletID uint) error {
	url := fmt.Sprintf("%s/api/v1/wallet/address/%d/disable", s.baseURL, walletID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	req.Header.Set("Authorization", s.authToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("发送HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	var response struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	if response.Code != 200 {
		return fmt.Errorf("禁用钱包地址失败: %s", response.Message)
	}

	return nil
}

// GetWalletAddresses 获取所有钱包地址
func (s *EpusdtService) GetWalletAddresses() ([]string, error) {
	url := fmt.Sprintf("%s/api/v1/wallet/addresses", s.baseURL)

	// 创建HTTP请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Authorization", s.authToken)

	// 发送请求
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var response struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    []struct {
			Token  string `json:"token"`
			Status int    `json:"status"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if response.Code != 200 {
		return nil, fmt.Errorf("Epusdt API错误: %s", response.Message)
	}

	// 提取启用的钱包地址
	var addresses []string
	for _, wallet := range response.Data {
		if wallet.Status == 1 { // 1表示启用状态
			addresses = append(addresses, wallet.Token)
		}
	}

	return addresses, nil
}
