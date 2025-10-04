package middleware

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/galaxyerp/galaxyErp/internal/utils"
)

// ValidationConfig 验证配置
type ValidationConfig struct {
	// 是否启用详细错误信息
	EnableDetailedErrors bool
	// 是否启用字段翻译
	EnableFieldTranslation bool
	// 自定义错误消息
	CustomMessages map[string]string
}

// DefaultValidationConfig 默认验证配置
func DefaultValidationConfig() ValidationConfig {
	return ValidationConfig{
		EnableDetailedErrors:   true,
		EnableFieldTranslation: true,
		CustomMessages:         getDefaultErrorMessages(),
	}
}

// getDefaultErrorMessages 获取默认错误消息
func getDefaultErrorMessages() map[string]string {
	return map[string]string{
		"required":  "字段 {field} 是必填的",
		"email":     "字段 {field} 必须是有效的邮箱地址",
		"min":       "字段 {field} 的最小长度为 {param}",
		"max":       "字段 {field} 的最大长度为 {param}",
		"len":       "字段 {field} 的长度必须为 {param}",
		"numeric":   "字段 {field} 必须是数字",
		"alpha":     "字段 {field} 只能包含字母",
		"alphanum":  "字段 {field} 只能包含字母和数字",
		"url":       "字段 {field} 必须是有效的URL",
		"uuid":      "字段 {field} 必须是有效的UUID",
		"datetime":  "字段 {field} 必须是有效的日期时间格式",
		"phone":     "字段 {field} 必须是有效的手机号码",
		"idcard":    "字段 {field} 必须是有效的身份证号码",
		"positive":  "字段 {field} 必须是正数",
		"nonnegative": "字段 {field} 必须是非负数",
	}
}

// Validator 验证器实例
var Validator *validator.Validate

// InitValidator 初始化验证器
func InitValidator() {
	Validator = validator.New()
	
	// 注册自定义验证规则
	registerCustomValidations()
	
	// 注册字段名翻译
	registerFieldTranslations()
}

// registerCustomValidations 注册自定义验证规则
func registerCustomValidations() {
	// 手机号验证
	Validator.RegisterValidation("phone", validatePhone)
	
	// 身份证号验证
	Validator.RegisterValidation("idcard", validateIDCard)
	
	// 正数验证
	Validator.RegisterValidation("positive", validatePositive)
	
	// 非负数验证
	Validator.RegisterValidation("nonnegative", validateNonNegative)
	
	// 日期时间验证
	Validator.RegisterValidation("datetime", validateDateTime)
	
	// 密码强度验证
	Validator.RegisterValidation("password", validatePassword)
	
	// 用户名验证
	Validator.RegisterValidation("username", validateUsername)
}

// registerFieldTranslations 注册字段名翻译
func registerFieldTranslations() {
	// 使用结构体标签获取字段名
	Validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		
		// 如果有中文标签，使用中文标签
		if label := fld.Tag.Get("label"); label != "" {
			return label
		}
		
		return name
	})
}

// ValidationMiddleware 验证中间件
func ValidationMiddleware(config ...ValidationConfig) gin.HandlerFunc {
	cfg := DefaultValidationConfig()
	if len(config) > 0 {
		cfg = config[0]
	}
	
	return func(c *gin.Context) {
		// 设置验证配置到上下文
		c.Set("validation_config", cfg)
		c.Next()
	}
}

// ValidateRequest 验证请求数据
func ValidateRequest(c *gin.Context, req interface{}) error {
	// 绑定请求数据
	if err := c.ShouldBindJSON(req); err != nil {
		logBindingError(c, err)
		return err
	}
	
	// 验证数据
	if err := Validator.Struct(req); err != nil {
		logValidationError(c, err)
		return err
	}
	
	return nil
}

// ValidateQuery 验证查询参数
func ValidateQuery(c *gin.Context, req interface{}) error {
	// 绑定查询参数
	if err := c.ShouldBindQuery(req); err != nil {
		logBindingError(c, err)
		return err
	}
	
	// 验证数据
	if err := Validator.Struct(req); err != nil {
		logValidationError(c, err)
		return err
	}
	
	return nil
}

// ValidateURI 验证URI参数
func ValidateURI(c *gin.Context, req interface{}) error {
	// 绑定URI参数
	if err := c.ShouldBindUri(req); err != nil {
		logBindingError(c, err)
		return err
	}
	
	// 验证数据
	if err := Validator.Struct(req); err != nil {
		logValidationError(c, err)
		return err
	}
	
	return nil
}

// logBindingError 记录绑定错误
func logBindingError(c *gin.Context, err error) {
	requestID := getRequestID(c)
	
	utils.GlobalLogger.Warn("Request binding failed",
		utils.String("request_id", requestID),
		utils.String("error", err.Error()),
		utils.String("method", c.Request.Method),
		utils.String("path", c.Request.URL.Path),
	)
}

// logValidationError 记录验证错误
func logValidationError(c *gin.Context, err error) {
	requestID := getRequestID(c)
	config := getValidationConfig(c)
	
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		utils.GlobalLogger.Warn("Unknown validation error",
			utils.String("request_id", requestID),
			utils.String("error", err.Error()),
		)
		return
	}
	
	// 构建详细错误信息
	errorDetails := buildValidationErrorDetails(validationErrors, config)
	
	utils.GlobalLogger.Warn("Request validation failed",
		utils.String("request_id", requestID),
		utils.String("method", c.Request.Method),
		utils.String("path", c.Request.URL.Path),
		utils.Any("validation_errors", errorDetails),
	)
}

// buildValidationErrorDetails 构建验证错误详情
func buildValidationErrorDetails(validationErrors validator.ValidationErrors, config ValidationConfig) []map[string]interface{} {
	var errors []map[string]interface{}
	
	for _, err := range validationErrors {
		errorDetail := map[string]interface{}{
			"field": err.Field(),
			"tag":   err.Tag(),
			"value": err.Value(),
		}
		
		if config.EnableDetailedErrors {
			errorDetail["message"] = getErrorMessage(err, config.CustomMessages)
			errorDetail["param"] = err.Param()
		}
		
		errors = append(errors, errorDetail)
	}
	
	return errors
}

// getErrorMessage 获取错误消息
func getErrorMessage(err validator.FieldError, customMessages map[string]string) string {
	tag := err.Tag()
	field := err.Field()
	param := err.Param()
	
	// 查找自定义消息
	if message, exists := customMessages[tag]; exists {
		message = strings.ReplaceAll(message, "{field}", field)
		message = strings.ReplaceAll(message, "{param}", param)
		return message
	}
	
	// 默认消息
	switch tag {
	case "required":
		return fmt.Sprintf("字段 %s 是必填的", field)
	case "email":
		return fmt.Sprintf("字段 %s 必须是有效的邮箱地址", field)
	case "min":
		return fmt.Sprintf("字段 %s 的最小长度为 %s", field, param)
	case "max":
		return fmt.Sprintf("字段 %s 的最大长度为 %s", field, param)
	case "len":
		return fmt.Sprintf("字段 %s 的长度必须为 %s", field, param)
	default:
		return fmt.Sprintf("字段 %s 验证失败: %s", field, tag)
	}
}

// getRequestID 获取请求ID
func getRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return "unknown"
}

// getValidationConfig 获取验证配置
func getValidationConfig(c *gin.Context) ValidationConfig {
	if config, exists := c.Get("validation_config"); exists {
		if cfg, ok := config.(ValidationConfig); ok {
			return cfg
		}
	}
	return DefaultValidationConfig()
}

// 自定义验证函数

// validatePhone 验证手机号
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if phone == "" {
		return true // 空值由required标签处理
	}
	
	// 中国手机号正则表达式
	phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return phoneRegex.MatchString(phone)
}

// validateIDCard 验证身份证号
func validateIDCard(fl validator.FieldLevel) bool {
	idCard := fl.Field().String()
	if idCard == "" {
		return true // 空值由required标签处理
	}
	
	// 18位身份证号正则表达式
	idCardRegex := regexp.MustCompile(`^[1-9]\d{5}(18|19|20)\d{2}((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$`)
	return idCardRegex.MatchString(idCard)
}

// validatePositive 验证正数
func validatePositive(fl validator.FieldLevel) bool {
	switch fl.Field().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fl.Field().Int() > 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fl.Field().Uint() > 0
	case reflect.Float32, reflect.Float64:
		return fl.Field().Float() > 0
	case reflect.String:
		val, err := strconv.ParseFloat(fl.Field().String(), 64)
		return err == nil && val > 0
	}
	return false
}

// validateNonNegative 验证非负数
func validateNonNegative(fl validator.FieldLevel) bool {
	switch fl.Field().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fl.Field().Int() >= 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true // uint类型天然非负
	case reflect.Float32, reflect.Float64:
		return fl.Field().Float() >= 0
	case reflect.String:
		val, err := strconv.ParseFloat(fl.Field().String(), 64)
		return err == nil && val >= 0
	}
	return false
}

// validateDateTime 验证日期时间格式
func validateDateTime(fl validator.FieldLevel) bool {
	dateTime := fl.Field().String()
	if dateTime == "" {
		return true // 空值由required标签处理
	}
	
	// 支持多种日期时间格式
	formats := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02",
		"15:04:05",
	}
	
	for _, format := range formats {
		if _, err := time.Parse(format, dateTime); err == nil {
			return true
		}
	}
	
	return false
}

// validatePassword 验证密码强度
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if password == "" {
		return true // 空值由required标签处理
	}
	
	// 密码长度至少8位
	if len(password) < 8 {
		return false
	}
	
	// 至少包含一个数字、一个小写字母、一个大写字母
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	
	return hasNumber && hasLower && hasUpper
}

// validateUsername 验证用户名
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	if username == "" {
		return true // 空值由required标签处理
	}
	
	// 用户名只能包含字母、数字、下划线，长度3-20位
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	return usernameRegex.MatchString(username)
}