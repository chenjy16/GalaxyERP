package repositories

import (
	"context"
	"errors"
	"github.com/galaxyerp/galaxyErp/internal/models"
	"gorm.io/gorm"
)

// ProjectRepository 项目仓储接口
type ProjectRepository interface {
	Create(ctx context.Context, project *models.Project) error
	GetByID(ctx context.Context, id uint) (*models.Project, error)
	Update(ctx context.Context, project *models.Project) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*models.Project, int64, error)
	GetByStatus(ctx context.Context, status string, offset, limit int) ([]*models.Project, int64, error)
}

// ProjectRepositoryImpl 项目仓储实现
type ProjectRepositoryImpl struct {
	db *gorm.DB
}

// NewProjectRepository 创建项目仓储实例
func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &ProjectRepositoryImpl{
		db: db,
	}
}

// Create 创建项目
func (r *ProjectRepositoryImpl) Create(ctx context.Context, project *models.Project) error {
	return r.db.WithContext(ctx).Create(project).Error
}

// GetByID 根据ID获取项目
func (r *ProjectRepositoryImpl) GetByID(ctx context.Context, id uint) (*models.Project, error) {
	var project models.Project
	err := r.db.WithContext(ctx).Preload("Tasks").Preload("Milestones").First(&project, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}

// Update 更新项目
func (r *ProjectRepositoryImpl) Update(ctx context.Context, project *models.Project) error {
	return r.db.WithContext(ctx).Save(project).Error
}

// Delete 删除项目
func (r *ProjectRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Project{}, id).Error
}

// List 获取项目列表
func (r *ProjectRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.Project, int64, error) {
	var projects []*models.Project
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Project{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&projects).Error
	if err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

// GetByStatus 根据状态获取项目
func (r *ProjectRepositoryImpl) GetByStatus(ctx context.Context, status string, offset, limit int) ([]*models.Project, int64, error) {
	var projects []*models.Project
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Project{}).
		Where("status = ?", status).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Offset(offset).Limit(limit).Find(&projects).Error
	if err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

// TaskRepository 任务仓储接口
type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	GetByID(ctx context.Context, id uint) (*models.Task, error)
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*models.Task, int64, error)
	GetByProjectID(ctx context.Context, projectID uint, offset, limit int) ([]*models.Task, int64, error)
	GetByStatus(ctx context.Context, status string, offset, limit int) ([]*models.Task, int64, error)
}

// TaskRepositoryImpl 任务仓储实现
type TaskRepositoryImpl struct {
	db *gorm.DB
}

// NewTaskRepository 创建任务仓储实例
func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &TaskRepositoryImpl{
		db: db,
	}
}

// Create 创建任务
func (r *TaskRepositoryImpl) Create(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

// GetByID 根据ID获取任务
func (r *TaskRepositoryImpl) GetByID(ctx context.Context, id uint) (*models.Task, error) {
	var task models.Task
	err := r.db.WithContext(ctx).Preload("Project").Preload("TimeEntries").First(&task, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

// Update 更新任务
func (r *TaskRepositoryImpl) Update(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).Save(task).Error
}

// Delete 删除任务
func (r *TaskRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Task{}, id).Error
}

// List 获取任务列表
func (r *TaskRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.Task, int64, error) {
	var tasks []*models.Task
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Task{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Project").Offset(offset).Limit(limit).Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// GetByProjectID 根据项目ID获取任务
func (r *TaskRepositoryImpl) GetByProjectID(ctx context.Context, projectID uint, offset, limit int) ([]*models.Task, int64, error) {
	var tasks []*models.Task
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Task{}).
		Where("project_id = ?", projectID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Project").
		Where("project_id = ?", projectID).
		Offset(offset).Limit(limit).Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// GetByStatus 根据状态获取任务
func (r *TaskRepositoryImpl) GetByStatus(ctx context.Context, status string, offset, limit int) ([]*models.Task, int64, error) {
	var tasks []*models.Task
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Task{}).
		Where("status = ?", status).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Project").
		Where("status = ?", status).
		Offset(offset).Limit(limit).Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// MilestoneRepository 里程碑仓储接口
type MilestoneRepository interface {
	Create(ctx context.Context, milestone *models.Milestone) error
	GetByID(ctx context.Context, id uint) (*models.Milestone, error)
	Update(ctx context.Context, milestone *models.Milestone) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*models.Milestone, int64, error)
	GetByProjectID(ctx context.Context, projectID uint, offset, limit int) ([]*models.Milestone, int64, error)
}

// MilestoneRepositoryImpl 里程碑仓储实现
type MilestoneRepositoryImpl struct {
	db *gorm.DB
}

// NewMilestoneRepository 创建里程碑仓储实例
func NewMilestoneRepository(db *gorm.DB) MilestoneRepository {
	return &MilestoneRepositoryImpl{
		db: db,
	}
}

// Create 创建里程碑
func (r *MilestoneRepositoryImpl) Create(ctx context.Context, milestone *models.Milestone) error {
	return r.db.WithContext(ctx).Create(milestone).Error
}

// GetByID 根据ID获取里程碑
func (r *MilestoneRepositoryImpl) GetByID(ctx context.Context, id uint) (*models.Milestone, error) {
	var milestone models.Milestone
	err := r.db.WithContext(ctx).Preload("Project").First(&milestone, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &milestone, nil
}

// Update 更新里程碑
func (r *MilestoneRepositoryImpl) Update(ctx context.Context, milestone *models.Milestone) error {
	return r.db.WithContext(ctx).Save(milestone).Error
}

// Delete 删除里程碑
func (r *MilestoneRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Milestone{}, id).Error
}

// List 获取里程碑列表
func (r *MilestoneRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.Milestone, int64, error) {
	var milestones []*models.Milestone
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Milestone{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Project").Offset(offset).Limit(limit).Find(&milestones).Error
	if err != nil {
		return nil, 0, err
	}

	return milestones, total, nil
}

// GetByProjectID 根据项目ID获取里程碑
func (r *MilestoneRepositoryImpl) GetByProjectID(ctx context.Context, projectID uint, offset, limit int) ([]*models.Milestone, int64, error) {
	var milestones []*models.Milestone
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Milestone{}).
		Where("project_id = ?", projectID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Project").
		Where("project_id = ?", projectID).
		Offset(offset).Limit(limit).Find(&milestones).Error
	if err != nil {
		return nil, 0, err
	}

	return milestones, total, nil
}

// TimeEntryRepository 时间记录仓储接口
type TimeEntryRepository interface {
	Create(ctx context.Context, timeEntry *models.TimeEntry) error
	GetByID(ctx context.Context, id uint) (*models.TimeEntry, error)
	Update(ctx context.Context, timeEntry *models.TimeEntry) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*models.TimeEntry, int64, error)
	GetByTaskID(ctx context.Context, taskID uint, offset, limit int) ([]*models.TimeEntry, int64, error)
	GetByUserID(ctx context.Context, userID uint, offset, limit int) ([]*models.TimeEntry, int64, error)
}

// TimeEntryRepositoryImpl 时间记录仓储实现
type TimeEntryRepositoryImpl struct {
	db *gorm.DB
}

// NewTimeEntryRepository 创建时间记录仓储实例
func NewTimeEntryRepository(db *gorm.DB) TimeEntryRepository {
	return &TimeEntryRepositoryImpl{
		db: db,
	}
}

// Create 创建时间记录
func (r *TimeEntryRepositoryImpl) Create(ctx context.Context, timeEntry *models.TimeEntry) error {
	return r.db.WithContext(ctx).Create(timeEntry).Error
}

// GetByID 根据ID获取时间记录
func (r *TimeEntryRepositoryImpl) GetByID(ctx context.Context, id uint) (*models.TimeEntry, error) {
	var timeEntry models.TimeEntry
	err := r.db.WithContext(ctx).Preload("Task").Preload("User").First(&timeEntry, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &timeEntry, nil
}

// Update 更新时间记录
func (r *TimeEntryRepositoryImpl) Update(ctx context.Context, timeEntry *models.TimeEntry) error {
	return r.db.WithContext(ctx).Save(timeEntry).Error
}

// Delete 删除时间记录
func (r *TimeEntryRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.TimeEntry{}, id).Error
}

// List 获取时间记录列表
func (r *TimeEntryRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.TimeEntry, int64, error) {
	var timeEntries []*models.TimeEntry
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.TimeEntry{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Task").Preload("User").Offset(offset).Limit(limit).Find(&timeEntries).Error
	if err != nil {
		return nil, 0, err
	}

	return timeEntries, total, nil
}

// GetByTaskID 根据任务ID获取时间记录
func (r *TimeEntryRepositoryImpl) GetByTaskID(ctx context.Context, taskID uint, offset, limit int) ([]*models.TimeEntry, int64, error) {
	var timeEntries []*models.TimeEntry
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.TimeEntry{}).
		Where("task_id = ?", taskID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Task").Preload("User").
		Where("task_id = ?", taskID).
		Offset(offset).Limit(limit).Find(&timeEntries).Error
	if err != nil {
		return nil, 0, err
	}

	return timeEntries, total, nil
}

// GetByUserID 根据用户ID获取时间记录
func (r *TimeEntryRepositoryImpl) GetByUserID(ctx context.Context, userID uint, offset, limit int) ([]*models.TimeEntry, int64, error) {
	var timeEntries []*models.TimeEntry
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.TimeEntry{}).
		Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Task").Preload("User").
		Where("user_id = ?", userID).
		Offset(offset).Limit(limit).Find(&timeEntries).Error
	if err != nil {
		return nil, 0, err
	}

	return timeEntries, total, nil
}
