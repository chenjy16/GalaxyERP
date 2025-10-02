package controllers

import (
	"net/http"
	"strconv"

	"github.com/galaxyerp/galaxyErp/internal/dto"
	"github.com/gin-gonic/gin"
)

// ControllerUtils 控制器工具类
type ControllerUtils struct{}

// NewControllerUtils 创建控制器工具实例
func NewControllerUtils() *ControllerUtils {
	return &ControllerUtils{}
}

// ParseIDParam 解析路径中的ID参数
func (u *ControllerUtils) ParseIDParam(ctx *gin.Context, paramName string) (uint, bool) {
	idStr := ctx.Param(paramName)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		u.RespondBadRequest(ctx, "ID格式错误")
		return 0, false
	}
	return uint(id), true
}

// BindJSON 绑定JSON请求体
func (u *ControllerUtils) BindJSON(ctx *gin.Context, req interface{}) bool {
	if err := ctx.ShouldBindJSON(req); err != nil {
		u.RespondBadRequest(ctx, "请求参数错误")
		return false
	}
	return true
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

// RespondBadRequest 返回400错误响应
func (u *ControllerUtils) RespondBadRequest(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
		Success:    false,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	})
}

// RespondInternalError 返回500错误响应
func (u *ControllerUtils) RespondInternalError(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
		Success:    false,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	})
}

// RespondUnauthorized 返回401错误响应
func (u *ControllerUtils) RespondUnauthorized(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
		Success:    false,
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	})
}

// RespondNotFound 返回404错误响应
func (u *ControllerUtils) RespondNotFound(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
		Success:    false,
		Message:    message,
		StatusCode: http.StatusNotFound,
	})
}

// RespondSuccess 返回成功响应
func (u *ControllerUtils) RespondSuccess(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Message: message,
	})
}

// RespondCreated 返回201创建成功响应
func (u *ControllerUtils) RespondCreated(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusCreated, data)
}

// RespondOK 返回200成功响应
func (u *ControllerUtils) RespondOK(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, data)
}

// RespondNotImplemented 返回501未实现响应
func (u *ControllerUtils) RespondNotImplemented(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusNotImplemented, gin.H{
		"message": message,
		"success": false,
	})
}
