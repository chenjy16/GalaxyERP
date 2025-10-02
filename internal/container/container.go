package container

import (
	"gorm.io/gorm"

	"github.com/galaxyerp/galaxyErp/internal/controllers"
	"github.com/galaxyerp/galaxyErp/internal/repositories"
	"github.com/galaxyerp/galaxyErp/internal/services"
)

// Container 依赖注入容器
type Container struct {
	DB *gorm.DB

	// Services
	UserService         services.UserService
	ItemService         services.ItemService
	StockService        services.StockService
	WarehouseService    services.WarehouseService
	MovementService     services.MovementService
	CustomerService     services.CustomerService
	SalesOrderService   services.SalesOrderService
	QuotationService    services.QuotationService
	SalesInvoiceService services.SalesInvoiceService
	ProductService      services.ProductService

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

	// Controllers
	UserController       *controllers.UserController
	InventoryController  *controllers.InventoryController
	SalesController      *controllers.SalesController
	ProductionController *controllers.ProductionController
	SystemController     *controllers.SystemController
	PurchaseController   *controllers.PurchaseController
	ProjectController    *controllers.ProjectController
	AccountingController *controllers.AccountingController
	HRController         *controllers.HRController
}

// NewContainer 创建新的依赖注入容器
func NewContainer(db *gorm.DB, jwtSecret string, jwtExpiryHours int) *Container {
	container := &Container{
		DB: db,
	}

	// 初始化服务层
	container.initServices(jwtSecret, jwtExpiryHours)

	// 初始化控制器层
	container.initControllers()

	return container
}

// initServices 初始化服务层
func (c *Container) initServices(jwtSecret string, jwtExpiryHours int) {
	// 创建仓库实例
	userRepo := repositories.NewUserRepository(c.DB)
	itemRepo := repositories.NewItemRepository(c.DB)
	stockRepo := repositories.NewStockRepository(c.DB)
	warehouseRepo := repositories.NewWarehouseRepository(c.DB)
	movementRepo := repositories.NewMovementRepository(c.DB)
	customerRepo := repositories.NewCustomerRepository(c.DB)
	salesOrderRepo := repositories.NewSalesOrderRepository(c.DB)
	quotationRepo := repositories.NewQuotationRepository(c.DB)
	productRepo := repositories.NewProductRepository(c.DB)

	// Purchase repositories
	supplierRepo := repositories.NewSupplierRepository(c.DB)
	purchaseRequestRepo := repositories.NewPurchaseRequestRepository(c.DB)
	purchaseOrderRepo := repositories.NewPurchaseOrderRepository(c.DB)

	// Project repositories
	projectRepo := repositories.NewProjectRepository(c.DB)
	taskRepo := repositories.NewTaskRepository(c.DB)
	milestoneRepo := repositories.NewMilestoneRepository(c.DB)
	timeEntryRepo := repositories.NewTimeEntryRepository(c.DB)

	// HR repositories
	employeeRepo := repositories.NewEmployeeRepository(c.DB)
	attendanceRepo := repositories.NewAttendanceRepository(c.DB)
	payrollRepo := repositories.NewPayrollRepository(c.DB)
	leaveRepo := repositories.NewLeaveRepository(c.DB)

	// Accounting repositories
	accountRepo := repositories.NewAccountRepository(c.DB)
	journalEntryRepo := repositories.NewJournalEntryRepository(c.DB)

	// 初始化服务
	c.UserService = services.NewUserService(userRepo, jwtSecret, jwtExpiryHours)
	c.ItemService = services.NewItemService(itemRepo)
	c.StockService = services.NewStockService(stockRepo)
	c.WarehouseService = services.NewWarehouseService(warehouseRepo)
	c.MovementService = services.NewMovementService(movementRepo, stockRepo, itemRepo, warehouseRepo)
	c.CustomerService = services.NewCustomerService(customerRepo)
	c.SalesOrderService = services.NewSalesOrderService(salesOrderRepo, customerRepo)
	c.QuotationService = services.NewQuotationService(quotationRepo, customerRepo)
	c.SalesInvoiceService = services.NewSalesInvoiceService(c.DB)
	c.ProductService = services.NewProductService(productRepo)

	// Purchase services
	c.SupplierService = services.NewSupplierService(supplierRepo)
	c.PurchaseRequestService = services.NewPurchaseRequestService(purchaseRequestRepo)
	c.PurchaseOrderService = services.NewPurchaseOrderService(purchaseOrderRepo)

	// Project services
	c.ProjectService = services.NewProjectService(projectRepo)
	c.TaskService = services.NewTaskService(taskRepo)
	c.MilestoneService = services.NewMilestoneService(milestoneRepo)
	c.TimeEntryService = services.NewTimeEntryService(timeEntryRepo)

	// HR services
	c.EmployeeService = services.NewEmployeeService(employeeRepo)
	c.AttendanceService = services.NewAttendanceService(attendanceRepo, employeeRepo)
	c.PayrollService = services.NewPayrollService(payrollRepo, employeeRepo)
	c.LeaveService = services.NewLeaveService(leaveRepo, employeeRepo)

	// Accounting services
	c.AccountService = services.NewAccountService(accountRepo)
	c.JournalEntryService = services.NewJournalEntryService(journalEntryRepo, accountRepo)
}

// initControllers 初始化控制器层
func (c *Container) initControllers() {
	// 创建ControllerUtils实例
	utils := controllers.NewControllerUtils()

	c.UserController = controllers.NewUserController(c.UserService)
	c.InventoryController = controllers.NewInventoryController(c.ItemService, c.StockService, c.WarehouseService, c.MovementService)
	c.SalesController = controllers.NewSalesController(c.CustomerService, c.SalesOrderService, c.QuotationService, c.SalesInvoiceService)
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
