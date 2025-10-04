package dto

import (
	"reflect"
	"time"

	"github.com/galaxyerp/galaxyErp/internal/models"
)

// Converter DTO转换器接口
type Converter[T any, D any] interface {
	ToDTO(model T) D
	ToModel(dto D) T
	ToDTOList(models []T) []D
	ToModelList(dtos []D) []T
}

// BaseConverter 基础转换器
type BaseConverter[T any, D any] struct{}

// ToDTO 转换为DTO
func (c *BaseConverter[T, D]) ToDTO(model T) D {
	var dto D
	copyFields(model, &dto)
	return dto
}

// ToModel 转换为模型
func (c *BaseConverter[T, D]) ToModel(dto D) T {
	var model T
	copyFields(dto, &model)
	return model
}

// ToDTOList 转换为DTO列表
func (c *BaseConverter[T, D]) ToDTOList(models []T) []D {
	dtos := make([]D, len(models))
	for i, model := range models {
		dtos[i] = c.ToDTO(model)
	}
	return dtos
}

// ToModelList 转换为模型列表
func (c *BaseConverter[T, D]) ToModelList(dtos []D) []T {
	models := make([]T, len(dtos))
	for i, dto := range dtos {
		models[i] = c.ToModel(dto)
	}
	return models
}

// UserConverter 用户转换器
type UserConverter struct {
	BaseConverter[models.User, UserResponse]
}

// NewUserConverter 创建用户转换器
func NewUserConverter() *UserConverter {
	return &UserConverter{}
}

// ToDTO 转换用户模型为DTO
func (c *UserConverter) ToDTO(user models.User) UserResponse {
	return UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Phone:       user.Phone,
		IsActive:    user.IsActive,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

// ToModel 转换DTO为用户模型
func (c *UserConverter) ToModel(dto UserResponse) models.User {
	return models.User{
		BaseModel: models.BaseModel{
			ID:        dto.ID,
			CreatedAt: dto.CreatedAt,
			UpdatedAt: dto.UpdatedAt,
		},
		Username:    dto.Username,
		Email:       dto.Email,
		FirstName:   dto.FirstName,
		LastName:    dto.LastName,
		Phone:       dto.Phone,
		IsActive:    dto.IsActive,
		LastLoginAt: dto.LastLoginAt,
	}
}

// ToCreateModel 转换创建请求为模型
func (c *UserConverter) ToCreateModel(dto UserCreateRequest) models.User {
	return models.User{
		Username:     dto.Username,
		Email:        dto.Email,
		Password:     dto.Password,
		FirstName:    dto.FirstName,
		LastName:     dto.LastName,
		Phone:        dto.Phone,
		DepartmentID: dto.DepartmentID,
		IsActive:     true,
	}
}

// ToUpdateModel 转换更新请求为模型
func (c *UserConverter) ToUpdateModel(dto UserUpdateRequest, existing models.User) models.User {
	if dto.Email != "" {
		existing.Email = dto.Email
	}
	if dto.FirstName != "" {
		existing.FirstName = dto.FirstName
	}
	if dto.LastName != "" {
		existing.LastName = dto.LastName
	}
	if dto.Phone != "" {
		existing.Phone = dto.Phone
	}
	if dto.DepartmentID != nil {
		existing.DepartmentID = dto.DepartmentID
	}
	if dto.IsActive != nil {
		existing.IsActive = *dto.IsActive
	}
	existing.UpdatedAt = time.Now()
	return existing
}

// SalesInvoiceConverter 销售发票转换器
type SalesInvoiceConverter struct {
	BaseConverter[models.SalesInvoice, SalesInvoiceResponse]
}

// NewSalesInvoiceConverter 创建销售发票转换器
func NewSalesInvoiceConverter() *SalesInvoiceConverter {
	return &SalesInvoiceConverter{}
}

// ToDTO 转换销售发票模型为DTO
func (c *SalesInvoiceConverter) ToDTO(invoice models.SalesInvoice) SalesInvoiceResponse {
	return SalesInvoiceResponse{
		ID:                invoice.ID,
		InvoiceNumber:     invoice.InvoiceNumber,
		CustomerID:        invoice.CustomerID,
		SalesOrderID:      invoice.SalesOrderID,
		DeliveryNoteID:    invoice.DeliveryNoteID,
		InvoiceDate:       invoice.InvoiceDate,
		DueDate:           invoice.DueDate,
		PostingDate:       invoice.PostingDate,
		DocStatus:         invoice.DocStatus,
		PaymentStatus:     invoice.PaymentStatus,
		Currency:          invoice.Currency,
		ExchangeRate:      invoice.ExchangeRate,
		SubTotal:          invoice.SubTotal,
		DiscountAmount:    invoice.DiscountAmount,
		TaxAmount:         invoice.TaxAmount,
		ShippingAmount:    invoice.ShippingAmount,
		GrandTotal:        invoice.GrandTotal,
		OutstandingAmount: invoice.OutstandingAmount,
		PaidAmount:        invoice.PaidAmount,
		BillingAddress:    invoice.BillingAddress,
		ShippingAddress:   invoice.ShippingAddress,
		PaymentTerms:      invoice.PaymentTerms,
		PaymentTermsDays:  invoice.PaymentTermsDays,
		SalesPersonID:     invoice.SalesPersonID,
		Territory:         invoice.Territory,
		CustomerPONumber:  invoice.CustomerPONumber,
		Project:           invoice.Project,
		CostCenter:        invoice.CostCenter,
		Terms:             invoice.Terms,
		Notes:             invoice.Notes,
		CreatedBy:         invoice.CreatedBy,
		SubmittedBy:       invoice.SubmittedBy,
		SubmittedAt:       invoice.SubmittedAt,
		CreatedAt:         invoice.CreatedAt,
		UpdatedAt:         invoice.UpdatedAt,
	}
}

// ToCreateModel 转换创建请求为模型
func (c *SalesInvoiceConverter) ToCreateModel(dto SalesInvoiceCreateRequest) models.SalesInvoice {
	return models.SalesInvoice{
		CustomerID:       dto.CustomerID,
		SalesOrderID:     dto.SalesOrderID,
		DeliveryNoteID:   dto.DeliveryNoteID,
		InvoiceDate:      dto.InvoiceDate,
		DueDate:          dto.DueDate,
		PostingDate:      dto.PostingDate,
		Currency:         dto.Currency,
		ExchangeRate:     dto.ExchangeRate,
		BillingAddress:   dto.BillingAddress,
		ShippingAddress:  dto.ShippingAddress,
		PaymentTerms:     dto.PaymentTerms,
		PaymentTermsDays: dto.PaymentTermsDays,
		SalesPersonID:    dto.SalesPersonID,
		Territory:        dto.Territory,
		CustomerPONumber: dto.CustomerPONumber,
		Project:          dto.Project,
		CostCenter:       dto.CostCenter,
		Terms:            dto.Terms,
		Notes:            dto.Notes,
		DocStatus:        "Draft",
		PaymentStatus:    "Unpaid",
	}
}

// ToUpdateModel 转换更新请求为模型
func (c *SalesInvoiceConverter) ToUpdateModel(dto SalesInvoiceUpdateRequest, existing models.SalesInvoice) models.SalesInvoice {
	if dto.CustomerID != nil {
		existing.CustomerID = *dto.CustomerID
	}
	if dto.InvoiceDate != nil {
		existing.InvoiceDate = *dto.InvoiceDate
	}
	if dto.DueDate != nil {
		existing.DueDate = *dto.DueDate
	}
	if dto.PostingDate != nil {
		existing.PostingDate = *dto.PostingDate
	}
	if dto.Currency != nil {
		existing.Currency = *dto.Currency
	}
	if dto.ExchangeRate != nil {
		existing.ExchangeRate = *dto.ExchangeRate
	}
	if dto.BillingAddress != nil {
		existing.BillingAddress = *dto.BillingAddress
	}
	if dto.ShippingAddress != nil {
		existing.ShippingAddress = *dto.ShippingAddress
	}
	if dto.PaymentTerms != nil {
		existing.PaymentTerms = *dto.PaymentTerms
	}
	if dto.PaymentTermsDays != nil {
		existing.PaymentTermsDays = *dto.PaymentTermsDays
	}
	if dto.SalesPersonID != nil {
		existing.SalesPersonID = dto.SalesPersonID
	}
	if dto.Territory != nil {
		existing.Territory = *dto.Territory
	}
	if dto.CustomerPONumber != nil {
		existing.CustomerPONumber = *dto.CustomerPONumber
	}
	if dto.Project != nil {
		existing.Project = *dto.Project
	}
	if dto.CostCenter != nil {
		existing.CostCenter = *dto.CostCenter
	}
	if dto.Terms != nil {
		existing.Terms = *dto.Terms
	}
	if dto.Notes != nil {
		existing.Notes = *dto.Notes
	}
	existing.UpdatedAt = time.Now()
	return existing
}

// DeliveryNoteConverter 交货单转换器
type DeliveryNoteConverter struct {
	BaseConverter[models.DeliveryNote, DeliveryNoteResponse]
}

// NewDeliveryNoteConverter 创建交货单转换器
func NewDeliveryNoteConverter() *DeliveryNoteConverter {
	return &DeliveryNoteConverter{}
}

// ToDTO 转换送货单模型为DTO
func (c *DeliveryNoteConverter) ToDTO(note models.DeliveryNote) DeliveryNoteResponse {
	return DeliveryNoteResponse{
		ID:             note.ID,
		DeliveryNumber: note.DeliveryNumber,
		Date:           note.Date,
		Status:         note.Status,
	}
}

// ToCreateModel 转换创建请求为模型
func (c *DeliveryNoteConverter) ToCreateModel(dto DeliveryNoteCreateRequest) models.DeliveryNote {
	return models.DeliveryNote{
		CustomerID:     dto.CustomerID,
		SalesOrderID:   dto.SalesOrderID,
		Date:           dto.Date,
		Status:         "Draft",
		Transporter:    dto.Transporter,
		DriverName:     dto.DriverName,
		VehicleNumber:  dto.VehicleNumber,
		Destination:    dto.Destination,
		Notes:          dto.Notes,
	}
}

// ToUpdateModel 转换更新请求为模型
func (c *DeliveryNoteConverter) ToUpdateModel(dto DeliveryNoteUpdateRequest, existing models.DeliveryNote) models.DeliveryNote {
	if dto.CustomerID != nil {
		existing.CustomerID = *dto.CustomerID
	}
	if dto.Date != nil {
		existing.Date = *dto.Date
	}
	if dto.Status != "" {
		existing.Status = dto.Status
	}
	if dto.Transporter != "" {
		existing.Transporter = dto.Transporter
	}
	if dto.DriverName != "" {
		existing.DriverName = dto.DriverName
	}
	if dto.VehicleNumber != "" {
		existing.VehicleNumber = dto.VehicleNumber
	}
	if dto.Destination != "" {
		existing.Destination = dto.Destination
	}
	if dto.Notes != "" {
		existing.Notes = dto.Notes
	}
	existing.UpdatedAt = time.Now()
	return existing
}

// copyFields 复制字段（简单实现）
func copyFields(src, dst interface{}) {
	srcVal := reflect.ValueOf(src)
	dstVal := reflect.ValueOf(dst)

	if srcVal.Kind() == reflect.Ptr {
		srcVal = srcVal.Elem()
	}
	if dstVal.Kind() == reflect.Ptr {
		dstVal = dstVal.Elem()
	}

	if !dstVal.CanSet() {
		return
	}

	srcType := srcVal.Type()
	dstType := dstVal.Type()

	for i := 0; i < srcType.NumField(); i++ {
		srcField := srcType.Field(i)
		srcFieldValue := srcVal.Field(i)

		if dstField, found := dstType.FieldByName(srcField.Name); found {
			dstFieldValue := dstVal.FieldByName(srcField.Name)
			if dstFieldValue.CanSet() && srcFieldValue.Type().AssignableTo(dstField.Type) {
				dstFieldValue.Set(srcFieldValue)
			}
		}
	}
}