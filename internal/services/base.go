package services

import (
	"context"
	"time"

	"github.com/galaxyerp/galaxyErp/internal/common"
	"github.com/galaxyerp/galaxyErp/internal/dto"
)

// CRUDService 通用CRUD服务接口
type CRUDService[T any, CreateReq any, UpdateReq any, Response any] interface {
	Create(ctx context.Context, req *CreateReq) (*dto.CreateResponse, error)
	GetByID(ctx context.Context, id uint) (*Response, error)
	Update(ctx context.Context, id uint, req *UpdateReq) (*dto.UpdateResponse, error)
	Delete(ctx context.Context, id uint) (*dto.DeleteResponse, error)
	List(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[Response], error)
	Search(ctx context.Context, req *dto.SearchRequest) (*dto.PaginatedResponse[Response], error)
}

// ListService 列表查询服务接口
type ListService[Response any] interface {
	List(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[Response], error)
	Search(ctx context.Context, req *dto.SearchRequest) (*dto.PaginatedResponse[Response], error)
	GetAll(ctx context.Context) ([]Response, error)
}

// BatchService 批量操作服务接口
type BatchService interface {
	BatchCreate(ctx context.Context, req *dto.BatchOperationRequest) (*dto.SuccessResponseDTO, error)
	BatchUpdate(ctx context.Context, req *dto.BatchOperationRequest) (*dto.SuccessResponseDTO, error)
	BatchDelete(ctx context.Context, req *dto.BatchOperationRequest) (*dto.DeleteResponse, error)
}

// ExportService 导出服务接口
type ExportService interface {
	Export(ctx context.Context, req *dto.ExportRequest) (*dto.FileUploadResponse, error)
	Import(ctx context.Context, req *dto.ImportRequest) (*dto.SuccessResponseDTO, error)
}

// StatisticsService 统计服务接口
type StatisticsService interface {
	GetStatistics(ctx context.Context, req *dto.DateRangeRequest) (*dto.StatisticsResponse, error)
	GetDashboardData(ctx context.Context) (*dto.SuccessResponseDTO, error)
}

// AuditableService 可审计服务接口
type AuditableService interface {
	GetAuditLogs(ctx context.Context, resourceID string, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.AuditLogResponse], error)
}

// BaseServiceConfig 基础服务配置
type BaseServiceConfig struct {
	EnableAudit      bool
	EnableValidation bool
	EnableCache      bool
	CacheExpiry      time.Duration // 缓存过期时间
	EnableMetrics    bool
}

// ServiceContext 服务上下文
type ServiceContext struct {
	UserID    uint
	Username  string
	IPAddress string
	UserAgent string
}

// GetServiceContext 从context中获取服务上下文
func GetServiceContext(ctx context.Context) *ServiceContext {
	if sctx, ok := ctx.Value("service_context").(*ServiceContext); ok {
		return sctx
	}
	return &ServiceContext{}
}

// WithServiceContext 将服务上下文添加到context中
func WithServiceContext(ctx context.Context, sctx *ServiceContext) context.Context {
	return context.WithValue(ctx, "service_context", sctx)
}

// 使用 common 包中的类型定义
type (
	ValidationGroup = common.ValidationGroup
	SortOrder       = common.SortOrder
	FilterOperator  = common.FilterOperator
	FilterCondition = common.FilterCondition
	SortCondition   = common.SortCondition
	QueryOptions    = common.QueryOptions
)

// 重新导出常量
const (
	ValidationGroupCreate = common.ValidationGroupCreate
	ValidationGroupUpdate = common.ValidationGroupUpdate
	ValidationGroupQuery  = common.ValidationGroupQuery
	ValidationGroupDelete = common.ValidationGroupDelete

	SortOrderAsc  = common.SortOrderAsc
	SortOrderDesc = common.SortOrderDesc

	FilterOperatorEq      = common.FilterOperatorEq
	FilterOperatorNe      = common.FilterOperatorNe
	FilterOperatorGt      = common.FilterOperatorGt
	FilterOperatorGte     = common.FilterOperatorGte
	FilterOperatorLt      = common.FilterOperatorLt
	FilterOperatorLte     = common.FilterOperatorLte
	FilterOperatorLike    = common.FilterOperatorLike
	FilterOperatorIn      = common.FilterOperatorIn
	FilterOperatorNotIn   = common.FilterOperatorNotIn
	FilterOperatorBetween = common.FilterOperatorBetween
	FilterOperatorIsNull  = common.FilterOperatorIsNull
	FilterOperatorNotNull = common.FilterOperatorNotNull
)

// BaseService 基础服务实现
type BaseService struct {
	config           *BaseServiceConfig
	validator        *common.BusinessRuleValidator
	paginationHelper *common.PaginationHelper
	searchBuilder    *common.SearchBuilder
	cacheManager     common.CacheManager
	responseHelper   *common.ResponseHelper
}

// NewBaseService 创建基础服务实例
func NewBaseService(config *BaseServiceConfig) *BaseService {
	return &BaseService{
		config:           config,
		validator:        common.NewBusinessRuleValidator(),
		paginationHelper: common.NewPaginationHelper(20, 100), // 默认页大小20，最大100
		searchBuilder:    common.NewSearchBuilder(),
		cacheManager:     common.NewMemoryCacheManager(),
		responseHelper:   common.NewResponseHelper(),
	}
}

// ValidateRequest 验证请求
func (s *BaseService) ValidateRequest(ctx context.Context, entityType string, req interface{}) error {
	if !s.config.EnableValidation {
		return nil
	}
	if appErr := s.validator.Validate(ctx, entityType, req); appErr != nil {
		return appErr
	}
	return nil
}

// BuildPaginationResponse 构建分页响应
func (s *BaseService) BuildPaginationResponse(data interface{}, total int64, req *dto.PaginationRequest) *dto.PaginationResponse {
	return s.paginationHelper.BuildPaginationResponse(data, total, req)
}

// BuildSearchConditions 构建搜索条件
func (s *BaseService) BuildSearchConditions(filters []FilterCondition, sorts []SortCondition) *QueryOptions {
	for _, filter := range filters {
		s.searchBuilder.AddFilter(filter.Field, filter.Operator, filter.Value)
	}
	for _, sort := range sorts {
		s.searchBuilder.AddSort(sort.Field, sort.Order)
	}
	return s.searchBuilder.Build(nil)
}

// GetFromCache 从缓存获取数据
func (s *BaseService) GetFromCache(ctx context.Context, key string) (interface{}, bool) {
	if !s.config.EnableCache {
		return nil, false
	}
	value, err := s.cacheManager.Get(ctx, key)
	if err != nil {
		return nil, false
	}
	return value, true
}

// SetToCache 设置缓存数据
func (s *BaseService) SetToCache(ctx context.Context, key string, value interface{}) error {
	if !s.config.EnableCache {
		return nil
	}
	return s.cacheManager.Set(ctx, key, value, s.config.CacheExpiry)
}

// DeleteFromCache 删除缓存数据
func (s *BaseService) DeleteFromCache(ctx context.Context, key string) error {
	if !s.config.EnableCache {
		return nil
	}
	return s.cacheManager.Delete(ctx, key)
}

// CreateSuccessResponse 创建成功响应
func (s *BaseService) CreateSuccessResponse(data interface{}, message ...string) *dto.SuccessResponseDTO {
	return s.responseHelper.Success(data, message...)
}

// CreateErrorResponse 创建错误响应
func (s *BaseService) CreateErrorResponse(err error, statusCode ...int) *dto.ErrorResponse {
	return s.responseHelper.Error(err, statusCode...)
}

// CreateListResponse 创建列表响应
func (s *BaseService) CreateListResponse(data interface{}, total int64) *dto.ListResponse {
	return s.responseHelper.List(data, total)
}

// CreateCreateResponse 创建创建响应
func (s *BaseService) CreateCreateResponse(id uint, data interface{}, message ...string) *dto.CreateResponse {
	return s.responseHelper.Create(id, data, message...)
}

// CreateUpdateResponse 创建更新响应
func (s *BaseService) CreateUpdateResponse(data interface{}, message ...string) *dto.UpdateResponse {
	return s.responseHelper.Update(data, message...)
}

// CreateDeleteResponse 创建删除响应
func (s *BaseService) CreateDeleteResponse(message ...string) *dto.DeleteResponse {
	return s.responseHelper.Delete(message...)
}