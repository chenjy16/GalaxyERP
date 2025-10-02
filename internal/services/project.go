package services

import (
	"context"
	"time"

	"github.com/galaxyerp/galaxyErp/internal/dto"
	"github.com/galaxyerp/galaxyErp/internal/models"
	"github.com/galaxyerp/galaxyErp/internal/repositories"
)

// ProjectService 项目服务接口
type ProjectService interface {
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
	projectRepo repositories.ProjectRepository
}

// NewProjectService 创建项目服务
func NewProjectService(projectRepo repositories.ProjectRepository) ProjectService {
	return &ProjectServiceImpl{
		projectRepo: projectRepo,
	}
}

// CreateProject 创建项目
func (s *ProjectServiceImpl) CreateProject(ctx context.Context, req *dto.ProjectCreateRequest) (*dto.ProjectResponse, error) {
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

	return s.convertToProjectResponse(project), nil
}

// GetProject 获取项目
func (s *ProjectServiceImpl) GetProject(ctx context.Context, id uint) (*dto.ProjectResponse, error) {
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.convertToProjectResponse(project), nil
}

// UpdateProject 更新项目
func (s *ProjectServiceImpl) UpdateProject(ctx context.Context, id uint, req *dto.ProjectUpdateRequest) (*dto.ProjectResponse, error) {
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
	if req.ManagerID != 0 {
		project.ManagerID = req.ManagerID
	}
	if req.ClientID != nil {
		project.CustomerID = req.ClientID
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

	return s.convertToProjectResponse(project), nil
}

// DeleteProject 删除项目
func (s *ProjectServiceImpl) DeleteProject(ctx context.Context, id uint) error {
	return s.projectRepo.Delete(ctx, id)
}

// ListProjects 列出项目
func (s *ProjectServiceImpl) ListProjects(ctx context.Context, filter *dto.ProjectFilter) (*dto.PaginatedResponse[dto.ProjectResponse], error) {
	offset := 0
	limit := 100
	if filter != nil {
		if filter.Page > 0 && filter.PageSize > 0 {
			offset = (filter.Page - 1) * filter.PageSize
			limit = filter.PageSize
		}
	}

	projects, total, err := s.projectRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ProjectResponse, len(projects))
	for i, project := range projects {
		responses[i] = *s.convertToProjectResponse(project)
	}

	page := 1
	pageSize := 100
	if filter != nil {
		page = filter.Page
		pageSize = filter.PageSize
	}

	return &dto.PaginatedResponse[dto.ProjectResponse]{
		Data:       responses,
		Total:      total,
		Page:       page,
		Limit:      pageSize,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}, nil
}

// SearchProjects 搜索项目
func (s *ProjectServiceImpl) SearchProjects(ctx context.Context, keyword string) ([]*dto.ProjectResponse, error) {
	// 暂时使用List方法，后续可以添加Search方法到repository
	projects, _, err := s.projectRepo.List(ctx, 0, 100)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.ProjectResponse, len(projects))
	for i, project := range projects {
		responses[i] = s.convertToProjectResponse(project)
	}

	return responses, nil
}

// convertToProjectResponse 转换为项目响应
func (s *ProjectServiceImpl) convertToProjectResponse(project *models.Project) *dto.ProjectResponse {
	return &dto.ProjectResponse{
		BaseModel: dto.BaseModel{
			ID:        project.ID,
			CreatedAt: project.CreatedAt,
			UpdatedAt: project.UpdatedAt,
		},
		ProjectName:     project.ProjectName,
		ProjectCode:     project.ProjectNumber,
		Description:     project.Description,
		ManagerID:       project.ManagerID,
		ClientID:        project.CustomerID,
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
	task := &models.Task{
		TaskNumber:     req.TaskNumber,
		TaskName:       req.TaskName,
		Description:    req.Description,
		ProjectID:      req.ProjectID,
		AssigneeID:     req.AssigneeID,
		ParentTaskID:   req.ParentTaskID,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
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
	offset := 0
	limit := 100
	if filter != nil {
		if filter.Page > 0 && filter.PageSize > 0 {
			offset = (filter.Page - 1) * filter.PageSize
			limit = filter.PageSize
		}
	}

	tasks, total, err := s.taskRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.TaskResponse, len(tasks))
	for i, task := range tasks {
		responses[i] = *s.convertToTaskResponse(task)
	}

	page := 1
	pageSize := 100
	if filter != nil {
		page = filter.Page
		pageSize = filter.PageSize
	}

	return &dto.PaginatedResponse[dto.TaskResponse]{
		Data:       responses,
		Total:      total,
		Page:       page,
		Limit:      pageSize,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}, nil
}

// convertToTaskResponse 转换为任务响应
func (s *TaskServiceImpl) convertToTaskResponse(task *models.Task) *dto.TaskResponse {
	return &dto.TaskResponse{
		BaseModel: dto.BaseModel{
			ID:        task.ID,
			CreatedAt: task.CreatedAt,
			UpdatedAt: task.UpdatedAt,
		},
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

	if req.MilestoneName != "" {
		milestone.MilestoneName = req.MilestoneName
	}
	if req.Description != "" {
		milestone.Description = req.Description
	}
	if req.DueDate != nil {
		milestone.DueDate = *req.DueDate
	}
	if req.Status != "" {
		milestone.Status = req.Status
		if req.Status == "completed" && milestone.CompletedDate == nil {
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
		BaseModel: dto.BaseModel{
			ID:        milestone.ID,
			CreatedAt: milestone.CreatedAt,
			UpdatedAt: milestone.UpdatedAt,
		},
		ProjectID:     milestone.ProjectID,
		MilestoneName: milestone.MilestoneName,
		Description:   milestone.Description,
		DueDate:       milestone.DueDate,
		CompletedDate: milestone.CompletedDate,
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
	if req.Description != "" {
		timeEntry.Description = req.Description
	}
	if req.IsBillable != nil {
		timeEntry.IsBillable = *req.IsBillable
	}
	if req.HourlyRate != nil {
		timeEntry.HourlyRate = *req.HourlyRate
	}
	if req.Status != "" {
		timeEntry.Status = req.Status
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
	var timeEntries []*models.TimeEntry
	var err error

	if employeeID != nil {
		timeEntries, _, err = s.timeEntryRepo.GetByUserID(ctx, *employeeID, 0, 100)
	} else {
		// 暂时使用List方法，后续可以添加GetByProjectID方法到repository
		timeEntries, _, err = s.timeEntryRepo.List(ctx, 0, 100)
	}

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
		BaseModel: dto.BaseModel{
			ID:        timeEntry.ID,
			CreatedAt: timeEntry.CreatedAt,
			UpdatedAt: timeEntry.UpdatedAt,
		},
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
