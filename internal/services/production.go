package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/galaxyerp/galaxyErp/internal/dto"
	"github.com/galaxyerp/galaxyErp/internal/models"
	"github.com/galaxyerp/galaxyErp/internal/repositories"
)

// ProductService 产品服务接口
type ProductService interface {
	CreateProduct(ctx context.Context, req *dto.ProductCreateRequest) (*dto.ProductResponse, error)
	GetProduct(ctx context.Context, id uint) (*dto.ProductResponse, error)
	UpdateProduct(ctx context.Context, id uint, req *dto.ProductUpdateRequest) (*dto.ProductResponse, error)
	DeleteProduct(ctx context.Context, id uint) error
	GetProducts(ctx context.Context, req *dto.PaginationRequest) (*dto.BaseResponse, error)
	SearchProducts(ctx context.Context, req *dto.ProductSearchRequest) (*dto.BaseResponse, error)
}

// ProductServiceImpl 产品服务实现
type ProductServiceImpl struct {
	productRepo repositories.ProductRepository
}

// NewProductService 创建产品服务实例
func NewProductService(productRepo repositories.ProductRepository) ProductService {
	return &ProductServiceImpl{
		productRepo: productRepo,
	}
}

// CreateProduct 创建产品
func (s *ProductServiceImpl) CreateProduct(ctx context.Context, req *dto.ProductCreateRequest) (*dto.ProductResponse, error) {
	product := &models.Product{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Unit:        req.Unit,
		Price:       req.Price,
		Cost:        req.Cost,
		Status:      req.Status,
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("创建产品失败: %w", err)
	}

	return s.toProductResponse(product), nil
}

// GetProduct 获取产品
func (s *ProductServiceImpl) GetProduct(ctx context.Context, id uint) (*dto.ProductResponse, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取产品失败: %w", err)
	}
	if product == nil {
		return nil, errors.New("产品不存在")
	}

	return s.toProductResponse(product), nil
}

// UpdateProduct 更新产品
func (s *ProductServiceImpl) UpdateProduct(ctx context.Context, id uint, req *dto.ProductUpdateRequest) (*dto.ProductResponse, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取产品失败: %w", err)
	}
	if product == nil {
		return nil, errors.New("产品不存在")
	}

	// 更新字段
	if req.Code != "" {
		product.Code = req.Code
	}
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Category != "" {
		product.Category = req.Category
	}
	if req.Unit != "" {
		product.Unit = req.Unit
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Cost != nil {
		product.Cost = *req.Cost
	}
	if req.Status != "" {
		product.Status = req.Status
	}

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("更新产品失败: %w", err)
	}

	return s.toProductResponse(product), nil
}

// DeleteProduct 删除产品
func (s *ProductServiceImpl) DeleteProduct(ctx context.Context, id uint) error {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取产品失败: %w", err)
	}
	if product == nil {
		return errors.New("产品不存在")
	}

	if err := s.productRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("删除产品失败: %w", err)
	}

	return nil
}

// GetProducts 获取产品列表
func (s *ProductServiceImpl) GetProducts(ctx context.Context, req *dto.PaginationRequest) (*dto.BaseResponse, error) {
	products, _, err := s.productRepo.List(ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, fmt.Errorf("获取产品列表失败: %w", err)
	}

	var productResponses []*dto.ProductResponse
	for _, product := range products {
		productResponses = append(productResponses, s.toProductResponse(product))
	}

	return &dto.BaseResponse{
		Success: true,
		Message: "获取产品列表成功",
		Data:    productResponses,
	}, nil
}

// SearchProducts 搜索产品
func (s *ProductServiceImpl) SearchProducts(ctx context.Context, req *dto.ProductSearchRequest) (*dto.BaseResponse, error) {
	products, _, err := s.productRepo.Search(ctx, req.Keyword, req.Page, req.PageSize)
	if err != nil {
		return nil, fmt.Errorf("搜索产品失败: %w", err)
	}

	var productResponses []*dto.ProductResponse
	for _, product := range products {
		productResponses = append(productResponses, s.toProductResponse(product))
	}

	return &dto.BaseResponse{
		Success: true,
		Message: "搜索产品成功",
		Data:    productResponses,
	}, nil
}

// toProductResponse 转换为产品响应
func (s *ProductServiceImpl) toProductResponse(product *models.Product) *dto.ProductResponse {
	return &dto.ProductResponse{
		ID:          product.ID,
		Code:        product.Code,
		Name:        product.Name,
		Description: product.Description,
		Category:    product.Category,
		Unit:        product.Unit,
		Price:       product.Price,
		Cost:        product.Cost,
		Status:      product.Status,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}
