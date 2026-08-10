package handlers

import (
	"net/http"

	stdErrors "errors"

	apperrors "music-player-backend/internal/errors"
	"music-player-backend/internal/models"
	"music-player-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type PlaylistHandler struct {
	service *services.PlaylistService
}

func NewPlaylistHandler(service *services.PlaylistService) *PlaylistHandler {
	return &PlaylistHandler{service}
}

// respondAppError memetakan apperrors.AppError ke HTTP status yang sesuai.
// Dipakai berulang di semua handler playlist supaya konsisten.
func respondAppError(c *gin.Context, err error, fallbackMessage string) {
	var appErr *apperrors.AppError
	if stdErrors.As(err, &appErr) {
		switch appErr.Code {
		case "VALIDATION_ERROR":
			c.JSON(http.StatusBadRequest, gin.H{"error": appErr.Code, "message": appErr.Message})
		case "NOT_FOUND":
			c.JSON(http.StatusNotFound, gin.H{"error": appErr.Code, "message": appErr.Message})
		case "DUPLICATE_ERROR":
			c.JSON(http.StatusConflict, gin.H{"error": appErr.Code, "message": appErr.Message})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": "An internal server error occurred"})
		}
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": fallbackMessage})
}

// ======================
// POST /api/playlists
// ======================
func (h *PlaylistHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")

	var req models.CreatePlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Invalid input format",
			"details": err.Error(),
		})
		return
	}

	playlist, err := h.service.Create(c.Request.Context(), userID, &req)
	if err != nil {
		respondAppError(c, err, "Failed to create playlist")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Playlist created successfully",
		"data":    gin.H{"playlist": playlist},
	})
}

// ======================
// GET /api/playlists
// ======================
func (h *PlaylistHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")

	playlists, err := h.service.ListByUser(c.Request.Context(), userID)
	if err != nil {
		respondAppError(c, err, "Failed to fetch playlists")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{"playlists": playlists},
	})
}

// ======================
// GET /api/playlists/:id
// ======================
func (h *PlaylistHandler) GetDetail(c *gin.Context) {
	userID := c.GetString("user_id")
	playlistID := c.Param("id")

	detail, err := h.service.GetDetail(c.Request.Context(), userID, playlistID)
	if err != nil {
		respondAppError(c, err, "Failed to fetch playlist")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{"playlist": detail},
	})
}

// ======================
// DELETE /api/playlists/:id
// ======================
func (h *PlaylistHandler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	playlistID := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), userID, playlistID); err != nil {
		respondAppError(c, err, "Failed to delete playlist")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Playlist deleted successfully"})
}

// ======================
// POST /api/playlists/:id/songs
// ======================
func (h *PlaylistHandler) AddSong(c *gin.Context) {
	userID := c.GetString("user_id")
	playlistID := c.Param("id")

	var req models.AddSongToPlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Invalid input format",
			"details": err.Error(),
		})
		return
	}

	if err := h.service.AddSong(c.Request.Context(), userID, playlistID, &req); err != nil {
		respondAppError(c, err, "Failed to add song to playlist")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song added to playlist"})
}

// ======================
// DELETE /api/playlists/:id/songs/:songId
// ======================
func (h *PlaylistHandler) RemoveSong(c *gin.Context) {
	userID := c.GetString("user_id")
	playlistID := c.Param("id")
	songID := c.Param("songId")

	if err := h.service.RemoveSong(c.Request.Context(), userID, playlistID, songID); err != nil {
		respondAppError(c, err, "Failed to remove song from playlist")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song removed from playlist"})
}