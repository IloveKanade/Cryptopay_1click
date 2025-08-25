package model

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 全局数据库连接
var DB *gorm.DB

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func InitDB(config DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)

	// 设置GORM日志级别
	logLevel := logger.Info
	if os.Getenv("GIN_MODE") == "release" {
		logLevel = logger.Error
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	// 获取底层的sql.DB对象
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	DB = db
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	// 自动迁移数据库表
	err := db.AutoMigrate(
		&User{},
		&PaymentLink{},
		&PaymentOrder{},
		&MiddlemanWallet{},
		&TransferRecord{},
		&FeeConfig{},
	)
	if err != nil {
		return err
	}

	// 初始化默认费率配置
	var feeConfig FeeConfig
	if err := db.First(&feeConfig).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 创建默认费率配置
			defaultFeeConfig := FeeConfig{
				FeeRate:      0.01,  // 1%
				MinFee:       0.1,   // 0.1 USDT
				MaxFee:       10,    // 10 USDT
				NetworkFeeMin: 0.001, // 0.001 USDT
				NetworkFeeMax: 0.01,  // 0.01 USDT
				IsActive:     true,
			}
			if err := db.Create(&defaultFeeConfig).Error; err != nil {
				log.Printf("创建默认费率配置失败: %v", err)
			}
		}
	}

	// 初始化默认管理员用户
	var adminUser User
	if err := db.Where("email = ?", "admin@example.com").First(&adminUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 创建默认管理员用户
			adminUser := User{
				Email: "admin@example.com",
				Name:  "系统管理员",
				Role:  "admin",
			}

			// 设置默认密码
			if err := adminUser.SetPassword("admin123"); err != nil {
				log.Printf("设置管理员密码失败: %v", err)
			} else {
				if err := db.Create(&adminUser).Error; err != nil {
					log.Printf("创建默认管理员用户失败: %v", err)
				} else {
					log.Println("已创建默认管理员用户: admin@example.com / admin123")
				}
			}
		}
	}

	log.Println("数据库表迁移完成")
	return nil
}
