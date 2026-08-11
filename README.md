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

| | |
|---|---|
| **Backend** | Go, Gin, PostgreSQL (pgx), JWT, bcrypt, zap logger |
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
│   ├── migration/            # file SQL migrasi (jalankan berurutan)
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
git clone <url-repo-ini>
cd spotify
```

### 2. Setup Backend

```bash
cd backend
```

Buat file `.env` di folder `backend/`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=melodify
JWT_SECRET=ganti_dengan_secret_key_yang_panjang_dan_acak
SERVER_PORT=8080
FRONTEND_ORIGIN=http://localhost:3000
```

Jalankan seluruh file migrasi secara **berurutan** ke database PostgreSQL kamu (lewat `psql`, DBeaver, atau tools lain):

```
migration/001_refresh_tokens.sql
migration/002_songs.sql
migration/003_playlists.sql
migration/004_admin.sql
```

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

Buat file `.env.local` di folder `frontend/`:

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

| Method | Endpoint | Auth | Keterangan |
|---|---|---|---|
| POST | `/api/auth/register` | – | Registrasi user baru |
| POST | `/api/auth/login` | – | Login |
| POST | `/api/auth/refresh` | – | Refresh access token |
| POST | `/api/auth/logout` | – | Logout |
| GET | `/api/me` | ✅ | Data user yang login |
| PUT | `/api/profile` | ✅ | Update profil |
| GET | `/api/songs` | ✅ | List/cari lagu (`?q=`) |
| GET | `/api/playlists` | ✅ | List playlist milik sendiri |
| POST | `/api/playlists` | ✅ | Buat playlist |
| POST | `/api/playlists/:id/songs` | ✅ | Tambah lagu ke playlist |
| POST | `/api/admin/songs` | 🛡️ Admin | Tambah lagu |
| POST | `/api/admin/upload` | 🛡️ Admin | Upload file cover/audio |
| GET | `/api/admin/users` | 🛡️ Admin | List semua user |
| GET | `/api/admin/playlists` | 🛡️ Admin | List semua playlist |

---

## 📌 Catatan Development

- File yang diupload lewat panel admin disimpan lokal di `backend/uploads/` — cukup untuk development, tapi untuk production sebaiknya dipindah ke object storage (Cloudflare R2 / AWS S3 / Supabase Storage) supaya tidak hilang saat redeploy.
- Refresh token disimpan sebagai cookie `httpOnly`, sehingga frontend dan backend perlu berjalan dengan origin yang terdaftar di `FRONTEND_ORIGIN` agar CORS berfungsi.

---

## 📄 Lisensi

Proyek ini dibuat untuk keperluan pembelajaran/portofolio.