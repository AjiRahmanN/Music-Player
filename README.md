# Melodify

Melodify adalah aplikasi music player berbasis web dengan arsitektur **full-stack**:

- **Frontend:** Next.js 16 + React 19 + TypeScript + Tailwind CSS + Zustand
- **Backend:** Go 1.24.5 + Gin
- **Database:** PostgreSQL
- **Authentication:** JWT access token + refresh token melalui cookie
- **HTTP client:** Axios

Project ini menyediakan autentikasi pengguna, katalog/pencarian lagu, audio player, profile pengguna, dan playlist pribadi. Berdasarkan source yang dipindai, halaman frontend menggunakan branding **Melodify - Your Music Library** dan backend menggunakan module `music-player-backend`. fileciteturn2file3L410-L441

## Fitur

### Authentication
- Register
- Login
- Logout
- Refresh access token
- Protected route dengan JWT
- Profile/current-user
- Update profile

Backend menggunakan middleware untuk membaca `Authorization: Bearer <token>`, memvalidasi JWT, memeriksa token type `access`, lalu memasukkan `user_id` ke request context. fileciteturn2file0L24-L25 fileciteturn2file0L80-L105

### Music
- Daftar lagu
- Detail lagu
- Pencarian berdasarkan judul atau artis
- Pemutaran audio
- Progress / seek
- Volume
- Repeat
- Queue / shuffle

Frontend memiliki `songService` untuk mengambil daftar/detail lagu dan `PlayerAudioEngine` yang menjembatani state player dengan elemen HTML `<audio>`. fileciteturn2file5L626-L659 fileciteturn2file11L1259-L1338

### Playlist
- Membuat playlist
- Melihat playlist milik user
- Melihat detail playlist
- Menghapus playlist
- Menambahkan lagu ke playlist
- Menghapus lagu dari playlist

Backend juga memeriksa kepemilikan playlist sebelum operasi detail/delete/add/remove dilakukan. fileciteturn3file4L452-L481 fileciteturn3file4L484-L523

## Tech Stack

| Bagian | Teknologi |
|---|---|
| Frontend | Next.js 16.1.1 |
| UI | React 19.2.3 |
| Language | TypeScript |
| Styling | Tailwind CSS 4 |
| State management | Zustand 5 |
| HTTP client | Axios |
| Icons | Lucide React |
| Backend | Go 1.24.5 |
| Web framework | Gin 1.11 |
| PostgreSQL driver | pgx/v5 |
| Authentication | golang-jwt/jwt/v5 |
| Password hashing | `golang.org/x/crypto` |
| Logging | Zap |
| Environment | godotenv |

Versi frontend dan dependencies tersebut berasal dari `frontend/package.json`. fileciteturn2file3L410-L441 Backend mendeklarasikan Go 1.24.5 serta Gin, JWT, UUID, pgx, crypto, dan dependency terkait di `go.mod`. fileciteturn1file1L226-L234

## Arsitektur

```text
┌──────────────────────┐
│      Next.js         │
│      Frontend        │
│                      │
│  App Router          │
│  Zustand             │
│  Axios               │
│  HTML Audio Player   │
└──────────┬───────────┘
           │ HTTP / JSON
           │ Bearer JWT
           ▼
┌──────────────────────┐
│       Go + Gin       │
│       Backend        │
│                      │
│ Routes               │
│ Middleware            │
│ Handlers             │
│ Services             │
│ Repositories         │
└──────────┬───────────┘
           │ pgx
           ▼
┌──────────────────────┐
│     PostgreSQL       │
│                      │
│ users                │
│ refresh_tokens       │
│ songs                │
│ playlists            │
│ playlist_songs       │
└──────────────────────┘
```

Backend menggunakan pemisahan **handlers → services → repositories → PostgreSQL**, dengan dependency injection dilakukan ketika router dibuat. fileciteturn3file0L36-L51

## Struktur Project

Struktur berikut disusun dari hasil scan project yang diberikan:

```text
spotify/
├── backend/
├── frontend/
├── scan_project.py
├── README.md
└── .gitignore
```

Struktur utama backend/frontend yang terdeteksi:

```text
├── backend/.env
├── backend/cmd/server/main.go
├── backend/go.mod
├── backend/internal/config/config.go
├── backend/internal/config/database.go
├── backend/internal/errors/errors.go
├── backend/internal/http/handlers/auth_handler.go
├── backend/internal/http/handlers/health_handler.go
├── backend/internal/http/handlers/playlist_handler.go
├── backend/internal/http/handlers/song_handler.go
├── backend/internal/http/handlers/user_handler.go
├── backend/internal/http/middleware/auth_middleware.go
├── backend/internal/http/routes/router.go
├── backend/internal/logger/logger.go
├── backend/internal/models/auth_request.go
├── backend/internal/models/playlist.go
├── backend/internal/models/refresh_token.go
├── backend/internal/models/song.go
├── backend/internal/models/user.go
├── backend/internal/repositories/playlist_repository.go
├── backend/internal/repositories/refresh_token_repository.go
├── backend/internal/repositories/song_repository.go
├── backend/internal/repositories/user_repository.go
├── backend/internal/services/auth_service.go
├── backend/internal/services/playlist_service.go
├── backend/internal/services/song_service.go
├── backend/internal/services/user_service.go
├── backend/internal/utils/cookie.go
├── backend/internal/utils/jwt.go
├── backend/internal/utils/password.go
├── backend/internal/utils/token.go
├── backend/migration/001_refresh_tokens.sql
├── backend/migration/002_songs.sql
├── backend/migration/003_playlists.sql
├── frontend/README.md
├── frontend/eslint.config.mjs
├── frontend/next-env.d.ts
├── frontend/next.config.ts
├── frontend/package.json
├── frontend/postcss.config.mjs
├── frontend/public/file.svg
├── frontend/public/globe.svg
├── frontend/public/next.svg
├── frontend/public/vercel.svg
├── frontend/public/window.svg
├── frontend/src/app/(auth)/login/page.tsx
├── frontend/src/app/(auth)/register/page.tsx
├── frontend/src/app/(main)/layout.tsx
├── frontend/src/app/(main)/page.tsx
├── frontend/src/app/(main)/test-player/page.tsx
├── frontend/src/app/globals.css
├── frontend/src/app/layout.tsx
├── frontend/src/app/profile/page.tsx
├── frontend/src/components/auth/AuthGuard.tsx
├── frontend/src/components/layout/Navbar.tsx
├── frontend/src/components/player/PlayerAudioEngine.tsx
├── frontend/src/components/player/PlayerBar.tsx
├── frontend/src/components/player/TrackList.tsx
├── frontend/src/lib/axios.ts
├── frontend/src/services/authService.ts
├── frontend/src/services/songService.ts
├── frontend/src/store/useAuthStore.ts
├── frontend/src/store/usePlayerStore.ts
├── frontend/src/types/auth.ts
├── frontend/src/types/player.ts
├── frontend/src/types/song.ts
├── frontend/tsconfig.json
├── frontend/tsconfig.tsbuildinfo
```

> Catatan: hasil scan adalah snapshot source code, sehingga struktur di atas hanya mencantumkan file yang terdeteksi oleh scanner dan tidak mengasumsikan file yang tidak muncul di hasil scan.

## Backend API

Base URL default saat development:

```text
http://localhost:8080/api
```

### Public endpoints

| Method | Endpoint | Keterangan |
|---|---|---|
| GET | `/health` | Health check backend + database |
| GET | `/ping` | Simple ping |
| POST | `/auth/register` | Registrasi user |
| POST | `/auth/login` | Login |
| POST | `/auth/refresh` | Refresh access token |
| POST | `/auth/logout` | Logout |
| GET | `/auth/health` | Health check authentication |

Route tersebut didefinisikan sebagai public route di router Gin. fileciteturn3file0L53-L67

### Protected endpoints

Semua endpoint berikut membutuhkan access token:

| Method | Endpoint | Keterangan |
|---|---|---|
| GET | `/me` | Current user |
| GET | `/profile` | Profile user |
| PUT | `/profile` | Update profile |
| GET | `/songs` | Daftar lagu |
| GET | `/songs?q=keyword` | Cari lagu berdasarkan judul/artis |
| GET | `/songs/:id` | Detail lagu |
| POST | `/playlists` | Buat playlist |
| GET | `/playlists` | Daftar playlist user |
| GET | `/playlists/:id` | Detail playlist |
| DELETE | `/playlists/:id` | Hapus playlist |
| POST | `/playlists/:id/songs` | Tambahkan lagu |
| DELETE | `/playlists/:id/songs/:songId` | Hapus lagu |

Route protected dan endpoint playlist/song tersebut terlihat langsung pada router backend. fileciteturn3file3L353-L378

## Database

Project menggunakan PostgreSQL melalui `pgxpool`. Backend membuat connection pool, melakukan `Ping`, dan menutup pool saat server shutdown. fileciteturn2file6L779-L795

Dari repository dan migration yang terdeteksi, data utama mencakup:

- `users`
- `refresh_tokens`
- `songs`
- `playlists`
- `playlist_songs`

Relasi playlist menggunakan `user_id` dan tabel penghubung `playlist_songs`; penghapusan playlist/song terkait menggunakan foreign key cascade pada migration playlist. fileciteturn2file1L154-L174

## Environment Variables

Backend membaca konfigurasi dari environment variable, termasuk:

```env
ENV=development
PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_database_password
DB_NAME=musicp

JWT_SECRET=your_secure_random_secret
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=720h

COOKIE_SECURE=false
COOKIE_DOMAIN=localhost
COOKIE_SAME_SITE=Lax

FRONTEND_ORIGIN=http://localhost:3000
```

**Jangan commit file `.env`.** Hasil scan menunjukkan bahwa repository lokal saat ini memiliki `backend/.env` yang berisi credential database dan JWT secret. fileciteturn1file1L201-L223

Jika secret tersebut pernah masuk Git history atau repository remote, **rotate JWT secret dan credential database** sebelum project dipublikasikan.

> Nilai `COOKIE_*` dan `FRONTEND_ORIGIN` di atas adalah contoh konfigurasi development. Sesuaikan dengan konfigurasi deployment sebenarnya.

## Menjalankan Backend

Masuk ke folder backend:

```bash
cd backend
```

Download dependency:

```bash
go mod download
```

Siapkan `.env`:

```bash
cp .env.example .env
```

Windows PowerShell:

```powershell
Copy-Item .env.example .env
```

Jalankan server:

```bash
go run ./cmd/server
```

Backend membaca `.env` saat environment bukan `production`, kemudian melakukan load config, inisialisasi logger, koneksi PostgreSQL, setup router, dan menjalankan HTTP server. fileciteturn2file6L757-L786

Default port pada konfigurasi project adalah `8080`.

Tes:

```text
GET http://localhost:8080/api/ping
```

## Menjalankan Frontend

Masuk ke folder frontend:

```bash
cd frontend
```

Install dependency:

```bash
npm install
```

Jalankan development server:

```bash
npm run dev
```

Buka:

```text
http://localhost:3000
```

Script yang tersedia pada frontend:

```bash
npm run dev
npm run build
npm run start
npm run lint
```

Script tersebut sesuai dengan `frontend/package.json`. fileciteturn2file3L410-L440

## Alur Authentication

```text
Register
   │
   ▼
POST /api/auth/register
   │
   ▼
User dibuat di PostgreSQL
   │
   ▼
Login
   │
   ▼
POST /api/auth/login
   │
   ├── Access Token
   │
   └── Refresh Token → cookie + database
             │
             ▼
        Protected API
             │
             ▼
     Authorization: Bearer <access>
```

Saat access token tidak valid/expired, frontend memiliki interceptor Axios yang mencoba refresh token sebelum mengarahkan user kembali ke `/login`. fileciteturn2file10L1144-L1163

## Security Notes

Beberapa hal penting sebelum repository dipublikasikan:

1. **Jangan commit `.env`.**
2. Gunakan JWT secret yang random dan berbeda untuk setiap environment.
3. Jangan memasukkan password database ke source code.
4. Pastikan konfigurasi cookie production menggunakan HTTPS.
5. Review CORS `FRONTEND_ORIGIN` sebelum deployment.
6. Jika `.env` pernah ter-commit, menghapus file dari working tree saja tidak cukup; secret perlu di-rotate dan Git history perlu dibersihkan bila diperlukan.
7. Jangan memasukkan hasil scan project yang berisi source code/configuration sensitif ke repository publik.

## Development

Frontend menggunakan Next.js App Router. Authentication state dan player state dikelola dengan Zustand. fileciteturn2file5L664-L714

Backend menerapkan beberapa lapisan:

```text
HTTP Request
    ↓
Gin Router
    ↓
Middleware
    ↓
Handler
    ↓
Service
    ↓
Repository
    ↓
PostgreSQL
```

Pendekatan ini membuat business logic tidak langsung ditempatkan di HTTP handler dan repository bertanggung jawab terhadap akses data.

## Status

Project saat ini memiliki fondasi full-stack music player dengan:

- Authentication
- JWT access/refresh token
- PostgreSQL
- Song catalog
- Song search
- Audio player
- User profile
- Personal playlists
- Protected API
- Health checks
- Structured backend logging

README ini dibuat berdasarkan hasil scan source project yang diberikan, bukan berdasarkan asumsi terhadap fitur yang tidak ditemukan di scan.
