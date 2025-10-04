package controllers

import (
	"github.com/galaxyerp/galaxyErp/internal/dto"
	"github.com/galaxyerp/galaxyErp/internal/models"
	"github.com/galaxyerp/galaxyErp/internal/services"
	"github.com/gin-gonic/gin"
)

// AccountingController 会计控制器
type AccountingController struct {
	accountService      services.AccountService
	journalEntryService services.JournalEntryService
	utils               *ControllerUtils
}

// NewAccountingController 创建会计控制器实例
func NewAccountingController(
	accountService services.AccountService,
	journalEntryService services.JournalEntryService,
) *AccountingController {
	return &AccountingController{
		accountService:      accountService,
		journalEntryService: journalEntryService,
		utils:               NewControllerUtils(),
	}
}

// CreateAccount 创建会计科目
// @Summary 创建会计科目
// @Description 创建新的会计科目
// @Tags 会计科目
// @Accept json
// @Produce json
// @Param account body dto.AccountCreateRequest true "会计科目信息"
// @Success 201 {object} dto.SuccessResponse{data=dto.AccountResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/accounting/accounts [post]
func (c *AccountingController) CreateAccount(ctx *gin.Context) {
	var req dto.AccountCreateRequest

	// 使用新的统一验证方法
	if !c.utils.BindAndValidateJSON(ctx, &req) {
		return
	}

	// 转换为模型
	account := models.Account{
		Code:        req.Code,
		Name:        req.Name,
		AccountType: req.Type,
		ParentID:    req.ParentID,
		Balance:     req.Balance,
		IsActive:    req.Status == "active",
	}

	if err := c.accountService.CreateAccount(ctx.Request.Context(), &account); err != nil {
		c.utils.RespondInternalError(ctx, "创建科目失败: "+err.Error())
		return
	}

	// 转换为响应 DTO
	response := dto.AccountResponse{
		ID:        account.ID,
		Code:      account.Code,
		Name:      account.Name,
		Type:      account.AccountType,
		ParentID:  account.ParentID,
		Balance:   account.Balance,
		Status:    map[bool]string{true: "active", false: "inactive"}[account.IsActive],
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}

	c.utils.RespondCreated(ctx, response)
}

// GetAccountList 获取会计科目列表
// @Summary 获取会计科目列表
// @Description 获取会计科目列表，支持按类型筛选和关键字搜索
// @Tags 会计科目
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param account_type query string false "科目类型"
// @Param keyword query string false "搜索关键字"
// @Success 200 {object} dto.PaginatedResponse{data=[]models.Account}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/accounting/accounts [get]
func (c *AccountingController) GetAccountList(ctx *gin.Context) {
	pagination := c.utils.ParsePaginationParams(ctx)
	accountType := ctx.Query("account_type")
	keyword := ctx.Query("keyword")

	var accounts []*models.Account
	var total int64
	var err error

	if keyword != "" {
		// 搜索科目
		accounts, total, err = c.accountService.SearchAccounts(ctx.Request.Context(), keyword, pagination.Page, pagination.PageSize)
	} else if accountType != "" {
		// 按类型筛选
		accounts, total, err = c.accountService.GetAccountsByType(ctx.Request.Context(), accountType, pagination.Page, pagination.PageSize)
	} else {
		// 获取所有科目
		accounts, total, err = c.accountService.ListAccounts(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	}

	if err != nil {
		c.utils.RespondInternalError(ctx, "获取科目列表失败: "+err.Error())
		return
	}

	// 转换为统一的分页响应格式
	pagination2 := c.utils.CreatePagination(pagination.Page, pagination.PageSize, total)
	c.utils.RespondPaginated(ctx, accounts, pagination2, "获取科目列表成功")
}

// GetAccount 获取会计科目详情
// @Summary 获取会计科目详情
// @Description 根据ID获取会计科目详情
// @Tags 会计科目
// @Accept json
// @Produce json
// @Param id path int true "科目ID"
// @Success 200 {object} dto.SuccessResponse{data=models.Account}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/accounting/accounts/{id} [get]
func (c *AccountingController) GetAccount(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	account, err := c.accountService.GetAccount(ctx.Request.Context(), id)
	if err != nil {
		if err.Error() == "科目不存在" {
			c.utils.RespondNotFound(ctx, "科目不存在")
		} else {
			c.utils.RespondInternalError(ctx, "获取科目失败: "+err.Error())
		}
		return
	}

	c.utils.RespondOK(ctx, account)
}

// UpdateAccount 更新会计科目
// @Summary 更新会计科目
// @Description 更新会计科目信息
// @Tags 会计科目
// @Accept json
// @Produce json
// @Param id path int true "科目ID"
// @Param account body models.Account true "会计科目信息"
// @Success 200 {object} dto.SuccessResponse{data=models.Account}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/accounting/accounts/{id} [put]
func (c *AccountingController) UpdateAccount(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	var account models.Account
	if !c.utils.BindAndValidateJSON(ctx, &account) {
		return
	}

	account.ID = id

	// 验证必填字段
	if account.Code == "" {
		c.utils.RespondBadRequest(ctx, "科目编码不能为空")
		return
	}
	if account.Name == "" {
		c.utils.RespondBadRequest(ctx, "科目名称不能为空")
		return
	}
	if account.AccountType == "" {
		c.utils.RespondBadRequest(ctx, "科目类型不能为空")
		return
	}

	if err := c.accountService.UpdateAccount(ctx.Request.Context(), &account); err != nil {
		if err.Error() == "科目不存在" {
			c.utils.RespondNotFound(ctx, "科目不存在")
		} else {
			c.utils.RespondInternalError(ctx, "更新科目失败: "+err.Error())
		}
		return
	}

	c.utils.RespondOK(ctx, account)
}

// DeleteAccount 删除会计科目
// @Summary 删除会计科目
// @Description 删除会计科目
// @Tags 会计科目
// @Accept json
// @Produce json
// @Param id path int true "科目ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/accounting/accounts/{id} [delete]
func (c *AccountingController) DeleteAccount(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	if err := c.accountService.DeleteAccount(ctx.Request.Context(), id); err != nil {
		if err.Error() == "科目不存在" {
			c.utils.RespondNotFound(ctx, "科目不存在")
		} else {
			c.utils.RespondInternalError(ctx, "删除科目失败: "+err.Error())
		}
		return
	}

	c.utils.RespondSuccess(ctx, "科目删除成功")
}

// GetAccountByCode 根据编码获取会计科目
// @Summary 根据编码获取会计科目
// @Description 根据科目编码获取会计科目详情
// @Tags 会计科目
// @Accept json
// @Produce json
// @Param code path string true "科目编码"
// @Success 200 {object} dto.SuccessResponse{data=models.Account}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/accounting/accounts/code/{code} [get]
func (c *AccountingController) GetAccountByCode(ctx *gin.Context) {
	code := ctx.Param("code")
	if code == "" {
		c.utils.RespondBadRequest(ctx, "科目编码不能为空")
		return
	}

	account, err := c.accountService.GetAccountByCode(ctx.Request.Context(), code)
	if err != nil {
		if err.Error() == "科目不存在" {
			c.utils.RespondNotFound(ctx, "科目不存在")
		} else {
			c.utils.RespondInternalError(ctx, "获取科目失败: "+err.Error())
		}
		return
	}

	c.utils.RespondOK(ctx, account)
}

// GetAccountChildren 获取子科目
// @Summary 获取子科目
// @Description 获取指定科目的所有子科目
// @Tags 会计科目
// @Accept json
// @Produce json
// @Param id path int true "父科目ID"
// @Success 200 {object} dto.SuccessResponse{data=[]models.Account}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/accounting/accounts/{id}/children [get]
func (c *AccountingController) GetAccountChildren(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	children, err := c.accountService.GetAccountChildren(ctx.Request.Context(), id)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取子科目失败: "+err.Error())
		return
	}

	c.utils.RespondOK(ctx, children)
}

// CreateJournalEntry 创建会计分录
// @Summary 创建会计分录
// @Description 创建新的会计分录
// @Tags 会计分录
// @Accept json
// @Produce json
// @Param entry body models.JournalEntry true "会计分录信息"
// @Success 201 {object} dto.SuccessResponse{data=models.JournalEntry}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/accounting/journal-entries [post]
func (c *AccountingController) CreateJournalEntry(ctx *gin.Context) {
	var req dto.JournalEntryCreateRequest
	if !c.utils.BindAndValidateJSON(ctx, &req) {
		return
	}

	// 验证借贷平衡
	var totalDebit, totalCredit float64
	for _, item := range req.Items {
		totalDebit += item.DebitAmount
		totalCredit += item.CreditAmount
	}

	if totalDebit != totalCredit {
		c.utils.RespondBadRequest(ctx, "借贷金额不平衡")
		return
	}

	entry, err := c.journalEntryService.CreateJournalEntryFromDTO(ctx.Request.Context(), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "创建分录失败: "+err.Error())
		return
	}

	c.utils.RespondCreated(ctx, entry)
}

// GetJournalEntryList 获取会计分录列表
// @Summary 获取会计分录列表
// @Description 获取会计分录列表，支持按日期范围筛选
// @Tags 会计分录
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param start_date query string false "开始日期"
// @Param end_date query string false "结束日期"
// @Success 200 {object} dto.PaginatedResponse{data=[]models.JournalEntry}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/accounting/journal-entries [get]
func (c *AccountingController) GetJournalEntryList(ctx *gin.Context) {
	pagination := c.utils.ParsePaginationParams(ctx)
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	var entries []*models.JournalEntry
	var total int64
	var err error

	if startDate != "" && endDate != "" {
		// 按日期范围筛选
		entries, total, err = c.journalEntryService.GetJournalEntriesByDateRange(ctx.Request.Context(), startDate, endDate, pagination.Page, pagination.PageSize)
	} else {
		// 获取所有分录
		entries, total, err = c.journalEntryService.ListJournalEntries(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	}

	if err != nil {
		c.utils.RespondInternalError(ctx, "获取分录列表失败: "+err.Error())
		return
	}

	// 转换为统一的分页响应格式
	pagination2 := c.utils.CreatePagination(pagination.Page, pagination.PageSize, total)
	c.utils.RespondPaginated(ctx, entries, pagination2, "获取分录列表成功")
}

// GetJournalEntry 获取会计分录详情
// @Summary 获取会计分录详情
// @Description 根据ID获取会计分录详情
// @Tags 会计分录
// @Accept json
// @Produce json
// @Param id path int true "分录ID"
// @Success 200 {object} dto.SuccessResponse{data=models.JournalEntry}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/accounting/journal-entries/{id} [get]
func (c *AccountingController) GetJournalEntry(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	entry, err := c.journalEntryService.GetJournalEntry(ctx.Request.Context(), id)
	if err != nil {
		if err.Error() == "分录不存在" {
			c.utils.RespondNotFound(ctx, "分录不存在")
		} else {
			c.utils.RespondInternalError(ctx, "获取分录失败: "+err.Error())
		}
		return
	}

	c.utils.RespondOK(ctx, entry)
}

// UpdateJournalEntry 更新会计分录
// @Summary 更新会计分录
// @Description 更新会计分录信息
// @Tags 会计分录
// @Accept json
// @Produce json
// @Param id path int true "分录ID"
// @Param entry body models.JournalEntry true "会计分录信息"
// @Success 200 {object} dto.SuccessResponse{data=models.JournalEntry}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/accounting/journal-entries/{id} [put]
func (c *AccountingController) UpdateJournalEntry(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	var entry models.JournalEntry
	if !c.utils.BindAndValidateJSON(ctx, &entry) {
		return
	}

	entry.ID = id

	// 验证必填字段
	if entry.TransactionID == 0 {
		c.utils.RespondBadRequest(ctx, "交易ID不能为空")
		return
	}
	if entry.AccountID == 0 {
		c.utils.RespondBadRequest(ctx, "科目ID不能为空")
		return
	}

	if err := c.journalEntryService.UpdateJournalEntry(ctx.Request.Context(), &entry); err != nil {
		if err.Error() == "分录不存在" {
			c.utils.RespondNotFound(ctx, "分录不存在")
		} else {
			c.utils.RespondInternalError(ctx, "更新分录失败: "+err.Error())
		}
		return
	}

	c.utils.RespondOK(ctx, entry)
}

// DeleteJournalEntry 删除会计分录
// @Summary 删除会计分录
// @Description 删除会计分录
// @Tags 会计分录
// @Accept json
// @Produce json
// @Param id path int true "分录ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/accounting/journal-entries/{id} [delete]
func (c *AccountingController) DeleteJournalEntry(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	if err := c.journalEntryService.DeleteJournalEntry(ctx.Request.Context(), id); err != nil {
		if err.Error() == "分录不存在" {
			c.utils.RespondNotFound(ctx, "分录不存在")
		} else {
			c.utils.RespondInternalError(ctx, "删除分录失败: "+err.Error())
		}
		return
	}

	c.utils.RespondSuccess(ctx, "分录删除成功")
}

// GetAccountTypes 获取科目类型列表
// @Summary 获取科目类型列表
// @Description 获取所有可用的科目类型
// @Tags 会计科目
// @Accept json
// @Produce json
// @Success 200 {object} dto.SuccessResponse{data=[]string}
// @Router /api/v1/accounting/account-types [get]
func (c *AccountingController) GetAccountTypes(ctx *gin.Context) {
	accountTypes := []string{
		"ASSET",     // 资产
		"LIABILITY", // 负债
		"EQUITY",    // 所有者权益
		"REVENUE",   // 收入
		"EXPENSE",   // 费用
	}

	c.utils.RespondOK(ctx, accountTypes)
}
