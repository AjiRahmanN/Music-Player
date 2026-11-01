-- Hak akses admin. Tidak ada endpoint publik untuk menjadi admin; promosikan
-- user lewat SQL manual, mis.:
--   UPDATE users SET is_admin = true WHERE email = 'email_kamu@contoh.com';
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_admin BOOLEAN NOT NULL DEFAULT FALSE;

-- Mempercepat COUNT(*) admin (dipakai untuk mencegah admin terakhir dicabut).
CREATE INDEX IF NOT EXISTS idx_users_is_admin ON users (is_admin) WHERE is_admin = true;
