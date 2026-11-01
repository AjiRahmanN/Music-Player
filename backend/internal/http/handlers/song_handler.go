package handlers

import (
	"net/http"
	"strconv"
	"strings"

	stdErrors "errors"

	apperrors "music-player-backend/internal/errors"
	"music-player-backend/internal/models"
	"music-player-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type SongHandler struct {
	service *services.SongService
}

func NewSongHandler(service *services.SongService) *SongHandler {
	return &SongHandler{service}
}

// ======================
// GET /api/songs
// GET /api/songs?q=keyword  -> pencarian berdasarkan judul/artis
// Parameter opsional: limit (default 100, maks 200) dan offset (default 0).
// ======================
func (h *SongHandler) List(c *gin.Context) {
	query := strings.TrimSpace(c.Query("q"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))

	var songs []models.Song
	var err error

	if query != "" {
		songs, err = h.service.Search(c.Request.Context(), query, limit, offset)
	} else {
		songs, err = h.service.List(c.Request.Context(), limit, offset)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Failed to fetch songs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"songs": songs,
		},
	})
}

// ======================
// GET /api/songs/:id
// ======================
func (h *SongHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	song, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		var appErr *apperrors.AppError
		if stdErrors.As(err, &appErr) && appErr.Code == "NOT_FOUND" {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   appErr.Code,
				"message": appErr.Message,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Failed to fetch song",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"song": song,
		},
	})
}

// ======================
// POST /api/admin/songs
// ======================
func (h *SongHandler) Create(c *gin.Context) {
	var req models.CreateSongRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Invalid input format",
			"details": err.Error(),
		})
		return
	}

	song, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		var appErr *apperrors.AppError
		if stdErrors.As(err, &appErr) && appErr.Code == "VALIDATION_ERROR" {
			c.JSON(http.StatusBadRequest, gin.H{"error": appErr.Code, "message": appErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": "Failed to create song"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Song created successfully",
		"data":    gin.H{"song": song},
	})
}

// ======================
// PUT /api/admin/songs/:id
// ======================
func (h *SongHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateSongRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Invalid input format",
			"details": err.Error(),
		})
		return
	}

	song, err := h.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		var appErr *apperrors.AppError
		if stdErrors.As(err, &appErr) {
			switch appErr.Code {
			case "VALIDATION_ERROR":
				c.JSON(http.StatusBadRequest, gin.H{"error": appErr.Code, "message": appErr.Message})
			case "NOT_FOUND":
				c.JSON(http.StatusNotFound, gin.H{"error": appErr.Code, "message": appErr.Message})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": "Failed to update song"})
			}
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": "Failed to update song"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Song updated successfully",
		"data":    gin.H{"song": song},
	})
}

// ======================
// DELETE /api/admin/songs/:id
// ======================
func (h *SongHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		var appErr *apperrors.AppError
		if stdErrors.As(err, &appErr) && appErr.Code == "NOT_FOUND" {
			c.JSON(http.StatusNotFound, gin.H{"error": appErr.Code, "message": appErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": "Failed to delete song"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song deleted successfully"})
}
