package models

import (
	"time"
)

// LegacyItem 兼容旧数据库结构的物料模型
type LegacyItem struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`
	Code        string    `json:"code" gorm:"uniqueIndex;size:100;not null"`
	Name        string    `json:"name" gorm:"size:255;not null"`
	Description string    `json:"description,omitempty" gorm:"type:text"`
	Category    string    `json:"category" gorm:"size:100"`
	Unit        string    `json:"unit" gorm:"size:50"`
	Cost        float64   `json:"cost" gorm:"default:0"`
	Price       float64   `json:"price" gorm:"default:0"`
	ReorderLevel int      `json:"reorder_level" gorm:"default:0"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
}

// TableName 指定表名
func (LegacyItem) TableName() string {
	return "items"
}