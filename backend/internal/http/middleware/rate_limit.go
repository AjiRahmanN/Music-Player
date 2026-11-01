package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter adalah limiter fixed-window sederhana di memori, per kunci
// (default: IP klien). Cukup untuk satu instance; kalau backend di-scale ke
// banyak instance, pindahkan ke store bersama (mis. Redis).
type RateLimiter struct {
	mu      sync.Mutex
	limit   int
	window  time.Duration
	entries map[string]*rlEntry
	now     func() time.Time
}

type rlEntry struct {
	count   int
	resetAt time.Time
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		limit:   limit,
		window:  window,
		entries: make(map[string]*rlEntry),
		now:     time.Now,
	}
}

// Allow mencatat satu percobaan untuk key. Mengembalikan false (dan sisa
// waktu tunggu) jika batas terlampaui.
func (l *RateLimiter) Allow(key string) (bool, time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	e, ok := l.entries[key]
	if !ok || !now.Before(e.resetAt) {
		l.entries[key] = &rlEntry{count: 1, resetAt: now.Add(l.window)}
		l.gcLocked(now)
		return true, 0
	}
	if e.count >= l.limit {
		return false, e.resetAt.Sub(now)
	}
	e.count++
	return true, 0
}

// gcLocked membuang entri kedaluwarsa agar map tidak tumbuh tanpa batas.
func (l *RateLimiter) gcLocked(now time.Time) {
	if len(l.entries) < 1024 {
		return
	}
	for k, e := range l.entries {
		if !now.Before(e.resetAt) {
			delete(l.entries, k)
		}
	}
}

// Middleware membatasi request per IP klien.
func (l *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, wait := l.Allow(c.ClientIP())
		if !ok {
			secs := int(wait.Seconds()) + 1
			c.Header("Retry-After", strconv.Itoa(secs))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "rate_limited",
				"message": "Terlalu banyak percobaan. Coba lagi beberapa saat lagi.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
