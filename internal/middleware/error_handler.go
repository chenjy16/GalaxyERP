package middleware

import (
	"errors"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"github.com/galaxyerp/galaxyErp/internal/common"
	"github.com/galaxyerp/galaxyErp/internal/utils"
)

// ErrorHandler 错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			handleError(c, errors.New(err))
		} else if err, ok := recovered.(error); ok {
			handleError(c, err)
		} else {
			handleError(c, fmt.Errorf("未知错误: %v", recovered))
		}
	})
}

// ErrorHandlerMiddleware 错误处理中间件（用于处理业务逻辑错误）
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			handleError(c, err)
			return
		}
	}
}

// handleError 统一错误处理函数
func handleError(c *gin.Context, err error) {
	// 记录错误日志
	logError(c, err)

	// 创建响应助手
	resp := common.NewAPIResponseHelper(c)

	// 根据错误类型进行处理
	switch {
	case isAppError(err):
		handleAppError(resp, err.(*common.AppError))
	case isValidationError(err):
		handleValidationError(resp, err)
	case isGormError(err):
		handleGormError(resp, err)
	case isBindingError(err):
		handleBindingError(resp, err)
	default:
		handleGenericError(resp, err)
	}
}

// isAppError 检查是否为应用错误
func isAppError(err error) bool {
	return common.IsAppError(err)
}

// isValidationError 检查是否为验证错误
func isValidationError(err error) bool {
	var validationErrors validator.ValidationErrors
	return errors.As(err, &validationErrors)
}

// isGormError 检查是否为GORM错误
func isGormError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound) ||
		errors.Is(err, gorm.ErrInvalidTransaction) ||
		errors.Is(err, gorm.ErrNotImplemented) ||
		errors.Is(err, gorm.ErrMissingWhereClause) ||
		errors.Is(err, gorm.ErrUnsupportedRelation) ||
		errors.Is(err, gorm.ErrPrimaryKeyRequired) ||
		errors.Is(err, gorm.ErrModelValueRequired) ||
		errors.Is(err, gorm.ErrInvalidData) ||
		errors.Is(err, gorm.ErrUnsupportedDriver) ||
		errors.Is(err, gorm.ErrRegistered) ||
		errors.Is(err, gorm.ErrInvalidField) ||
		errors.Is(err, gorm.ErrEmptySlice) ||
		errors.Is(err, gorm.ErrDryRunModeUnsupported) ||
		errors.Is(err, gorm.ErrInvalidDB) ||
		errors.Is(err, gorm.ErrInvalidValue) ||
		errors.Is(err, gorm.ErrInvalidValueOfLength)
}

// isBindingError 检查是否为绑定错误
func isBindingError(err error) bool {
	return strings.Contains(err.Error(), "bind") ||
		strings.Contains(err.Error(), "unmarshal") ||
		strings.Contains(err.Error(), "invalid character")
}

// handleAppError 处理应用错误
func handleAppError(resp *common.APIResponseHelper, appErr *common.AppError) {
	details := make(map[string]interface{})
	if appErr.Details != "" {
		details["details"] = appErr.Details
	}
	if appErr.Cause != nil {
		details["cause"] = appErr.Cause.Error()
	}

	resp.Error(appErr.StatusCode, string(appErr.Code), appErr.Message, details)
}

// handleValidationError 处理验证错误
func handleValidationError(resp *common.APIResponseHelper, err error) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		details := make(map[string]interface{})
		fieldErrors := make(map[string]string)

		for _, fieldError := range validationErrors {
			fieldName := getFieldName(fieldError)
			fieldErrors[fieldName] = getValidationErrorMessage(fieldError)
		}

		details["field_errors"] = fieldErrors
		resp.ValidationError("输入数据验证失败", details)
	} else {
		resp.BadRequest("数据验证失败: " + err.Error())
	}
}

// handleGormError 处理GORM错误
func handleGormError(resp *common.APIResponseHelper, err error) {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		resp.NotFound("记录")
	case errors.Is(err, gorm.ErrInvalidTransaction):
		resp.InternalError("数据库事务错误")
	case errors.Is(err, gorm.ErrMissingWhereClause):
		resp.BadRequest("缺少查询条件")
	case errors.Is(err, gorm.ErrInvalidData):
		resp.BadRequest("无效的数据格式")
	default:
		resp.InternalError("数据库操作失败")
	}
}

// handleBindingError 处理绑定错误
func handleBindingError(resp *common.APIResponseHelper, err error) {
	details := map[string]interface{}{
		"error": err.Error(),
	}
	resp.BadRequest("请求数据格式错误", details)
}

// handleGenericError 处理通用错误
func handleGenericError(resp *common.APIResponseHelper, err error) {
	// 对于未知错误，不暴露具体错误信息给客户端
	resp.InternalError("服务器内部错误")
}

// getFieldName 获取字段名称
func getFieldName(fieldError validator.FieldError) string {
	// 可以根据需要进行字段名称转换
	fieldName := fieldError.Field()
	
	// 转换为下划线命名
	return toSnakeCase(fieldName)
}

// getValidationErrorMessage 获取验证错误消息
func getValidationErrorMessage(fieldError validator.FieldError) string {
	fieldName := fieldError.Field()
	tag := fieldError.Tag()
	param := fieldError.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("%s是必填字段", fieldName)
	case "email":
		return fmt.Sprintf("%s必须是有效的邮箱地址", fieldName)
	case "min":
		return fmt.Sprintf("%s最小值为%s", fieldName, param)
	case "max":
		return fmt.Sprintf("%s最大值为%s", fieldName, param)
	case "len":
		return fmt.Sprintf("%s长度必须为%s", fieldName, param)
	case "oneof":
		return fmt.Sprintf("%s必须是以下值之一: %s", fieldName, param)
	case "numeric":
		return fmt.Sprintf("%s必须是数字", fieldName)
	case "alpha":
		return fmt.Sprintf("%s只能包含字母", fieldName)
	case "alphanum":
		return fmt.Sprintf("%s只能包含字母和数字", fieldName)
	case "url":
		return fmt.Sprintf("%s必须是有效的URL", fieldName)
	case "uuid":
		return fmt.Sprintf("%s必须是有效的UUID", fieldName)
	case "datetime":
		return fmt.Sprintf("%s必须是有效的日期时间格式", fieldName)
	default:
		return fmt.Sprintf("%s验证失败", fieldName)
	}
}

// toSnakeCase 转换为下划线命名
func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}

// logError 记录错误日志
func logError(c *gin.Context, err error) {
	// 获取请求信息
	method := c.Request.Method
	path := c.Request.URL.Path
	clientIP := c.ClientIP()
	userAgent := c.Request.UserAgent()
	requestID := c.GetString("request_id")

	// 构建基础日志字段
	baseFields := []interface{}{
		utils.String("method", method),
		utils.String("path", path),
		utils.String("client_ip", clientIP),
		utils.String("user_agent", userAgent),
		utils.String("request_id", requestID),
		utils.String("error", err.Error()),
	}

	// 如果是应用错误，添加更多信息
	var fields []interface{}
	fields = append(fields, baseFields...)
	
	if appErr := common.GetAppError(err); appErr != nil {
		fields = append(fields,
			utils.String("error_code", string(appErr.Code)),
			utils.Int("status_code", appErr.StatusCode),
		)
		if appErr.Details != "" {
			fields = append(fields, utils.String("details", appErr.Details))
		}
		if appErr.Cause != nil {
			fields = append(fields, utils.String("cause", appErr.Cause.Error()))
		}
	}

	// 转换为zap.Field切片
	zapFields := make([]interface{}, len(fields))
	copy(zapFields, fields)

	// 如果是panic错误，记录堆栈信息
	if strings.Contains(err.Error(), "runtime error") {
		zapFields = append(zapFields, utils.String("stack", string(debug.Stack())))
		// 直接使用字符串格式记录日志
		utils.GlobalLogger.Error(fmt.Sprintf("Panic recovered - %s %s: %s", method, path, err.Error()))
	} else if common.IsAppError(err) {
		// 应用错误使用Warn级别
		utils.GlobalLogger.Warn(fmt.Sprintf("Application error - %s %s: %s", method, path, err.Error()))
	} else {
		// 其他错误使用Error级别
		utils.GlobalLogger.Error(fmt.Sprintf("Request error - %s %s: %s", method, path, err.Error()))
	}
}