-- Tabel users. Harus dijalankan PERTAMA: refresh_tokens, playlists, dll.
-- mereferensikan users(id).
--
-- Catatan: kolom is_admin ditambahkan di 004_admin.sql.
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Unik tanpa membedakan huruf besar/kecil. Nama index sengaja mengandung
-- "email" / "username": kode Go (user_repository, admin_user_service)
-- memetakan pelanggaran unique ke DuplicateError berdasarkan nama constraint.
CREATE UNIQUE INDEX IF NOT EXISTS users_email_lower_key ON users (LOWER(email));
CREATE UNIQUE INDEX IF NOT EXISTS users_username_lower_key ON users (LOWER(username));
