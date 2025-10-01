package dto

import (
	"time"
)

// ProductCreateRequest 产品创建请求
type ProductCreateRequest struct {
	Code        string  `json:"code" binding:"required,max=50"`
	Name        string  `json:"name" binding:"required,max=100"`
	Description string  `json:"description,omitempty"`
	Category    string  `json:"category" binding:"required,max=50"`
	Unit        string  `json:"unit" binding:"required,max=20"`
	Price       float64 `json:"price" binding:"required,min=0"`
	Cost        float64 `json:"cost" binding:"required,min=0"`
	Status      string  `json:"status" binding:"required,oneof=active inactive"`
}

// ProductUpdateRequest 产品更新请求
type ProductUpdateRequest struct {
	Code        string   `json:"code,omitempty" binding:"omitempty,max=50"`
	Name        string   `json:"name,omitempty" binding:"omitempty,max=100"`
	Description string   `json:"description,omitempty"`
	Category    string   `json:"category,omitempty" binding:"omitempty,max=50"`
	Unit        string   `json:"unit,omitempty" binding:"omitempty,max=20"`
	Price       *float64 `json:"price,omitempty" binding:"omitempty,min=0"`
	Cost        *float64 `json:"cost,omitempty" binding:"omitempty,min=0"`
	Status      string   `json:"status,omitempty" binding:"omitempty,oneof=active inactive"`
}

// ProductResponse 产品响应
type ProductResponse struct {
	ID          uint      `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Category    string    `json:"category"`
	Unit        string    `json:"unit"`
	Price       float64   `json:"price"`
	Cost        float64   `json:"cost"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProductListResponse 产品列表响应
type ProductListResponse struct {
	ID       uint    `json:"id"`
	Code     string  `json:"code"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Unit     string  `json:"unit"`
	Price    float64 `json:"price"`
	Cost     float64 `json:"cost"`
	Status   string  `json:"status"`
}

// ProductSearchRequest 产品搜索请求
type ProductSearchRequest struct {
	SearchRequest
	Category string `json:"category,omitempty" form:"category"`
	Status   string `json:"status,omitempty" form:"status"`
	MinPrice *float64 `json:"min_price,omitempty" form:"min_price"`
	MaxPrice *float64 `json:"max_price,omitempty" form:"max_price"`
}

// ProductFilter 产品过滤器
type ProductFilter struct {
	PaginationRequest
	Category string `json:"category,omitempty" form:"category"`
	Status   string `json:"status,omitempty" form:"status"`
	Name     string `json:"name,omitempty" form:"name"`
	Code     string `json:"code,omitempty" form:"code"`
}