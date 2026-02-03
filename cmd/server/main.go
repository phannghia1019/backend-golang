package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hbc/app-hrm-api-go/internal/config"
	"github.com/hbc/app-hrm-api-go/internal/handlers"
	"github.com/hbc/app-hrm-api-go/internal/middleware"
	"github.com/hbc/app-hrm-api-go/internal/utils"
)

var log = utils.NewLogger("Server")
var buildTime = time.Now()

func main() {
	// Load configuration
	envFile := getEnvFile()
	cfg, err := config.Load(envFile)
	if err != nil {
		log.Error("Failed to load config: %v", err)
		os.Exit(1)
	}

	// Set Gin mode
	if cfg.NodeEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to MongoDB
	if err := config.ConnectMongoDB(); err != nil {
		log.Warn("Failed to connect to MongoDB: %v (continuing without MongoDB)", err)
	}

	// Create Gin router
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LanguageMiddleware())

	// Health check endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"env":       cfg.NodeEnv,
			"buildDate": buildTime.Format("02/01/2006 15:04:05"),
			"fromNow":   time.Since(buildTime).String(),
			"version":   "1.0.0-go",
		})
	})

	r.GET("/maintaince", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Hệ thống đang bảo trì")
	})

	// Setup routes
	setupRoutes(r)

	// Start server
	port := cfg.Port
	if port == "" {
		port = "3000"
	}

	log.Info("Express server started on url: http://localhost:%s/v1", port)
	if err := r.Run(":" + port); err != nil {
		log.Error("Failed to start server: %v", err)
		os.Exit(1)
	}
}

// setupRoutes configures all API routes
func setupRoutes(r *gin.Engine) {
	authHandler := handlers.NewAuthHandler()

	// V1 routes
	v1 := r.Group("/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.GET("/permissions", authHandler.GetPermissions)
		}
	}

	// Multi-company routes (with database_name prefix)
	dbV1 := r.Group("/:database_name/v1")
	{
		// Auth routes with database resolver
		auth := dbV1.Group("/auth")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.GET("/permissions", authHandler.GetPermissions)
		}
	}
}

// getEnvFile returns the appropriate .env file path
func getEnvFile() string {
	env := os.Getenv("NODE_ENV")
	if env == "" {
		env = "development"
	}
	return fmt.Sprintf("env/.env.%s", env)
}
