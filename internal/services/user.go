package services

import (
	"context"
	"errors"
	"fmt"
	"time"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/galaxyerp/galaxyErp/internal/models"
	"github.com/galaxyerp/galaxyErp/internal/repositories"
	"github.com/galaxyerp/galaxyErp/internal/dto"
)

// UserService 用户服务接口
type UserService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.LoginResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	GetProfile(ctx context.Context, userID uint) (*dto.UserProfileResponse, error)
	UpdateProfile(ctx context.Context, userID uint, req *dto.UserUpdateRequest) error
	ChangePassword(ctx context.Context, userID uint, req *dto.ChangePasswordRequest) error
	GetUsers(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.UserResponse], error)
	SearchUsers(ctx context.Context, req *dto.SearchRequest) (*dto.PaginatedResponse[dto.UserResponse], error)
	DeleteUser(ctx context.Context, userID uint) error
}

// UserServiceImpl 用户服务实现
type UserServiceImpl struct {
	userRepo repositories.UserRepository
	jwtSecret string
	jwtExpiryHours int
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo repositories.UserRepository, jwtSecret string, jwtExpiryHours int) UserService {
	return &UserServiceImpl{
		userRepo: userRepo,
		jwtSecret: jwtSecret,
		jwtExpiryHours: jwtExpiryHours,
	}
}

// Register 用户注册
func (s *UserServiceImpl) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.LoginResponse, error) {
	// 检查用户是否已存在
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("检查用户失败: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("邮箱已被注册")
	}

	// 检查用户名是否已存在
	existingUser, err = s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("检查用户名失败: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("用户名已被使用")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	// 创建用户
	user := &models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		IsActive:  true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	// 生成JWT令牌
	token, err := s.generateJWT(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("生成令牌失败: %w", err)
	}

	// 构建响应
	userResponse := dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	expiryDuration := time.Duration(s.jwtExpiryHours) * time.Hour
	return &dto.LoginResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(expiryDuration),
		User:      userResponse,
	}, nil
}

// Login 用户登录
func (s *UserServiceImpl) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// 根据用户名查找用户
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("查找用户失败: %w", err)
	}
	if user == nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查用户是否激活
	if !user.IsActive {
		return nil, errors.New("用户账户已被禁用")
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLoginAt = &now
	if err := s.userRepo.Update(ctx, user); err != nil {
		// 记录错误但不影响登录
		fmt.Printf("更新最后登录时间失败: %v\n", err)
	}

	// 生成JWT令牌
	token, err := s.generateJWT(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("生成令牌失败: %w", err)
	}

	// 构建响应
	userResponse := dto.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Phone:       user.Phone,
		IsActive:    user.IsActive,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	return &dto.LoginResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		User:      userResponse,
	}, nil
}

// GetProfile 获取用户资料
func (s *UserServiceImpl) GetProfile(ctx context.Context, userID uint) (*dto.UserProfileResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}

	// 构建用户资料响应
	profile := &dto.UserProfileResponse{
		UserResponse: dto.UserResponse{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Phone:       user.Phone,
			IsActive:    user.IsActive,
			LastLoginAt: user.LastLoginAt,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		},
		Permissions: []string{}, // TODO: 实现权限获取
		MenuItems:   []string{}, // TODO: 实现菜单项获取
	}

	return profile, nil
}

// UpdateProfile 更新用户资料
func (s *UserServiceImpl) UpdateProfile(ctx context.Context, userID uint, req *dto.UserUpdateRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("获取用户失败: %w", err)
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	// 更新字段
	if req.Username != "" {
		// 检查用户名是否已被其他用户使用
		existingUser, err := s.userRepo.GetByUsername(ctx, req.Username)
		if err != nil {
			return fmt.Errorf("检查用户名失败: %w", err)
		}
		if existingUser != nil && existingUser.ID != userID {
			return errors.New("用户名已被使用")
		}
		user.Username = req.Username
	}

	if req.Email != "" {
		// 检查邮箱是否已被其他用户使用
		existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
		if err != nil {
			return fmt.Errorf("检查邮箱失败: %w", err)
		}
		if existingUser != nil && existingUser.ID != userID {
			return errors.New("邮箱已被使用")
		}
		user.Email = req.Email
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.DepartmentID != nil {
		user.DepartmentID = req.DepartmentID
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	return s.userRepo.Update(ctx, user)
}

// ChangePassword 修改密码
func (s *UserServiceImpl) ChangePassword(ctx context.Context, userID uint, req *dto.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("获取用户失败: %w", err)
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return errors.New("旧密码错误")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	user.Password = string(hashedPassword)
	return s.userRepo.Update(ctx, user)
}

// GetUsers 获取用户列表
func (s *UserServiceImpl) GetUsers(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.UserResponse], error) {
	offset := req.GetOffset()
	limit := req.GetLimit()

	users, total, err := s.userRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("获取用户列表失败: %w", err)
	}

	// 转换为响应格式
	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = dto.UserResponse{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Phone:       user.Phone,
			IsActive:    user.IsActive,
			LastLoginAt: user.LastLoginAt,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		}
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &dto.PaginatedResponse[dto.UserResponse]{
		Data:       userResponses,
		Total:      total,
		Page:       req.Page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// SearchUsers 搜索用户
func (s *UserServiceImpl) SearchUsers(ctx context.Context, req *dto.SearchRequest) (*dto.PaginatedResponse[dto.UserResponse], error) {
	offset := req.GetOffset()
	limit := req.GetLimit()

	users, total, err := s.userRepo.Search(ctx, req.Keyword, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("搜索用户失败: %w", err)
	}

	// 转换为响应格式
	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = dto.UserResponse{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Phone:       user.Phone,
			IsActive:    user.IsActive,
			LastLoginAt: user.LastLoginAt,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		}
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &dto.PaginatedResponse[dto.UserResponse]{
		Data:       userResponses,
		Total:      total,
		Page:       req.Page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// DeleteUser 删除用户
func (s *UserServiceImpl) DeleteUser(ctx context.Context, userID uint) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("获取用户失败: %w", err)
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	return s.userRepo.Delete(ctx, userID)
}

// generateJWT 生成JWT令牌
func (s *UserServiceImpl) generateJWT(userID uint, username string) (string, error) {
	expiryDuration := time.Duration(s.jwtExpiryHours) * time.Hour
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(expiryDuration).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}