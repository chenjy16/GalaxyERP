package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	// Register custom tag name function to use json tag as field name
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Register custom validators
	registerCustomValidators()
}

// registerCustomValidators registers custom validators
func registerCustomValidators() {
	// Chinese mobile number validation
	validate.RegisterValidation("chinese_mobile", validateChineseMobile)

	// Chinese ID card validation
	validate.RegisterValidation("chinese_id_card", validateChineseIDCard)

	// Bank card number validation
	validate.RegisterValidation("bank_card", validateBankCard)

	// Account code validation
	validate.RegisterValidation("account_code", validateAccountCode)

	// Currency amount validation (up to 2 decimal places)
	validate.RegisterValidation("currency", validateCurrency)

	// Date format validation
	validate.RegisterValidation("date_format", validateDateFormat)

	// Strong password validation
	validate.RegisterValidation("strong_password", validateStrongPassword)
}

// validateChineseMobile validates Chinese mobile number
func validateChineseMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, mobile)
	return matched
}

// validateChineseIDCard validates Chinese ID card number
func validateChineseIDCard(fl validator.FieldLevel) bool {
	idCard := fl.Field().String()
	// 18-digit ID card number validation
	matched, _ := regexp.MatchString(`^[1-9]\d{5}(18|19|20)\d{2}((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$`, idCard)
	return matched
}

// validateBankCard validates bank card number
func validateBankCard(fl validator.FieldLevel) bool {
	cardNumber := fl.Field().String()
	// Bank card numbers are typically 16-19 digits
	matched, _ := regexp.MatchString(`^\d{16,19}$`, cardNumber)
	return matched
}

// validateAccountCode validates account code
func validateAccountCode(fl validator.FieldLevel) bool {
	code := fl.Field().String()
	// Account code format: starts with digit, may contain dot separator, e.g. 1001.01
	matched, _ := regexp.MatchString(`^[1-9]\d{3}(\.\d{2})?$`, code)
	return matched
}

// validateCurrency validates currency amount (up to 2 decimal places)
func validateCurrency(fl validator.FieldLevel) bool {
	field := fl.Field()

	// Handle different types of numeric fields
	switch field.Kind() {
	case reflect.Float32, reflect.Float64:
		value := field.Float()
		// Check if negative
		if value < 0 {
			return false
		}
		// Check if decimal places exceed 2
		valueStr := fmt.Sprintf("%.10f", value)
		parts := strings.Split(valueStr, ".")
		if len(parts) > 1 {
			// Remove trailing zeros
			decimal := strings.TrimRight(parts[1], "0")
			if len(decimal) > 2 {
				return false
			}
		}
		return true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value := field.Int()
		return value >= 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true // uint types are inherently non-negative
	case reflect.String:
		value := field.String()
		matched, _ := regexp.MatchString(`^\d+(\.\d{1,2})?$`, value)
		return matched
	default:
		return false
	}
}

// validateDateFormat validates date format (YYYY-MM-DD)
func validateDateFormat(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

// validateStrongPassword validates strong password (at least 8 characters with uppercase, lowercase, numbers and special characters)
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}

	hasUpper, _ := regexp.MatchString(`[A-Z]`, password)
	hasLower, _ := regexp.MatchString(`[a-z]`, password)
	hasDigit, _ := regexp.MatchString(`\d`, password)
	hasSpecial, _ := regexp.MatchString(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`, password)

	return hasUpper && hasLower && hasDigit && hasSpecial
}

// ValidateStruct validates a struct and returns formatted error messages
func ValidateStruct(s interface{}) map[string]string {
	errors := make(map[string]string)

	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := err.Field()
			errors[fieldName] = getErrorMessage(err)
		}
	}

	return errors
}

// BindAndValidate binds JSON request and validates it
func BindAndValidate(c *gin.Context, obj interface{}) bool {
	// Use standard JSON decoding, completely skip Gin's validator
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		BadRequestResponse(c, "Failed to read request body: "+err.Error())
		return false
	}

	if err := json.Unmarshal(body, obj); err != nil {
		BadRequestResponse(c, "Invalid JSON format: "+err.Error())
		return false
	}

	// Use our custom validator for validation
	if errors := ValidateStruct(obj); len(errors) > 0 {
		ValidationErrorResponse(c, errors)
		return false
	}

	return true
}

// BindAndValidateWithCustomError binds JSON request and validates with custom error response
func BindAndValidateWithCustomError(ctx *gin.Context, req interface{}, customMessage string) bool {
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": customMessage,
			"errors":  err.Error(),
		})
		return false
	}

	if errors := ValidateStruct(req); len(errors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": customMessage,
			"errors":  errors,
		})
		return false
	}

	return true
}

// ValidateStructWithGroups validates a struct with validation groups
func ValidateStructWithGroups(s interface{}, groups ...string) map[string]string {
	errors := make(map[string]string)

	for _, group := range groups {
		if err := validate.StructPartial(s, group); err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				fieldName := err.Field()
				if _, exists := errors[fieldName]; !exists {
					errors[fieldName] = getErrorMessage(err)
				}
			}
		}
	}

	return errors
}

// ValidateField validates a single field
func ValidateField(field interface{}, tag string) error {
	return validate.Var(field, tag)
}

// ValidateFields validates multiple fields with their respective tags
func ValidateFields(fields map[string]interface{}, tags map[string]string) map[string]string {
	errors := make(map[string]string)

	for fieldName, fieldValue := range fields {
		if tag, exists := tags[fieldName]; exists {
			if err := validate.Var(fieldValue, tag); err != nil {
				if validationErrors, ok := err.(validator.ValidationErrors); ok {
					for _, validationError := range validationErrors {
						errors[fieldName] = getErrorMessage(validationError)
						break // 只取第一个错误
					}
				} else {
					errors[fieldName] = err.Error()
				}
			}
		}
	}

	return errors
}

// ValidateConditional validates a struct only if condition is met
func ValidateConditional(s interface{}, condition bool) map[string]string {
	if !condition {
		return make(map[string]string)
	}
	return ValidateStruct(s)
}

// getErrorMessage returns a user-friendly error message for validation errors
func getErrorMessage(err validator.FieldError) string {
	field := err.Field()
	tag := err.Tag()
	param := err.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, param)
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", field, param)
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", field, param)
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, param)
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, param)
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, param)
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, param)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, param)
	case "unique":
		return fmt.Sprintf("%s already exists", field)
	case "alphanum":
		return fmt.Sprintf("%s must contain only alphanumeric characters", field)
	case "alpha":
		return fmt.Sprintf("%s must contain only alphabetic characters", field)
	case "numeric":
		return fmt.Sprintf("%s must be numeric", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "uri":
		return fmt.Sprintf("%s must be a valid URI", field)
	case "datetime":
		return fmt.Sprintf("%s must be a valid datetime", field)
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	// Custom validator error messages
	case "chinese_mobile":
		return fmt.Sprintf("%s must be a valid Chinese mobile number", field)
	case "chinese_id_card":
		return fmt.Sprintf("%s must be a valid Chinese ID card number", field)
	case "bank_card":
		return fmt.Sprintf("%s must be a valid bank card number", field)
	case "account_code":
		return fmt.Sprintf("%s must be a valid account code", field)
	case "currency":
		return fmt.Sprintf("%s must be a valid currency amount (up to 2 decimal places)", field)
	case "date_format":
		return fmt.Sprintf("%s must be a valid date format (YYYY-MM-DD)", field)
	case "strong_password":
		return fmt.Sprintf("%s must be a strong password (at least 8 characters with uppercase, lowercase, numbers and special characters)", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// ValidateID validates that an ID parameter is valid
func ValidateID(c *gin.Context, paramName string) (uint, bool) {
	idStr := c.Param(paramName)
	if idStr == "" {
		BadRequestResponse(c, fmt.Sprintf("%s parameter is required", paramName))
		return 0, false
	}

	var id uint
	if err := c.ShouldBindUri(&struct {
		ID uint `uri:"id" binding:"required,min=1"`
	}{ID: 0}); err != nil {
		BadRequestResponse(c, fmt.Sprintf("Invalid %s parameter", paramName))
		return 0, false
	}

	// Parse the ID manually since ShouldBindUri doesn't work as expected
	fmt.Sscanf(idStr, "%d", &id)
	if id == 0 {
		BadRequestResponse(c, fmt.Sprintf("Invalid %s parameter", paramName))
		return 0, false
	}

	return id, true
}

// ValidatePagination validates pagination parameters
func ValidatePagination(c *gin.Context) (int, int, bool) {
	page := 1
	pageSize := 10

	if pageStr := c.Query("page"); pageStr != "" {
		fmt.Sscanf(pageStr, "%d", &page)
		if page < 1 {
			page = 1
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		fmt.Sscanf(pageSizeStr, "%d", &pageSize)
		if pageSize < 1 {
			pageSize = 10
		}
		if pageSize > 100 {
			pageSize = 100
		}
	}

	return page, pageSize, true
}
