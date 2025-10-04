package repositories

import (
	"context"
	"github.com/galaxyerp/galaxyErp/internal/common"
	"gorm.io/gorm"
)

// BaseRepository 基础仓储接口
type BaseRepository[T any] interface {
	Create(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id uint) (*T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, options *common.QueryOptions) ([]*T, int64, error)
	Exists(ctx context.Context, id uint) (bool, error)
	Count(ctx context.Context, filters []common.FilterCondition) (int64, error)
}

// TransactionRepository 事务仓储接口
type TransactionRepository interface {
	BeginTx(ctx context.Context) (Transaction, error)
	WithTx(tx Transaction) TransactionRepository
}

// Transaction 事务接口
type Transaction interface {
	Commit() error
	Rollback() error
	GetDB() *gorm.DB
}

// QueryBuilder 查询构建器接口
type QueryBuilder interface {
	Where(condition string, args ...interface{}) QueryBuilder
	WhereIn(field string, values []interface{}) QueryBuilder
	WhereNotIn(field string, values []interface{}) QueryBuilder
	WhereBetween(field string, start, end interface{}) QueryBuilder
	WhereNull(field string) QueryBuilder
	WhereNotNull(field string) QueryBuilder
	OrderBy(field string, order common.SortOrder) QueryBuilder
	Limit(limit int) QueryBuilder
	Offset(offset int) QueryBuilder
	Preload(associations ...string) QueryBuilder
	Build() *gorm.DB
}

// BaseRepositoryImpl 基础仓储实现
type BaseRepositoryImpl[T any] struct {
	db *gorm.DB
}

// NewBaseRepository 创建基础仓储实例
func NewBaseRepository[T any](db *gorm.DB) BaseRepository[T] {
	return &BaseRepositoryImpl[T]{db: db}
}

// Create 创建实体
func (r *BaseRepositoryImpl[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// GetByID 根据ID获取实体
func (r *BaseRepositoryImpl[T]) GetByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// Update 更新实体
func (r *BaseRepositoryImpl[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

// Delete 删除实体
func (r *BaseRepositoryImpl[T]) Delete(ctx context.Context, id uint) error {
	var entity T
	return r.db.WithContext(ctx).Delete(&entity, id).Error
}

// List 列表查询
func (r *BaseRepositoryImpl[T]) List(ctx context.Context, options *common.QueryOptions) ([]*T, int64, error) {
	var entities []*T
	var total int64

	query := r.buildQuery(options)
	
	// 计算总数
	countQuery := r.buildQuery(&common.QueryOptions{Filters: options.Filters})
	if err := countQuery.WithContext(ctx).Model(new(T)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询数据
	if err := query.WithContext(ctx).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

// Exists 检查实体是否存在
func (r *BaseRepositoryImpl[T]) Exists(ctx context.Context, id uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(new(T)).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// Count 计算数量
func (r *BaseRepositoryImpl[T]) Count(ctx context.Context, filters []common.FilterCondition) (int64, error) {
	var count int64
	query := r.buildQuery(&common.QueryOptions{Filters: filters})
	err := query.WithContext(ctx).Model(new(T)).Count(&count).Error
	return count, err
}

// buildQuery 构建查询
func (r *BaseRepositoryImpl[T]) buildQuery(options *common.QueryOptions) *gorm.DB {
	query := r.db.Model(new(T))

	if options == nil {
		return query
	}

	// 应用过滤条件
	for _, filter := range options.Filters {
		query = r.applyFilter(query, filter)
	}

	// 应用排序
	for _, sort := range options.Sorts {
		order := "ASC"
		if sort.Order == common.SortOrderDesc {
			order = "DESC"
		}
		query = query.Order(sort.Field + " " + order)
	}

	// 应用分页
	if options.Pagination != nil {
		query = query.Offset(options.Pagination.GetOffset()).Limit(options.Pagination.GetLimit())
	}

	// 应用关联查询
	for _, include := range options.Includes {
		query = query.Preload(include)
	}

	return query
}

// applyFilter 应用过滤条件
func (r *BaseRepositoryImpl[T]) applyFilter(query *gorm.DB, filter common.FilterCondition) *gorm.DB {
	switch filter.Operator {
	case common.FilterOperatorEq:
		return query.Where(filter.Field+" = ?", filter.Value)
	case common.FilterOperatorNe:
		return query.Where(filter.Field+" != ?", filter.Value)
	case common.FilterOperatorGt:
		return query.Where(filter.Field+" > ?", filter.Value)
	case common.FilterOperatorGte:
		return query.Where(filter.Field+" >= ?", filter.Value)
	case common.FilterOperatorLt:
		return query.Where(filter.Field+" < ?", filter.Value)
	case common.FilterOperatorLte:
		return query.Where(filter.Field+" <= ?", filter.Value)
	case common.FilterOperatorLike:
		return query.Where(filter.Field+" LIKE ?", "%"+filter.Value.(string)+"%")
	case common.FilterOperatorIn:
		return query.Where(filter.Field+" IN ?", filter.Values)
	case common.FilterOperatorNotIn:
		return query.Where(filter.Field+" NOT IN ?", filter.Values)
	case common.FilterOperatorBetween:
		if len(filter.Values) >= 2 {
			return query.Where(filter.Field+" BETWEEN ? AND ?", filter.Values[0], filter.Values[1])
		}
	case common.FilterOperatorIsNull:
		return query.Where(filter.Field + " IS NULL")
	case common.FilterOperatorNotNull:
		return query.Where(filter.Field + " IS NOT NULL")
	}
	return query
}

// TransactionRepositoryImpl 事务仓储实现
type TransactionRepositoryImpl struct {
	db *gorm.DB
}

// NewTransactionRepository 创建事务仓储实例
func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &TransactionRepositoryImpl{db: db}
}

// BeginTx 开始事务
func (r *TransactionRepositoryImpl) BeginTx(ctx context.Context) (Transaction, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &TransactionImpl{tx: tx}, nil
}

// WithTx 使用事务
func (r *TransactionRepositoryImpl) WithTx(tx Transaction) TransactionRepository {
	return &TransactionRepositoryImpl{db: tx.GetDB()}
}

// TransactionImpl 事务实现
type TransactionImpl struct {
	tx *gorm.DB
}

// Commit 提交事务
func (t *TransactionImpl) Commit() error {
	return t.tx.Commit().Error
}

// Rollback 回滚事务
func (t *TransactionImpl) Rollback() error {
	return t.tx.Rollback().Error
}

// GetDB 获取数据库连接
func (t *TransactionImpl) GetDB() *gorm.DB {
	return t.tx
}

// QueryBuilderImpl 查询构建器实现
type QueryBuilderImpl struct {
	db *gorm.DB
}

// NewQueryBuilder 创建查询构建器
func NewQueryBuilder(db *gorm.DB) QueryBuilder {
	return &QueryBuilderImpl{db: db}
}

// Where 添加WHERE条件
func (qb *QueryBuilderImpl) Where(condition string, args ...interface{}) QueryBuilder {
	qb.db = qb.db.Where(condition, args...)
	return qb
}

// WhereIn 添加IN条件
func (qb *QueryBuilderImpl) WhereIn(field string, values []interface{}) QueryBuilder {
	qb.db = qb.db.Where(field+" IN ?", values)
	return qb
}

// WhereNotIn 添加NOT IN条件
func (qb *QueryBuilderImpl) WhereNotIn(field string, values []interface{}) QueryBuilder {
	qb.db = qb.db.Where(field+" NOT IN ?", values)
	return qb
}

// WhereBetween 添加BETWEEN条件
func (qb *QueryBuilderImpl) WhereBetween(field string, start, end interface{}) QueryBuilder {
	qb.db = qb.db.Where(field+" BETWEEN ? AND ?", start, end)
	return qb
}

// WhereNull 添加IS NULL条件
func (qb *QueryBuilderImpl) WhereNull(field string) QueryBuilder {
	qb.db = qb.db.Where(field + " IS NULL")
	return qb
}

// WhereNotNull 添加IS NOT NULL条件
func (qb *QueryBuilderImpl) WhereNotNull(field string) QueryBuilder {
	qb.db = qb.db.Where(field + " IS NOT NULL")
	return qb
}

// OrderBy 添加排序
func (qb *QueryBuilderImpl) OrderBy(field string, order common.SortOrder) QueryBuilder {
	orderStr := "ASC"
	if order == common.SortOrderDesc {
		orderStr = "DESC"
	}
	qb.db = qb.db.Order(field + " " + orderStr)
	return qb
}

// Limit 设置限制数量
func (qb *QueryBuilderImpl) Limit(limit int) QueryBuilder {
	qb.db = qb.db.Limit(limit)
	return qb
}

// Offset 设置偏移量
func (qb *QueryBuilderImpl) Offset(offset int) QueryBuilder {
	qb.db = qb.db.Offset(offset)
	return qb
}

// Preload 预加载关联
func (qb *QueryBuilderImpl) Preload(associations ...string) QueryBuilder {
	for _, assoc := range associations {
		qb.db = qb.db.Preload(assoc)
	}
	return qb
}

// Build 构建查询
func (qb *QueryBuilderImpl) Build() *gorm.DB {
	return qb.db
}