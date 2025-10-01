package dto

import (
	"time"
)

// BaseResponse 基础响应结构
type BaseResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PaginationRequest 分页请求
type PaginationRequest struct {
	Page     int    `json:"page" form:"page" binding:"min=1"`
	PageSize int    `json:"page_size" form:"page_size" binding:"min=1,max=100"`
	SortBy   string `json:"sort_by,omitempty" form:"sort_by"`
	SortDesc bool   `json:"sort_desc,omitempty" form:"sort_desc"`
}

// GetOffset 获取偏移量
func (p *PaginationRequest) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return (p.Page - 1) * p.GetLimit()
}

// GetLimit 获取限制数量
func (p *PaginationRequest) GetLimit() int {
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	return p.PageSize
}

// PaginationResponse 分页响应
type PaginationResponse struct {
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
	Data       interface{} `json:"data"`
}

// IDRequest ID请求
type IDRequest struct {
	ID uint `json:"id" uri:"id" binding:"required,min=1"`
}

// IDsRequest 多个ID请求
type IDsRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"`
}

// StatusRequest 状态更新请求
type StatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// DateRangeRequest 日期范围请求
type DateRangeRequest struct {
	StartDate time.Time `json:"start_date" form:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" form:"end_date" binding:"required"`
}

// SearchRequest 搜索请求
type SearchRequest struct {
	PaginationRequest
	Keyword string `json:"keyword,omitempty" form:"keyword"`
	Status  string `json:"status,omitempty" form:"status"`
}

// BulkOperationRequest 批量操作请求
type BulkOperationRequest struct {
	IDs       []uint `json:"ids" binding:"required,min=1"`
	Operation string `json:"operation" binding:"required"`
}

// PaginatedResponse 泛型分页响应
type PaginatedResponse[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
}

// BaseModel 基础模型
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

// SuccessResponse 成功响应
type SuccessResponse struct {
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
	Format    string    `json:"format" form:"format" binding:"required,oneof=excel pdf csv"`
	StartDate time.Time `json:"start_date,omitempty" form:"start_date"`
	EndDate   time.Time `json:"end_date,omitempty" form:"end_date"`
	Filters   map[string]interface{} `json:"filters,omitempty"`
}

// ImportRequest 导入请求
type ImportRequest struct {
	FilePath string `json:"file_path" binding:"required"`
	Options  map[string]interface{} `json:"options,omitempty"`
}