package routes

import (
	"net/http"
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

	// Batas ukuran form-data upload (25MB, samakan dengan maxFileSize terbesar di UploadHandler)
	r.MaxMultipartMemory = 25 << 20

	// Header keamanan dasar (nosniff untuk file upload) + batas ukuran body.
	// Server tidak punya Read/WriteTimeout global, jadi body dibatasi di sini:
	// 1 MB untuk request JSON, 26 MB khusus endpoint upload.
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.BodyLimit(1<<20, map[string]int64{
		"/api/admin/upload": 26 << 20,
	}))

	// Rate limit endpoint auth yang rawan brute force / abuse.
	loginLimiter := middleware.NewRateLimiter(10, time.Minute)   // 10 percobaan login/menit/IP
	registerLimiter := middleware.NewRateLimiter(5, time.Minute) // 5 registrasi/menit/IP
	refreshLimiter := middleware.NewRateLimiter(60, time.Minute) // refresh lebih longgar

	// Layani file yang sudah diupload secara statis, contoh:
	// http://localhost:8080/uploads/covers/169999999.jpg
	// (didaftarkan SETELAH SecurityHeaders supaya ikut mendapat nosniff)
	r.Static("/uploads", "./uploads")

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
	adminUserRepo := repositories.NewAdminUserRepository(db)

	authService := services.NewAuthService(userRepo, rtRepo)
	userService := services.NewUserService(userRepo)
	songService := services.NewSongService(songRepo)
	playlistService := services.NewPlaylistService(playlistRepo, songRepo)
	adminUserService := services.NewAdminUserService(adminUserRepo)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	songHandler := handlers.NewSongHandler(songService)
	playlistHandler := handlers.NewPlaylistHandler(playlistService)
	adminUserHandler := handlers.NewAdminUserHandler(adminUserService)
	uploadHandler := handlers.NewUploadHandler()
	healthHandler := handlers.NewHealthHandler()

	// Public routes
	public := r.Group("/api")
	{
		public.GET("/health", healthHandler.Check)
		public.GET("/ping", healthHandler.Ping)

		authGroup := public.Group("/auth")
		{
			authGroup.POST("/register", registerLimiter.Middleware(), authHandler.Register)
			authGroup.POST("/login", loginLimiter.Middleware(), authHandler.Login)
			authGroup.POST("/refresh", refreshLimiter.Middleware(), authHandler.Refresh)
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
	}

	// Admin routes (butuh login DAN is_admin = true)
	admin := r.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware(db))
	{
		// Kalau request ini sampai ke handler (lolos AdminMiddleware),
		// berarti user itu pasti admin — dipakai frontend buat cek akses.
		admin.GET("/check", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"isAdmin": true})
		})

		admin.POST("/upload", uploadHandler.Upload)
		admin.POST("/songs", songHandler.Create)
		admin.PUT("/songs/:id", songHandler.Update)
		admin.DELETE("/songs/:id", songHandler.Delete)

		admin.GET("/users", adminUserHandler.List)
		admin.PUT("/users/:id", adminUserHandler.Update)
		admin.DELETE("/users/:id", adminUserHandler.Delete)

		admin.GET("/playlists", playlistHandler.AdminListAll)
		admin.DELETE("/playlists/:id", playlistHandler.AdminDelete)

		// Future routes can be added here
		// protected.GET("/users", userHandler.List)
	}

	return r
}
