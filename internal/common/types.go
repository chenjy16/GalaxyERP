package common

import (
	"github.com/galaxyerp/galaxyErp/internal/dto"
)

// SortOrder 排序方向
type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

// FilterOperator 过滤操作符
type FilterOperator string

const (
	FilterOperatorEq       FilterOperator = "eq"       // 等于
	FilterOperatorNe       FilterOperator = "ne"       // 不等于
	FilterOperatorGt       FilterOperator = "gt"       // 大于
	FilterOperatorGte      FilterOperator = "gte"      // 大于等于
	FilterOperatorLt       FilterOperator = "lt"       // 小于
	FilterOperatorLte      FilterOperator = "lte"      // 小于等于
	FilterOperatorLike     FilterOperator = "like"     // 模糊匹配
	FilterOperatorIn       FilterOperator = "in"       // 包含
	FilterOperatorNotIn    FilterOperator = "not_in"   // 不包含
	FilterOperatorBetween  FilterOperator = "between"  // 范围
	FilterOperatorIsNull   FilterOperator = "is_null"  // 为空
	FilterOperatorNotNull  FilterOperator = "not_null" // 不为空
)

// FilterCondition 过滤条件
type FilterCondition struct {
	Field    string         `json:"field"`
	Operator FilterOperator `json:"operator"`
	Value    interface{}    `json:"value"`
	Values   []interface{}  `json:"values,omitempty"` // 用于 in, not_in, between 操作
}

// SortCondition 排序条件
type SortCondition struct {
	Field string    `json:"field"`
	Order SortOrder `json:"order"`
}

// QueryOptions 查询选项
type QueryOptions struct {
	Filters    []FilterCondition      `json:"filters,omitempty"`
	Sorts      []SortCondition        `json:"sorts,omitempty"`
	Pagination *dto.PaginationRequest `json:"pagination,omitempty"`
	Includes   []string               `json:"includes,omitempty"` // 关联查询
}

// ValidationGroup 验证组枚举
type ValidationGroup string

const (
	ValidationGroupCreate ValidationGroup = "create"
	ValidationGroupUpdate ValidationGroup = "update"
	ValidationGroupQuery  ValidationGroup = "query"
	ValidationGroupDelete ValidationGroup = "delete"
)