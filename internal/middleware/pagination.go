package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/galaxyerp/galaxyErp/internal/utils"
)

// PaginationMiddleware 分页中间件
func PaginationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析分页参数
		pagination := utils.ParsePaginationFromQuery(c)
		
		// 验证分页参数
		pagination.Validate()
		
		// 将分页参数存储到上下文中
		c.Set("pagination", pagination)
		
		c.Next()
	}
}

// ListMiddleware 列表中间件（包含分页、排序、搜索）
func ListMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析列表请求参数
		listReq := utils.ParseListRequestFromQuery(c)
		
		// 验证参数
		listReq.Validate(nil) // 可以在具体使用时传入允许的排序字段
		
		// 将参数存储到上下文中
		c.Set("list_request", listReq)
		c.Set("pagination", listReq.PaginationRequest)
		c.Set("sort", listReq.SortRequest)
		c.Set("search", listReq.SearchRequest)
		
		c.Next()
	}
}

// GetPaginationFromContext 从上下文获取分页参数
func GetPaginationFromContext(c *gin.Context) utils.PaginationRequest {
	if pagination, exists := c.Get("pagination"); exists {
		if p, ok := pagination.(utils.PaginationRequest); ok {
			return p
		}
	}
	return utils.DefaultPaginationRequest()
}

// GetListRequestFromContext 从上下文获取列表请求参数
func GetListRequestFromContext(c *gin.Context) utils.ListRequest {
	if listReq, exists := c.Get("list_request"); exists {
		if req, ok := listReq.(utils.ListRequest); ok {
			return req
		}
	}
	return utils.ListRequest{
		PaginationRequest: utils.DefaultPaginationRequest(),
		SortRequest:       utils.DefaultSortRequest(),
		SearchRequest:     utils.SearchRequest{},
	}
}

// GetSortFromContext 从上下文获取排序参数
func GetSortFromContext(c *gin.Context) utils.SortRequest {
	if sort, exists := c.Get("sort"); exists {
		if s, ok := sort.(utils.SortRequest); ok {
			return s
		}
	}
	return utils.DefaultSortRequest()
}

// GetSearchFromContext 从上下文获取搜索参数
func GetSearchFromContext(c *gin.Context) utils.SearchRequest {
	if search, exists := c.Get("search"); exists {
		if s, ok := search.(utils.SearchRequest); ok {
			return s
		}
	}
	return utils.SearchRequest{}
}