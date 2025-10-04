package common

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/galaxyerp/galaxyErp/internal/dto"
)

// PaginationHelper 分页处理器
type PaginationHelper struct {
	DefaultPageSize int
	MaxPageSize     int
}

// NewPaginationHelper 创建分页处理器
func NewPaginationHelper(defaultPageSize, maxPageSize int) *PaginationHelper {
	return &PaginationHelper{
		DefaultPageSize: defaultPageSize,
		MaxPageSize:     maxPageSize,
	}
}

// NormalizePagination 标准化分页参数
func (h *PaginationHelper) NormalizePagination(req *dto.PaginationRequest) *dto.PaginationRequest {
	if req == nil {
		req = &dto.PaginationRequest{}
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	if req.PageSize <= 0 {
		req.PageSize = h.DefaultPageSize
	}

	if req.PageSize > h.MaxPageSize {
		req.PageSize = h.MaxPageSize
	}

	return req
}

// BuildPaginationResponse 构建分页响应
func (h *PaginationHelper) BuildPaginationResponse(data interface{}, total int64, req *dto.PaginationRequest) *dto.PaginationResponse {
	req = h.NormalizePagination(req)
	
	totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	
	return &dto.PaginationResponse{
		Page:       req.Page,
		PageSize:   req.PageSize,
		Total:      total,
		TotalPages: totalPages,
		Data:       data,
	}
}

// SearchBuilder 搜索条件构建器
type SearchBuilder struct {
	filters []FilterCondition
	sorts   []SortCondition
}

// NewSearchBuilder 创建搜索构建器
func NewSearchBuilder() *SearchBuilder {
	return &SearchBuilder{
		filters: make([]FilterCondition, 0),
		sorts:   make([]SortCondition, 0),
	}
}

// AddFilter 添加过滤条件
func (b *SearchBuilder) AddFilter(field string, operator FilterOperator, value interface{}) *SearchBuilder {
	b.filters = append(b.filters, FilterCondition{
		Field:    field,
		Operator: operator,
		Value:    value,
	})
	return b
}

// AddFilterWithValues 添加多值过滤条件
func (b *SearchBuilder) AddFilterWithValues(field string, operator FilterOperator, values []interface{}) *SearchBuilder {
	b.filters = append(b.filters, FilterCondition{
		Field:    field,
		Operator: operator,
		Values:   values,
	})
	return b
}

// AddSort 添加排序条件
func (b *SearchBuilder) AddSort(field string, order SortOrder) *SearchBuilder {
	b.sorts = append(b.sorts, SortCondition{
		Field: field,
		Order: order,
	})
	return b
}

// AddKeywordSearch 添加关键词搜索
func (b *SearchBuilder) AddKeywordSearch(keyword string, fields []string) *SearchBuilder {
	if keyword != "" && len(fields) > 0 {
		for _, field := range fields {
			b.AddFilter(field, FilterOperatorLike, keyword)
		}
	}
	return b
}

// AddDateRangeFilter 添加日期范围过滤
func (b *SearchBuilder) AddDateRangeFilter(field string, startDate, endDate time.Time) *SearchBuilder {
	if !startDate.IsZero() && !endDate.IsZero() {
		b.AddFilterWithValues(field, FilterOperatorBetween, []interface{}{startDate, endDate})
	} else if !startDate.IsZero() {
		b.AddFilter(field, FilterOperatorGte, startDate)
	} else if !endDate.IsZero() {
		b.AddFilter(field, FilterOperatorLte, endDate)
	}
	return b
}

// Build 构建查询选项
func (b *SearchBuilder) Build(pagination *dto.PaginationRequest) *QueryOptions {
	return &QueryOptions{
		Filters:    b.filters,
		Sorts:      b.sorts,
		Pagination: pagination,
	}
}

// CacheManager 缓存管理器接口
type CacheManager interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, expiry time.Duration) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context, pattern string) error
	Exists(ctx context.Context, key string) (bool, error)
}

// MemoryCacheManager 内存缓存管理器
type MemoryCacheManager struct {
	cache map[string]*cacheItem
}

type cacheItem struct {
	value  interface{}
	expiry time.Time
}

// NewMemoryCacheManager 创建内存缓存管理器
func NewMemoryCacheManager() CacheManager {
	return &MemoryCacheManager{
		cache: make(map[string]*cacheItem),
	}
}

// Get 获取缓存
func (m *MemoryCacheManager) Get(ctx context.Context, key string) (interface{}, error) {
	item, exists := m.cache[key]
	if !exists {
		return nil, NewNotFoundError("cache key")
	}

	if !item.expiry.IsZero() && time.Now().After(item.expiry) {
		delete(m.cache, key)
		return nil, NewNotFoundError("cache key")
	}

	return item.value, nil
}

// Set 设置缓存
func (m *MemoryCacheManager) Set(ctx context.Context, key string, value interface{}, expiry time.Duration) error {
	var expiryTime time.Time
	if expiry > 0 {
		expiryTime = time.Now().Add(expiry)
	}

	m.cache[key] = &cacheItem{
		value:  value,
		expiry: expiryTime,
	}
	return nil
}

// Delete 删除缓存
func (m *MemoryCacheManager) Delete(ctx context.Context, key string) error {
	delete(m.cache, key)
	return nil
}

// Clear 清空匹配模式的缓存
func (m *MemoryCacheManager) Clear(ctx context.Context, pattern string) error {
	for key := range m.cache {
		if strings.Contains(key, pattern) {
			delete(m.cache, key)
		}
	}
	return nil
}

// Exists 检查缓存是否存在
func (m *MemoryCacheManager) Exists(ctx context.Context, key string) (bool, error) {
	_, exists := m.cache[key]
	return exists, nil
}

// FileUploadHandler 文件上传处理器
type FileUploadHandler struct {
	UploadDir     string
	MaxFileSize   int64
	AllowedTypes  []string
	GenerateNames bool
}

// NewFileUploadHandler 创建文件上传处理器
func NewFileUploadHandler(uploadDir string, maxFileSize int64, allowedTypes []string) *FileUploadHandler {
	return &FileUploadHandler{
		UploadDir:     uploadDir,
		MaxFileSize:   maxFileSize,
		AllowedTypes:  allowedTypes,
		GenerateNames: true,
	}
}

// UploadFile 上传文件
func (h *FileUploadHandler) UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader) (*dto.FileUploadResponse, error) {
	// 检查文件大小
	if header.Size > h.MaxFileSize {
		return nil, NewValidationError(fmt.Sprintf("File size exceeds maximum allowed size of %d bytes", h.MaxFileSize))
	}

	// 检查文件类型
	if !h.isAllowedType(header.Filename) {
		return nil, NewValidationError(fmt.Sprintf("File type not allowed. Allowed types: %v", h.AllowedTypes))
	}

	// 确保上传目录存在
	if err := os.MkdirAll(h.UploadDir, 0755); err != nil {
		return nil, NewInternalError("Failed to create upload directory")
	}

	// 生成文件名
	filename := h.generateFilename(header.Filename)
	filePath := filepath.Join(h.UploadDir, filename)

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, NewInternalError("Failed to create file")
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, file); err != nil {
		return nil, NewInternalError("Failed to save file")
	}

	return &dto.FileUploadResponse{
		FileName: filename,
		FilePath: filePath,
		FileSize: header.Size,
		FileType: filepath.Ext(header.Filename),
	}, nil
}

// DeleteFile 删除文件
func (h *FileUploadHandler) DeleteFile(ctx context.Context, filePath string) error {
	// 确保文件在上传目录内
	if !strings.HasPrefix(filePath, h.UploadDir) {
		return NewValidationError("Invalid file path")
	}

	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return NewNotFoundError("file")
		}
		return NewInternalError("Failed to delete file")
	}

	return nil
}

// isAllowedType 检查文件类型是否允许
func (h *FileUploadHandler) isAllowedType(filename string) bool {
	if len(h.AllowedTypes) == 0 {
		return true
	}

	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowedType := range h.AllowedTypes {
		if strings.ToLower(allowedType) == ext {
			return true
		}
	}
	return false
}

// generateFilename 生成文件名
func (h *FileUploadHandler) generateFilename(originalName string) string {
	if !h.GenerateNames {
		return originalName
	}

	ext := filepath.Ext(originalName)
	name := strings.TrimSuffix(originalName, ext)
	timestamp := time.Now().Unix()
	
	return fmt.Sprintf("%s_%d%s", name, timestamp, ext)
}

// ResponseHelper 响应助手
type ResponseHelper struct{}

// NewResponseHelper 创建响应助手
func NewResponseHelper() *ResponseHelper {
	return &ResponseHelper{}
}

// Success 创建成功响应
func (h *ResponseHelper) Success(data interface{}, message ...string) *dto.SuccessResponseDTO {
	msg := "Success"
	if len(message) > 0 {
		msg = message[0]
	}
	
	return &dto.SuccessResponseDTO{
		Success: true,
		Message: msg,
		Data:    data,
	}
}

// Error 创建错误响应
func (h *ResponseHelper) Error(err error, statusCode ...int) *dto.ErrorResponse {
	code := 500
	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	if appErr := GetAppError(err); appErr != nil {
		return &dto.ErrorResponse{
			Success:    false,
			Message:    appErr.Message,
			StatusCode: appErr.StatusCode,
		}
	}

	return &dto.ErrorResponse{
		Success:    false,
		Message:    err.Error(),
		StatusCode: code,
	}
}

// List 创建列表响应
func (h *ResponseHelper) List(data interface{}, total int64) *dto.ListResponse {
	return &dto.ListResponse{
		Success: true,
		Data:    data,
		Total:   total,
	}
}

// Create 创建创建响应
func (h *ResponseHelper) Create(id uint, data interface{}, message ...string) *dto.CreateResponse {
	msg := "Created successfully"
	if len(message) > 0 {
		msg = message[0]
	}

	return &dto.CreateResponse{
		Success: true,
		Message: msg,
		ID:      id,
		Data:    data,
	}
}

// Update 创建更新响应
func (h *ResponseHelper) Update(data interface{}, message ...string) *dto.UpdateResponse {
	msg := "Updated successfully"
	if len(message) > 0 {
		msg = message[0]
	}

	return &dto.UpdateResponse{
		Success: true,
		Message: msg,
		Data:    data,
	}
}

// Delete 创建删除响应
func (h *ResponseHelper) Delete(message ...string) *dto.DeleteResponse {
	msg := "Deleted successfully"
	if len(message) > 0 {
		msg = message[0]
	}

	return &dto.DeleteResponse{
		Success: true,
		Message: msg,
	}
}