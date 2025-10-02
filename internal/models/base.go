package models

import (
	"gorm.io/gorm"
	"time"
)

// BaseModel 基础模型，包含所有模型的公共字段
type BaseModel struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// AuditableModel 可审计模型，包含创建者和更新者信息
type AuditableModel struct {
	BaseModel
	CreatedBy uint `json:"created_by,omitempty" gorm:"index"`
	UpdatedBy uint `json:"updated_by,omitempty" gorm:"index"`
}

// StatusModel 状态模型，包含状态字段
type StatusModel struct {
	AuditableModel
	IsActive bool `json:"is_active" gorm:"default:true;index"`
}

// CodeModel 编码模型，包含编码字段
type CodeModel struct {
	StatusModel
	Code string `json:"code" gorm:"uniqueIndex;size:50;not null"`
	Name string `json:"name" gorm:"size:255;not null"`
}

// DescriptionModel 描述模型，包含描述字段
type DescriptionModel struct {
	CodeModel
	Description string `json:"description,omitempty" gorm:"type:text"`
}
