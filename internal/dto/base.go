package dto

import (
	"time"
)

// BaseResponse base response structure
type BaseResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PaginationRequest pagination request
type PaginationRequest struct {
	Page     int    `json:"page" form:"page" validate:"min=1"`
	PageSize int    `json:"page_size" form:"page_size" validate:"min=1,max=100"`
	Search   string `json:"search" form:"search"`
	SortBy   string `json:"sort_by" form:"sort_by"`
	SortDesc bool   `json:"sort_desc" form:"sort_desc"`
}

// GetOffset gets offset value
func (p *PaginationRequest) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return (p.Page - 1) * p.GetLimit()
}

// GetLimit gets limit value
func (p *PaginationRequest) GetLimit() int {
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	return p.PageSize
}

// PaginationResponse pagination response
type PaginationResponse struct {
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
	Data       interface{} `json:"data"`
}

// IDRequest ID request
type IDRequest struct {
	ID uint `json:"id" uri:"id" validate:"required,min=1"`
}

// IDsRequest multiple IDs request
type IDsRequest struct {
	IDs []uint `json:"ids" validate:"required,min=1"`
}

// StatusRequest status update request
type StatusRequest struct {
	Status string `json:"status" validate:"required"`
}

// DateRangeRequest date range request
type DateRangeRequest struct {
	StartDate time.Time `json:"start_date" form:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" form:"end_date" validate:"required"`
}

// SearchRequest search request
type SearchRequest struct {
	PaginationRequest
	Keyword string `json:"keyword,omitempty" form:"keyword"`
	Status  string `json:"status,omitempty" form:"status"`
}

// BulkOperationRequest bulk operation request
type BatchOperationRequest struct {
	IDs       []uint `json:"ids" validate:"required,min=1"`
	Operation string `json:"operation" validate:"required"`
}

// PaginatedResponse generic paginated response
type PaginatedResponse[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
}

// BaseModel base model
type BaseModel struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FileUploadResponse 文件上传响应
type FileUploadResponse struct {
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
	FileSize int64  `json:"file_size"`
	FileType string `json:"file_type"`
}

// SuccessResponse 成功响应
func SuccessResponse(data interface{}, message ...string) *BaseResponse {
	msg := "Success"
	if len(message) > 0 {
		msg = message[0]
	}
	return &BaseResponse{
		Success: true,
		Message: msg,
		Data:    data,
	}
}

// ValidationError 验证错误
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Success    bool              `json:"success"`
	Message    string            `json:"message"`
	Errors     []ValidationError `json:"errors,omitempty"`
	StatusCode int               `json:"status_code"`
}

// 通用响应结构
type SuccessResponseDTO struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ListResponse 列表响应
type ListResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Total   int64       `json:"total,omitempty"`
}

// CreateResponse 创建响应
type CreateResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	ID      uint        `json:"id"`
	Data    interface{} `json:"data,omitempty"`
}

// UpdateResponse 更新响应
type UpdateResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// DeleteResponse 删除响应
type DeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// StatisticsResponse 统计响应
type StatisticsResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Period  string      `json:"period,omitempty"`
}

// ExportRequest 导出请求
type ExportRequest struct {
	Format    string                 `json:"format" form:"format" validate:"required,oneof=excel pdf csv"`
	StartDate time.Time              `json:"start_date,omitempty" form:"start_date"`
	EndDate   time.Time              `json:"end_date,omitempty" form:"end_date"`
	Filters   map[string]interface{} `json:"filters,omitempty"`
}

// ImportRequest 导入请求
type ImportRequest struct {
	FilePath string                 `json:"file_path" validate:"required"`
	Options  map[string]interface{} `json:"options,omitempty"`
}

// AuditLogSearchRequest 审计日志搜索请求
type AuditLogSearchRequest struct {
	PaginationRequest
	UserID     *uint     `json:"user_id,omitempty" form:"user_id"`
	Username   string    `json:"username,omitempty" form:"username"`
	Action     string    `json:"action,omitempty" form:"action"`
	Resource   string    `json:"resource,omitempty" form:"resource"`
	ResourceID string    `json:"resource_id,omitempty" form:"resource_id"`
	Status     string    `json:"status,omitempty" form:"status"`
	IPAddress  string    `json:"ip_address,omitempty" form:"ip_address"`
	StartTime  time.Time `json:"start_time,omitempty" form:"start_time"`
	EndTime    time.Time `json:"end_time,omitempty" form:"end_time"`
}

// AuditLogResponse 审计日志响应
type AuditLogResponse struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	Username    string    `json:"username"`
	Action      string    `json:"action"`
	Resource    string    `json:"resource"`
	ResourceID  string    `json:"resource_id,omitempty"`
	Method      string    `json:"method,omitempty"`
	Path        string    `json:"path,omitempty"`
	Description string    `json:"description,omitempty"`
	IPAddress   string    `json:"ip_address,omitempty"`
	UserAgent   string    `json:"user_agent,omitempty"`
	OldValues   string    `json:"old_values,omitempty"`
	NewValues   string    `json:"new_values,omitempty"`
	Changes     string    `json:"changes,omitempty"`
	Status      string    `json:"status"`
	ErrorMsg    string    `json:"error_msg,omitempty"`
	Duration    int64     `json:"duration,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
