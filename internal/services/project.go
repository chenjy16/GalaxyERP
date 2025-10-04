package services

import (
	"context"
	"fmt"
	"time"

	"github.com/galaxyerp/galaxyErp/internal/common"
	"github.com/galaxyerp/galaxyErp/internal/dto"
	"github.com/galaxyerp/galaxyErp/internal/models"
	"github.com/galaxyerp/galaxyErp/internal/repositories"
)

// ProjectService 项目服务接口
type ProjectService interface {
	CRUDService[models.Project, dto.ProjectCreateRequest, dto.ProjectUpdateRequest, dto.ProjectResponse]
	CreateProject(ctx context.Context, req *dto.ProjectCreateRequest) (*dto.ProjectResponse, error)
	GetProject(ctx context.Context, id uint) (*dto.ProjectResponse, error)
	UpdateProject(ctx context.Context, id uint, req *dto.ProjectUpdateRequest) (*dto.ProjectResponse, error)
	DeleteProject(ctx context.Context, id uint) error
	ListProjects(ctx context.Context, filter *dto.ProjectFilter) (*dto.PaginatedResponse[dto.ProjectResponse], error)
	SearchProjects(ctx context.Context, keyword string) ([]*dto.ProjectResponse, error)
}

// TaskService 任务服务接口
type TaskService interface {
	CreateTask(ctx context.Context, req *dto.TaskCreateRequest) (*dto.TaskResponse, error)
	GetTask(ctx context.Context, id uint) (*dto.TaskResponse, error)
	UpdateTask(ctx context.Context, id uint, req *dto.TaskUpdateRequest) (*dto.TaskResponse, error)
	DeleteTask(ctx context.Context, id uint) error
	ListTasks(ctx context.Context, filter *dto.TaskFilter) (*dto.PaginatedResponse[dto.TaskResponse], error)
}

// MilestoneService 里程碑服务接口
type MilestoneService interface {
	CreateMilestone(ctx context.Context, req *dto.MilestoneCreateRequest) (*dto.MilestoneResponse, error)
	GetMilestone(ctx context.Context, id uint) (*dto.MilestoneResponse, error)
	UpdateMilestone(ctx context.Context, id uint, req *dto.MilestoneUpdateRequest) (*dto.MilestoneResponse, error)
	DeleteMilestone(ctx context.Context, id uint) error
	ListMilestones(ctx context.Context, projectID uint) ([]*dto.MilestoneResponse, error)
}

// TimeEntryService 工时记录服务接口
type TimeEntryService interface {
	CreateTimeEntry(ctx context.Context, req *dto.TimeEntryCreateRequest) (*dto.TimeEntryResponse, error)
	GetTimeEntry(ctx context.Context, id uint) (*dto.TimeEntryResponse, error)
	UpdateTimeEntry(ctx context.Context, id uint, req *dto.TimeEntryUpdateRequest) (*dto.TimeEntryResponse, error)
	DeleteTimeEntry(ctx context.Context, id uint) error
	ListTimeEntries(ctx context.Context, projectID uint, employeeID *uint) ([]*dto.TimeEntryResponse, error)
}

// ProjectServiceImpl 项目服务实现
type ProjectServiceImpl struct {
	*BaseService
	projectRepo repositories.ProjectRepository
}

// NewProjectService 创建项目服务
func NewProjectService(projectRepo repositories.ProjectRepository) ProjectService {
	config := &BaseServiceConfig{
		EnableAudit:      true,
		EnableValidation: true,
		EnableCache:      true,
		EnableMetrics:    true,
		CacheExpiry:      time.Hour,
	}
	
	return &ProjectServiceImpl{
		BaseService: NewBaseService(config),
		projectRepo: projectRepo,
	}
}

// Create 实现 CRUDService 接口的 Create 方法
func (s *ProjectServiceImpl) Create(ctx context.Context, req *dto.ProjectCreateRequest) (*dto.CreateResponse, error) {
	// 验证请求
	if err := s.ValidateRequest(ctx, "project", req); err != nil {
		return nil, err
	}

	project := &models.Project{
		ProjectNumber: req.ProjectCode,
		ProjectName:   req.ProjectName,
		Description:   req.Description,
		ManagerID:     req.ManagerID,
		CustomerID:    req.ClientID,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		Budget:        req.Budget,
		Status:        req.Status,
		Priority:      req.Priority,
	}

	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, err
	}

	// 清除相关缓存
	s.DeleteFromCache(ctx, "projects_list")

	return s.CreateCreateResponse(project.ID, s.convertToProjectResponse(project), "项目创建成功"), nil
}

// CreateProject 创建项目
func (s *ProjectServiceImpl) CreateProject(ctx context.Context, req *dto.ProjectCreateRequest) (*dto.ProjectResponse, error) {
	createResp, err := s.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return createResp.Data.(*dto.ProjectResponse), nil
}

// GetByID 实现 CRUDService 接口的 GetByID 方法
func (s *ProjectServiceImpl) GetByID(ctx context.Context, id uint) (*dto.ProjectResponse, error) {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("project_%d", id)
	if cached, found := s.GetFromCache(ctx, cacheKey); found {
		if project, ok := cached.(*dto.ProjectResponse); ok {
			return project, nil
		}
	}

	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := s.convertToProjectResponse(project)
	
	// 缓存结果
	s.SetToCache(ctx, cacheKey, response)

	return response, nil
}

// GetProject 获取项目
func (s *ProjectServiceImpl) GetProject(ctx context.Context, id uint) (*dto.ProjectResponse, error) {
	return s.GetByID(ctx, id)
}

// Update 实现 CRUDService 接口的 Update 方法
func (s *ProjectServiceImpl) Update(ctx context.Context, id uint, req *dto.ProjectUpdateRequest) (*dto.UpdateResponse, error) {
	// 验证请求
	if err := s.ValidateRequest(ctx, "project", req); err != nil {
		return nil, err
	}

	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.ProjectName != "" {
		project.ProjectName = req.ProjectName
	}
	if req.Description != "" {
		project.Description = req.Description
	}
	if req.ManagerID != nil && *req.ManagerID != 0 {
		project.ManagerID = *req.ManagerID
	}
	if req.ClientID != nil {
		customerID := *req.ClientID
		project.CustomerID = &customerID
	}
	if req.StartDate != nil {
		project.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		project.EndDate = *req.EndDate
	}
	if req.Budget != nil {
		project.Budget = *req.Budget
	}
	if req.Status != "" {
		project.Status = req.Status
	}
	if req.Priority != "" {
		project.Priority = req.Priority
	}

	if err := s.projectRepo.Update(ctx, project); err != nil {
		return nil, err
	}

	// 清除相关缓存
	cacheKey := fmt.Sprintf("project_%d", id)
	s.DeleteFromCache(ctx, cacheKey)
	s.DeleteFromCache(ctx, "projects_list")

	return s.CreateUpdateResponse(s.convertToProjectResponse(project), "项目更新成功"), nil
}

// Delete 实现 CRUDService 接口的 Delete 方法
func (s *ProjectServiceImpl) Delete(ctx context.Context, id uint) (*dto.DeleteResponse, error) {
	// 检查项目是否存在
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, fmt.Errorf("项目不存在")
	}

	if err := s.projectRepo.Delete(ctx, id); err != nil {
		return nil, err
	}

	// 清除相关缓存
	cacheKey := fmt.Sprintf("project_%d", id)
	s.DeleteFromCache(ctx, cacheKey)
	s.DeleteFromCache(ctx, "projects_list")

	return s.CreateDeleteResponse("项目删除成功"), nil
}

// UpdateProject 更新项目
func (s *ProjectServiceImpl) UpdateProject(ctx context.Context, id uint, req *dto.ProjectUpdateRequest) (*dto.ProjectResponse, error) {
	updateResp, err := s.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}
	return updateResp.Data.(*dto.ProjectResponse), nil
}

// DeleteProject 删除项目
func (s *ProjectServiceImpl) DeleteProject(ctx context.Context, id uint) error {
	_, err := s.Delete(ctx, id)
	return err
}

// List 实现 CRUDService 接口的 List 方法
func (s *ProjectServiceImpl) List(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.ProjectResponse], error) {
	projects, total, err := s.projectRepo.List(ctx, &common.QueryOptions{
		Pagination: req,
	})
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ProjectResponse, len(projects))
	for i, project := range projects {
		responses[i] = *s.convertToProjectResponse(project)
	}

	// 计算总页数
	limit := req.GetLimit()
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &dto.PaginatedResponse[dto.ProjectResponse]{
		Data:       responses,
		Total:      total,
		Page:       req.Page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// Search 实现 CRUDService 接口的 Search 方法
func (s *ProjectServiceImpl) Search(ctx context.Context, req *dto.SearchRequest) (*dto.PaginatedResponse[dto.ProjectResponse], error) {
	projects, total, err := s.projectRepo.List(ctx, &common.QueryOptions{
		Pagination: &req.PaginationRequest,
	})
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ProjectResponse, len(projects))
	for i, project := range projects {
		responses[i] = *s.convertToProjectResponse(project)
	}

	limit := req.GetLimit()
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &dto.PaginatedResponse[dto.ProjectResponse]{
		Data:       responses,
		Total:      total,
		Page:       req.Page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// ListProjects 列出项目
func (s *ProjectServiceImpl) ListProjects(ctx context.Context, filter *dto.ProjectFilter) (*dto.PaginatedResponse[dto.ProjectResponse], error) {
	// 转换为通用分页请求
	req := &dto.PaginationRequest{
		Page:     1,
		PageSize: 100,
	}
	if filter != nil {
		if filter.Page > 0 {
			req.Page = filter.Page
		}
		if filter.PageSize > 0 {
			req.PageSize = filter.PageSize
		}
	}
	return s.List(ctx, req)
}

// SearchProjects 搜索项目
func (s *ProjectServiceImpl) SearchProjects(ctx context.Context, keyword string) ([]*dto.ProjectResponse, error) {
	// 转换为通用搜索请求
	searchReq := &dto.SearchRequest{
		PaginationRequest: dto.PaginationRequest{
			Page:     1,
			PageSize: 100,
		},
		Keyword: keyword,
	}
	
	result, err := s.Search(ctx, searchReq)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.ProjectResponse, len(result.Data))
	for i, project := range result.Data {
		responses[i] = &project
	}

	return responses, nil
}

// convertToProjectResponse 转换为项目响应
func (s *ProjectServiceImpl) convertToProjectResponse(project *models.Project) *dto.ProjectResponse {
	var clientID uint
	if project.CustomerID != nil {
		clientID = *project.CustomerID
	}
	
	return &dto.ProjectResponse{
		ID:              project.ID,
		CreatedAt:       project.CreatedAt,
		UpdatedAt:       project.UpdatedAt,
		ProjectName:     project.ProjectName,
		ProjectCode:     project.ProjectNumber,
		Description:     project.Description,
		ManagerID:       project.ManagerID,
		ClientID:        clientID,
		StartDate:       project.StartDate,
		EndDate:         project.EndDate,
		ActualStartDate: project.ActualStartDate,
		ActualEndDate:   project.ActualEndDate,
		Budget:          project.Budget,
		ActualCost:      project.ActualCost,
		Progress:        project.Progress,
		Status:          project.Status,
		Priority:        project.Priority,
	}
}

// TaskServiceImpl 任务服务实现
type TaskServiceImpl struct {
	taskRepo repositories.TaskRepository
}

// NewTaskService 创建任务服务
func NewTaskService(taskRepo repositories.TaskRepository) TaskService {
	return &TaskServiceImpl{
		taskRepo: taskRepo,
	}
}

// CreateTask 创建任务
func (s *TaskServiceImpl) CreateTask(ctx context.Context, req *dto.TaskCreateRequest) (*dto.TaskResponse, error) {
	var startDate, endDate time.Time
	if req.StartDate != nil {
		startDate = *req.StartDate
	}
	if req.EndDate != nil {
		endDate = *req.EndDate
	}
	
	task := &models.Task{
		TaskNumber:     req.TaskNumber,
		TaskName:       req.TaskName,
		Description:    req.Description,
		ProjectID:      req.ProjectID,
		AssigneeID:     req.AssigneeID,
		ParentTaskID:   req.ParentTaskID,
		StartDate:      startDate,
		EndDate:        endDate,
		EstimatedHours: req.EstimatedHours,
		Priority:       req.Priority,
		Status:         "todo",
	}

	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, err
	}

	return s.convertToTaskResponse(task), nil
}

// GetTask 获取任务
func (s *TaskServiceImpl) GetTask(ctx context.Context, id uint) (*dto.TaskResponse, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.convertToTaskResponse(task), nil
}

// UpdateTask 更新任务
func (s *TaskServiceImpl) UpdateTask(ctx context.Context, id uint, req *dto.TaskUpdateRequest) (*dto.TaskResponse, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.TaskName != "" {
		task.TaskName = req.TaskName
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.AssigneeID != nil {
		task.AssigneeID = req.AssigneeID
	}
	if req.Status != "" {
		task.Status = req.Status
	}
	if req.Priority != "" {
		task.Priority = req.Priority
	}
	if req.StartDate != nil {
		task.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		task.EndDate = *req.EndDate
	}
	if req.EstimatedHours != nil {
		task.EstimatedHours = *req.EstimatedHours
	}
	if req.ActualHours != nil {
		task.ActualHours = *req.ActualHours
	}
	if req.Progress != nil {
		task.Progress = *req.Progress
	}

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}

	return s.convertToTaskResponse(task), nil
}

// DeleteTask 删除任务
func (s *TaskServiceImpl) DeleteTask(ctx context.Context, id uint) error {
	return s.taskRepo.Delete(ctx, id)
}

// ListTasks 列出任务
func (s *TaskServiceImpl) ListTasks(ctx context.Context, filter *dto.TaskFilter) (*dto.PaginatedResponse[dto.TaskResponse], error) {
	// 创建分页请求
	paginationReq := &dto.PaginationRequest{
		Page:     filter.Page,
		PageSize: filter.PageSize,
	}
	if paginationReq.Page <= 0 {
		paginationReq.Page = 1
	}
	if paginationReq.PageSize <= 0 {
		paginationReq.PageSize = 20
	}

	tasks, total, err := s.taskRepo.List(ctx, &common.QueryOptions{
		Pagination: paginationReq,
	})
	if err != nil {
		return nil, fmt.Errorf("获取任务列表失败: %w", err)
	}

	responses := make([]dto.TaskResponse, len(tasks))
	for i, task := range tasks {
		responses[i] = *s.convertToTaskResponse(task)
	}

	page := filter.Page
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &dto.PaginatedResponse[dto.TaskResponse]{
		Data:       responses,
		Total:      total,
		Page:       page,
		Limit:      pageSize,
		TotalPages: totalPages,
	}, nil
}

// convertToTaskResponse 转换为任务响应
func (s *TaskServiceImpl) convertToTaskResponse(task *models.Task) *dto.TaskResponse {
	return &dto.TaskResponse{
		ID:              task.ID,
		CreatedAt:       task.CreatedAt,
		UpdatedAt:       task.UpdatedAt,
		TaskNumber:      task.TaskNumber,
		TaskName:        task.TaskName,
		Description:     task.Description,
		ProjectID:       task.ProjectID,
		ParentTaskID:    task.ParentTaskID,
		AssigneeID:      task.AssigneeID,
		StartDate:       task.StartDate,
		EndDate:         task.EndDate,
		ActualStartDate: task.ActualStartDate,
		ActualEndDate:   task.ActualEndDate,
		EstimatedHours:  task.EstimatedHours,
		ActualHours:     task.ActualHours,
		Progress:        task.Progress,
		Priority:        task.Priority,
		Status:          task.Status,
	}
}

// MilestoneServiceImpl 里程碑服务实现
type MilestoneServiceImpl struct {
	milestoneRepo repositories.MilestoneRepository
}

// NewMilestoneService 创建里程碑服务
func NewMilestoneService(milestoneRepo repositories.MilestoneRepository) MilestoneService {
	return &MilestoneServiceImpl{
		milestoneRepo: milestoneRepo,
	}
}

// CreateMilestone 创建里程碑
func (s *MilestoneServiceImpl) CreateMilestone(ctx context.Context, req *dto.MilestoneCreateRequest) (*dto.MilestoneResponse, error) {
	milestone := &models.Milestone{
		ProjectID:     req.ProjectID,
		MilestoneName: req.MilestoneName,
		Description:   req.Description,
		DueDate:       req.DueDate,
		Status:        "pending",
	}

	if err := s.milestoneRepo.Create(ctx, milestone); err != nil {
		return nil, err
	}

	return s.convertToMilestoneResponse(milestone), nil
}

// GetMilestone 获取里程碑
func (s *MilestoneServiceImpl) GetMilestone(ctx context.Context, id uint) (*dto.MilestoneResponse, error) {
	milestone, err := s.milestoneRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.convertToMilestoneResponse(milestone), nil
}

// UpdateMilestone 更新里程碑
func (s *MilestoneServiceImpl) UpdateMilestone(ctx context.Context, id uint, req *dto.MilestoneUpdateRequest) (*dto.MilestoneResponse, error) {
	milestone, err := s.milestoneRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.MilestoneName != nil && *req.MilestoneName != "" {
		milestone.MilestoneName = *req.MilestoneName
	}
	if req.Description != nil && *req.Description != "" {
		milestone.Description = *req.Description
	}
	if req.DueDate != nil {
		milestone.DueDate = *req.DueDate
	}
	if req.Status != nil && *req.Status != "" {
		milestone.Status = *req.Status
		if *req.Status == "completed" && milestone.CompletedDate == nil {
			now := time.Now()
			milestone.CompletedDate = &now
		}
	}

	if err := s.milestoneRepo.Update(ctx, milestone); err != nil {
		return nil, err
	}

	return s.convertToMilestoneResponse(milestone), nil
}

// DeleteMilestone 删除里程碑
func (s *MilestoneServiceImpl) DeleteMilestone(ctx context.Context, id uint) error {
	return s.milestoneRepo.Delete(ctx, id)
}

// ListMilestones 列出里程碑
func (s *MilestoneServiceImpl) ListMilestones(ctx context.Context, projectID uint) ([]*dto.MilestoneResponse, error) {
	milestones, _, err := s.milestoneRepo.GetByProjectID(ctx, projectID, 0, 100)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.MilestoneResponse, len(milestones))
	for i, milestone := range milestones {
		responses[i] = s.convertToMilestoneResponse(milestone)
	}

	return responses, nil
}

// convertToMilestoneResponse 转换为里程碑响应
func (s *MilestoneServiceImpl) convertToMilestoneResponse(milestone *models.Milestone) *dto.MilestoneResponse {
	return &dto.MilestoneResponse{
		ID:            milestone.ID,
		CreatedAt:     milestone.CreatedAt,
		UpdatedAt:     milestone.UpdatedAt,
		ProjectID:     milestone.ProjectID,
		MilestoneName: milestone.MilestoneName,
		Name:          milestone.MilestoneName, // 使用 MilestoneName 作为 Name
		Description:   milestone.Description,
		DueDate:       milestone.DueDate,
		Status:        milestone.Status,
	}
}

// TimeEntryServiceImpl 工时记录服务实现
type TimeEntryServiceImpl struct {
	timeEntryRepo repositories.TimeEntryRepository
}

// NewTimeEntryService 创建工时记录服务
func NewTimeEntryService(timeEntryRepo repositories.TimeEntryRepository) TimeEntryService {
	return &TimeEntryServiceImpl{
		timeEntryRepo: timeEntryRepo,
	}
}

// CreateTimeEntry 创建工时记录
func (s *TimeEntryServiceImpl) CreateTimeEntry(ctx context.Context, req *dto.TimeEntryCreateRequest) (*dto.TimeEntryResponse, error) {
	// 计算工时
	hours := req.EndTime.Sub(req.StartTime).Hours()
	amount := hours * req.HourlyRate

	timeEntry := &models.TimeEntry{
		ProjectID:   req.ProjectID,
		TaskID:      req.TaskID,
		Date:        req.Date,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Hours:       hours,
		Description: req.Description,
		IsBillable:  req.IsBillable,
		HourlyRate:  req.HourlyRate,
		Amount:      amount,
		Status:      "draft",
	}

	if err := s.timeEntryRepo.Create(ctx, timeEntry); err != nil {
		return nil, err
	}

	return s.convertToTimeEntryResponse(timeEntry), nil
}

// GetTimeEntry 获取工时记录
func (s *TimeEntryServiceImpl) GetTimeEntry(ctx context.Context, id uint) (*dto.TimeEntryResponse, error) {
	timeEntry, err := s.timeEntryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.convertToTimeEntryResponse(timeEntry), nil
}

// UpdateTimeEntry 更新工时记录
func (s *TimeEntryServiceImpl) UpdateTimeEntry(ctx context.Context, id uint, req *dto.TimeEntryUpdateRequest) (*dto.TimeEntryResponse, error) {
	timeEntry, err := s.timeEntryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.TaskID != nil {
		timeEntry.TaskID = req.TaskID
	}
	if req.Date != nil {
		timeEntry.Date = *req.Date
	}
	if req.StartTime != nil {
		timeEntry.StartTime = *req.StartTime
	}
	if req.EndTime != nil {
		timeEntry.EndTime = *req.EndTime
	}
	if req.Description != nil && *req.Description != "" {
		timeEntry.Description = *req.Description
	}
	if req.IsBillable != nil {
		timeEntry.IsBillable = *req.IsBillable
	}
	if req.HourlyRate != nil {
		timeEntry.HourlyRate = *req.HourlyRate
	}
	if req.Status != nil && *req.Status != "" {
		timeEntry.Status = *req.Status
	}

	// 重新计算工时和金额
	timeEntry.Hours = timeEntry.EndTime.Sub(timeEntry.StartTime).Hours()
	timeEntry.Amount = timeEntry.Hours * timeEntry.HourlyRate

	if err := s.timeEntryRepo.Update(ctx, timeEntry); err != nil {
		return nil, err
	}

	return s.convertToTimeEntryResponse(timeEntry), nil
}

// DeleteTimeEntry 删除工时记录
func (s *TimeEntryServiceImpl) DeleteTimeEntry(ctx context.Context, id uint) error {
	return s.timeEntryRepo.Delete(ctx, id)
}

// ListTimeEntries 列出工时记录
func (s *TimeEntryServiceImpl) ListTimeEntries(ctx context.Context, projectID uint, employeeID *uint) ([]*dto.TimeEntryResponse, error) {
	req := &dto.PaginationRequest{
		Page:     1,
		PageSize: 100,
	}
	timeEntries, _, err := s.timeEntryRepo.List(ctx, &common.QueryOptions{
		Pagination: req,
	})

	if err != nil {
		return nil, err
	}

	responses := make([]*dto.TimeEntryResponse, len(timeEntries))
	for i, timeEntry := range timeEntries {
		responses[i] = s.convertToTimeEntryResponse(timeEntry)
	}

	return responses, nil
}

// convertToTimeEntryResponse 转换为工时记录响应
func (s *TimeEntryServiceImpl) convertToTimeEntryResponse(timeEntry *models.TimeEntry) *dto.TimeEntryResponse {
	return &dto.TimeEntryResponse{
		ID:          timeEntry.ID,
		CreatedAt:   timeEntry.CreatedAt,
		UpdatedAt:   timeEntry.UpdatedAt,
		EmployeeID:  timeEntry.EmployeeID,
		ProjectID:   timeEntry.ProjectID,
		TaskID:      timeEntry.TaskID,
		Date:        timeEntry.Date,
		StartTime:   timeEntry.StartTime,
		EndTime:     timeEntry.EndTime,
		Hours:       timeEntry.Hours,
		Description: timeEntry.Description,
		IsBillable:  timeEntry.IsBillable,
		HourlyRate:  timeEntry.HourlyRate,
		Amount:      timeEntry.Amount,
		Status:      timeEntry.Status,
	}
}
