package common

import (
	"fmt"
	"net/http"
	"runtime"
	"go.uber.org/zap"
	"github.com/galaxyerp/galaxyErp/internal/utils"
)

// ErrorCode 错误码类型
type ErrorCode string

// 业务错误码枚举
const (
	// 通用错误码
	ErrCodeSuccess           ErrorCode = "SUCCESS"
	ErrCodeInternalError     ErrorCode = "INTERNAL_ERROR"
	ErrCodeInvalidRequest    ErrorCode = "INVALID_REQUEST"
	ErrCodeUnauthorized      ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden         ErrorCode = "FORBIDDEN"
	ErrCodeNotFound          ErrorCode = "NOT_FOUND"
	ErrCodeConflict          ErrorCode = "CONFLICT"
	ErrCodeValidationFailed  ErrorCode = "VALIDATION_FAILED"
	ErrCodeDatabaseError     ErrorCode = "DATABASE_ERROR"

	// 用户相关错误码
	ErrCodeUserNotFound      ErrorCode = "USER_NOT_FOUND"
	ErrCodeUserExists        ErrorCode = "USER_EXISTS"
	ErrCodeInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
	ErrCodePasswordTooWeak   ErrorCode = "PASSWORD_TOO_WEAK"
	ErrCodeTokenExpired      ErrorCode = "TOKEN_EXPIRED"
	ErrCodeTokenInvalid      ErrorCode = "TOKEN_INVALID"

	// 销售相关错误码
	ErrCodeSalesOrderNotFound     ErrorCode = "SALES_ORDER_NOT_FOUND"
	ErrCodeSalesOrderCancelled    ErrorCode = "SALES_ORDER_CANCELLED"
	ErrCodeSalesOrderCompleted    ErrorCode = "SALES_ORDER_COMPLETED"
	ErrCodeInsufficientInventory  ErrorCode = "INSUFFICIENT_INVENTORY"
	ErrCodeInvalidOrderStatus     ErrorCode = "INVALID_ORDER_STATUS"

	// 采购相关错误码
	ErrCodePurchaseOrderNotFound  ErrorCode = "PURCHASE_ORDER_NOT_FOUND"
	ErrCodeSupplierNotFound       ErrorCode = "SUPPLIER_NOT_FOUND"
	ErrCodeInvalidPurchaseStatus  ErrorCode = "INVALID_PURCHASE_STATUS"

	// 库存相关错误码
	ErrCodeInventoryNotFound      ErrorCode = "INVENTORY_NOT_FOUND"
	ErrCodeInventoryInsufficient  ErrorCode = "INVENTORY_INSUFFICIENT"
	ErrCodeWarehouseNotFound      ErrorCode = "WAREHOUSE_NOT_FOUND"

	// 财务相关错误码
	ErrCodeAccountNotFound        ErrorCode = "ACCOUNT_NOT_FOUND"
	ErrCodeInvoiceNotFound        ErrorCode = "INVOICE_NOT_FOUND"
	ErrCodePaymentFailed          ErrorCode = "PAYMENT_FAILED"
	ErrCodeInvalidAmount          ErrorCode = "INVALID_AMOUNT"

	// 生产相关错误码
	ErrCodeProductionOrderNotFound ErrorCode = "PRODUCTION_ORDER_NOT_FOUND"
	ErrCodeBOMNotFound            ErrorCode = "BOM_NOT_FOUND"
	ErrCodeWorkCenterNotFound     ErrorCode = "WORK_CENTER_NOT_FOUND"

	// 项目相关错误码
	ErrCodeProjectNotFound        ErrorCode = "PROJECT_NOT_FOUND"
	ErrCodeTaskNotFound           ErrorCode = "TASK_NOT_FOUND"
	ErrCodeInvalidProjectStatus   ErrorCode = "INVALID_PROJECT_STATUS"

	// 人力资源相关错误码
	ErrCodeEmployeeNotFound       ErrorCode = "EMPLOYEE_NOT_FOUND"
	ErrCodeDepartmentNotFound     ErrorCode = "DEPARTMENT_NOT_FOUND"
	ErrCodePositionNotFound       ErrorCode = "POSITION_NOT_FOUND"
)

// AppError 应用错误结构
type AppError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	Details    string    `json:"details,omitempty"`
	StatusCode int       `json:"status_code"`
	Cause      error     `json:"-"`
	StackTrace string    `json:"stack_trace,omitempty"`
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap 支持错误链
func (e *AppError) Unwrap() error {
	return e.Cause
}

// NewAppError 创建应用错误
func NewAppError(code ErrorCode, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		StackTrace: getStackTrace(),
	}
}

// NewAppErrorWithDetails 创建带详情的应用错误
func NewAppErrorWithDetails(code ErrorCode, message, details string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Details:    details,
		StatusCode: statusCode,
		StackTrace: getStackTrace(),
	}
}

// NewAppErrorWithCause 创建带原因的应用错误
func NewAppErrorWithCause(code ErrorCode, message string, statusCode int, cause error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Cause:      cause,
		StackTrace: getStackTrace(),
	}
}

// WrapError 包装错误
func WrapError(code ErrorCode, message string, statusCode int, cause error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Cause:      cause,
		StackTrace: getStackTrace(),
	}
}

// LogAppError 记录应用错误日志
func LogAppError(err error, operation string, fields ...zap.Field) {
	if appErr, ok := err.(*AppError); ok {
		logFields := []zap.Field{
			utils.String("operation", operation),
			utils.String("error_code", string(appErr.Code)),
			utils.String("error_message", appErr.Message),
		}
		
		if appErr.Details != "" {
			logFields = append(logFields, utils.String("error_details", appErr.Details))
		}
		
		if appErr.Cause != nil {
			logFields = append(logFields, utils.String("cause", appErr.Cause.Error()))
		}
		
		if appErr.StackTrace != "" {
			logFields = append(logFields, utils.String("stack_trace", appErr.StackTrace))
		}
		
		// 添加额外的字段
		logFields = append(logFields, fields...)
		
		utils.LogError("Application error occurred", logFields...)
	} else {
		utils.LogError("Unknown error occurred", 
			utils.String("operation", operation),
			utils.String("error", err.Error()),
		)
	}
}

// LogAppErrorWithContext 记录带上下文的应用错误日志
func LogAppErrorWithContext(err error, operation string, userID uint, fields ...zap.Field) {
	contextFields := []zap.Field{utils.Uint("user_id", userID)}
	contextFields = append(contextFields, fields...)
	LogAppError(err, operation, contextFields...)
}

// getStackTrace 获取堆栈跟踪
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

// 预定义的错误创建函数

// NewInternalError 创建内部错误
func NewInternalError(message string) *AppError {
	return NewAppError(ErrCodeInternalError, message, http.StatusInternalServerError)
}

// NewValidationError 创建验证错误
func NewValidationError(message string) *AppError {
	return NewAppError(ErrCodeValidationFailed, message, http.StatusBadRequest)
}

// NewNotFoundError 创建未找到错误
func NewNotFoundError(resource string) *AppError {
	return NewAppError(ErrCodeNotFound, fmt.Sprintf("%s not found", resource), http.StatusNotFound)
}

// NewUnauthorizedError 创建未授权错误
func NewUnauthorizedError(message string) *AppError {
	return NewAppError(ErrCodeUnauthorized, message, http.StatusUnauthorized)
}

// NewForbiddenError 创建禁止访问错误
func NewForbiddenError(message string) *AppError {
	return NewAppError(ErrCodeForbidden, message, http.StatusForbidden)
}

// NewConflictError 创建冲突错误
func NewConflictError(message string) *AppError {
	return NewAppError(ErrCodeConflict, message, http.StatusConflict)
}

// NewDatabaseError 创建数据库错误
func NewDatabaseError(message string, cause error) *AppError {
	return NewAppErrorWithCause(ErrCodeDatabaseError, message, http.StatusInternalServerError, cause)
}

// 业务特定错误创建函数

// NewUserNotFoundError 创建用户未找到错误
func NewUserNotFoundError() *AppError {
	return NewAppError(ErrCodeUserNotFound, "User not found", http.StatusNotFound)
}

// NewUserExistsError 创建用户已存在错误
func NewUserExistsError(email string) *AppError {
	return NewAppErrorWithDetails(ErrCodeUserExists, "User already exists", fmt.Sprintf("Email: %s", email), http.StatusConflict)
}

// NewInvalidCredentialsError 创建无效凭据错误
func NewInvalidCredentialsError() *AppError {
	return NewAppError(ErrCodeInvalidCredentials, "Invalid credentials", http.StatusUnauthorized)
}

// NewTokenExpiredError 创建令牌过期错误
func NewTokenExpiredError() *AppError {
	return NewAppError(ErrCodeTokenExpired, "Token expired", http.StatusUnauthorized)
}

// NewInsufficientInventoryError 创建库存不足错误
func NewInsufficientInventoryError(productName string, available, required int) *AppError {
	details := fmt.Sprintf("Product: %s, Available: %d, Required: %d", productName, available, required)
	return NewAppErrorWithDetails(ErrCodeInsufficientInventory, "Insufficient inventory", details, http.StatusBadRequest)
}

// IsAppError 检查是否为应用错误
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetAppError 获取应用错误
func GetAppError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return nil
}

// 错误码映射函数，用于从旧的错误类型转换
func MapErrorTypeToCode(errorType string, code string) ErrorCode {
	switch code {
	case "USER_NOT_FOUND":
		return ErrCodeUserNotFound
	case "INVALID_PASSWORD":
		return ErrCodeInvalidCredentials
	case "USER_EXISTS", "EMAIL_EXISTS", "USERNAME_EXISTS":
		return ErrCodeUserExists
	case "INVALID_TOKEN":
		return ErrCodeTokenInvalid
	case "PERMISSION_DENIED":
		return ErrCodeForbidden
	case "INVALID_INPUT":
		return ErrCodeValidationFailed
	case "DATABASE_ERROR", "USER_CHECK_FAILED", "USERNAME_CHECK_FAILED", "PASSWORD_HASH_FAILED":
		return ErrCodeDatabaseError
	case "SYSTEM_ERROR":
		return ErrCodeInternalError
	case "SALES_ORDER_NOT_FOUND":
		return ErrCodeSalesOrderNotFound
	case "CUSTOMER_NOT_FOUND":
		return ErrCodeNotFound
	case "AUDIT_LOG_NOT_FOUND":
		return ErrCodeNotFound
	case "AUDIT_LOG_CREATE_FAILED", "AUDIT_LOG_GET_FAILED", "AUDIT_LOG_CLEANUP_FAILED":
		return ErrCodeDatabaseError
	default:
		// 根据错误类型进行映射
		switch errorType {
		case "validation":
			return ErrCodeValidationFailed
		case "database":
			return ErrCodeDatabaseError
		case "authentication":
			return ErrCodeUnauthorized
		case "permission":
			return ErrCodeForbidden
		case "business":
			return ErrCodeInvalidRequest
		case "system":
			return ErrCodeInternalError
		case "external":
			return ErrCodeInternalError
		default:
			return ErrCodeInternalError
		}
	}
}

// 便捷的错误创建函数，兼容旧的utils.ErrorType
func NewAppErrorFromType(errorType string, code, message string) *AppError {
	errorCode := MapErrorTypeToCode(errorType, code)
	statusCode := ErrorCodeToHTTPStatus(errorCode)
	return NewAppError(errorCode, message, statusCode)
}

func NewAppErrorFromTypeWithDetails(errorType string, code, message, details string) *AppError {
	errorCode := MapErrorTypeToCode(errorType, code)
	statusCode := ErrorCodeToHTTPStatus(errorCode)
	return NewAppErrorWithDetails(errorCode, message, details, statusCode)
}

func NewAppErrorFromTypeWithCause(errorType string, code, message string, cause error) *AppError {
	errorCode := MapErrorTypeToCode(errorType, code)
	statusCode := ErrorCodeToHTTPStatus(errorCode)
	return NewAppErrorWithCause(errorCode, message, statusCode, cause)
}

// ErrorCodeToHTTPStatus 错误码到HTTP状态码的映射
func ErrorCodeToHTTPStatus(code ErrorCode) int {
	switch code {
	case ErrCodeSuccess:
		return http.StatusOK
	case ErrCodeInvalidRequest, ErrCodeValidationFailed, ErrCodeInvalidAmount:
		return http.StatusBadRequest
	case ErrCodeUnauthorized, ErrCodeInvalidCredentials, ErrCodeTokenExpired, ErrCodeTokenInvalid:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeNotFound, ErrCodeUserNotFound, ErrCodeSalesOrderNotFound, ErrCodePurchaseOrderNotFound,
		 ErrCodeInventoryNotFound, ErrCodeAccountNotFound, ErrCodeInvoiceNotFound, ErrCodeProductionOrderNotFound,
		 ErrCodeBOMNotFound, ErrCodeWorkCenterNotFound, ErrCodeProjectNotFound, ErrCodeTaskNotFound,
		 ErrCodeEmployeeNotFound, ErrCodeDepartmentNotFound, ErrCodePositionNotFound, ErrCodeSupplierNotFound,
		 ErrCodeWarehouseNotFound:
		return http.StatusNotFound
	case ErrCodeConflict, ErrCodeUserExists:
		return http.StatusConflict
	case ErrCodeInternalError, ErrCodeDatabaseError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}