package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/galaxyerp/galaxyErp/internal/dto"
	"github.com/galaxyerp/galaxyErp/internal/models"
)

// DeliveryNoteRepository 交付单仓储接口
type DeliveryNoteRepository interface {
	BaseRepository[models.DeliveryNote]
	GetByDeliveryNumber(ctx context.Context, deliveryNumber string) (*models.DeliveryNote, error)
	GetBySalesOrderID(ctx context.Context, salesOrderID uint) ([]*models.DeliveryNote, error)
	GetByCustomerID(ctx context.Context, customerID uint) ([]*models.DeliveryNote, error)
	UpdateStatus(ctx context.Context, id uint, status string) error
	GetStatistics(ctx *gin.Context) (*dto.DeliveryNoteStatisticsResponse, error)
	GetDeliveryTrend(days int) ([]dto.DeliveryTrendData, error)
	CheckDeliveryNumberExists(deliveryNumber string, excludeID ...uint) (bool, error)
	GenerateDeliveryNumber() (string, error)
	ListWithFilters(req *dto.DeliveryNoteListRequest) ([]*models.DeliveryNote, int64, error)
}

// DeliveryNoteRepositoryImpl 交付单仓储实现
type DeliveryNoteRepositoryImpl struct {
	BaseRepository[models.DeliveryNote]
	db *gorm.DB
}

// NewDeliveryNoteRepository 创建交付单仓储实例
func NewDeliveryNoteRepository(db *gorm.DB) DeliveryNoteRepository {
	return &DeliveryNoteRepositoryImpl{
		BaseRepository: NewBaseRepository[models.DeliveryNote](db),
		db:             db,
	}
}

// GetByDeliveryNumber 根据发货单号获取发货单
func (r *DeliveryNoteRepositoryImpl) GetByDeliveryNumber(ctx context.Context, deliveryNumber string) (*models.DeliveryNote, error) {
	var deliveryNote models.DeliveryNote
	err := r.db.WithContext(ctx).Preload("Customer").
		Preload("SalesOrder").
		Preload("Items").
		Preload("Items.Product").
		Where("delivery_number = ?", deliveryNumber).
		First(&deliveryNote).Error
	if err != nil {
		return nil, err
	}
	return &deliveryNote, nil
}

// ListWithFilters 获取发货单列表（带过滤条件）
func (r *DeliveryNoteRepositoryImpl) ListWithFilters(req *dto.DeliveryNoteListRequest) ([]*models.DeliveryNote, int64, error) {
	var deliveryNotes []*models.DeliveryNote
	var total int64

	query := r.db.Model(&models.DeliveryNote{})

	// 应用过滤条件
	if req.CustomerID != nil {
		query = query.Where("customer_id = ?", *req.CustomerID)
	}
	if req.SalesOrderID != nil {
		query = query.Where("sales_order_id = ?", *req.SalesOrderID)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Search != "" {
		query = query.Where("delivery_number LIKE ? OR notes LIKE ?", "%"+req.Search+"%", "%"+req.Search+"%")
	}
	if req.DateFrom != "" {
		if startDate, err := time.Parse("2006-01-02", req.DateFrom); err == nil {
			query = query.Where("date >= ?", startDate)
		}
	}
	if req.DateTo != "" {
		if endDate, err := time.Parse("2006-01-02", req.DateTo); err == nil {
			query = query.Where("date <= ?", endDate)
		}
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := query.Preload("Customer").
		Preload("SalesOrder").
		Preload("Items").
		Preload("Items.Product").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Order("created_at DESC").
		Find(&deliveryNotes).Error

	if err != nil {
		return nil, 0, err
	}

	return deliveryNotes, total, nil
}

// GetBySalesOrderID 根据销售订单ID获取发货单
func (r *DeliveryNoteRepositoryImpl) GetBySalesOrderID(ctx context.Context, salesOrderID uint) ([]*models.DeliveryNote, error) {
	var deliveryNotes []*models.DeliveryNote
	err := r.db.WithContext(ctx).Preload("Customer").
		Preload("Items").
		Preload("Items.Product").
		Where("sales_order_id = ?", salesOrderID).
		Find(&deliveryNotes).Error
	return deliveryNotes, err
}

// GetByCustomerID 根据客户ID获取发货单
func (r *DeliveryNoteRepositoryImpl) GetByCustomerID(ctx context.Context, customerID uint) ([]*models.DeliveryNote, error) {
	var deliveryNotes []*models.DeliveryNote
	err := r.db.WithContext(ctx).Preload("SalesOrder").
		Preload("Items").
		Preload("Items.Product").
		Where("customer_id = ?", customerID).
		Find(&deliveryNotes).Error
	return deliveryNotes, err
}

// UpdateStatus 更新发货单状态
func (r *DeliveryNoteRepositoryImpl) UpdateStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).Model(&models.DeliveryNote{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// GetStatistics 获取发货单统计信息
func (r *DeliveryNoteRepositoryImpl) GetStatistics(ctx *gin.Context) (*dto.DeliveryNoteStatisticsResponse, error) {
	var stats dto.DeliveryNoteStatisticsResponse

	// 总发货单数量
	if err := r.db.Model(&models.DeliveryNote{}).Count(&stats.TotalDeliveries).Error; err != nil {
		return nil, err
	}

	// 待发货数量
	if err := r.db.Model(&models.DeliveryNote{}).
		Where("status = ?", "pending").
		Count(&stats.PendingDeliveries).Error; err != nil {
		return nil, err
	}

	// 已发货数量
	if err := r.db.Model(&models.DeliveryNote{}).
		Where("status = ?", "delivered").
		Count(&stats.CompletedDeliveries).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetDeliveryTrend 获取发货趋势数据
func (r *DeliveryNoteRepositoryImpl) GetDeliveryTrend(days int) ([]dto.DeliveryTrendData, error) {
	var trendData []dto.DeliveryTrendData

	query := `
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as count
		FROM delivery_notes 
		WHERE created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`

	rows, err := r.db.Raw(query, days).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var data dto.DeliveryTrendData
		if err := rows.Scan(&data.Date, &data.Count); err != nil {
			return nil, err
		}
		trendData = append(trendData, data)
	}

	return trendData, nil
}

// CheckDeliveryNumberExists 检查发货单号是否存在
func (r *DeliveryNoteRepositoryImpl) CheckDeliveryNumberExists(deliveryNumber string, excludeID ...uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.DeliveryNote{}).Where("delivery_number = ?", deliveryNumber)
	
	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}
	
	err := query.Count(&count).Error
	return count > 0, err
}

// GenerateDeliveryNumber 生成发货单号
func (r *DeliveryNoteRepositoryImpl) GenerateDeliveryNumber() (string, error) {
	now := time.Now()
	prefix := fmt.Sprintf("DN%s", now.Format("20060102"))
	
	var count int64
	err := r.db.Model(&models.DeliveryNote{}).
		Where("delivery_number LIKE ?", prefix+"%").
		Count(&count).Error
	if err != nil {
		return "", err
	}
	
	return fmt.Sprintf("%s%04d", prefix, count+1), nil
}
