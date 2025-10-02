package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/galaxyerp/galaxyErp/internal/dto"
	"github.com/galaxyerp/galaxyErp/internal/services"
)

// HRController HR控制器
type HRController struct {
	employeeService   services.EmployeeService
	attendanceService services.AttendanceService
	payrollService    services.PayrollService
	leaveService      services.LeaveService
	utils            *ControllerUtils
}

// NewHRController 创建HR控制器实例
func NewHRController(
	employeeService services.EmployeeService,
	attendanceService services.AttendanceService,
	payrollService services.PayrollService,
	leaveService services.LeaveService,
) *HRController {
	return &HRController{
		employeeService:   employeeService,
		attendanceService: attendanceService,
		payrollService:    payrollService,
		leaveService:      leaveService,
		utils:            NewControllerUtils(),
	}
}

// ===== 员工管理 =====

// CreateEmployee 创建员工
// @Summary 创建员工
// @Description 创建新员工
// @Tags 员工管理
// @Accept json
// @Produce json
// @Param employee body dto.EmployeeCreateRequest true "员工信息"
// @Success 200 {object} dto.EmployeeResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/hr/employees [post]
func (c *HRController) CreateEmployee(ctx *gin.Context) {
	var req dto.EmployeeCreateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	response, err := c.employeeService.Create(ctx.Request.Context(), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "创建员工失败: "+err.Error())
		return
	}

	c.utils.RespondOK(ctx, response)
}

// ===== 请假管理 =====

// CreateLeave 创建请假申请
// @Summary 创建请假申请
// @Description 创建新的请假申请
// @Tags 请假管理
// @Accept json
// @Produce json
// @Param leave body dto.LeaveCreateRequest true "请假申请信息"
// @Success 200 {object} dto.LeaveResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/hr/leaves [post]
func (c *HRController) CreateLeave(ctx *gin.Context) {
	var req dto.LeaveCreateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	response, err := c.leaveService.Create(ctx.Request.Context(), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "创建请假申请失败: "+err.Error())
		return
	}

	c.utils.RespondOK(ctx, response)
}

// GetLeave 获取请假申请详情
// @Summary 获取请假申请详情
// @Description 根据ID获取请假申请详情
// @Tags 请假管理
// @Accept json
// @Produce json
// @Param id path int true "请假申请ID"
// @Success 200 {object} dto.LeaveResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/hr/leaves/{id} [get]
func (c *HRController) GetLeave(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		c.utils.RespondBadRequest(ctx, "无效的请假申请ID")
		return
	}

	response, err := c.leaveService.GetByID(ctx.Request.Context(), uint(id))
	if err != nil {
		c.utils.RespondNotFound(ctx, "请假申请不存在")
		return
	}

	c.utils.RespondOK(ctx, response)
}

// UpdateLeave 更新请假申请
// @Summary 更新请假申请
// @Description 更新请假申请信息
// @Tags 请假管理
// @Accept json
// @Produce json
// @Param id path int true "请假申请ID"
// @Param leave body dto.LeaveUpdateRequest true "请假申请信息"
// @Success 200 {object} dto.LeaveResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/hr/leaves/{id} [put]
func (c *HRController) UpdateLeave(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		c.utils.RespondBadRequest(ctx, "无效的请假申请ID")
		return
	}

	var req dto.LeaveUpdateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	response, err := c.leaveService.Update(ctx.Request.Context(), uint(id), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "更新请假申请失败: "+err.Error())
		return
	}

	c.utils.RespondOK(ctx, response)
}

// DeleteLeave 删除请假申请
// @Summary 删除请假申请
// @Description 删除请假申请
// @Tags 请假管理
// @Accept json
// @Produce json
// @Param id path int true "请假申请ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/hr/leaves/{id} [delete]
func (c *HRController) DeleteLeave(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		c.utils.RespondBadRequest(ctx, "无效的请假申请ID")
		return
	}

	err = c.leaveService.Delete(ctx.Request.Context(), uint(id))
	if err != nil {
		c.utils.RespondInternalError(ctx, "删除请假申请失败: "+err.Error())
		return
	}

	c.utils.RespondOK(ctx, gin.H{"message": "请假申请删除成功"})
}

// GetLeaveList 获取请假申请列表
// @Summary 获取请假申请列表
// @Description 获取请假申请列表，支持分页和筛选
// @Tags 请假管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Param employee_id query int false "员工ID"
// @Param leave_type query string false "请假类型"
// @Param status query string false "状态"
// @Param start_date query string false "开始日期"
// @Param end_date query string false "结束日期"
// @Success 200 {object} dto.PaginatedResponse{data=[]dto.LeaveListResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/hr/leaves [get]
func (c *HRController) GetLeaveList(ctx *gin.Context) {
	pagination := c.utils.ParsePaginationParams(ctx)
	
	// 构建筛选条件
	filter := &dto.LeaveFilter{
		PaginationRequest: *pagination,
	}
	
	if employeeIDStr := ctx.Query("employee_id"); employeeIDStr != "" {
		if employeeID, err := strconv.ParseUint(employeeIDStr, 10, 32); err == nil {
			empID := uint(employeeID)
			filter.EmployeeID = &empID
		}
	}
	
	if leaveType := ctx.Query("leave_type"); leaveType != "" {
		filter.LeaveType = leaveType
	}
	
	if status := ctx.Query("status"); status != "" {
		filter.Status = status
	}
	
	// 注意：这里简化处理，实际应该解析时间格式
	// 在实际应用中应该使用 time.Parse 来解析日期字符串

	response, err := c.leaveService.ListWithFilters(ctx.Request.Context(), filter)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取请假申请列表失败: "+err.Error())
		return
	}

	c.utils.RespondOK(ctx, response)
}

// GetEmployeeLeaves 获取员工的请假申请列表
// @Summary 获取员工的请假申请列表
// @Description 根据员工ID获取其请假申请列表
// @Tags 请假管理
// @Accept json
// @Produce json
// @Param id path int true "员工ID"
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Success 200 {object} dto.PaginatedResponse{data=[]dto.LeaveListResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/hr/employees/{id}/leaves [get]
func (c *HRController) GetEmployeeLeaves(ctx *gin.Context) {
	employeeID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		c.utils.RespondBadRequest(ctx, "无效的员工ID")
		return
	}

	pagination := c.utils.ParsePaginationParams(ctx)

	response, err := c.leaveService.GetByEmployeeID(ctx.Request.Context(), uint(employeeID), pagination)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取员工请假申请列表失败: "+err.Error())
		return
	}

	c.utils.RespondOK(ctx, response)
}

// ApproveLeave 审批请假申请
// @Summary 审批请假申请
// @Description 审批或拒绝请假申请
// @Tags 请假管理
// @Accept json
// @Produce json
// @Param id path int true "请假申请ID"
// @Param approval body dto.LeaveApprovalRequest true "审批信息"
// @Success 200 {object} dto.LeaveResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/hr/leaves/{id}/approve [post]
func (c *HRController) ApproveLeave(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		c.utils.RespondBadRequest(ctx, "无效的请假申请ID")
		return
	}

	var req dto.LeaveApprovalRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	// 这里需要获取当前用户ID作为审批人ID，暂时使用1作为示例
	// 在实际应用中应该从JWT token或session中获取当前用户ID
	approverID := uint(1)

	response, err := c.leaveService.Approve(ctx.Request.Context(), uint(id), approverID, &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "审批请假申请失败: "+err.Error())
		return
	}

	c.utils.RespondOK(ctx, response)
}

// CancelLeave 取消请假申请
// @Summary 取消请假申请
// @Description 取消请假申请
// @Tags 请假管理
// @Accept json
// @Produce json
// @Param id path int true "请假申请ID"
// @Success 200 {object} dto.LeaveResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/hr/leaves/{id}/cancel [post]
func (c *HRController) CancelLeave(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		c.utils.RespondBadRequest(ctx, "无效的请假申请ID")
		return
	}

	// 这里需要获取当前用户ID，暂时使用1作为示例
	// 在实际应用中应该从JWT token或session中获取当前用户ID
	employeeID := uint(1)

	response, err := c.leaveService.Cancel(ctx.Request.Context(), uint(id), employeeID)
	if err != nil {
		c.utils.RespondInternalError(ctx, "取消请假申请失败: "+err.Error())
		return
	}

	c.utils.RespondOK(ctx, response)
}

// GetPendingLeaves 获取待审批的请假申请列表
// @Summary 获取待审批的请假申请列表
// @Description 获取所有待审批的请假申请列表
// @Tags 请假管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Success 200 {object} dto.PaginatedResponse{data=[]dto.LeaveListResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/hr/leaves/pending [get]
func (c *HRController) GetPendingLeaves(ctx *gin.Context) {
	pagination := c.utils.ParsePaginationParams(ctx)

	response, err := c.leaveService.GetPendingApprovals(ctx.Request.Context(), pagination)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取待审批请假申请列表失败: "+err.Error())
		return
	}

	c.utils.RespondOK(ctx, response)
}

// GetEmployee 获取员工详情
// @Summary 获取员工详情
// @Description 根据ID获取员工详情
// @Tags 员工管理
// @Accept json
// @Produce json
// @Param id path int true "员工ID"
// @Success 200 {object} dto.EmployeeResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/hr/employees/{id} [get]
func (c *HRController) GetEmployee(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		c.utils.RespondBadRequest(ctx, "无效的员工ID")
		return
	}

	response, err := c.employeeService.GetByID(ctx.Request.Context(), uint(id))
	if err != nil {
		c.utils.RespondNotFound(ctx, "员工不存在")
		return
	}

	c.utils.RespondOK(ctx, response)
}

// UpdateEmployee 更新员工信息
// @Summary 更新员工信息
// @Description 更新员工信息
// @Tags 员工管理
// @Accept json
// @Produce json
// @Param id path int true "员工ID"
// @Param employee body dto.EmployeeUpdateRequest true "员工信息"
// @Success 200 {object} dto.EmployeeResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/hr/employees/{id} [put]
func (c *HRController) UpdateEmployee(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		c.utils.RespondBadRequest(ctx, "无效的员工ID")
		return
	}

	var req dto.EmployeeUpdateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	response, err := c.employeeService.Update(ctx.Request.Context(), uint(id), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "更新员工失败: "+err.Error())
		return
	}

	c.utils.RespondOK(ctx, response)
}

// DeleteEmployee 删除员工
// @Summary 删除员工
// @Description 删除员工
// @Tags 员工管理
// @Accept json
// @Produce json
// @Param id path int true "员工ID"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/hr/employees/{id} [delete]
func (c *HRController) DeleteEmployee(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		c.utils.RespondBadRequest(ctx, "无效的员工ID")
		return
	}

	err = c.employeeService.Delete(ctx.Request.Context(), uint(id))
	if err != nil {
		c.utils.RespondInternalError(ctx, "删除员工失败: "+err.Error())
		return
	}

	c.utils.RespondOK(ctx, gin.H{"message": "员工删除成功"})
}

// GetEmployees 获取员工列表
// @Summary 获取员工列表
// @Description 获取员工列表
// @Tags 员工管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} dto.EmployeeListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/hr/employees [get]
func (c *HRController) GetEmployees(ctx *gin.Context) {
	pagination := c.utils.ParsePaginationParams(ctx)

	response, err := c.employeeService.List(ctx.Request.Context(), pagination)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取员工列表失败: "+err.Error())
		return
	}

	c.utils.RespondOK(ctx, response)
}

// SearchEmployees 搜索员工
// @Summary 搜索员工
// @Description 搜索员工
// @Tags 员工管理
// @Accept json
// @Produce json
// @Param search body dto.EmployeeSearchRequest true "搜索条件"
// @Success 200 {object} dto.EmployeeListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/hr/employees/search [post]
func (c *HRController) SearchEmployees(ctx *gin.Context) {
	var req dto.EmployeeSearchRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	response, err := c.employeeService.Search(ctx.Request.Context(), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "搜索员工失败: "+err.Error())
		return
	}

	c.utils.RespondOK(ctx, response)
}

// ===== 考勤管理 =====

// CreateAttendance 创建考勤记录
// @Summary 创建考勤记录
// @Description 创建新的考勤记录
// @Tags 考勤管理
// @Accept json
// @Produce json
// @Param attendance body dto.AttendanceCreateRequest true "考勤信息"
// @Success 200 {object} dto.AttendanceResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/hr/attendance [post]
func (c *HRController) CreateAttendance(ctx *gin.Context) {
	var req dto.AttendanceCreateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	response, err := c.attendanceService.Create(ctx.Request.Context(), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "创建考勤记录失败: "+err.Error())
		return
	}

	c.utils.RespondOK(ctx, response)
}

// GetAttendance 获取考勤记录详情
// @Summary 获取考勤记录详情
// @Description 根据ID获取考勤记录详情
// @Tags 考勤管理
// @Accept json
// @Produce json
// @Param id path int true "考勤记录ID"
// @Success 200 {object} dto.AttendanceResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/hr/attendance/{id} [get]
func (c *HRController) GetAttendance(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		c.utils.RespondBadRequest(ctx, "无效的考勤记录ID")
		return
	}

	response, err := c.attendanceService.GetByID(ctx.Request.Context(), uint(id))
	if err != nil {
		c.utils.RespondNotFound(ctx, "考勤记录不存在")
		return
	}

	c.utils.RespondOK(ctx, response)
}

// UpdateAttendance 更新考勤记录
// @Summary 更新考勤记录
// @Description 更新考勤记录
// @Tags 考勤管理
// @Accept json
// @Produce json
// @Param id path int true "考勤记录ID"
// @Param attendance body dto.AttendanceUpdateRequest true "考勤信息"
// @Success 200 {object} dto.AttendanceResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/hr/attendance/{id} [put]
func (c *HRController) UpdateAttendance(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		c.utils.RespondBadRequest(ctx, "无效的考勤记录ID")
		return
	}

	var req dto.AttendanceUpdateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	response, err := c.attendanceService.Update(ctx.Request.Context(), uint(id), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "更新考勤记录失败: "+err.Error())
		return
	}

	c.utils.RespondOK(ctx, response)
}

// DeleteAttendance 删除考勤记录
// @Summary 删除考勤记录
// @Description 删除考勤记录
// @Tags 考勤管理
// @Accept json
// @Produce json
// @Param id path int true "考勤记录ID"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/hr/attendance/{id} [delete]
func (c *HRController) DeleteAttendance(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		c.utils.RespondBadRequest(ctx, "无效的考勤记录ID")
		return
	}

	err = c.attendanceService.Delete(ctx.Request.Context(), uint(id))
	if err != nil {
		c.utils.RespondInternalError(ctx, "删除考勤记录失败: "+err.Error())
		return
	}

	c.utils.RespondOK(ctx, gin.H{"message": "考勤记录删除成功"})
}

// GetAttendanceList 获取考勤记录列表
// @Summary 获取考勤记录列表
// @Description 获取考勤记录列表
// @Tags 考勤管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} dto.PaginatedResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/hr/attendance [get]
func (c *HRController) GetAttendanceList(ctx *gin.Context) {
	pagination := c.utils.ParsePaginationParams(ctx)

	response, err := c.attendanceService.List(ctx.Request.Context(), pagination)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取考勤记录列表失败: "+err.Error())
		return
	}

	c.utils.RespondOK(ctx, response)
}

// ===== 薪资管理 =====

// CreatePayroll 创建薪资记录
// @Summary 创建薪资记录
// @Description 创建新的薪资记录
// @Tags 薪资管理
// @Accept json
// @Produce json
// @Param payroll body dto.PayrollCreateRequest true "薪资信息"
// @Success 200 {object} dto.PayrollResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/hr/payroll [post]
func (c *HRController) CreatePayroll(ctx *gin.Context) {
	var req dto.PayrollCreateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	response, err := c.payrollService.Create(ctx.Request.Context(), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "创建薪资记录失败: "+err.Error())
		return
	}

	c.utils.RespondOK(ctx, response)
}

// GetPayroll 获取薪资记录详情
// @Summary 获取薪资记录详情
// @Description 根据ID获取薪资记录详情
// @Tags 薪资管理
// @Accept json
// @Produce json
// @Param id path int true "薪资记录ID"
// @Success 200 {object} dto.PayrollResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/hr/payroll/{id} [get]
func (c *HRController) GetPayroll(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		c.utils.RespondBadRequest(ctx, "无效的薪资记录ID")
		return
	}

	response, err := c.payrollService.GetByID(ctx.Request.Context(), uint(id))
	if err != nil {
		c.utils.RespondNotFound(ctx, "薪资记录不存在")
		return
	}

	c.utils.RespondOK(ctx, response)
}

// UpdatePayroll 更新薪资记录
// @Summary 更新薪资记录
// @Description 更新薪资记录
// @Tags 薪资管理
// @Accept json
// @Produce json
// @Param id path int true "薪资记录ID"
// @Param payroll body dto.PayrollUpdateRequest true "薪资信息"
// @Success 200 {object} dto.PayrollResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/hr/payroll/{id} [put]
func (c *HRController) UpdatePayroll(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		c.utils.RespondBadRequest(ctx, "无效的薪资记录ID")
		return
	}

	var req dto.PayrollUpdateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	response, err := c.payrollService.Update(ctx.Request.Context(), uint(id), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "更新薪资记录失败: "+err.Error())
		return
	}

	c.utils.RespondOK(ctx, response)
}

// DeletePayroll 删除薪资记录
// @Summary 删除薪资记录
// @Description 删除薪资记录
// @Tags 薪资管理
// @Accept json
// @Produce json
// @Param id path int true "薪资记录ID"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/hr/payroll/{id} [delete]
func (c *HRController) DeletePayroll(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		c.utils.RespondBadRequest(ctx, "无效的薪资记录ID")
		return
	}

	err = c.payrollService.Delete(ctx.Request.Context(), uint(id))
	if err != nil {
		c.utils.RespondInternalError(ctx, "删除薪资记录失败: "+err.Error())
		return
	}

	c.utils.RespondOK(ctx, gin.H{"message": "薪资记录删除成功"})
}

// GetPayrollList 获取薪资记录列表
// @Summary 获取薪资记录列表
// @Description 获取薪资记录列表
// @Tags 薪资管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} dto.PaginatedResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/hr/payroll [get]
func (c *HRController) GetPayrollList(ctx *gin.Context) {
	pagination := c.utils.ParsePaginationParams(ctx)

	response, err := c.payrollService.List(ctx.Request.Context(), pagination)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取薪资记录列表失败: "+err.Error())
		return
	}

	c.utils.RespondOK(ctx, response)
}