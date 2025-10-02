package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/galaxyerp/galaxyErp/internal/container"
	"github.com/galaxyerp/galaxyErp/internal/middleware"
	"github.com/galaxyerp/galaxyErp/internal/models"
	"github.com/galaxyerp/galaxyErp/internal/routes"

	"github.com/galaxyerp/galaxyErp/internal/utils"
)

func main() {
	// Initialize configuration
	initConfig()

	// Initialize logger based on environment
	var logger *zap.Logger
	logLevel := viper.GetString("logging.level")
	logFormat := viper.GetString("logging.format")

	if logLevel == "debug" || logFormat == "console" {
		// Use development logger for debug mode or console format
		config := zap.NewDevelopmentConfig()
		if logLevel == "debug" {
			config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		} else if logLevel == "info" {
			config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		} else if logLevel == "error" {
			config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
		}
		logger, _ = config.Build()
	} else {
		// Use production logger
		logger, _ = zap.NewProduction()
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	// Connect to database
	utils.ConnectDatabase()

	// Auto migrate models
	utils.GetDB().AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.RolePermission{},
		&models.UserRole{},
		&models.DataPermission{},
		&models.Company{},
		&models.Department{},
		&models.Position{},
		&models.SystemConfig{},
		&models.ApprovalWorkflow{},
		&models.ApprovalStep{},
		&models.AuditLog{},
		&models.Backup{},
		&models.Account{},
		&models.Transaction{},
		&models.JournalEntry{},
		&models.Receivable{},
		&models.Payable{},
		&models.FixedAsset{},
		&models.DepreciationEntry{},
		&models.TaxRate{},
		&models.TaxEntry{},
		&models.Currency{},
		&models.FinancialReport{},
		&models.CostCenter{},
		&models.BankAccount{},
		&models.PaymentEntry{},
		&models.Budget{},
		&models.ExchangeRateHistory{},
		&models.TaxTemplate{},
		&models.FiscalYear{},
		&models.AccountingPeriod{},
		&models.Product{},
		&models.Item{},
		&models.Warehouse{},
		&models.Location{},
		&models.Stock{},
		&models.StockMovement{},
		&models.StockTransfer{},
		&models.Customer{},
		&models.Quotation{},
		&models.QuotationItem{},
		&models.SalesOrder{},
		&models.SalesOrderItem{},
		&models.DeliveryNote{},
		&models.DeliveryNoteItem{},
		&models.Supplier{},
		&models.PurchaseRequest{},
		&models.PurchaseRequestItem{},
		&models.PurchaseOrder{},
		&models.PurchaseOrderItem{},
		&models.PurchaseReceipt{},
		&models.PurchaseReceiptItem{},
		// Production Manufacturing models
		&models.ProductionPlan{},
		&models.MaterialRequirement{},
		&models.ProcessRoute{},
		&models.ProcessOperation{},
		&models.WorkCenter{},
		&models.ProductionOrder{},
		&models.ProductionOrderItem{},
		&models.ProductionProgress{},
		&models.QualityInspection{},
		&models.NonConformingProduct{},
		&models.Equipment{},
		&models.EquipmentMaintenance{},
		&models.EquipmentFailure{},
		// Project Management models
		&models.Project{},
		&models.ProjectMilestone{},
		&models.ProjectTask{},
		&models.TaskTimeRecord{},
		&models.ProjectResource{},
		&models.ProjectReport{},
		// Human Resources models
		&models.Employee{},
		&models.Attendance{}, // 修正为Attendance
		&models.Leave{},      // 修正为Leave
		&models.OvertimeRecord{},
		&models.Payroll{},
		&models.PerformanceGoal{},
		&models.PerformanceReview{},
		&models.Training{},            // 修正为Training
		&models.TrainingParticipant{}, // 修正为TrainingParticipant
		&models.Skill{},
		&models.EmployeeSkill{},
	)

	// 初始化依赖注入容器
	appContainer := container.NewContainer(utils.GetDB(), viper.GetString("jwt.secret"), viper.GetInt("jwt.expiry"))

	// Create server
	r := gin.Default()

	// Register routes
	registerRoutes(r, appContainer)

	// Get server port from config
	port := viper.GetString("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("Failed to start server", zap.Error(err))
		}
	}()

	zap.L().Info("Server started", zap.String("port", port))

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zap.L().Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server forced to shutdown", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}

func initConfig() {
	// Check for environment variable to determine config mode
	env := os.Getenv("GALAXYERP_ENV")
	if env == "" {
		env = "dev" // default to development
	}

	// Set config file based on environment
	switch env {
	case "dev":
		viper.SetConfigName("dev")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./configs")
	case "test":
		viper.SetConfigName("test")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./configs")
	case "prod":
		viper.SetConfigFile(".env")
	default:
		viper.SetConfigName("dev")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./configs")
	}

	// Enable viper to read Environment Variables
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; try to read .env file for prod or dev config for dev
			env := os.Getenv("GALAXYERP_ENV")
			if env == "" {
				env = "dev"
			}

			if env == "prod" {
				viper.SetConfigFile(".env")
				if err := viper.ReadInConfig(); err != nil {
					fmt.Println("Warning: No .env file found for production configuration")
				} else {
					fmt.Println("Using .env configuration file")
				}
			} else {
				viper.SetConfigName("dev")
				viper.SetConfigType("yaml")
				viper.AddConfigPath("./configs")
				if err := viper.ReadInConfig(); err != nil {
					fmt.Println("Warning: No config file found (configs/dev.yaml)")
				} else {
					fmt.Println("Using dev configuration file")
				}
			}
		} else {
			// Config file was found but another error was produced
			log.Fatalf("Error reading config file: %s", err)
		}
	} else {
		fmt.Printf("Using %s configuration\n", env)
	}
}

func registerRoutes(r *gin.Engine, appContainer *container.Container) {
	// Add CORS middleware
	r.Use(middleware.CORSMiddleware())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "GalaxyERP is running",
		})
	})

	// API version 1 group
	v1 := r.Group("/api/v1")
	{
		v1.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "GalaxyERP API v1",
			})
		})

		// Authentication routes (public for login/register)
		authGroup := v1.Group("/auth")
		{
			// Use real user controller handlers
			authGroup.POST("/register", appContainer.UserController.Register)
			authGroup.POST("/login", appContainer.UserController.Login)
		}

		// Authenticated auth routes (require token)
		authProtected := v1.Group("/auth")
		authProtected.Use(middleware.AuthMiddleware(viper.GetString("jwt.secret")))
		{
			authProtected.GET("/me", appContainer.UserController.GetProfile)
			// Simple logout endpoint
			authProtected.POST("/logout", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"success": true,
					"message": "已登出",
				})
			})
			// Token refresh placeholder
			authProtected.POST("/refresh", func(c *gin.Context) {
				c.JSON(501, gin.H{"error": "功能暂未实现"})
			})
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(viper.GetString("jwt.secret")))
		{
			// Register modular routes
			routes.RegisterUserRoutes(protected, appContainer)
			routes.RegisterAccountingRoutes(protected, appContainer)
			routes.RegisterInventoryRoutes(protected, appContainer)
			routes.RegisterSalesRoutes(protected, appContainer)
			routes.RegisterPurchaseRoutes(protected, appContainer)
			routes.RegisterProductionRoutes(protected, appContainer)
			routes.RegisterHRRoutes(protected, appContainer)
			routes.RegisterProjectRoutes(protected, appContainer)
			// System management routes - modularized
			routes.RegisterSystemRoutes(protected, appContainer)
		}
	}
}
