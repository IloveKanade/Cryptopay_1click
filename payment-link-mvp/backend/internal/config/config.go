package config

import (
	"os"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Epusdt   EpusdtConfig
	Callback CallbackConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type JWTConfig struct {
	Secret string
}

type EpusdtConfig struct {
	BaseURL   string
	AuthToken string
}

type CallbackConfig struct {
	TransferSuccessURL string // 转账成功回调地址
	TransferFailURL    string // 转账失败回调地址
	NotifySecret       string // 回调签名密钥
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "epusdt"),
			Password: getEnv("DB_PASSWORD", "epusdt123456"),
			DBName:   getEnv("DB_NAME", "epusdt"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your-secret-key"),
		},
		Epusdt: EpusdtConfig{
			BaseURL:   getEnv("EPUSDT_BASE_URL", "http://localhost:8000"),
			AuthToken: getEnv("EPUSDT_AUTH_TOKEN", "epusdt_password_xasddawqe"),
		},
		Callback: CallbackConfig{
			TransferSuccessURL: getEnv("CALLBACK_SUCCESS_URL", ""),
			TransferFailURL:    getEnv("CALLBACK_FAIL_URL", ""),
			NotifySecret:       getEnv("CALLBACK_SECRET", "callback_secret_key"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
