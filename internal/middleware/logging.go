package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/galaxyerp/galaxyErp/internal/utils"
)

// responseWriter 包装 gin.ResponseWriter 以捕获响应内容
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// LoggingMiddleware 创建HTTP请求日志中间件
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 读取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 包装响应写入器
		w := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = w

		// 处理请求
		c.Next()

		// 计算处理时间
		duration := time.Since(start)

		// 构建日志字段
		fields := []zap.Field{
			utils.String("method", c.Request.Method),
			utils.String("path", c.Request.URL.Path),
			utils.String("query", c.Request.URL.RawQuery),
			utils.String("ip", c.ClientIP()),
			utils.String("user_agent", c.Request.UserAgent()),
			utils.Int("status", c.Writer.Status()),
			utils.Int64("response_time_ms", duration.Milliseconds()),
			utils.Int("response_size", c.Writer.Size()),
		}

		// 添加请求体（仅对非GET请求且内容不为空时）
		if c.Request.Method != "GET" && len(requestBody) > 0 && len(requestBody) < 1024 {
			fields = append(fields, utils.String("request_body", string(requestBody)))
		}

		// 添加响应体（仅对小响应）
		if w.body.Len() > 0 && w.body.Len() < 1024 {
			fields = append(fields, utils.String("response_body", w.body.String()))
		}

		// 添加用户信息（如果存在）
		if userID, exists := c.Get("user_id"); exists {
			fields = append(fields, utils.Any("user_id", userID))
		}

		// 根据状态码选择日志级别
		switch {
		case c.Writer.Status() >= 500:
			utils.LogError("HTTP Request", fields...)
		case c.Writer.Status() >= 400:
			utils.Warn("HTTP Request", fields...)
		default:
			utils.Info("HTTP Request", fields...)
		}
	}
}
