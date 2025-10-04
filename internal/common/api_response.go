package common

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// StandardAPIResponse 统一的API响应结构
type StandardAPIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     *ErrorInfo  `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

// ErrorInfo 错误信息结构
type ErrorInfo struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// PaginatedAPIResponse 分页响应结构
type PaginatedAPIResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination"`
	Error      *ErrorInfo  `json:"error,omitempty"`
	Timestamp  time.Time   `json:"timestamp"`
	RequestID  string      `json:"request_id,omitempty"`
}

// Pagination 分页信息
type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// APIResponseHelper API响应助手
type APIResponseHelper struct {
	ctx *gin.Context
}

// NewAPIResponseHelper 创建API响应助手
func NewAPIResponseHelper(ctx *gin.Context) *APIResponseHelper {
	return &APIResponseHelper{ctx: ctx}
}

// Success 成功响应
func (r *APIResponseHelper) Success(data interface{}, message ...string) {
	msg := "操作成功"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}

	response := StandardAPIResponse{
		Success:   true,
		Message:   msg,
		Data:      data,
		Timestamp: time.Now(),
		RequestID: r.getRequestID(),
	}

	r.ctx.JSON(http.StatusOK, response)
}

// Created 创建成功响应
func (r *APIResponseHelper) Created(data interface{}, message ...string) {
	msg := "创建成功"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}

	response := StandardAPIResponse{
		Success:   true,
		Message:   msg,
		Data:      data,
		Timestamp: time.Now(),
		RequestID: r.getRequestID(),
	}

	r.ctx.JSON(http.StatusCreated, response)
}

// Updated 更新成功响应
func (r *APIResponseHelper) Updated(data interface{}, message ...string) {
	msg := "更新成功"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}

	response := StandardAPIResponse{
		Success:   true,
		Message:   msg,
		Data:      data,
		Timestamp: time.Now(),
		RequestID: r.getRequestID(),
	}

	r.ctx.JSON(http.StatusOK, response)
}

// Deleted 删除成功响应
func (r *APIResponseHelper) Deleted(message ...string) {
	msg := "删除成功"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}

	response := StandardAPIResponse{
		Success:   true,
		Message:   msg,
		Timestamp: time.Now(),
		RequestID: r.getRequestID(),
	}

	r.ctx.JSON(http.StatusOK, response)
}

// Error 错误响应
func (r *APIResponseHelper) Error(statusCode int, errorCode, message string, details ...map[string]interface{}) {
	var errorDetails map[string]interface{}
	if len(details) > 0 {
		errorDetails = details[0]
	}

	response := StandardAPIResponse{
		Success: false,
		Error: &ErrorInfo{
			Code:    errorCode,
			Message: message,
			Details: errorDetails,
		},
		Timestamp: time.Now(),
		RequestID: r.getRequestID(),
	}

	r.ctx.JSON(statusCode, response)
}

// BadRequest 400错误响应
func (r *APIResponseHelper) BadRequest(message string, details ...map[string]interface{}) {
	r.Error(http.StatusBadRequest, "BAD_REQUEST", message, details...)
}

// Unauthorized 401错误响应
func (r *APIResponseHelper) Unauthorized(message ...string) {
	msg := "未授权访问"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	r.Error(http.StatusUnauthorized, "UNAUTHORIZED", msg)
}

// Forbidden 403错误响应
func (r *APIResponseHelper) Forbidden(message ...string) {
	msg := "禁止访问"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	r.Error(http.StatusForbidden, "FORBIDDEN", msg)
}

// NotFound 404错误响应
func (r *APIResponseHelper) NotFound(resource string) {
	r.Error(http.StatusNotFound, "NOT_FOUND", resource+"不存在")
}

// Conflict 409错误响应
func (r *APIResponseHelper) Conflict(message string) {
	r.Error(http.StatusConflict, "CONFLICT", message)
}

// ValidationError 422验证错误响应
func (r *APIResponseHelper) ValidationError(message string, details map[string]interface{}) {
	r.Error(http.StatusUnprocessableEntity, "VALIDATION_ERROR", message, details)
}

// InternalError 500内部错误响应
func (r *APIResponseHelper) InternalError(message ...string) {
	msg := "内部服务器错误"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	r.Error(http.StatusInternalServerError, "INTERNAL_ERROR", msg)
}

// NotImplemented 501未实现响应
func (r *APIResponseHelper) NotImplemented(message ...string) {
	msg := "功能暂未实现"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	r.Error(http.StatusNotImplemented, "NOT_IMPLEMENTED", msg)
}

// Paginated 分页响应
func (r *APIResponseHelper) Paginated(data interface{}, pagination *Pagination, message ...string) {
	msg := "获取成功"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}

	response := PaginatedAPIResponse{
		Success:    true,
		Message:    msg,
		Data:       data,
		Pagination: pagination,
		Timestamp:  time.Now(),
		RequestID:  r.getRequestID(),
	}

	r.ctx.JSON(http.StatusOK, response)
}

// getRequestID 获取请求ID
func (r *APIResponseHelper) getRequestID() string {
	if requestID := r.ctx.GetString("request_id"); requestID != "" {
		return requestID
	}
	return ""
}

// 全局响应函数，用于向后兼容

// APISuccessResponse 全局成功响应函数
func APISuccessResponse(ctx *gin.Context, data interface{}, message ...string) {
	NewAPIResponseHelper(ctx).Success(data, message...)
}

// APICreatedResponse 全局创建成功响应函数
func APICreatedResponse(ctx *gin.Context, data interface{}, message ...string) {
	NewAPIResponseHelper(ctx).Created(data, message...)
}

// APIErrorResponse 全局错误响应函数
func APIErrorResponse(ctx *gin.Context, statusCode int, errorCode, message string, details ...map[string]interface{}) {
	NewAPIResponseHelper(ctx).Error(statusCode, errorCode, message, details...)
}

// APIBadRequestResponse 全局400错误响应函数
func APIBadRequestResponse(ctx *gin.Context, message string, details ...map[string]interface{}) {
	NewAPIResponseHelper(ctx).BadRequest(message, details...)
}

// APINotFoundResponse 全局404错误响应函数
func APINotFoundResponse(ctx *gin.Context, resource string) {
	NewAPIResponseHelper(ctx).NotFound(resource)
}

// APIInternalErrorResponse 全局500错误响应函数
func APIInternalErrorResponse(ctx *gin.Context, message ...string) {
	NewAPIResponseHelper(ctx).InternalError(message...)
}

// APIPaginatedResponse 全局分页响应函数
func APIPaginatedResponse(ctx *gin.Context, data interface{}, pagination *Pagination, message ...string) {
	NewAPIResponseHelper(ctx).Paginated(data, pagination, message...)
}

// NewPagination 创建分页信息
func NewPagination(page, pageSize int, total int64) *Pagination {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	
	return &Pagination{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}