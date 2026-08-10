package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"music-player-backend/internal/config"
	"music-player-backend/internal/http/routes"
	"music-player-backend/internal/logger"
	"music-player-backend/internal/utils"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Load environment variables
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("Warning: .env not loaded")
		}
	}

	// Load configuration
	if _, err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger.InitLogger(config.AppConfig.Server.Env)

	// Initialize JWT secret
	utils.SetJWTSecret(config.AppConfig.JWT.Secret)
	// Verify JWT secret is set
	if len(utils.JwtSecret) == 0 {
		logger.Log.Fatal("JWT_SECRET is not properly initialized")
	}
	// Connect to database
	if err := config.ConnectDB(); err != nil {
		logger.Log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer config.CloseDB()

	// Setup router with dependencies
	r := routes.SetupRouter(config.DB)

	// Create server with graceful shutdown
	server := &http.Server{
		Addr:         ":" + config.AppConfig.Server.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on :%s in %s mode", config.AppConfig.Server.Port, config.AppConfig.Server.Env)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Log.Info("Server exited properly")
}