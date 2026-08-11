package handlers

import (
	"net/http"

	"music-player-backend/internal/models"
	"music-player-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type AdminUserHandler struct {
	service *services.AdminUserService
}

func NewAdminUserHandler(service *services.AdminUserService) *AdminUserHandler {
	return &AdminUserHandler{service}
}

// ======================
// GET /api/admin/users
// ======================
func (h *AdminUserHandler) List(c *gin.Context) {
	users, err := h.service.List(c.Request.Context())
	if err != nil {
		respondAppError(c, err, "Failed to fetch users")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"users": users}})
}

// ======================
// PUT /api/admin/users/:id
// ======================
func (h *AdminUserHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateUserByAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Invalid input format",
			"details": err.Error(),
		})
		return
	}

	user, err := h.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		respondAppError(c, err, "Failed to update user")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"data":    gin.H{"user": user},
	})
}

// ======================
// DELETE /api/admin/users/:id
// ======================
func (h *AdminUserHandler) Delete(c *gin.Context) {
	requesterID := c.GetString("user_id")
	targetID := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), requesterID, targetID); err != nil {
		respondAppError(c, err, "Failed to delete user")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}