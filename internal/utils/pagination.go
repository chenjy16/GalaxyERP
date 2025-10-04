package utils

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PaginationRequest 分页请求参数
type PaginationRequest struct {
	Page     int `json:"page" form:"page" binding:"min=1" label:"页码"`
	PageSize int `json:"page_size" form:"page_size" binding:"min=1,max=100" label:"每页大小"`
}

// PaginationResponse 分页响应信息
type PaginationResponse struct {
	Page       int   `json:"page"`        // 当前页码
	PageSize   int   `json:"page_size"`   // 每页大小
	Total      int64 `json:"total"`       // 总记录数
	TotalPages int   `json:"total_pages"` // 总页数
	HasNext    bool  `json:"has_next"`    // 是否有下一页
	HasPrev    bool  `json:"has_prev"`    // 是否有上一页
}

// PaginatedResult 分页结果
type PaginatedResult struct {
	Data       interface{}        `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}

// DefaultPaginationRequest 默认分页请求
func DefaultPaginationRequest() PaginationRequest {
	return PaginationRequest{
		Page:     1,
		PageSize: 20,
	}
}

// ParsePaginationFromQuery 从查询参数解析分页信息
func ParsePaginationFromQuery(c *gin.Context) PaginationRequest {
	pagination := DefaultPaginationRequest()

	// 解析页码
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			pagination.Page = page
		}
	}

	// 解析每页大小
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 && pageSize <= 100 {
			pagination.PageSize = pageSize
		}
	}

	return pagination
}

// Validate 验证分页参数
func (p *PaginationRequest) Validate() error {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 20
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	return nil
}

// GetOffset 获取偏移量
func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// GetLimit 获取限制数量
func (p *PaginationRequest) GetLimit() int {
	return p.PageSize
}

// BuildPaginationResponse 构建分页响应
func BuildPaginationResponse(req PaginationRequest, total int64) PaginationResponse {
	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))
	
	return PaginationResponse{
		Page:       req.Page,
		PageSize:   req.PageSize,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    req.Page < totalPages,
		HasPrev:    req.Page > 1,
	}
}

// PaginateQuery 对GORM查询应用分页
func PaginateQuery(db *gorm.DB, req PaginationRequest) *gorm.DB {
	req.Validate()
	offset := req.GetOffset()
	limit := req.GetLimit()
	
	return db.Offset(offset).Limit(limit)
}

// PaginateWithCount 分页查询并计算总数
func PaginateWithCount(db *gorm.DB, req PaginationRequest, result interface{}) (PaginationResponse, error) {
	req.Validate()
	
	// 计算总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return PaginationResponse{}, err
	}
	
	// 应用分页查询
	if err := PaginateQuery(db, req).Find(result).Error; err != nil {
		return PaginationResponse{}, err
	}
	
	// 构建分页响应
	pagination := BuildPaginationResponse(req, total)
	
	return pagination, nil
}

// PaginateWithCountAndOrder 分页查询并计算总数（带排序）
func PaginateWithCountAndOrder(db *gorm.DB, req PaginationRequest, orderBy string, result interface{}) (PaginationResponse, error) {
	req.Validate()
	
	// 计算总数（不包含排序，提高性能）
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return PaginationResponse{}, err
	}
	
	// 应用排序和分页查询
	query := db
	if orderBy != "" {
		query = query.Order(orderBy)
	}
	
	if err := PaginateQuery(query, req).Find(result).Error; err != nil {
		return PaginationResponse{}, err
	}
	
	// 构建分页响应
	pagination := BuildPaginationResponse(req, total)
	
	return pagination, nil
}

// PaginationHelper 分页助手
type PaginationHelper struct {
	db *gorm.DB
}

// NewPaginationHelper 创建分页助手
func NewPaginationHelper(db *gorm.DB) *PaginationHelper {
	return &PaginationHelper{db: db}
}

// Paginate 执行分页查询
func (h *PaginationHelper) Paginate(req PaginationRequest, result interface{}, scopes ...func(*gorm.DB) *gorm.DB) (PaginationResponse, error) {
	req.Validate()
	
	// 应用作用域
	query := h.db
	for _, scope := range scopes {
		query = scope(query)
	}
	
	return PaginateWithCount(query, req, result)
}

// PaginateWithOrder 执行带排序的分页查询
func (h *PaginationHelper) PaginateWithOrder(req PaginationRequest, orderBy string, result interface{}, scopes ...func(*gorm.DB) *gorm.DB) (PaginationResponse, error) {
	req.Validate()
	
	// 应用作用域
	query := h.db
	for _, scope := range scopes {
		query = scope(query)
	}
	
	return PaginateWithCountAndOrder(query, req, orderBy, result)
}

// PaginateWithSearch 执行带搜索的分页查询
func (h *PaginationHelper) PaginateWithSearch(req PaginationRequest, searchFields []string, searchTerm string, result interface{}, scopes ...func(*gorm.DB) *gorm.DB) (PaginationResponse, error) {
	req.Validate()
	
	// 应用作用域
	query := h.db
	for _, scope := range scopes {
		query = scope(query)
	}
	
	// 应用搜索条件
	if searchTerm != "" && len(searchFields) > 0 {
		searchQuery := query.Where("FALSE") // 初始化为假条件
		for _, field := range searchFields {
			searchQuery = searchQuery.Or(field+" LIKE ?", "%"+searchTerm+"%")
		}
		query = searchQuery
	}
	
	return PaginateWithCount(query, req, result)
}

// SortRequest 排序请求参数
type SortRequest struct {
	SortBy    string `json:"sort_by" form:"sort_by" label:"排序字段"`
	SortOrder string `json:"sort_order" form:"sort_order" binding:"oneof=asc desc" label:"排序方向"`
}

// DefaultSortRequest 默认排序请求
func DefaultSortRequest() SortRequest {
	return SortRequest{
		SortBy:    "id",
		SortOrder: "desc",
	}
}

// ParseSortFromQuery 从查询参数解析排序信息
func ParseSortFromQuery(c *gin.Context) SortRequest {
	sort := DefaultSortRequest()

	if sortBy := c.Query("sort_by"); sortBy != "" {
		sort.SortBy = sortBy
	}

	if sortOrder := c.Query("sort_order"); sortOrder != "" {
		if sortOrder == "asc" || sortOrder == "desc" {
			sort.SortOrder = sortOrder
		}
	}

	return sort
}

// GetOrderBy 获取排序字符串
func (s *SortRequest) GetOrderBy() string {
	if s.SortBy == "" {
		return ""
	}
	return s.SortBy + " " + s.SortOrder
}

// Validate 验证排序参数
func (s *SortRequest) Validate(allowedFields []string) error {
	// 检查排序字段是否在允许列表中
	if s.SortBy != "" && len(allowedFields) > 0 {
		allowed := false
		for _, field := range allowedFields {
			if s.SortBy == field {
				allowed = true
				break
			}
		}
		if !allowed {
			s.SortBy = "id" // 回退到默认字段
		}
	}

	// 验证排序方向
	if s.SortOrder != "asc" && s.SortOrder != "desc" {
		s.SortOrder = "desc"
	}

	return nil
}

// SearchRequest 搜索请求参数
type SearchRequest struct {
	Search string `json:"search" form:"search" label:"搜索关键词"`
}

// ParseSearchFromQuery 从查询参数解析搜索信息
func ParseSearchFromQuery(c *gin.Context) SearchRequest {
	return SearchRequest{
		Search: c.Query("search"),
	}
}

// ListRequest 列表请求参数（包含分页、排序、搜索）
type ListRequest struct {
	PaginationRequest
	SortRequest
	SearchRequest
}

// ParseListRequestFromQuery 从查询参数解析列表请求
func ParseListRequestFromQuery(c *gin.Context) ListRequest {
	return ListRequest{
		PaginationRequest: ParsePaginationFromQuery(c),
		SortRequest:       ParseSortFromQuery(c),
		SearchRequest:     ParseSearchFromQuery(c),
	}
}

// Validate 验证列表请求参数
func (r *ListRequest) Validate(allowedSortFields []string) error {
	r.PaginationRequest.Validate()
	r.SortRequest.Validate(allowedSortFields)
	return nil
}

// BuildListResponse 构建列表响应
func BuildListResponse(data interface{}, pagination PaginationResponse) PaginatedResult {
	return PaginatedResult{
		Data:       data,
		Pagination: pagination,
	}
}

// PaginateList 执行完整的列表查询（分页+排序+搜索）
func PaginateList(db *gorm.DB, req ListRequest, searchFields []string, result interface{}, scopes ...func(*gorm.DB) *gorm.DB) (PaginatedResult, error) {
	req.Validate(nil) // 这里可以传入允许的排序字段
	
	// 应用作用域
	query := db
	for _, scope := range scopes {
		query = scope(query)
	}
	
	// 应用搜索条件
	if req.Search != "" && len(searchFields) > 0 {
		searchQuery := query.Where("FALSE") // 初始化为假条件
		for _, field := range searchFields {
			searchQuery = searchQuery.Or(field+" LIKE ?", "%"+req.Search+"%")
		}
		query = searchQuery
	}
	
	// 应用排序
	if orderBy := req.SortRequest.GetOrderBy(); orderBy != "" {
		query = query.Order(orderBy)
	}
	
	// 执行分页查询
	pagination, err := PaginateWithCount(query, req.PaginationRequest, result)
	if err != nil {
		return PaginatedResult{}, err
	}
	
	return BuildListResponse(result, pagination), nil
}