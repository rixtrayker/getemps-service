package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/yourusername/getemps-service/internal/cache"
	"github.com/yourusername/getemps-service/internal/config"
	"github.com/yourusername/getemps-service/internal/database"
	"github.com/yourusername/getemps-service/internal/handler"
	"github.com/yourusername/getemps-service/internal/middleware"
	"github.com/yourusername/getemps-service/internal/repository"
	"github.com/yourusername/getemps-service/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Setup logger
	logger := logrus.New()
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetLevel(logrus.DebugLevel)
	}

	// Connect to database
	dbConfig := database.Config{
		Host:         cfg.Database.Host,
		Port:         cfg.Database.Port,
		User:         cfg.Database.User,
		Password:     cfg.Database.Password,
		DBName:       cfg.Database.DBName,
		SSLMode:      cfg.Database.SSLMode,
		MaxOpenConns: cfg.Database.MaxOpenConns,
		MaxIdleConns: cfg.Database.MaxIdleConns,
	}

	db, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	logger.Info("Database connection established")

	// Initialize cache
	var appCache cache.Cache
	if cfg.Cache.Enabled {
		appCache = cache.NewInMemoryCache(
			time.Duration(cfg.Cache.TTL)*time.Second,
			10*time.Minute, // cleanup interval
		)
		logger.Info("Cache enabled")
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	salaryRepo := repository.NewSalaryRepository(db)

	// Initialize services
	processStatusService := service.NewProcessStatusService(
		userRepo, 
		salaryRepo, 
		appCache, 
		time.Duration(cfg.Cache.TTL)*time.Second,
	)

	// Initialize handlers
	employeeHandler := handler.NewEmployeeHandler(processStatusService)

	// Setup router
	router := setupRouter(employeeHandler, logger, cfg)

	// Setup HTTP server
	server := &http.Server{
		Addr:    ":" + cfg.App.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		logger.Infof("Server starting on port %s", cfg.App.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited")
}

func setupRouter(employeeHandler *handler.EmployeeHandler, logger *logrus.Logger, cfg *config.Config) *gin.Engine {
	router := gin.New()

	// Add middleware
	router.Use(middleware.Logger(logger))
	router.Use(middleware.Recovery())

	// Health check endpoint
	router.GET("/health", employeeHandler.HealthCheck)

	// API routes with optional authentication
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(cfg.Security.APISecretKey, cfg.Security.TokenRequired))
	{
		api.POST("/GetEmpStatus", employeeHandler.GetEmployeeStatus)
	}

	return router
}