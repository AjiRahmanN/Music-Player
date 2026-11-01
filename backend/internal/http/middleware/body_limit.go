package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BodyLimit membatasi ukuran body request. Karena server tidak memakai
// ReadTimeout global (supaya upload besar di koneksi lambat tidak terputus),
// batas ukuran inilah yang melindungi server dari body tak terbatas.
//
// overrides diindeks dengan c.FullPath() (mis. "/api/admin/upload").
func BodyLimit(defaultLimit int64, overrides map[string]int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := defaultLimit
		if v, ok := overrides[c.FullPath()]; ok {
			limit = v
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit)
		c.Next()
	}
}

// SecurityHeaders memasang header dasar. nosniff penting karena file upload
// disajikan dari origin yang sama dengan API.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.Writer.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Referrer-Policy", "no-referrer")
		c.Next()
	}
}
