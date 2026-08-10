package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"music-player-backend/internal/errors"
	"music-player-backend/internal/models"
	"music-player-backend/internal/services"
	"music-player-backend/internal/utils"

	stdErrors "errors"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()
	var req models.RegisterRequest

	// Bind JSON dengan validasi
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Invalid input format",
			"details": err.Error(),
		})
		return
	}

	// Call service
	user, err := h.service.Register(ctx, &req)
	if err != nil {
		// Handle specific error types
		if stdErrors.Is(err, context.DeadlineExceeded) {
			c.JSON(http.StatusGatewayTimeout, gin.H{
				"error":   "timeout_error",
				"message": "Request processing timeout",
			})
			return
		}
		var appErr *errors.AppError
		if stdErrors.As(err, &appErr) {
			switch appErr.Code {
			case "VALIDATION_ERROR":
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   appErr.Code,
					"message": appErr.Message,
				})
			case "DUPLICATE_ERROR":
				c.JSON(http.StatusConflict, gin.H{
					"error":   appErr.Code,
					"message": appErr.Message,
				})
			case "AUTHENTICATION_ERROR":
				c.JSON(http.StatusUnauthorized, gin.H{
					"error":   appErr.Code,
					"message": appErr.Message,
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "internal_error",
					"message": "An internal server error occurred",
				})
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal_error",
				"message": "Failed to register user",
			})
		}
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"data": gin.H{
			"user": gin.H{
				"id":        user.ID,
				"username":  user.Username,
				"email":     user.Email,
				"createdAt": user.CreatedAt,
			},
		},
	})
}

// Login handles user authentication
// Login handler menggunakan cookie helper
func (h *AuthHandler) Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Invalid input format",
			"details": err.Error(),
		})
		return
	}
	accessToken, refreshToken, err := h.service.Login(ctx, &req)
	if err != nil {
		if stdErrors.Is(err, context.DeadlineExceeded) {
			c.JSON(http.StatusGatewayTimeout, gin.H{
				"error":   "timeout_error",
				"message": "Request processing timeout",
			})
			return
		}
		var appErr *errors.AppError
		if stdErrors.As(err, &appErr) && appErr.Code == "INTERNAL_ERROR" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": "Unable to sign in"})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "authentication_error",
			"message": "Invalid email or password",
		})
		return
	}

	// Set cookie dengan config yang benar
	utils.SetRefreshTokenCookie(c, refreshToken, 60*60*24*30) // 30 days
	user, err := h.service.GetUserFromToken(ctx, accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": "Failed to load user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data": gin.H{
			"accessToken": accessToken,
			"tokenType":   "Bearer",
			"expiresIn":   900,
			"user": gin.H{
				"id": user.ID, "username": user.Username, "email": user.Email, "createdAt": user.CreatedAt,
			},
		},
	})
}

// Refresh generates new access token using refresh token
func (h *AuthHandler) Refresh(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()
	// Get refresh token from cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "authentication_error",
			"message": "Refresh token not found",
		})
		return
	}

	if refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "authentication_error",
			"message": "Refresh token is empty",
		})
		return
	}

	// Call service
	newAccessToken, newRefreshToken, err := h.service.Refresh(ctx, refreshToken)
	if err != nil {
		if stdErrors.Is(err, context.DeadlineExceeded) {
			c.JSON(http.StatusGatewayTimeout, gin.H{
				"error":   "timeout_error",
				"message": "Request processing timeout",
			})
			return
		}
		// Clear invalid cookie
		utils.ClearRefreshTokenCookie(c)

		// Return appropriate error based on error type
		errorMsg := "Invalid or expired refresh token"
		if strings.Contains(err.Error(), "expired") {
			errorMsg = "Refresh token expired"
		}

		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "authentication_error",
			"message": errorMsg,
		})
		return
	}

	utils.SetRefreshTokenCookie(c, newRefreshToken, 60*60*24*30)

	// Return new access token
	c.JSON(http.StatusOK, gin.H{
		"message": "Token refreshed successfully",
		"data": gin.H{
			"accessToken": newAccessToken,
			"tokenType":   "Bearer",
			"expiresIn":   900, // 15 minutes in seconds
		},
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()
	// Get refresh token from cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err == nil && refreshToken != "" {
		// Invalidate all tokens for this user
		if err := h.service.LogoutAll(ctx, refreshToken); err != nil {
			if stdErrors.Is(err, context.DeadlineExceeded) {
				c.JSON(http.StatusGatewayTimeout, gin.H{
					"error":   "timeout_error",
					"message": "Request processing timeout",
				})
				return
			}
			// Log error but still proceed to clear cookie
			// Don't expose internal errors to client
			_ = err // Ignore for security reasons
		}
	}

	utils.ClearRefreshTokenCookie(c)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

// GetCurrentUser returns current user info (protected endpoint)
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "authentication_error",
			"message": "User not authenticated",
		})
		return
	}

	// Get access token from header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "authentication_error",
			"message": "Authorization header missing",
		})
		return
	}

	// Extract token from "Bearer <token>"
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "authentication_error",
			"message": "Invalid authorization header format",
		})
		return
	}

	accessToken := tokenParts[1]

	// Verify token and get user info
	user, err := h.service.GetUserFromToken(ctx, accessToken)
	if err != nil {
		if stdErrors.Is(err, context.DeadlineExceeded) {
			c.JSON(http.StatusGatewayTimeout, gin.H{
				"error":   "timeout_error",
				"message": "Request processing timeout",
			})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "authentication_error",
			"message": "Invalid or expired token",
		})
		return
	}

	// Return user info (excluding sensitive data)
	c.JSON(http.StatusOK, gin.H{
		"message": "User retrieved successfully",
		"data": gin.H{
			"user": gin.H{
				"id":        user.ID,
				"username":  user.Username,
				"email":     user.Email,
				"createdAt": user.CreatedAt,
			},
		},
	})
}

// HealthCheck provides API health status
func (h *AuthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "auth",
		"message": "Authentication service is running",
	})
}
