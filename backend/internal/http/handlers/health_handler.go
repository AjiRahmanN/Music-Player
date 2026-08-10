package handlers

import (
	"context"
	"net/http"
	"sync"
	"time"

	"music-player-backend/internal/config"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
    startTime time.Time
    mu        sync.RWMutex
}

func NewHealthHandler() *HealthHandler {
    return &HealthHandler{
        startTime: time.Now(),
    }
}

func (h *HealthHandler) GetUptime() time.Duration {
    h.mu.RLock()
    defer h.mu.RUnlock()
    return time.Since(h.startTime)
}

func (h *HealthHandler) ResetStartTime() {
    h.mu.Lock()
    defer h.mu.Unlock()
    h.startTime = time.Now()
}

// Check provides comprehensive health status
func (h *HealthHandler) Check(c *gin.Context) {
    healthStatus := gin.H{
        "status":    "healthy",
        "timestamp": time.Now().UTC(),
        "services":  make(map[string]string),
    }

    // Check database connection
    if config.DB == nil {
        healthStatus["services"].(map[string]string)["database"] = "disconnected"
        healthStatus["status"] = "degraded"
    } else {
        ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
        defer cancel()

        if err := config.DB.Ping(ctx); err != nil {
            healthStatus["services"].(map[string]string)["database"] = "unhealthy"
            healthStatus["status"] = "unhealthy"
            healthStatus["error"] = err.Error()
        } else {
            healthStatus["services"].(map[string]string)["database"] = "connected"
        }
    }

    // Add server info
    if config.AppConfig != nil {
        healthStatus["environment"] = config.AppConfig.Server.Env
        healthStatus["uptime"] = h.GetUptime().String()
    }

    statusCode := http.StatusOK
    if healthStatus["status"] == "unhealthy" {
        statusCode = http.StatusServiceUnavailable
    }

    c.JSON(statusCode, healthStatus)
}

// Simple health check for load balancers
func (h *HealthHandler) Ping(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status":  "ok",
        "message": "pong",
        "uptime":  h.GetUptime().String(),
    })
}