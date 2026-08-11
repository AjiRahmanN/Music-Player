package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

var allowedExtensions = map[string][]string{
	"cover": {".jpg", ".jpeg", ".png", ".webp"},
	"audio": {".mp3", ".wav", ".ogg", ".m4a"},
}

var maxFileSize = map[string]int64{
	"cover": 5 << 20,  // 5 MB
	"audio": 25 << 20, // 25 MB
}

// ======================
// POST /api/admin/upload
// Form-data: file=<binary>, type=cover|audio
// ======================
func (h *UploadHandler) Upload(c *gin.Context) {
	fileType := c.PostForm("type")
	if fileType != "cover" && fileType != "audio" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Field 'type' harus 'cover' atau 'audio'",
		})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "File tidak ditemukan di request",
		})
		return
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !isAllowedExt(fileType, ext) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": fmt.Sprintf("Ekstensi %s tidak diizinkan untuk tipe %s", ext, fileType),
		})
		return
	}

	if fileHeader.Size > maxFileSize[fileType] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Ukuran file terlalu besar",
		})
		return
	}

	// Folder tujuan: ./uploads/covers atau ./uploads/audio
	dir := filepath.Join("uploads", fileType+"s")
	if err := os.MkdirAll(dir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Gagal menyiapkan folder upload",
		})
		return
	}

	// Nama file unik supaya tidak bentrok antar upload
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	destPath := filepath.Join(dir, filename)

	src, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": "Gagal membaca file"})
		return
	}
	defer src.Close()

	dst, err := os.Create(destPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": "Gagal menyimpan file"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": "Gagal menyimpan file"})
		return
	}

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	url := fmt.Sprintf("%s://%s/uploads/%ss/%s", scheme, c.Request.Host, fileType, filename)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{"url": url},
	})
}

func isAllowedExt(fileType, ext string) bool {
	for _, allowed := range allowedExtensions[fileType] {
		if allowed == ext {
			return true
		}
	}
	return false
}