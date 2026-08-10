package routes

import (
	"os"
	"strings"
	"time"

	"music-player-backend/internal/http/handlers"
	"music-player-backend/internal/http/middleware"
	"music-player-backend/internal/logger"
	"music-player-backend/internal/repositories"
	"music-player-backend/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var Log *zap.Logger

func SetupRouter(db *pgxpool.Pool) *gin.Engine {
	r := gin.Default()
	allowedOrigins := strings.Split(os.Getenv("FRONTEND_ORIGIN"), ",")
	if allowedOrigins[0] == "" {
		allowedOrigins = []string{"http://localhost:3000", "http://localhost:8080"}
	}

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "X-Total-Count"},
		AllowCredentials: true,
		MaxAge:           12 * 3600,
	}))

	// Request logging middleware
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		// Log request details
		logger.Log.Info("Request completed",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
		)
	})

	// Dependency injection
	userRepo := repositories.NewUserRepository(db)
	rtRepo := repositories.NewRefreshTokenRepository(db)
	songRepo := repositories.NewSongRepository(db)
	playlistRepo := repositories.NewPlaylistRepository(db)

	authService := services.NewAuthService(userRepo, rtRepo)
	userService := services.NewUserService(userRepo)
	songService := services.NewSongService(songRepo)
	playlistService := services.NewPlaylistService(playlistRepo, songRepo)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	songHandler := handlers.NewSongHandler(songService)
	playlistHandler := handlers.NewPlaylistHandler(playlistService)
	healthHandler := handlers.NewHealthHandler()

	// Public routes
	public := r.Group("/api")
	{
		public.GET("/health", healthHandler.Check)
		public.GET("/ping", healthHandler.Ping)

		authGroup := public.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/refresh", authHandler.Refresh)
			authGroup.POST("/logout", authHandler.Logout)
			authGroup.GET("/health", authHandler.HealthCheck)
		}
	}

	// Protected routes (require authentication)
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// User routes
		protected.GET("/me", userHandler.Me)
		protected.GET("/profile", authHandler.GetCurrentUser)
		protected.PUT("/profile", userHandler.UpdateProfile)

		// Song routes
		protected.GET("/songs", songHandler.List)
		protected.GET("/songs/:id", songHandler.GetByID)

		// Playlist routes
		protected.POST("/playlists", playlistHandler.Create)
		protected.GET("/playlists", playlistHandler.List)
		protected.GET("/playlists/:id", playlistHandler.GetDetail)
		protected.DELETE("/playlists/:id", playlistHandler.Delete)
		protected.POST("/playlists/:id/songs", playlistHandler.AddSong)
		protected.DELETE("/playlists/:id/songs/:songId", playlistHandler.RemoveSong)

		// Future routes can be added here
		// protected.GET("/users", userHandler.List)
	}

	return r
}