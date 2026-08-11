package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AdminMiddleware harus dipasang SETELAH AuthMiddleware, karena butuh
// "user_id" yang sudah di-set AuthMiddleware di context.
//
// Sengaja query langsung ke DB (bukan lewat UserRepository yang sudah ada)
// supaya tidak perlu mengubah struct/queries User yang mungkin dipakai
// di banyak tempat lain — mengurangi risiko pecah di kode yang sudah jalan.
func AdminMiddleware(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "User not authenticated",
			})
			c.Abort()
			return
		}

		var isAdmin bool
		err := db.QueryRow(c.Request.Context(), "SELECT is_admin FROM users WHERE id = $1", userID).Scan(&isAdmin)
		if err != nil || !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "Admin access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}