package utils

import (
	"go.uber.org/zap"
)

// Logger 统一日志接口
type Logger struct {
	zap *zap.Logger
}

// NewLogger 创建新的日志实例
func NewLogger() *Logger {
	return &Logger{
		zap: zap.L(),
	}
}

// Debug 记录调试信息
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.zap.Debug(msg, fields...)
}

// Info 记录信息
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.zap.Info(msg, fields...)
}

// Warn 记录警告
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.zap.Warn(msg, fields...)
}

// Error 记录错误
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.zap.Error(msg, fields...)
}

// Fatal 记录致命错误并退出
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.zap.Fatal(msg, fields...)
}

// With 添加字段到日志上下文
func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{
		zap: l.zap.With(fields...),
	}
}

// 全局日志实例
var GlobalLogger = NewLogger()

// 便捷函数
func Debug(msg string, fields ...zap.Field) {
	GlobalLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	GlobalLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	GlobalLogger.Warn(msg, fields...)
}

func LogError(msg string, fields ...zap.Field) {
	GlobalLogger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	GlobalLogger.Fatal(msg, fields...)
}

// WithFields 创建带字段的日志实例
func WithFields(fields ...zap.Field) *Logger {
	return GlobalLogger.With(fields...)
}

// 常用字段构造函数
func String(key, val string) zap.Field {
	return zap.String(key, val)
}

func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

func Int64(key string, val int64) zap.Field {
	return zap.Int64(key, val)
}

func Float64(key string, val float64) zap.Field {
	return zap.Float64(key, val)
}

func Bool(key string, val bool) zap.Field {
	return zap.Bool(key, val)
}

func Uint(key string, val uint) zap.Field {
	return zap.Uint(key, val)
}

func ErrorField(err error) zap.Field {
	return zap.Error(err)
}

func Any(key string, val interface{}) zap.Field {
	return zap.Any(key, val)
}
