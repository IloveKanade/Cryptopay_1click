package service

import (
	"errors"
	"payment-link-mvp/internal/config"
	"payment-link-mvp/internal/model"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}

func (s *UserService) Register(req RegisterRequest) error {
	// 检查邮箱是否已存在
	var existingUser model.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return errors.New("邮箱已被注册")
	}

	// 创建新用户
	user := model.User{
		Email: req.Email,
		Name:  req.Name,
		Role:  "user", // 默认为普通用户
	}

	if err := user.SetPassword(req.Password); err != nil {
		return err
	}

	return s.db.Create(&user).Error
}

func (s *UserService) Login(req LoginRequest) (*LoginResponse, error) {
	var user model.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	if !user.CheckPassword(req.Password) {
		return nil, errors.New("密码错误")
	}

	// 生成JWT令牌，包含用户角色信息
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role, // 添加角色信息
	})

	tokenString, err := token.SignedString([]byte(config.Load().JWT.Secret))
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: tokenString,
		User:  user,
	}, nil
}

func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateAdminUser 创建管理员用户
func (s *UserService) CreateAdminUser(req RegisterRequest) error {
	// 检查邮箱是否已存在
	var existingUser model.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return errors.New("邮箱已被注册")
	}

	// 创建管理员用户
	user := model.User{
		Email: req.Email,
		Name:  req.Name,
		Role:  "admin", // 设置为管理员角色
	}

	if err := user.SetPassword(req.Password); err != nil {
		return err
	}

	return s.db.Create(&user).Error
}

// GetUserByEmail 根据邮箱获取用户
func (s *UserService) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
