package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestRateLimiterAllowAndReset(t *testing.T) {
	l := NewRateLimiter(3, time.Minute)
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	l.now = func() time.Time { return now }

	for i := 0; i < 3; i++ {
		if ok, _ := l.Allow("1.1.1.1"); !ok {
			t.Fatalf("percobaan %d seharusnya diizinkan", i+1)
		}
	}
	ok, wait := l.Allow("1.1.1.1")
	if ok || wait <= 0 {
		t.Fatalf("percobaan ke-4 seharusnya ditolak dengan wait > 0 (ok=%v wait=%v)", ok, wait)
	}
	if ok, _ := l.Allow("2.2.2.2"); !ok {
		t.Fatal("IP lain tidak boleh terkena limit")
	}

	now = now.Add(time.Minute + time.Second)
	if ok, _ := l.Allow("1.1.1.1"); !ok {
		t.Fatal("setelah window lewat, percobaan seharusnya diizinkan lagi")
	}
}

func TestRateLimiterMiddleware429(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/login", NewRateLimiter(1, time.Minute).Middleware(), func(c *gin.Context) { c.Status(http.StatusOK) })

	do := func() *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/login", nil)
		req.RemoteAddr = "9.9.9.9:1234"
		r.ServeHTTP(w, req)
		return w
	}
	if w := do(); w.Code != http.StatusOK {
		t.Fatalf("request pertama = %d", w.Code)
	}
	w := do()
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("request kedua = %d, want 429", w.Code)
	}
	if w.Header().Get("Retry-After") == "" {
		t.Fatal("header Retry-After kosong")
	}
}

func TestBodyLimitDefaultAndOverride(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(BodyLimit(10, map[string]int64{"/big": 100}))
	read := func(c *gin.Context) {
		buf := make([]byte, 200)
		n, err := c.Request.Body.Read(buf)
		_ = n
		if err != nil && n < 200 && strings.Contains(err.Error(), "too large") {
			c.Status(http.StatusRequestEntityTooLarge)
			return
		}
		c.Status(http.StatusOK)
	}
	r.POST("/small", read)
	r.POST("/big", read)

	post := func(path string, size int) int {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, path, strings.NewReader(strings.Repeat("a", size))))
		return w.Code
	}
	if got := post("/small", 50); got != http.StatusRequestEntityTooLarge {
		t.Fatalf("/small 50 byte = %d, want 413", got)
	}
	if got := post("/big", 50); got != http.StatusOK {
		t.Fatalf("/big 50 byte = %d, want 200", got)
	}
	if got := post("/big", 150); got != http.StatusRequestEntityTooLarge {
		t.Fatalf("/big 150 byte = %d, want 413", got)
	}
}

func TestSecurityHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(SecurityHeaders())
	r.GET("/", func(c *gin.Context) { c.Status(http.StatusOK) })
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
	if w.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatal("nosniff tidak terpasang")
	}
}
