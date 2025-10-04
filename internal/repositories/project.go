package repositories

import (
	"context"
	"github.com/galaxyerp/galaxyErp/internal/models"
	"gorm.io/gorm"
)

// ProjectRepository 项目仓储接口
type ProjectRepository interface {
	BaseRepository[models.Project]
	GetByStatus(ctx context.Context, status string, offset, limit int) ([]*models.Project, int64, error)
}

// ProjectRepositoryImpl 项目仓储实现
type ProjectRepositoryImpl struct {
	BaseRepository[models.Project]
	db *gorm.DB
}

// NewProjectRepository 创建项目仓储实例
func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &ProjectRepositoryImpl{
		BaseRepository: NewBaseRepository[models.Project](db),
		db:             db,
	}
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
	BaseRepository[models.Task]
	GetByProjectID(ctx context.Context, projectID uint, offset, limit int) ([]*models.Task, int64, error)
	GetByStatus(ctx context.Context, status string, offset, limit int) ([]*models.Task, int64, error)
}

// TaskRepositoryImpl 任务仓储实现
type TaskRepositoryImpl struct {
	BaseRepository[models.Task]
	db *gorm.DB
}

// NewTaskRepository 创建任务仓储实例
func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &TaskRepositoryImpl{
		BaseRepository: NewBaseRepository[models.Task](db),
		db:             db,
	}
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
	BaseRepository[models.Milestone]
	GetByProjectID(ctx context.Context, projectID uint, offset, limit int) ([]*models.Milestone, int64, error)
}

// MilestoneRepositoryImpl 里程碑仓储实现
type MilestoneRepositoryImpl struct {
	BaseRepository[models.Milestone]
	db *gorm.DB
}

// NewMilestoneRepository 创建里程碑仓储实例
func NewMilestoneRepository(db *gorm.DB) MilestoneRepository {
	return &MilestoneRepositoryImpl{
		BaseRepository: NewBaseRepository[models.Milestone](db),
		db:             db,
	}
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
	BaseRepository[models.TimeEntry]
	GetByTaskID(ctx context.Context, taskID uint, offset, limit int) ([]*models.TimeEntry, int64, error)
	GetByUserID(ctx context.Context, userID uint, offset, limit int) ([]*models.TimeEntry, int64, error)
}

// TimeEntryRepositoryImpl 时间记录仓储实现
type TimeEntryRepositoryImpl struct {
	BaseRepository[models.TimeEntry]
	db *gorm.DB
}

// NewTimeEntryRepository 创建时间记录仓储实例
func NewTimeEntryRepository(db *gorm.DB) TimeEntryRepository {
	return &TimeEntryRepositoryImpl{
		BaseRepository: NewBaseRepository[models.TimeEntry](db),
		db:             db,
	}
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
