package services

import (
	"context"
	"fmt"
	"github.com/galaxyerp/galaxyErp/internal/dto"
	"github.com/galaxyerp/galaxyErp/internal/models"
	"github.com/galaxyerp/galaxyErp/internal/repositories"
	"github.com/galaxyerp/galaxyErp/internal/common"
	"github.com/galaxyerp/galaxyErp/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
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
	*BaseService
	userRepo        repositories.UserRepository
	auditLogService AuditLogService
	jwtSecret       string
	jwtExpiryHours  int
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo repositories.UserRepository, auditLogService AuditLogService, jwtSecret string, jwtExpiryHours int) UserService {
	config := &BaseServiceConfig{
		EnableValidation: true,
		EnableCache:      true,
		CacheExpiry:      time.Hour * 24, // 24小时缓存
		EnableAudit:      true,
		EnableMetrics:    true,
	}
	
	return &UserServiceImpl{
		BaseService:     NewBaseService(config),
		userRepo:        userRepo,
		auditLogService: auditLogService,
		jwtSecret:       jwtSecret,
		jwtExpiryHours:  jwtExpiryHours,
	}
}

// Register 用户注册
func (s *UserServiceImpl) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.LoginResponse, error) {
	// 使用BaseService验证请求
	if err := s.ValidateRequest(ctx, "user", req); err != nil {
		return nil, err
	}

	// 检查用户是否已存在
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		appErr := common.NewAppErrorFromTypeWithCause("database", "USER_CHECK_FAILED", "检查用户失败", err)
		common.LogAppError(appErr, "user_register", utils.String("email", req.Email))
		return nil, appErr
	}
	if existingUser != nil {
		appErr := common.NewAppErrorFromTypeWithDetails("business", "EMAIL_EXISTS", "邮箱已被注册", fmt.Sprintf("邮箱 %s 已被注册", req.Email))
		common.LogAppError(appErr, "user_register", utils.String("email", req.Email))
		return nil, appErr
	}

	// 检查用户名是否已存在
	existingUser, err = s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		appErr := common.NewAppErrorFromTypeWithCause("database", "USERNAME_CHECK_FAILED", "检查用户名失败", err)
		common.LogAppError(appErr, "user_register", utils.String("username", req.Username))
		return nil, appErr
	}
	if existingUser != nil {
		appErr := common.NewAppErrorFromTypeWithDetails("business", "USERNAME_EXISTS", "用户名已被使用", fmt.Sprintf("用户名 %s 已被使用", req.Username))
		common.LogAppError(appErr, "user_register", utils.String("username", req.Username))
		return nil, appErr
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		appErr := common.NewAppErrorFromTypeWithCause("system", "PASSWORD_HASH_FAILED", "密码加密失败", err)
		common.LogAppError(appErr, "user_register", utils.String("username", req.Username))
		return nil, appErr
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
		appErr := common.NewAppErrorFromTypeWithCause("database", "USER_CREATE_FAILED", "创建用户失败", err)
		common.LogAppError(appErr, "user_register",
			utils.String("username", req.Username),
			utils.String("email", req.Email))
		return nil, appErr
	}

	// 记录用户注册成功的业务日志
	utils.Info("用户注册成功",
		utils.Uint("user_id", user.ID),
		utils.String("username", user.Username),
		utils.String("email", user.Email),
		utils.String("operation", "user_register"),
	)

	// 记录审计日志
	if err := s.auditLogService.LogAction(ctx, user.ID, user.Username, "CREATE", "USER", fmt.Sprintf("%d", user.ID), fmt.Sprintf("用户注册: %s", user.Username), nil, user); err != nil {
		utils.LogError("记录审计日志失败", utils.ErrorField(err))
	}

	// 生成JWT令牌
	token, err := s.generateJWT(user.ID, user.Username)
	if err != nil {
		appErr := common.NewAppErrorFromTypeWithCause("system", "JWT_GENERATE_FAILED", "生成令牌失败", err)
		common.LogAppError(appErr, "user_register", utils.Uint("user_id", user.ID))
		return nil, appErr
	}

	// 记录用户登录成功的业务日志
	utils.Info("用户登录成功",
		utils.Uint("user_id", user.ID),
		utils.String("username", user.Username),
		utils.String("operation", "user_login"),
	)

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
		appErr := common.NewAppErrorFromTypeWithCause("database", "USER_LOOKUP_FAILED", "查找用户失败", err)
		common.LogAppError(appErr, "user_login", utils.String("username", req.Username))
		return nil, appErr
	}
	if user == nil {
		appErr := common.NewAppErrorFromType("auth", "INVALID_CREDENTIALS", "用户名或密码错误")
		common.LogAppError(appErr, "user_login", utils.String("username", req.Username))
		return nil, appErr
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		appErr := common.NewAppErrorFromType("auth", "INVALID_CREDENTIALS", "用户名或密码错误")
		common.LogAppError(appErr, "user_login",
			utils.String("username", req.Username),
			utils.Uint("user_id", user.ID))
		return nil, appErr
	}

	// 检查用户是否激活
	if !user.IsActive {
		appErr := common.NewAppErrorFromType("auth", "USER_DISABLED", "用户账户已被禁用")
		common.LogAppError(appErr, "user_login",
			utils.String("username", req.Username),
			utils.Uint("user_id", user.ID))
		return nil, appErr
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLoginAt = &now
	if err := s.userRepo.Update(ctx, user); err != nil {
		// 记录错误但不影响登录
		appErr := common.NewAppErrorFromTypeWithCause("database", "UPDATE_LOGIN_TIME_FAILED", "更新最后登录时间失败", err)
		common.LogAppError(appErr, "user_login", utils.Uint("user_id", user.ID))
	}

	// 生成JWT令牌
	token, err := s.generateJWT(user.ID, user.Username)
	if err != nil {
		appErr := common.NewAppErrorFromTypeWithCause("system", "JWT_GENERATE_FAILED", "生成令牌失败", err)
		common.LogAppError(appErr, "user_login", utils.Uint("user_id", user.ID))
		return nil, appErr
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
	// 获取用户信息
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		appErr := common.NewAppErrorFromTypeWithCause("database", "USER_GET_FAILED", "获取用户失败", err)
		common.LogAppError(appErr, "user_get_profile", utils.Uint("user_id", userID))
		return nil, appErr
	}
	if user == nil {
		appErr := common.NewAppErrorFromType("business", "USER_NOT_FOUND", "用户不存在")
		common.LogAppError(appErr, "user_get_profile", utils.Uint("user_id", userID))
		return nil, appErr
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
	// 获取原始用户数据
	oldUser, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		appErr := common.NewAppErrorFromTypeWithCause("database", "USER_GET_FAILED", "获取用户失败", err)
		common.LogAppError(appErr, "user_update_profile", utils.Uint("user_id", userID))
		return appErr
	}
	if oldUser == nil {
		appErr := common.NewAppErrorFromType("business", "USER_NOT_FOUND", "用户不存在")
		common.LogAppError(appErr, "user_update_profile", utils.Uint("user_id", userID))
		return appErr
	}

	// 创建用户副本用于更新
	user := *oldUser

	// 更新字段
	if req.Username != "" {
		// 检查用户名是否已被其他用户使用
		existingUser, err := s.userRepo.GetByUsername(ctx, req.Username)
		if err != nil {
			appErr := common.NewAppErrorFromTypeWithCause("database", "USERNAME_CHECK_FAILED", "检查用户名失败", err)
			common.LogAppError(appErr, "user_update_profile", 
				utils.Uint("user_id", userID),
				utils.String("username", req.Username))
			return appErr
		}
		if existingUser != nil && existingUser.ID != userID {
			appErr := common.NewAppErrorFromTypeWithDetails("business", "USERNAME_EXISTS", "用户名已被使用", 
				fmt.Sprintf("用户名 %s 已被其他用户使用", req.Username))
			common.LogAppError(appErr, "user_update_profile", 
				utils.Uint("user_id", userID),
				utils.String("username", req.Username))
			return appErr
		}
		user.Username = req.Username
	}

	if req.Email != "" {
		// 检查邮箱是否已被其他用户使用
		existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
		if err != nil {
			appErr := common.NewAppErrorFromTypeWithCause("database", "EMAIL_CHECK_FAILED", "检查邮箱失败", err)
			common.LogAppError(appErr, "user_update_profile", 
				utils.Uint("user_id", userID),
				utils.String("email", req.Email))
			return appErr
		}
		if existingUser != nil && existingUser.ID != userID {
			appErr := common.NewAppErrorFromTypeWithDetails("business", "EMAIL_EXISTS", "邮箱已被使用", 
				fmt.Sprintf("邮箱 %s 已被其他用户使用", req.Email))
			common.LogAppError(appErr, "user_update_profile", 
				utils.Uint("user_id", userID),
				utils.String("email", req.Email))
			return appErr
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

	// 更新用户
	if err := s.userRepo.Update(ctx, &user); err != nil {
		appErr := common.NewAppErrorFromTypeWithCause("database", "USER_UPDATE_FAILED", "更新用户失败", err)
		common.LogAppError(appErr, "user_update_profile", utils.Uint("user_id", userID))
		return appErr
	}

	// 记录审计日志
	if err := s.auditLogService.LogAction(ctx, userID, user.Username, "UPDATE", "USER", fmt.Sprintf("%d", userID), fmt.Sprintf("用户信息更新: %s", user.Username), oldUser, user); err != nil {
		utils.LogError("记录审计日志失败", utils.ErrorField(err))
	}

	return nil
}

// ChangePassword 修改密码
func (s *UserServiceImpl) ChangePassword(ctx context.Context, userID uint, req *dto.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		appErr := common.NewAppErrorFromTypeWithCause("database", "USER_GET_FAILED", "获取用户失败", err)
		common.LogAppError(appErr, "user_change_password", utils.Uint("user_id", userID))
		return appErr
	}
	if user == nil {
		appErr := common.NewAppErrorFromType("business", "USER_NOT_FOUND", "用户不存在")
		common.LogAppError(appErr, "user_change_password", utils.Uint("user_id", userID))
		return appErr
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		appErr := common.NewAppErrorFromType("auth", "INVALID_OLD_PASSWORD", "旧密码错误")
		common.LogAppError(appErr, "user_change_password", utils.Uint("user_id", userID))
		return appErr
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		appErr := common.NewAppErrorFromTypeWithCause("system", "PASSWORD_HASH_FAILED", "密码加密失败", err)
		common.LogAppError(appErr, "user_change_password", utils.Uint("user_id", userID))
		return appErr
	}

	user.Password = string(hashedPassword)
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		appErr := common.NewAppErrorFromTypeWithCause("database", "PASSWORD_UPDATE_FAILED", "更新密码失败", err)
		common.LogAppError(appErr, "user_change_password", utils.Uint("user_id", userID))
		return appErr
	}

	// 记录密码修改成功的业务日志
	utils.Info("用户密码修改成功",
		utils.Uint("user_id", userID),
		utils.String("username", user.Username),
		utils.String("operation", "change_password"),
	)

	return nil
}

// GetUsers 获取用户列表
func (s *UserServiceImpl) GetUsers(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.UserResponse], error) {
	// 使用BaseService标准化分页参数
	normalizedReq := s.paginationHelper.NormalizePagination(req)
	
	// 检查缓存
	cacheKey := fmt.Sprintf("users:list:%d:%d", normalizedReq.Page, normalizedReq.PageSize)
	if cached, found := s.GetFromCache(ctx, cacheKey); found {
		if result, ok := cached.(*dto.PaginatedResponse[dto.UserResponse]); ok {
			return result, nil
		}
	}

	offset := normalizedReq.GetOffset()
	limit := normalizedReq.GetLimit()

	users, total, err := s.userRepo.ListUsers(ctx, offset, limit)
	if err != nil {
		appErr := common.NewAppErrorFromTypeWithCause("database", "USER_LIST_FAILED", "获取用户列表失败", err)
		common.LogAppError(appErr, "user_get_users", 
			utils.Int("page", normalizedReq.Page),
			utils.Int("page_size", normalizedReq.PageSize))
		return nil, appErr
	}

	// 使用转换器转换为响应格式
	userResponses := make([]dto.UserResponse, len(users))
	converter := &dto.UserConverter{}
	for i, user := range users {
		userResponses[i] = converter.ToDTO(*user)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	result := &dto.PaginatedResponse[dto.UserResponse]{
		Data:       userResponses,
		Total:      total,
		Page:       normalizedReq.Page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	// 设置缓存
	s.SetToCache(ctx, cacheKey, result)

	return result, nil
}

// SearchUsers 搜索用户
func (s *UserServiceImpl) SearchUsers(ctx context.Context, req *dto.SearchRequest) (*dto.PaginatedResponse[dto.UserResponse], error) {
	// 检查缓存
	cacheKey := fmt.Sprintf("users:search:%s:%d:%d", req.Keyword, req.Page, req.PageSize)
	if cached, found := s.GetFromCache(ctx, cacheKey); found {
		if result, ok := cached.(*dto.PaginatedResponse[dto.UserResponse]); ok {
			return result, nil
		}
	}

	offset := req.GetOffset()
	limit := req.GetLimit()

	users, total, err := s.userRepo.Search(ctx, req.Keyword, offset, limit)
	if err != nil {
		appErr := common.NewAppErrorFromTypeWithCause("database", "USER_SEARCH_FAILED", "搜索用户失败", err)
		common.LogAppError(appErr, "user_search_users", 
			utils.String("keyword", req.Keyword),
			utils.Int("page", req.Page),
			utils.Int("page_size", req.PageSize))
		return nil, appErr
	}

	// 使用转换器转换为响应格式
	userResponses := make([]dto.UserResponse, len(users))
	converter := &dto.UserConverter{}
	for i, user := range users {
		userResponses[i] = converter.ToDTO(*user)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	result := &dto.PaginatedResponse[dto.UserResponse]{
		Data:       userResponses,
		Total:      total,
		Page:       req.Page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	// 设置缓存
	s.SetToCache(ctx, cacheKey, result)

	return result, nil
}

// DeleteUser 删除用户
func (s *UserServiceImpl) DeleteUser(ctx context.Context, userID uint) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		appErr := common.NewAppErrorFromTypeWithCause("database", "USER_GET_FAILED", "获取用户失败", err)
		common.LogAppError(appErr, "user_delete_user", utils.Uint("user_id", userID))
		return appErr
	}
	if user == nil {
		appErr := common.NewAppErrorFromType("business", "USER_NOT_FOUND", "用户不存在")
		common.LogAppError(appErr, "user_delete_user", utils.Uint("user_id", userID))
		return appErr
	}

	err = s.userRepo.Delete(ctx, userID)
	if err != nil {
		appErr := common.NewAppErrorFromTypeWithCause("database", "USER_DELETE_FAILED", "删除用户失败", err)
		common.LogAppError(appErr, "user_delete_user", utils.Uint("user_id", userID))
		return appErr
	}

	// 清理相关缓存
	s.DeleteFromCache(ctx, fmt.Sprintf("user:%d", userID))
	s.DeleteFromCache(ctx, "users:list:*")
	s.DeleteFromCache(ctx, "users:search:*")

	// 记录用户删除成功的业务日志
	utils.Info("用户删除成功",
		utils.Uint("user_id", userID),
		utils.String("username", user.Username),
		utils.String("operation", "delete_user"),
	)

	// 记录审计日志
	if err := s.auditLogService.LogAction(ctx, userID, user.Username, "DELETE", "USER", fmt.Sprintf("%d", userID), fmt.Sprintf("用户删除: %s", user.Username), user, nil); err != nil {
		utils.LogError("记录审计日志失败", utils.ErrorField(err))
	}

	return nil
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
