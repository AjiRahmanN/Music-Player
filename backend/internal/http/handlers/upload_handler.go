package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"music-player-backend/internal/models"

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

// allowedMIME memetakan ekstensi -> MIME hasil sniffing isi file
// (http.DetectContentType). Ekstensi saja mudah dipalsukan, jadi isi file
// ikut diperiksa dan harus cocok dengan ekstensinya.
var allowedMIME = map[string]map[string][]string{
	"cover": {
		".jpg":  {"image/jpeg"},
		".jpeg": {"image/jpeg"},
		".png":  {"image/png"},
		".webp": {"image/webp"},
	},
	"audio": {
		".mp3": {"audio/mpeg"},
		".wav": {"audio/wave", "audio/wav", "audio/x-wav"},
		".ogg": {"application/ogg", "audio/ogg"},
		// Go mendeteksi container MP4/M4A sebagai video/mp4.
		".m4a": {"video/mp4", "audio/mp4", "audio/x-m4a"},
	},
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
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error":   "validation_error",
				"message": "Ukuran file terlalu besar",
			})
			return
		}
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

	src, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": "Gagal membaca file"})
		return
	}
	defer src.Close()

	// Periksa isi file (bukan cuma ekstensi) dari 512 byte pertama.
	head := make([]byte, 512)
	n, readErr := io.ReadFull(src, head)
	if readErr != nil && !errors.Is(readErr, io.ErrUnexpectedEOF) && !errors.Is(readErr, io.EOF) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": "Gagal membaca file"})
		return
	}
	if !sniffAllowed(fileType, ext, head[:n]) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Isi file tidak sesuai dengan ekstensi/tipe yang diizinkan",
		})
		return
	}
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": "Gagal membaca file"})
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

	// Simpan/kirim PATH RELATIF; base URL berasal dari konfigurasi server
	// (PUBLIC_BASE_URL), bukan dari header Host yang bisa dimanipulasi.
	// "url" berisi versi absolut untuk preview di client, "path" versi relatif.
	relPath := fmt.Sprintf("/uploads/%ss/%s", fileType, filename)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"url":  models.ResolveAssetURL(relPath),
			"path": relPath,
		},
	})
}

// sniffAllowed memastikan isi file (dari header-nya) cocok dengan ekstensi.
func sniffAllowed(fileType, ext string, head []byte) bool {
	allowed := allowedMIME[fileType][ext]
	if len(allowed) == 0 || len(head) == 0 {
		return false
	}
	detected := http.DetectContentType(head)
	if i := strings.Index(detected, ";"); i >= 0 {
		detected = strings.TrimSpace(detected[:i])
	}
	for _, m := range allowed {
		if m == detected {
			return true
		}
	}
	return false
}

func isAllowedExt(fileType, ext string) bool {
	for _, allowed := range allowedExtensions[fileType] {
		if allowed == ext {
			return true
		}
	}
	return false
}
