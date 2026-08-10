package handlers

import (
	"net/http"

	apperrors "music-player-backend/internal/errors"
	"music-player-backend/internal/models"
	"music-player-backend/internal/services"

	stdErrors "errors"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service}
}

// ======================
// GET /api/me
// ======================
func (h *UserHandler) Me(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.service.GetMe(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"created":  user.CreatedAt,
	})
}

// ======================
// PUT /api/profile
// ======================
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "authentication_error",
			"message": "User not authenticated",
		})
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Invalid input format",
			"details": err.Error(),
		})
		return
	}

	user, err := h.service.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		var appErr *apperrors.AppError
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
			case "NOT_FOUND":
				c.JSON(http.StatusNotFound, gin.H{
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
				"message": "Failed to update profile",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
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
