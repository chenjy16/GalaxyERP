package middleware

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/galaxyerp/galaxyErp/internal/utils"
)

// RequestLoggerConfig 请求日志配置
type RequestLoggerConfig struct {
	// 是否记录请求体
	LogRequestBody bool
	// 是否记录响应体
	LogResponseBody bool
	// 最大请求体大小（字节）
	MaxRequestBodySize int64
	// 最大响应体大小（字节）
	MaxResponseBodySize int64
	// 跳过记录的路径
	SkipPaths []string
	// 跳过记录的方法
	SkipMethods []string
}

// DefaultRequestLoggerConfig 默认请求日志配置
func DefaultRequestLoggerConfig() RequestLoggerConfig {
	return RequestLoggerConfig{
		LogRequestBody:      true,
		LogResponseBody:     false, // 默认不记录响应体，避免日志过大
		MaxRequestBodySize:  1024 * 10, // 10KB
		MaxResponseBodySize: 1024 * 10, // 10KB
		SkipPaths: []string{
			"/health",
			"/metrics",
			"/favicon.ico",
		},
		SkipMethods: []string{
			"OPTIONS",
		},
	}
}

// RequestLogger 请求日志中间件
func RequestLogger(config ...RequestLoggerConfig) gin.HandlerFunc {
	cfg := DefaultRequestLoggerConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	return func(c *gin.Context) {
		// 检查是否跳过记录
		if shouldSkip(c, cfg) {
			c.Next()
			return
		}

		// 生成请求ID
		requestID := generateRequestID()
		c.Set("request_id", requestID)

		// 记录开始时间
		startTime := time.Now()

		// 记录请求信息
		logRequestStart(c, cfg, requestID, startTime)

		// 如果需要记录响应体，包装ResponseWriter
		var responseBody *bytes.Buffer
		if cfg.LogResponseBody {
			responseBody = &bytes.Buffer{}
			c.Writer = &responseBodyWriter{
				ResponseWriter: c.Writer,
				body:          responseBody,
			}
		}

		// 处理请求
		c.Next()

		// 记录请求结束信息
		logRequestEnd(c, cfg, requestID, startTime, responseBody)
	}
}

// shouldSkip 检查是否应该跳过记录
func shouldSkip(c *gin.Context, cfg RequestLoggerConfig) bool {
	path := c.Request.URL.Path
	method := c.Request.Method

	// 检查跳过的路径
	for _, skipPath := range cfg.SkipPaths {
		if path == skipPath {
			return true
		}
	}

	// 检查跳过的方法
	for _, skipMethod := range cfg.SkipMethods {
		if method == skipMethod {
			return true
		}
	}

	return false
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	return uuid.New().String()
}

// logRequestStart 记录请求开始
func logRequestStart(c *gin.Context, cfg RequestLoggerConfig, requestID string, startTime time.Time) {
	method := c.Request.Method
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery
	clientIP := c.ClientIP()
	userAgent := c.Request.UserAgent()
	contentType := c.Request.Header.Get("Content-Type")
	contentLength := c.Request.ContentLength

	// 构建日志消息
	logMsg := fmt.Sprintf("Request started - %s %s", method, path)

	// 记录基础信息
	utils.GlobalLogger.Info(logMsg,
		utils.String("request_id", requestID),
		utils.String("method", method),
		utils.String("path", path),
		utils.String("query", query),
		utils.String("client_ip", clientIP),
		utils.String("user_agent", userAgent),
		utils.String("content_type", contentType),
		utils.Int64("content_length", contentLength),
		utils.String("timestamp", startTime.Format(time.RFC3339)),
	)

	// 记录请求体（如果配置允许）
	if cfg.LogRequestBody && method != "GET" && method != "DELETE" && contentLength > 0 {
		logRequestBody(c, cfg, requestID)
	}
}

// logRequestBody 记录请求体
func logRequestBody(c *gin.Context, cfg RequestLoggerConfig, requestID string) {
	if c.Request.Body == nil {
		return
	}

	// 读取请求体
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		utils.GlobalLogger.Warn("Failed to read request body",
			utils.String("request_id", requestID),
			utils.String("error", err.Error()),
		)
		return
	}

	// 恢复请求体，以便后续处理
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// 检查大小限制
	if int64(len(bodyBytes)) > cfg.MaxRequestBodySize {
		utils.GlobalLogger.Info("Request body logged (truncated)",
			utils.String("request_id", requestID),
			utils.String("body", string(bodyBytes[:cfg.MaxRequestBodySize])+"..."),
			utils.Int("actual_size", len(bodyBytes)),
			utils.Int64("max_size", cfg.MaxRequestBodySize),
		)
	} else {
		utils.GlobalLogger.Info("Request body logged",
			utils.String("request_id", requestID),
			utils.String("body", string(bodyBytes)),
		)
	}
}

// logRequestEnd 记录请求结束
func logRequestEnd(c *gin.Context, cfg RequestLoggerConfig, requestID string, startTime time.Time, responseBody *bytes.Buffer) {
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	statusCode := c.Writer.Status()
	responseSize := c.Writer.Size()

	method := c.Request.Method
	path := c.Request.URL.Path

	// 构建日志消息
	logMsg := fmt.Sprintf("Request completed - %s %s [%d]", method, path, statusCode)

	// 选择日志级别
	var logLevel string
	switch {
	case statusCode >= 500:
		logLevel = "error"
	case statusCode >= 400:
		logLevel = "warn"
	default:
		logLevel = "info"
	}

	// 根据状态码选择日志级别
	switch logLevel {
	case "error":
		utils.GlobalLogger.Error(logMsg,
			utils.String("request_id", requestID),
			utils.String("method", method),
			utils.String("path", path),
			utils.Int("status_code", statusCode),
			utils.Int("response_size", responseSize),
			utils.String("duration", duration.String()),
			utils.Int64("duration_ms", duration.Milliseconds()),
			utils.String("timestamp", endTime.Format(time.RFC3339)),
		)
	case "warn":
		utils.GlobalLogger.Warn(logMsg,
			utils.String("request_id", requestID),
			utils.String("method", method),
			utils.String("path", path),
			utils.Int("status_code", statusCode),
			utils.Int("response_size", responseSize),
			utils.String("duration", duration.String()),
			utils.Int64("duration_ms", duration.Milliseconds()),
			utils.String("timestamp", endTime.Format(time.RFC3339)),
		)
	default:
		utils.GlobalLogger.Info(logMsg,
			utils.String("request_id", requestID),
			utils.String("method", method),
			utils.String("path", path),
			utils.Int("status_code", statusCode),
			utils.Int("response_size", responseSize),
			utils.String("duration", duration.String()),
			utils.Int64("duration_ms", duration.Milliseconds()),
			utils.String("timestamp", endTime.Format(time.RFC3339)),
		)
	}

	// 记录响应体（如果配置允许）
	if cfg.LogResponseBody && responseBody != nil {
		logResponseBody(cfg, requestID, responseBody)
	}
}

// logResponseBody 记录响应体
func logResponseBody(cfg RequestLoggerConfig, requestID string, responseBody *bytes.Buffer) {
	bodyBytes := responseBody.Bytes()
	
	if len(bodyBytes) == 0 {
		return
	}

	// 检查大小限制
	if int64(len(bodyBytes)) > cfg.MaxResponseBodySize {
		utils.GlobalLogger.Info("Response body logged (truncated)",
			utils.String("request_id", requestID),
			utils.String("body", string(bodyBytes[:cfg.MaxResponseBodySize])+"..."),
			utils.Int("actual_size", len(bodyBytes)),
			utils.Int64("max_size", cfg.MaxResponseBodySize),
		)
	} else {
		utils.GlobalLogger.Info("Response body logged",
			utils.String("request_id", requestID),
			utils.String("body", string(bodyBytes)),
		)
	}
}

// responseBodyWriter 响应体写入器
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 写入响应体
func (w *responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// WriteString 写入字符串响应体
func (w *responseBodyWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}