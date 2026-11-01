# 🎵 Melodify

Aplikasi web music streaming full-stack terinspirasi dari Spotify — dibangun dari nol menggunakan **Go (Gin)** di backend dan **Next.js 15** di frontend. Mendukung autentikasi, pemutar musik lengkap (queue, shuffle, repeat, seek), playlist, pencarian lagu, dan panel admin untuk kelola konten.

![Go](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go&logoColor=white)
![Next.js](https://img.shields.io/badge/Next.js-15-black?logo=next.js&logoColor=white)
![TypeScript](https://img.shields.io/badge/TypeScript-3178C6?logo=typescript&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?logo=postgresql&logoColor=white)
![TailwindCSS](https://img.shields.io/badge/TailwindCSS-38B2AC?logo=tailwind-css&logoColor=white)

---

## ✨ Fitur

**Autentikasi**

- Register, login, logout dengan JWT access token + refresh token (httpOnly cookie, rotasi otomatis)
- Access token hanya disimpan di memori; setelah reload, sesi dipulihkan lewat refresh token
- Refresh token ber-single-flight di frontend (request paralel berbagi satu proses refresh)
- Logout hanya mengakhiri sesi di perangkat yang sedang dipakai
- Rate limiting pada login, register, dan refresh (per IP)
- Edit profil (username & email)

**Pemutar Musik**

- Play / pause / next / previous
- Progress bar interaktif (klik & geser untuk seek)
- Kontrol volume + mute
- Shuffle & repeat (off / all / one)
- Antrian lagu mengikuti playlist yang sedang diputar

**Lagu & Playlist**

- Pencarian lagu (judul/artis) dengan debounce
- Buat, lihat, dan hapus playlist pribadi
- Tambah/hapus lagu dari playlist

**Panel Admin**

- Kelola lagu — tambah/edit/hapus, upload file cover & audio langsung (durasi terdeteksi otomatis dari file)
- Kelola user — edit data & status admin, hapus user
- Kelola playlist — lihat & hapus playlist milik semua user

---

## 🛠️ Tech Stack

|              |                                                                                 |
| ------------ | ------------------------------------------------------------------------------- |
| **Backend**  | Go, Gin, PostgreSQL (pgx), JWT, bcrypt, zap logger                              |
| **Frontend** | Next.js 15 (App Router), TypeScript, Zustand, Tailwind CSS, Axios, Lucide Icons |

---

## 📁 Struktur Proyek

```
spotify/
├── backend/
│   ├── cmd/server/           # entry point
│   ├── internal/
│   │   ├── config/           # konfigurasi & koneksi database
│   │   ├── models/           # struct data & validasi request
│   │   ├── repositories/     # query database
│   │   ├── services/         # logic bisnis
│   │   ├── http/
│   │   │   ├── handlers/     # HTTP handler
│   │   │   ├── middleware/   # auth & admin middleware
│   │   │   └── routes/       # routing
│   │   └── utils/            # JWT, password hashing, dll
│   ├── migration/            # file SQL migrasi 000-004 (jalankan berurutan)
│   └── uploads/              # file cover & audio hasil upload (dev only)
│
└── frontend/
    └── src/
        ├── app/
        │   ├── (auth)/        # login, register — tanpa Navbar/AuthGuard
        │   └── (main)/        # home, playlists, admin, profile — dijaga AuthGuard
        ├── components/        # komponen React (player, admin, layout)
        ├── services/          # pemanggil API (axios)
        ├── store/             # Zustand (auth, player)
        └── types/             # TypeScript types
```

---

## 🚀 Cara Menjalankan

### Prasyarat

- [Go](https://go.dev/) 1.24+
- [Node.js](https://nodejs.org/) 18+
- [PostgreSQL](https://www.postgresql.org/) 14+

### 1. Clone repository

```bash
git clone https://github.com/AjiRahmanN/Music-Player.git
cd Music-Player
```

### 2. Setup Backend

```bash
cd backend
```

Salin `.env.example` menjadi `.env` di folder `backend/` lalu isi nilainya:

```bash
cp .env.example .env
```

| Variabel | Wajib | Keterangan |
| --- | --- | --- |
| `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` | ✅ | Koneksi PostgreSQL |
| `JWT_SECRET` | ✅ | String acak panjang (`openssl rand -base64 48`) |
| `FRONTEND_ORIGIN` | ✅ | Origin frontend untuk CORS (pisahkan dengan koma) |
| `PORT` | – | Default `8080` (`SERVER_PORT` masih diterima sebagai fallback) |
| `PUBLIC_BASE_URL` | – | URL publik backend untuk URL file upload. Default `http://localhost:<PORT>` |
| `JWT_ACCESS_EXPIRY` | – | Default `15m` |
| `JWT_REFRESH_EXPIRY` | – | Default `720h` (30 hari) |
| `ENV` | – | `development` (default), `staging`, atau `production` (cookie Secure + SameSite=None, butuh HTTPS) |
| `COOKIE_DOMAIN` | – | Domain cookie refresh token (staging/production) |

Jalankan seluruh file migrasi secara **berurutan** ke database PostgreSQL kamu. Dengan `psql`:

```bash
for f in migration/*.sql; do psql "postgres://USER:PASSWORD@localhost:5432/melodify" -f "$f"; done
```

Urutannya: `000_users.sql` (tabel users), `001_refresh_tokens.sql`, `002_songs.sql`, `003_playlists.sql`, `004_admin.sql` (kolom `is_admin`). Semua file idempotent (`IF NOT EXISTS`).

Jalankan server:

```bash
go run cmd/server/main.go
```

Backend berjalan di `http://localhost:8080`.

### 3. Setup Frontend

```bash
cd frontend
npm install
```

Salin `.env.example` menjadi `.env.local` di folder `frontend/`:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api
```

Jalankan dev server:

```bash
npm run dev
```

Frontend berjalan di `http://localhost:3000`.

### 4. Jadikan akun sebagai admin

Tidak ada endpoint publik untuk membuat admin (demi keamanan). Setelah register, jalankan query manual ke database:

```sql
UPDATE users SET is_admin = true WHERE email = 'email_kamu@contoh.com';
```

Login ulang setelah itu supaya sesi ter-refresh, lalu akses `/admin`.

---

## 📡 Ringkasan API

| Method | Endpoint                   | Auth     | Keterangan                  |
| ------ | -------------------------- | -------- | --------------------------- |
| POST   | `/api/auth/register`       | –        | Registrasi user baru        |
| POST   | `/api/auth/login`          | –        | Login                       |
| POST   | `/api/auth/refresh`        | –        | Refresh access token        |
| POST   | `/api/auth/logout`         | –        | Logout                      |
| GET    | `/api/me`                  | ✅       | Data user yang login        |
| PUT    | `/api/profile`             | ✅       | Update profil               |
| GET    | `/api/songs`               | ✅       | List/cari lagu (`?q=`, `?limit=` maks 200, `?offset=`) |
| GET    | `/api/playlists`           | ✅       | List playlist milik sendiri |
| POST   | `/api/playlists`           | ✅       | Buat playlist               |
| POST   | `/api/playlists/:id/songs` | ✅       | Tambah lagu ke playlist     |
| POST   | `/api/admin/songs`         | 🛡️ Admin | Tambah lagu                 |
| POST   | `/api/admin/upload`        | 🛡️ Admin | Upload file cover/audio     |
| GET    | `/api/admin/users`         | 🛡️ Admin | List semua user             |
| GET    | `/api/admin/playlists`     | 🛡️ Admin | List semua playlist         |

---

## 🧪 Testing

```bash
# Backend (auth, token, playlist ownership, admin rules, upload, middleware, config)
cd backend && go test ./...

# Frontend (player store, auth store, interceptor refresh token)
cd frontend && npm test
npm run typecheck
```

---

## 📌 Catatan Development

- Server tidak memakai `ReadTimeout`/`WriteTimeout` global (akan memutus upload 25 MB dan streaming audio di koneksi lambat). Perlindungannya: `ReadHeaderTimeout`, batas ukuran body per route (1 MB untuk JSON, 26 MB untuk upload), dan context timeout di handler.
- Lagu menyimpan path file upload **relatif** (`/uploads/audios/x.mp3`); backend mengubahnya menjadi URL absolut memakai `PUBLIC_BASE_URL` saat mengirim response. Pindah domain cukup mengubah variabel itu.
- Rate limiter bersifat in-memory (cukup untuk satu instance). Jika backend di-scale ke beberapa instance, pindahkan ke store bersama seperti Redis.
- File yang diupload lewat panel admin disimpan lokal di `backend/uploads/` — cukup untuk development, tapi untuk production sebaiknya dipindah ke object storage (Cloudflare R2 / AWS S3 / Supabase Storage) supaya tidak hilang saat redeploy.
- Refresh token disimpan sebagai cookie `httpOnly`, sehingga frontend dan backend perlu berjalan dengan origin yang terdaftar di `FRONTEND_ORIGIN` agar CORS berfungsi.

---

## 📄 Lisensi

Proyek ini dibuat untuk keperluan pembelajaran/portofolio.
