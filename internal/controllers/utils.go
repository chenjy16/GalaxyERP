package controllers

import (
	"strconv"

	"github.com/galaxyerp/galaxyErp/internal/common"
	"github.com/galaxyerp/galaxyErp/internal/dto"
	"github.com/galaxyerp/galaxyErp/internal/utils"
	"github.com/gin-gonic/gin"
)

// ControllerUtils controller utility class
type ControllerUtils struct{}

// NewControllerUtils creates a controller utility instance
func NewControllerUtils() *ControllerUtils {
	return &ControllerUtils{}
}

// ParseIDParam parses ID parameter from path
func (u *ControllerUtils) ParseIDParam(ctx *gin.Context, paramName string) (uint, bool) {
	idStr := ctx.Param(paramName)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		u.RespondBadRequest(ctx, "Invalid ID format")
		return 0, false
	}
	return uint(id), true
}

// BindJSON binds JSON request body
func (u *ControllerUtils) BindJSON(ctx *gin.Context, req interface{}) bool {
	if err := ctx.ShouldBindJSON(req); err != nil {
		u.RespondBadRequest(ctx, "Invalid request parameters: "+err.Error())
		return false
	}
	return true
}

// BindAndValidateJSON binds and validates JSON request body
func (u *ControllerUtils) BindAndValidateJSON(ctx *gin.Context, req interface{}) bool {
	return utils.BindAndValidate(ctx, req)
}

// BindAndValidateJSONWithMessage binds and validates JSON request body with custom error message
func (u *ControllerUtils) BindAndValidateJSONWithMessage(ctx *gin.Context, req interface{}, message string) bool {
	return utils.BindAndValidateWithCustomError(ctx, req, message)
}

// BindAndValidateQuery binds and validates query parameters
func (u *ControllerUtils) BindAndValidateQuery(ctx *gin.Context, req interface{}) bool {
	if err := ctx.ShouldBindQuery(req); err != nil {
		u.RespondBadRequest(ctx, "查询参数错误: "+err.Error())
		return false
	}
	
	// Validate the bound struct
	if validationErrors := u.ValidateStruct(req); len(validationErrors) > 0 {
		u.RespondBadRequest(ctx, "查询参数验证失败")
		return false
	}
	
	return true
}

// ValidateStruct validates struct
func (u *ControllerUtils) ValidateStruct(req interface{}) map[string]string {
	return utils.ValidateStruct(req)
}

// ValidateField validates single field
func (u *ControllerUtils) ValidateField(field interface{}, tag string) error {
	return utils.ValidateField(field, tag)
}

// ParsePaginationParams 解析分页参数
func (u *ControllerUtils) ParsePaginationParams(ctx *gin.Context) *dto.PaginationRequest {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	return &dto.PaginationRequest{
		Page:     page,
		PageSize: pageSize,
	}
}

// RespondBadRequest returns 400 error response
func (u *ControllerUtils) RespondBadRequest(ctx *gin.Context, message string) {
	common.APIBadRequestResponse(ctx, message)
}

// RespondInternalError returns 500 error response
func (u *ControllerUtils) RespondInternalError(ctx *gin.Context, message string) {
	common.APIInternalErrorResponse(ctx, message)
}

// RespondUnauthorized returns 401 error response
func (u *ControllerUtils) RespondUnauthorized(ctx *gin.Context, message string) {
	helper := common.NewAPIResponseHelper(ctx)
	helper.Unauthorized(message)
}

// RespondNotFound returns 404 error response
func (u *ControllerUtils) RespondNotFound(ctx *gin.Context, message string) {
	common.APINotFoundResponse(ctx, message)
}

// RespondSuccess returns success response
func (u *ControllerUtils) RespondSuccess(ctx *gin.Context, message string) {
	common.APISuccessResponse(ctx, nil, message)
}

// RespondCreated 返回201创建成功响应
func (u *ControllerUtils) RespondCreated(ctx *gin.Context, data interface{}) {
	common.APICreatedResponse(ctx, data)
}

// RespondOK 返回200成功响应
func (u *ControllerUtils) RespondOK(ctx *gin.Context, data interface{}) {
	common.APISuccessResponse(ctx, data)
}

// RespondNotImplemented 返回501未实现响应
func (u *ControllerUtils) RespondNotImplemented(ctx *gin.Context, message string) {
	helper := common.NewAPIResponseHelper(ctx)
	helper.NotImplemented(message)
}

// RespondPaginated 返回分页响应
func (u *ControllerUtils) RespondPaginated(ctx *gin.Context, data interface{}, pagination *common.Pagination, message ...string) {
	common.APIPaginatedResponse(ctx, data, pagination, message...)
}

// CreatePagination 创建分页对象
func (u *ControllerUtils) CreatePagination(page, pageSize int, total int64) *common.Pagination {
	return common.NewPagination(page, pageSize, total)
}
