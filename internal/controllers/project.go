package controllers

import (
	"strconv"

	"github.com/galaxyerp/galaxyErp/internal/dto"
	"github.com/galaxyerp/galaxyErp/internal/services"
	"github.com/gin-gonic/gin"
)

// ProjectController 项目管理控制器
type ProjectController struct {
	projectService   services.ProjectService
	taskService      services.TaskService
	milestoneService services.MilestoneService
	timeEntryService services.TimeEntryService
	utils            *ControllerUtils
}

// NewProjectController 创建项目控制器实例
func NewProjectController(
	projectService services.ProjectService,
	taskService services.TaskService,
	milestoneService services.MilestoneService,
	timeEntryService services.TimeEntryService,
) *ProjectController {
	return &ProjectController{
		projectService:   projectService,
		taskService:      taskService,
		milestoneService: milestoneService,
		timeEntryService: timeEntryService,
		utils:            NewControllerUtils(),
	}
}

// 项目相关接口

// CreateProject 创建项目
func (c *ProjectController) CreateProject(ctx *gin.Context) {
	var req dto.ProjectCreateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	project, err := c.projectService.CreateProject(ctx.Request.Context(), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "创建项目失败")
		return
	}

	c.utils.RespondCreated(ctx, project)
}

// GetProject 获取项目详情
func (c *ProjectController) GetProject(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	project, err := c.projectService.GetProject(ctx.Request.Context(), id)
	if err != nil {
		c.utils.RespondNotFound(ctx, "项目不存在")
		return
	}

	c.utils.RespondOK(ctx, project)
}

// UpdateProject 更新项目
func (c *ProjectController) UpdateProject(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	var req dto.ProjectUpdateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	project, err := c.projectService.UpdateProject(ctx.Request.Context(), id, &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "更新项目失败")
		return
	}

	c.utils.RespondOK(ctx, project)
}

// DeleteProject 删除项目
func (c *ProjectController) DeleteProject(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	err := c.projectService.DeleteProject(ctx.Request.Context(), id)
	if err != nil {
		c.utils.RespondInternalError(ctx, "删除项目失败")
		return
	}

	c.utils.RespondSuccess(ctx, "项目删除成功")
}

// ListProjects 获取项目列表
func (c *ProjectController) ListProjects(ctx *gin.Context) {
	var filter dto.ProjectFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		c.utils.RespondBadRequest(ctx, "查询参数错误")
		return
	}

	projects, err := c.projectService.ListProjects(ctx.Request.Context(), &filter)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取项目列表失败")
		return
	}

	c.utils.RespondOK(ctx, projects)
}

// 任务相关接口

// CreateTask 创建任务
func (c *ProjectController) CreateTask(ctx *gin.Context) {
	var req dto.TaskCreateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	task, err := c.taskService.CreateTask(ctx.Request.Context(), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "创建任务失败")
		return
	}

	c.utils.RespondCreated(ctx, task)
}

// GetTask 获取任务详情
func (c *ProjectController) GetTask(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	task, err := c.taskService.GetTask(ctx.Request.Context(), id)
	if err != nil {
		c.utils.RespondNotFound(ctx, "任务不存在")
		return
	}

	c.utils.RespondOK(ctx, task)
}

// UpdateTask 更新任务
func (c *ProjectController) UpdateTask(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	var req dto.TaskUpdateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	task, err := c.taskService.UpdateTask(ctx.Request.Context(), id, &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "更新任务失败")
		return
	}

	c.utils.RespondOK(ctx, task)
}

// DeleteTask 删除任务
func (c *ProjectController) DeleteTask(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	err := c.taskService.DeleteTask(ctx.Request.Context(), id)
	if err != nil {
		c.utils.RespondInternalError(ctx, "删除任务失败")
		return
	}

	c.utils.RespondSuccess(ctx, "任务删除成功")
}

// ListTasks 获取任务列表
func (c *ProjectController) ListTasks(ctx *gin.Context) {
	var filter dto.TaskFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		c.utils.RespondBadRequest(ctx, "查询参数错误")
		return
	}

	tasks, err := c.taskService.ListTasks(ctx.Request.Context(), &filter)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取任务列表失败")
		return
	}

	c.utils.RespondOK(ctx, tasks)
}

// 里程碑相关接口

// CreateMilestone 创建里程碑
func (c *ProjectController) CreateMilestone(ctx *gin.Context) {
	var req dto.MilestoneCreateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	milestone, err := c.milestoneService.CreateMilestone(ctx.Request.Context(), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "创建里程碑失败")
		return
	}

	c.utils.RespondCreated(ctx, milestone)
}

// GetMilestone 获取里程碑详情
func (c *ProjectController) GetMilestone(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	milestone, err := c.milestoneService.GetMilestone(ctx.Request.Context(), id)
	if err != nil {
		c.utils.RespondNotFound(ctx, "里程碑不存在")
		return
	}

	c.utils.RespondOK(ctx, milestone)
}

// UpdateMilestone 更新里程碑
func (c *ProjectController) UpdateMilestone(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	var req dto.MilestoneUpdateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	milestone, err := c.milestoneService.UpdateMilestone(ctx.Request.Context(), id, &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "更新里程碑失败")
		return
	}

	c.utils.RespondOK(ctx, milestone)
}

// DeleteMilestone 删除里程碑
func (c *ProjectController) DeleteMilestone(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	err := c.milestoneService.DeleteMilestone(ctx.Request.Context(), id)
	if err != nil {
		c.utils.RespondInternalError(ctx, "删除里程碑失败")
		return
	}

	c.utils.RespondSuccess(ctx, "里程碑删除成功")
}

// ListMilestones 获取项目里程碑列表
func (c *ProjectController) ListMilestones(ctx *gin.Context) {
	projectID, ok := c.utils.ParseIDParam(ctx, "project_id")
	if !ok {
		return
	}

	milestones, err := c.milestoneService.ListMilestones(ctx.Request.Context(), projectID)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取里程碑列表失败")
		return
	}

	c.utils.RespondOK(ctx, milestones)
}

// 工时记录相关接口

// CreateTimeEntry 创建工时记录
func (c *ProjectController) CreateTimeEntry(ctx *gin.Context) {
	var req dto.TimeEntryCreateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	timeEntry, err := c.timeEntryService.CreateTimeEntry(ctx.Request.Context(), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "创建工时记录失败")
		return
	}

	c.utils.RespondCreated(ctx, timeEntry)
}

// GetTimeEntry 获取工时记录详情
func (c *ProjectController) GetTimeEntry(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	timeEntry, err := c.timeEntryService.GetTimeEntry(ctx.Request.Context(), id)
	if err != nil {
		c.utils.RespondNotFound(ctx, "工时记录不存在")
		return
	}

	c.utils.RespondOK(ctx, timeEntry)
}

// UpdateTimeEntry 更新工时记录
func (c *ProjectController) UpdateTimeEntry(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	var req dto.TimeEntryUpdateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	timeEntry, err := c.timeEntryService.UpdateTimeEntry(ctx.Request.Context(), id, &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "更新工时记录失败")
		return
	}

	c.utils.RespondOK(ctx, timeEntry)
}

// DeleteTimeEntry 删除工时记录
func (c *ProjectController) DeleteTimeEntry(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	err := c.timeEntryService.DeleteTimeEntry(ctx.Request.Context(), id)
	if err != nil {
		c.utils.RespondInternalError(ctx, "删除工时记录失败")
		return
	}

	c.utils.RespondSuccess(ctx, "工时记录删除成功")
}

// ListTimeEntries 获取项目工时记录列表
func (c *ProjectController) ListTimeEntries(ctx *gin.Context) {
	projectID, ok := c.utils.ParseIDParam(ctx, "project_id")
	if !ok {
		return
	}

	// 可选的员工ID过滤
	var employeeID *uint
	if empIDStr := ctx.Query("employee_id"); empIDStr != "" {
		// 手动解析员工ID
		if empID, err := strconv.ParseUint(empIDStr, 10, 32); err == nil {
			id := uint(empID)
			employeeID = &id
		}
	}

	timeEntries, err := c.timeEntryService.ListTimeEntries(ctx.Request.Context(), projectID, employeeID)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取工时记录列表失败")
		return
	}

	c.utils.RespondOK(ctx, timeEntries)
}
