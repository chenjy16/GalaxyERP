package controllers

import (
	"github.com/galaxyerp/galaxyErp/internal/dto"
	"github.com/galaxyerp/galaxyErp/internal/services"
	"github.com/gin-gonic/gin"
)

// ProductionController 生产控制器
type ProductionController struct {
	productService services.ProductService
	utils          *ControllerUtils
}

// NewProductionController 创建生产控制器实例
func NewProductionController(productService services.ProductService) *ProductionController {
	return &ProductionController{
		productService: productService,
		utils:          NewControllerUtils(),
	}
}

// CreateProduct 创建产品
// @Summary 创建产品
// @Description 创建新产品
// @Tags 产品管理
// @Accept json
// @Produce json
// @Param product body dto.ProductCreateRequest true "产品信息"
// @Success 201 {object} dto.ProductResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/products [post]
func (c *ProductionController) CreateProduct(ctx *gin.Context) {
	var req dto.ProductCreateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	product, err := c.productService.CreateProduct(ctx.Request.Context(), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "创建产品失败")
		return
	}

	c.utils.RespondCreated(ctx, product)
}

// GetProduct 获取产品
// @Summary 获取产品
// @Description 根据ID获取产品信息
// @Tags 产品管理
// @Accept json
// @Produce json
// @Param id path int true "产品ID"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/products/{id} [get]
func (c *ProductionController) GetProduct(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	product, err := c.productService.GetProduct(ctx.Request.Context(), id)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取产品失败")
		return
	}

	c.utils.RespondOK(ctx, product)
}

// UpdateProduct 更新产品
// @Summary 更新产品
// @Description 更新产品信息
// @Tags 产品管理
// @Accept json
// @Produce json
// @Param id path int true "产品ID"
// @Param product body dto.ProductUpdateRequest true "产品信息"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/products/{id} [put]
func (c *ProductionController) UpdateProduct(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	var req dto.ProductUpdateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	_, err := c.productService.UpdateProduct(ctx.Request.Context(), id, &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "更新产品失败")
		return
	}

	c.utils.RespondSuccess(ctx, "更新产品成功")
}

// DeleteProduct 删除产品
// @Summary 删除产品
// @Description 删除产品
// @Tags 产品管理
// @Accept json
// @Produce json
// @Param id path int true "产品ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/products/{id} [delete]
func (c *ProductionController) DeleteProduct(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	err := c.productService.DeleteProduct(ctx.Request.Context(), id)
	if err != nil {
		c.utils.RespondInternalError(ctx, "删除产品失败")
		return
	}

	c.utils.RespondSuccess(ctx, "删除产品成功")
}

// ListProducts 获取产品列表
// @Summary 获取产品列表
// @Description 获取产品列表
// @Tags 产品管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} dto.ProductListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/products [get]
func (c *ProductionController) ListProducts(ctx *gin.Context) {
	req := c.utils.ParsePaginationParams(ctx)

	response, err := c.productService.GetProducts(ctx.Request.Context(), req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取产品列表失败")
		return
	}

	c.utils.RespondOK(ctx, response)
}

// SearchProducts 搜索产品
// @Summary 搜索产品
// @Description 搜索产品
// @Tags 产品管理
// @Accept json
// @Produce json
// @Param keyword query string false "搜索关键词"
// @Param category query string false "产品分类"
// @Param status query string false "产品状态"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} dto.ProductListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/products/search [get]
func (c *ProductionController) SearchProducts(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	category := ctx.Query("category")
	status := ctx.Query("status")
	pagination := c.utils.ParsePaginationParams(ctx)

	req := &dto.ProductSearchRequest{
		SearchRequest: dto.SearchRequest{
			PaginationRequest: *pagination,
			Keyword:           keyword,
			Status:            status,
		},
		Category: category,
	}

	response, err := c.productService.SearchProducts(ctx.Request.Context(), req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "搜索产品失败")
		return
	}

	c.utils.RespondOK(ctx, response)
}
