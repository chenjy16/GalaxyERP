package container

import (
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/galaxyerp/galaxyErp/internal/controllers"
	"github.com/galaxyerp/galaxyErp/internal/handlers"
	"github.com/galaxyerp/galaxyErp/internal/repositories"
	"github.com/galaxyerp/galaxyErp/internal/services"
)

// Container 依赖注入容器
type Container struct {
	DB *gorm.DB

	// Repository Interfaces (仓储层接口)
	UserRepository         repositories.UserRepository
	ItemRepository         repositories.ItemRepository
	StockRepository        repositories.StockRepository
	WarehouseRepository    repositories.WarehouseRepository
	MovementRepository     repositories.MovementRepository
	CustomerRepository     repositories.CustomerRepository
	SalesOrderRepository   repositories.SalesOrderRepository
	QuotationRepository    repositories.QuotationRepository
	SupplierRepository     repositories.SupplierRepository
	PurchaseRequestRepository repositories.PurchaseRequestRepository
	PurchaseOrderRepository   repositories.PurchaseOrderRepository
	ProjectRepository      repositories.ProjectRepository
	TaskRepository         repositories.TaskRepository
	EmployeeRepository     repositories.EmployeeRepository
	AccountRepository      repositories.AccountRepository
	AuditLogRepository     repositories.AuditLogRepository
	ProductRepository      repositories.ProductRepository
	SalesInvoiceRepository repositories.SalesInvoiceRepository
	DeliveryNoteRepository repositories.DeliveryNoteRepository

	// Service Interfaces (服务层接口)
	AuditLogService          services.AuditLogService
	UserService              services.UserService
	ItemService              services.ItemService
	StockService             services.StockService
	WarehouseService         services.WarehouseService
	MovementService          services.MovementService
	CustomerService          services.CustomerService
	SalesOrderService        services.SalesOrderService
	QuotationService         services.QuotationService
	QuotationTemplateService services.QuotationTemplateService
	QuotationVersionService  services.QuotationVersionService
	SalesInvoiceService      services.SalesInvoiceService
	DeliveryNoteService      services.DeliveryNoteServiceInterface
	ProductService           services.ProductService

	// Purchase Services
	SupplierService        services.SupplierService
	PurchaseRequestService services.PurchaseRequestService
	PurchaseOrderService   services.PurchaseOrderService

	// Project Services
	ProjectService   services.ProjectService
	TaskService      services.TaskService
	MilestoneService services.MilestoneService
	TimeEntryService services.TimeEntryService

	// HR Services
	EmployeeService   services.EmployeeService
	AttendanceService services.AttendanceService
	PayrollService    services.PayrollService
	LeaveService      services.LeaveService

	// Accounting Services
	AccountService      services.AccountService
	JournalEntryService services.JournalEntryService
	PaymentEntryService services.PaymentEntryService

	// Handlers
	AuditLogHandler *handlers.AuditLogHandler

	// Controllers
	UserController         *controllers.UserController
	InventoryController    *controllers.InventoryController
	SalesController        *controllers.SalesController
	DeliveryNoteController *controllers.DeliveryNoteController
	ProductionController   *controllers.ProductionController
	SystemController       *controllers.SystemController
	PurchaseController     *controllers.PurchaseController
	ProjectController      *controllers.ProjectController
	AccountingController   *controllers.AccountingController
	HRController           *controllers.HRController
}

// NewContainer 创建新的依赖注入容器
func NewContainer(db *gorm.DB, jwtSecret string, jwtExpiryHours int) *Container {
	container := &Container{
		DB: db,
	}

	// 初始化仓储层
	container.initRepositories()

	// 初始化服务层
	container.initServices(jwtSecret, jwtExpiryHours)

	// 初始化控制器层
	container.initControllers()

	return container
}

// initRepositories 初始化仓储层
func (c *Container) initRepositories() {
	// 创建仓库实例并存储到容器中
	c.UserRepository = repositories.NewUserRepository(c.DB)
	c.ItemRepository = repositories.NewItemRepository(c.DB)
	c.StockRepository = repositories.NewStockRepository(c.DB)
	c.WarehouseRepository = repositories.NewWarehouseRepository(c.DB)
	c.MovementRepository = repositories.NewMovementRepository(c.DB)
	c.CustomerRepository = repositories.NewCustomerRepository(c.DB)
	c.SalesOrderRepository = repositories.NewSalesOrderRepository(c.DB)
	c.QuotationRepository = repositories.NewQuotationRepository(c.DB)
	c.SalesInvoiceRepository = repositories.NewSalesInvoiceRepository(c.DB)
	c.DeliveryNoteRepository = repositories.NewDeliveryNoteRepository(c.DB)
	c.ProductRepository = repositories.NewProductRepository(c.DB)

	// Purchase repositories
	c.SupplierRepository = repositories.NewSupplierRepository(c.DB)
	c.PurchaseRequestRepository = repositories.NewPurchaseRequestRepository(c.DB)
	c.PurchaseOrderRepository = repositories.NewPurchaseOrderRepository(c.DB)

	// Project repositories
	c.ProjectRepository = repositories.NewProjectRepository(c.DB)
	c.TaskRepository = repositories.NewTaskRepository(c.DB)

	// HR repositories
	c.EmployeeRepository = repositories.NewEmployeeRepository(c.DB)

	// Accounting repositories
	c.AccountRepository = repositories.NewAccountRepository(c.DB)

	// Audit log repository
	c.AuditLogRepository = repositories.NewAuditLogRepository(c.DB)
}

// initServices 初始化服务层
func (c *Container) initServices(jwtSecret string, jwtExpiryHours int) {
	// 创建其他仓库实例（暂时未接口化的）
	quotationTemplateRepo := repositories.NewQuotationTemplateRepository(c.DB)
	quotationVersionRepo := repositories.NewQuotationVersionRepository(c.DB)
	milestoneRepo := repositories.NewMilestoneRepository(c.DB)
	timeEntryRepo := repositories.NewTimeEntryRepository(c.DB)
	attendanceRepo := repositories.NewAttendanceRepository(c.DB)
	payrollRepo := repositories.NewPayrollRepository(c.DB)
	leaveRepo := repositories.NewLeaveRepository(c.DB)
	journalEntryRepo := repositories.NewJournalEntryRepository(c.DB)
	paymentEntryRepo := repositories.NewPaymentEntryRepository(c.DB)

	// 初始化审计日志服务
	c.AuditLogService = services.NewAuditLogService(c.AuditLogRepository, zap.L())

	// 初始化服务（使用容器中的仓储接口）
	c.UserService = services.NewUserService(c.UserRepository, c.AuditLogService, jwtSecret, jwtExpiryHours)
	c.ItemService = services.NewItemService(c.ItemRepository)
	c.StockService = services.NewStockService(c.StockRepository)
	c.WarehouseService = services.NewWarehouseService(c.WarehouseRepository)
	c.MovementService = services.NewMovementService(c.MovementRepository, c.StockRepository, c.ItemRepository, c.WarehouseRepository)
	c.CustomerService = services.NewCustomerService(c.CustomerRepository)
	c.ProductService = services.NewProductService(c.ProductRepository)

	// Accounting services (需要先初始化，因为其他服务可能依赖)
	c.AccountService = services.NewAccountService(c.AccountRepository)
	c.JournalEntryService = services.NewJournalEntryService(journalEntryRepo, c.AccountRepository)
	c.PaymentEntryService = services.NewPaymentEntryService(paymentEntryRepo)

	// Sales services (依赖会计服务)
	c.SalesOrderService = services.NewSalesOrderService(c.SalesOrderRepository, c.CustomerRepository)
	c.QuotationService = services.NewQuotationService(c.QuotationRepository, c.CustomerRepository)
	c.QuotationTemplateService = services.NewQuotationTemplateService(quotationTemplateRepo, c.QuotationRepository)
	c.QuotationVersionService = services.NewQuotationVersionService(quotationVersionRepo, c.QuotationRepository)
	c.SalesInvoiceService = services.NewSalesInvoiceService(c.SalesInvoiceRepository, c.CustomerRepository, c.SalesOrderRepository, c.PaymentEntryService)
	c.DeliveryNoteService = services.NewDeliveryNoteService(c.DeliveryNoteRepository, c.SalesOrderRepository, c.CustomerRepository)

	// Purchase services
	c.SupplierService = services.NewSupplierService(c.SupplierRepository)
	c.PurchaseRequestService = services.NewPurchaseRequestService(c.PurchaseRequestRepository)
	c.PurchaseOrderService = services.NewPurchaseOrderService(c.PurchaseOrderRepository)

	// Project services
	c.ProjectService = services.NewProjectService(c.ProjectRepository)
	c.TaskService = services.NewTaskService(c.TaskRepository)
	c.MilestoneService = services.NewMilestoneService(milestoneRepo)
	c.TimeEntryService = services.NewTimeEntryService(timeEntryRepo)

	// HR services
	c.EmployeeService = services.NewEmployeeService(c.EmployeeRepository)
	c.AttendanceService = services.NewAttendanceService(attendanceRepo, c.EmployeeRepository)
	c.PayrollService = services.NewPayrollService(payrollRepo, c.EmployeeRepository)
	c.LeaveService = services.NewLeaveService(leaveRepo, c.EmployeeRepository)
}

// initControllers 初始化控制器层
func (c *Container) initControllers() {
	// 创建ControllerUtils实例
	utils := controllers.NewControllerUtils()

	// 初始化处理器
	c.AuditLogHandler = handlers.NewAuditLogHandler(c.AuditLogService, zap.L())

	c.UserController = controllers.NewUserController(c.UserService)
	c.InventoryController = controllers.NewInventoryController(c.ItemService, c.StockService, c.WarehouseService, c.MovementService)
	c.SalesController = controllers.NewSalesController(c.CustomerService, c.SalesOrderService, c.QuotationService, c.QuotationTemplateService, c.SalesInvoiceService, c.QuotationVersionService)
	c.DeliveryNoteController = controllers.NewDeliveryNoteController(c.DeliveryNoteService)
	c.ProductionController = controllers.NewProductionController(c.ProductService)
	c.SystemController = controllers.NewSystemController()

	// Purchase Controller
	c.PurchaseController = controllers.NewPurchaseController(
		utils,
		c.SupplierService,
		c.PurchaseRequestService,
		c.PurchaseOrderService,
	)

	// Project Controller
	c.ProjectController = controllers.NewProjectController(
		c.ProjectService,
		c.TaskService,
		c.MilestoneService,
		c.TimeEntryService,
	)

	// Accounting Controller
	c.AccountingController = controllers.NewAccountingController(
		c.AccountService,
		c.JournalEntryService,
	)

	// HR Controller
	c.HRController = controllers.NewHRController(
		c.EmployeeService,
		c.AttendanceService,
		c.PayrollService,
		c.LeaveService,
	)
}
