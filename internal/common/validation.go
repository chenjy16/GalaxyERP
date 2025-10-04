package common

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ValidationRule 验证规则接口
type ValidationRule interface {
	Validate(ctx context.Context, value interface{}) *AppError
	GetMessage() string
}

// ValidationRuleFunc 验证规则函数类型
type ValidationRuleFunc func(ctx context.Context, value interface{}) *AppError

// Validate 实现 ValidationRule 接口
func (f ValidationRuleFunc) Validate(ctx context.Context, value interface{}) *AppError {
	return f(ctx, value)
}

// GetMessage 获取错误消息
func (f ValidationRuleFunc) GetMessage() string {
	return "Validation failed"
}

// ConditionalRule 条件验证规则
type ConditionalRule struct {
	Condition func(ctx context.Context, value interface{}) bool
	Rule      ValidationRule
	Message   string
}

// Validate 验证
func (r *ConditionalRule) Validate(ctx context.Context, value interface{}) *AppError {
	if r.Condition(ctx, value) {
		return r.Rule.Validate(ctx, value)
	}
	return nil
}

// GetMessage 获取错误消息
func (r *ConditionalRule) GetMessage() string {
	if r.Message != "" {
		return r.Message
	}
	return r.Rule.GetMessage()
}

// RuleComposer 验证规则组合器
type RuleComposer struct {
	rules []ValidationRule
}

// NewRuleComposer 创建验证规则组合器
func NewRuleComposer() *RuleComposer {
	return &RuleComposer{
		rules: make([]ValidationRule, 0),
	}
}

// AddRule 添加验证规则
func (c *RuleComposer) AddRule(rule ValidationRule) *RuleComposer {
	c.rules = append(c.rules, rule)
	return c
}

// AddConditionalRule 添加条件验证规则
func (c *RuleComposer) AddConditionalRule(condition func(ctx context.Context, value interface{}) bool, rule ValidationRule, message string) *RuleComposer {
	c.rules = append(c.rules, &ConditionalRule{
		Condition: condition,
		Rule:      rule,
		Message:   message,
	})
	return c
}

// Validate 执行验证
func (c *RuleComposer) Validate(ctx context.Context, value interface{}) *AppError {
	for _, rule := range c.rules {
		if err := rule.Validate(ctx, value); err != nil {
			return err
		}
	}
	return nil
}

// 基础验证规则

// RequiredRule 必填验证规则
type RequiredRule struct {
	Message string
}

// Validate 验证
func (r *RequiredRule) Validate(ctx context.Context, value interface{}) *AppError {
	if value == nil {
		return NewValidationError(r.GetMessage())
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		if strings.TrimSpace(v.String()) == "" {
			return NewValidationError(r.GetMessage())
		}
	case reflect.Slice, reflect.Array, reflect.Map:
		if v.Len() == 0 {
			return NewValidationError(r.GetMessage())
		}
	case reflect.Ptr:
		if v.IsNil() {
			return NewValidationError(r.GetMessage())
		}
	}
	return nil
}

// GetMessage 获取错误消息
func (r *RequiredRule) GetMessage() string {
	if r.Message != "" {
		return r.Message
	}
	return "Field is required"
}

// MinLengthRule 最小长度验证规则
type MinLengthRule struct {
	MinLength int
	Message   string
}

// Validate 验证
func (r *MinLengthRule) Validate(ctx context.Context, value interface{}) *AppError {
	if value == nil {
		return nil
	}

	str, ok := value.(string)
	if !ok {
		return NewValidationError("Value must be a string")
	}

	if len(str) < r.MinLength {
		return NewValidationError(r.GetMessage())
	}
	return nil
}

// GetMessage 获取错误消息
func (r *MinLengthRule) GetMessage() string {
	if r.Message != "" {
		return r.Message
	}
	return fmt.Sprintf("Minimum length is %d", r.MinLength)
}

// MaxLengthRule 最大长度验证规则
type MaxLengthRule struct {
	MaxLength int
	Message   string
}

// Validate 验证
func (r *MaxLengthRule) Validate(ctx context.Context, value interface{}) *AppError {
	if value == nil {
		return nil
	}

	str, ok := value.(string)
	if !ok {
		return NewValidationError("Value must be a string")
	}

	if len(str) > r.MaxLength {
		return NewValidationError(r.GetMessage())
	}
	return nil
}

// GetMessage 获取错误消息
func (r *MaxLengthRule) GetMessage() string {
	if r.Message != "" {
		return r.Message
	}
	return fmt.Sprintf("Maximum length is %d", r.MaxLength)
}

// EmailRule 邮箱验证规则
type EmailRule struct {
	Message string
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Validate 验证
func (r *EmailRule) Validate(ctx context.Context, value interface{}) *AppError {
	if value == nil {
		return nil
	}

	str, ok := value.(string)
	if !ok {
		return NewValidationError("Value must be a string")
	}

	if !emailRegex.MatchString(str) {
		return NewValidationError(r.GetMessage())
	}
	return nil
}

// GetMessage 获取错误消息
func (r *EmailRule) GetMessage() string {
	if r.Message != "" {
		return r.Message
	}
	return "Invalid email format"
}

// RangeRule 范围验证规则
type RangeRule struct {
	Min     float64
	Max     float64
	Message string
}

// Validate 验证
func (r *RangeRule) Validate(ctx context.Context, value interface{}) *AppError {
	if value == nil {
		return nil
	}

	var num float64
	var err error

	switch v := value.(type) {
	case int:
		num = float64(v)
	case int32:
		num = float64(v)
	case int64:
		num = float64(v)
	case float32:
		num = float64(v)
	case float64:
		num = v
	case string:
		num, err = strconv.ParseFloat(v, 64)
		if err != nil {
			return NewValidationError("Value must be a number")
		}
	default:
		return NewValidationError("Value must be a number")
	}

	if num < r.Min || num > r.Max {
		return NewValidationError(r.GetMessage())
	}
	return nil
}

// GetMessage 获取错误消息
func (r *RangeRule) GetMessage() string {
	if r.Message != "" {
		return r.Message
	}
	return fmt.Sprintf("Value must be between %.2f and %.2f", r.Min, r.Max)
}

// DateRangeRule 日期范围验证规则
type DateRangeRule struct {
	MinDate time.Time
	MaxDate time.Time
	Message string
}

// Validate 验证
func (r *DateRangeRule) Validate(ctx context.Context, value interface{}) *AppError {
	if value == nil {
		return nil
	}

	date, ok := value.(time.Time)
	if !ok {
		return NewValidationError("Value must be a time.Time")
	}

	if !r.MinDate.IsZero() && date.Before(r.MinDate) {
		return NewValidationError(r.GetMessage())
	}

	if !r.MaxDate.IsZero() && date.After(r.MaxDate) {
		return NewValidationError(r.GetMessage())
	}

	return nil
}

// GetMessage 获取错误消息
func (r *DateRangeRule) GetMessage() string {
	if r.Message != "" {
		return r.Message
	}
	return "Date is out of range"
}

// BusinessRule 业务规则接口
type BusinessRule interface {
	Validate(ctx context.Context, entity interface{}) *AppError
	GetName() string
}

// BusinessRuleValidator 业务规则验证器
type BusinessRuleValidator struct {
	rules map[string][]BusinessRule
}

// NewBusinessRuleValidator 创建业务规则验证器
func NewBusinessRuleValidator() *BusinessRuleValidator {
	return &BusinessRuleValidator{
		rules: make(map[string][]BusinessRule),
	}
}

// AddRule 添加业务规则
func (v *BusinessRuleValidator) AddRule(entityType string, rule BusinessRule) {
	if v.rules[entityType] == nil {
		v.rules[entityType] = make([]BusinessRule, 0)
	}
	v.rules[entityType] = append(v.rules[entityType], rule)
}

// Validate 验证业务规则
func (v *BusinessRuleValidator) Validate(ctx context.Context, entityType string, entity interface{}) *AppError {
	rules, exists := v.rules[entityType]
	if !exists {
		return nil
	}

	for _, rule := range rules {
		if err := rule.Validate(ctx, entity); err != nil {
			return err
		}
	}
	return nil
}

// 便捷函数

// Required 创建必填规则
func Required(message ...string) ValidationRule {
	rule := &RequiredRule{}
	if len(message) > 0 {
		rule.Message = message[0]
	}
	return rule
}

// MinLength 创建最小长度规则
func MinLength(length int, message ...string) ValidationRule {
	rule := &MinLengthRule{MinLength: length}
	if len(message) > 0 {
		rule.Message = message[0]
	}
	return rule
}

// MaxLength 创建最大长度规则
func MaxLength(length int, message ...string) ValidationRule {
	rule := &MaxLengthRule{MaxLength: length}
	if len(message) > 0 {
		rule.Message = message[0]
	}
	return rule
}

// Email 创建邮箱验证规则
func Email(message ...string) ValidationRule {
	rule := &EmailRule{}
	if len(message) > 0 {
		rule.Message = message[0]
	}
	return rule
}

// Range 创建范围验证规则
func Range(min, max float64, message ...string) ValidationRule {
	rule := &RangeRule{Min: min, Max: max}
	if len(message) > 0 {
		rule.Message = message[0]
	}
	return rule
}

// DateRange 创建日期范围验证规则
func DateRange(minDate, maxDate time.Time, message ...string) ValidationRule {
	rule := &DateRangeRule{MinDate: minDate, MaxDate: maxDate}
	if len(message) > 0 {
		rule.Message = message[0]
	}
	return rule
}