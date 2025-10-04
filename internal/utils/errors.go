package utils

import (
	"fmt"
	"go.uber.org/zap"
	"runtime"
)

// ErrorType 错误类型枚举
type ErrorType string

const (
	ErrorTypeValidation ErrorType = "validation"
	ErrorTypeDatabase   ErrorType = "database"
	ErrorTypeAuth       ErrorType = "authentication"
	ErrorTypePermission ErrorType = "permission"
	ErrorTypeBusiness   ErrorType = "business"
	ErrorTypeSystem     ErrorType = "system"
	ErrorTypeExternal   ErrorType = "external"
)

// AppError 应用程序错误结构
type AppError struct {
	Type       ErrorType `json:"type"`
	Code       string    `json:"code"`
	Message    string    `json:"message"`
	Details    string    `json:"details,omitempty"`
	Cause      error     `json:"-"`
	StackTrace string    `json:"stack_trace,omitempty"`
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap 实现 errors.Unwrap 接口
func (e *AppError) Unwrap() error {
	return e.Cause
}

// NewAppError 创建新的应用程序错误
func NewAppError(errorType ErrorType, code, message string) *AppError {
	return &AppError{
		Type:       errorType,
		Code:       code,
		Message:    message,
		StackTrace: getStackTrace(),
	}
}

// NewAppErrorWithCause 创建带原因的应用程序错误
func NewAppErrorWithCause(errorType ErrorType, code, message string, cause error) *AppError {
	return &AppError{
		Type:       errorType,
		Code:       code,
		Message:    message,
		Cause:      cause,
		StackTrace: getStackTrace(),
	}
}

// NewAppErrorWithDetails 创建带详细信息的应用程序错误
func NewAppErrorWithDetails(errorType ErrorType, code, message, details string) *AppError {
	return &AppError{
		Type:       errorType,
		Code:       code,
		Message:    message,
		Details:    details,
		StackTrace: getStackTrace(),
	}
}

// LogAppError 记录应用程序错误日志
func LogAppError(err error, operation string, fields ...zap.Field) {
	if appErr, ok := err.(*AppError); ok {
		// 应用程序错误，使用结构化日志
		logFields := []zap.Field{
			zap.String("operation", operation),
			zap.String("error_type", string(appErr.Type)),
			zap.String("error_code", appErr.Code),
			zap.String("error_message", appErr.Message),
		}

		if appErr.Details != "" {
			logFields = append(logFields, zap.String("error_details", appErr.Details))
		}

		if appErr.Cause != nil {
			logFields = append(logFields, zap.String("cause", appErr.Cause.Error()))
		}

		if appErr.StackTrace != "" {
			logFields = append(logFields, zap.String("stack_trace", appErr.StackTrace))
		}

		// 添加额外字段
		logFields = append(logFields, fields...)

		// 根据错误类型选择日志级别
		switch appErr.Type {
		case ErrorTypeValidation, ErrorTypeBusiness:
			zap.L().Warn("业务错误", logFields...)
		case ErrorTypeAuth, ErrorTypePermission:
			zap.L().Warn("认证/权限错误", logFields...)
		case ErrorTypeDatabase, ErrorTypeSystem, ErrorTypeExternal:
			zap.L().Error("系统错误", logFields...)
		default:
			zap.L().Error("未知错误", logFields...)
		}
	} else {
		// 普通错误，使用简单日志
		logFields := []zap.Field{
			zap.String("operation", operation),
			zap.Error(err),
		}
		logFields = append(logFields, fields...)
		zap.L().Error("错误", logFields...)
	}
}

// LogAppErrorWithContext 记录带上下文的应用程序错误日志
func LogAppErrorWithContext(err error, operation string, userID uint, fields ...zap.Field) {
	contextFields := []zap.Field{
		zap.Uint("user_id", userID),
	}
	contextFields = append(contextFields, fields...)
	LogAppError(err, operation, contextFields...)
}

// getStackTrace 获取堆栈跟踪信息
func getStackTrace() string {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			return string(buf[:n])
		}
		buf = make([]byte, 2*len(buf))
	}
}

// 预定义的常见错误
var (
	ErrUserNotFound     = NewAppError(ErrorTypeBusiness, "USER_NOT_FOUND", "用户不存在")
	ErrInvalidPassword  = NewAppError(ErrorTypeAuth, "INVALID_PASSWORD", "密码错误")
	ErrUserExists       = NewAppError(ErrorTypeBusiness, "USER_EXISTS", "用户已存在")
	ErrInvalidToken     = NewAppError(ErrorTypeAuth, "INVALID_TOKEN", "无效的令牌")
	ErrPermissionDenied = NewAppError(ErrorTypePermission, "PERMISSION_DENIED", "权限不足")
	ErrInvalidInput     = NewAppError(ErrorTypeValidation, "INVALID_INPUT", "输入参数无效")
	ErrDatabaseError    = NewAppError(ErrorTypeDatabase, "DATABASE_ERROR", "数据库操作失败")
	ErrSystemError      = NewAppError(ErrorTypeSystem, "SYSTEM_ERROR", "系统内部错误")
)
