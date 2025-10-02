package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/galaxyerp/galaxyErp/internal/dto"
	"github.com/galaxyerp/galaxyErp/internal/models"
)

type DeliveryNoteRepository struct {
	db *gorm.DB
}

func NewDeliveryNoteRepository(db *gorm.DB) *DeliveryNoteRepository {
	return &DeliveryNoteRepository{db: db}
}

// Create 创建发货单
func (r *DeliveryNoteRepository) Create(deliveryNote *models.DeliveryNote) error {
	return r.db.Create(deliveryNote).Error
}

// GetByID 根据ID获取发货单
func (r *DeliveryNoteRepository) GetByID(id uint) (*models.DeliveryNote, error) {
	var deliveryNote models.DeliveryNote
	err := r.db.Preload("Customer").
		Preload("SalesOrder").
		Preload("Items").
		Preload("Items.Product").
		First(&deliveryNote, id).Error
	if err != nil {
		return nil, err
	}
	return &deliveryNote, nil
}

// GetByDeliveryNumber 根据发货单号获取发货单
func (r *DeliveryNoteRepository) GetByDeliveryNumber(deliveryNumber string) (*models.DeliveryNote, error) {
	var deliveryNote models.DeliveryNote
	err := r.db.Preload("Customer").
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

// Update 更新发货单
func (r *DeliveryNoteRepository) Update(deliveryNote *models.DeliveryNote) error {
	return r.db.Save(deliveryNote).Error
}

// Delete 删除发货单
func (r *DeliveryNoteRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除发货单明细
		if err := tx.Where("delivery_note_id = ?", id).Delete(&models.DeliveryNoteItem{}).Error; err != nil {
			return err
		}
		// 删除发货单
		return tx.Delete(&models.DeliveryNote{}, id).Error
	})
}

// List 获取发货单列表
func (r *DeliveryNoteRepository) List(req *dto.DeliveryNoteListRequest) ([]*models.DeliveryNote, int64, error) {
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

	// 应用排序
	if req.SortBy != "" {
		order := req.SortBy
		if req.SortDesc {
			order += " DESC"
		}
		query = query.Order(order)
	} else {
		query = query.Order("created_at DESC")
	}

	// 应用分页
	if req.Page > 0 && req.PageSize > 0 {
		offset := (req.Page - 1) * req.PageSize
		query = query.Offset(offset).Limit(req.PageSize)
	}

	// 预加载关联数据
	err := query.Preload("Customer").
		Preload("SalesOrder").
		Preload("Items").
		Find(&deliveryNotes).Error

	return deliveryNotes, total, err
}

// GetBySalesOrderID 根据销售订单ID获取发货单列表
func (r *DeliveryNoteRepository) GetBySalesOrderID(salesOrderID uint) ([]*models.DeliveryNote, error) {
	var deliveryNotes []*models.DeliveryNote
	err := r.db.Preload("Customer").
		Preload("Items").
		Where("sales_order_id = ?", salesOrderID).
		Find(&deliveryNotes).Error
	return deliveryNotes, err
}

// GetByCustomerID 根据客户ID获取发货单列表
func (r *DeliveryNoteRepository) GetByCustomerID(customerID uint) ([]*models.DeliveryNote, error) {
	var deliveryNotes []*models.DeliveryNote
	err := r.db.Preload("SalesOrder").
		Preload("Items").
		Where("customer_id = ?", customerID).
		Find(&deliveryNotes).Error
	return deliveryNotes, err
}

// UpdateStatus 更新发货单状态
func (r *DeliveryNoteRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.DeliveryNote{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// GetStatistics 获取发货单统计信息
func (r *DeliveryNoteRepository) GetStatistics(ctx *gin.Context) (*dto.DeliveryNoteStatisticsResponse, error) {
	var stats dto.DeliveryNoteStatisticsResponse

	// 总发货单数
	if err := r.db.Model(&models.DeliveryNote{}).Count(&stats.TotalDeliveries).Error; err != nil {
		return nil, err
	}

	// 待发货数量
	if err := r.db.Model(&models.DeliveryNote{}).
		Where("status = ?", "pending").
		Count(&stats.PendingDeliveries).Error; err != nil {
		return nil, err
	}

	// 已完成发货数量
	if err := r.db.Model(&models.DeliveryNote{}).
		Where("status = ?", "delivered").
		Count(&stats.CompletedDeliveries).Error; err != nil {
		return nil, err
	}

	// 总发货数量
	var totalQuantity sql.NullFloat64
	if err := r.db.Model(&models.DeliveryNote{}).
		Select("SUM(total_quantity)").
		Scan(&totalQuantity).Error; err != nil {
		return nil, err
	}
	if totalQuantity.Valid {
		stats.TotalQuantity = totalQuantity.Float64
	}

	return &stats, nil
}

// GetDeliveryTrend 获取发货趋势数据
func (r *DeliveryNoteRepository) GetDeliveryTrend(days int) ([]dto.DeliveryTrendData, error) {
	var trendData []dto.DeliveryTrendData

	query := `
		SELECT 
			DATE(date) as date,
			COUNT(*) as count,
			SUM(total_quantity) as total_quantity
		FROM delivery_notes 
		WHERE date >= ? 
		GROUP BY DATE(date) 
		ORDER BY date
	`

	startDate := time.Now().AddDate(0, 0, -days)
	err := r.db.Raw(query, startDate).Scan(&trendData).Error

	return trendData, err
}

// CheckDeliveryNumberExists 检查发货单号是否存在
func (r *DeliveryNoteRepository) CheckDeliveryNumberExists(deliveryNumber string, excludeID ...uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.DeliveryNote{}).Where("delivery_number = ?", deliveryNumber)
	
	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}
	
	err := query.Count(&count).Error
	return count > 0, err
}

// GenerateDeliveryNumber 生成发货单号
func (r *DeliveryNoteRepository) GenerateDeliveryNumber() (string, error) {
	var count int64
	today := time.Now().Format("20060102")
	
	// 获取今天的发货单数量
	err := r.db.Model(&models.DeliveryNote{}).
		Where("delivery_number LIKE ?", "DN"+today+"%").
		Count(&count).Error
	if err != nil {
		return "", err
	}
	
	return fmt.Sprintf("DN%s%04d", today, count+1), nil
}