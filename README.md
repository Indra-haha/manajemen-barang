# Manajemen Barang API

REST API manajemen inventaris barang yang dibangun dengan **Go + Gin**, **PostgreSQL 17**, autentikasi **JWT (HS256)**, dan dijalankan lewat **Docker Compose** (nginx + app + postgres). Dokumentasi endpoint lengkap ada di [`API.md`](API.md).

## Tech Stack

| Komponen | Teknologi |
|---|---|
| Bahasa | Go 1.25 |
| HTTP Framework | Gin |
| Database | PostgreSQL 17 (`pgx/v5` pgxpool, SQL mentah tanpa ORM) |
| Autentikasi | JWT HS256 (`golang-jwt/jwt/v5`) + bcrypt |
| Reverse Proxy | nginx 1.27 (port publik `8080`) |
| Deploy | Docker Compose (multi-stage build, image runtime Alpine non-root) |

## Menjalankan dengan Docker

### Prasyarat

- Docker & Docker Compose

### Langkah

1. Siapkan file `.env` di root (docker compose otomatis membacanya):

   ```env
   POSTGRES_USER=barang
   POSTGRES_PASSWORD=barang
   POSTGRES_DB=manajemen_barang

   DB_HOST=postgres
   DB_PORT=5432
   DB_SSLMODE=disable
   JWT_SECRET=ganti-dengan-secret-acak-minimal-32-karakter
   ```

   > `JWT_SECRET` wajib diisi — server akan **panic saat startup** bila kosong. Gunakan secret acak yang panjang (mis. hasil `openssl rand -hex 32`).

2. Build dan jalankan seluruh service:

   ```bash
   docker compose up --build
   ```

3. API siap diakses di **`http://localhost:8080`**.

### Arsitektur Container

```
Client ──▶ nginx (host:8080 → container:80)
              │  proxy_pass
              ▼
           app (internal :8080, Go server)
              │  pgxpool
              ▼
           postgres (:5432, tidak dipublikasikan ke host)
```

- Migrasi skema (`migrations/schema.sql`) dijalankan **otomatis** oleh entrypoint Postgres saat volume `postgres_data` pertama kali dibuat.
- Untuk reset database total: `docker compose down -v` lalu `docker compose up --build`.

### Perintah Berguna

| Perintah | Fungsi |
|---|---|
| `docker compose up --build` | Build & jalankan semua service |
| `docker compose down` | Hentikan container (data tetap tersimpan) |
| `docker compose down -v` | Hentikan & hapus volume (reset database) |
| `docker compose logs -f app` | Pantau log aplikasi |

### Menjalankan Tanpa Docker (Opsional)

Butuh PostgreSQL berjalan sendiri, lalu:

```bash
export DB_HOST=localhost DB_PORT=5432 DB_USER=barang DB_PASS=barang DB_NAME=manajemen_barang DB_SSLMODE=disable
export JWT_SECRET=rahasia
export PORT=8080
go run ./cmd/api
```

## Alur Request

### Alur Umum Setiap Request

```
1. Client mengirim HTTP request ke http://localhost:8080
        │
2. nginx (port 80 di container) menerima & meneruskan
   ke http://app:8080 dengan header X-Real-IP,
   X-Forwarded-For, X-Forwarded-Proto
        │
3. Gin router (internal/router/router.go) mencocokkan
   method + path ke handler yang sesuai
        │
4. [Untuk route /api/barang/*] Middleware AuthRequired
   mengecek JWT:
     - Ambil token: cookie "token" lebih dulu,
       lalu fallback header "Authorization: Bearer <token>"
     - Validasi signature HS256 + klaim exp
     - Gagal → 401 {"error": "..."} + Abort
     - Berhasil → user_id & username disimpan ke context
        │
5. Handler memvalidasi payload (binding + validasi bisnis),
   memanggil repository
        │
6. Repository mengeksekusi SQL mentah via pgxpool
   ke PostgreSQL
        │
7. Response dikirim sebagai JSON
   (sukses: objek/array; gagal: {"error": "pesan"})
```

### Alur Login → Akses Terproteksi

```
POST /api/login  { "username", "password" }
   → bcrypt.CompareHashAndPassword
   → JWT HS256 ditandatangani (klaim: sub, username, iat, exp — TTL 24 jam)
   → token dikirim DUA cara sekaligus:
       a. Set-Cookie: token=...; HttpOnly; SameSite=Lax; Path=/
       b. body JSON: { "token": "..." }

Client memilih salah satu untuk request berikutnya:
   - Browser: cookie dikirim otomatis oleh browser
   - Non-browser (curl/mobile): header Authorization: Bearer <token>
```

### Ringkasan Endpoint

| Method | Path | Auth | Fungsi |
|---|---|---|---|
| POST | `/api/register` | – | Registrasi user baru |
| POST | `/api/login` | – | Login, terima token |
| GET | `/api/barang` | JWT | Daftar barang (`?kategori=&q=`) |
| POST | `/api/barang` | JWT | Tambah barang |
| GET | `/api/barang/:id` | JWT | Detail barang |
| PUT | `/api/barang/:id` | JWT | Ubah data barang |
| PUT | `/api/barang/:id/stok` | JWT | Penyesuaian stok (delta atomik) |
| DELETE | `/api/barang/:id` | JWT | Hapus barang |

Contoh pemakaian curl lengkap ada di [`API.md`](API.md).

## Keamanan Token

### Mekanisme yang Digunakan: **httpOnly Cookie + Bearer (dual mode)**

Sistem ini **tidak memakai localStorage**. Saat login (`internal/handler/auth_handler.go`), token dikirim dua cara sekaligus:

1. **httpOnly Cookie** (`cookie name: token`)
   - `HttpOnly = true` — cookie **tidak bisa dibaca/diubah oleh JavaScript** di browser, sehingga token tidak mudah dicuri skrip XSS.
   - `SameSite = Lax` — mitigasi CSRF untuk request lintas situs.
   - `Secure` dikontrol env `COOKIE_SECURE` (default `false`; **set `true` di production dengan HTTPS**).
   - Path `/`, TTL sama dengan token (24 jam).

2. **JSON body** (`{ "token": "..." }`)
   - Untuk client non-browser (mobile, Postman, curl, service pihak ketiga) yang tidak punya mekanisme cookie.

Middleware (`internal/handler/middleware.go`) mengestrak token dengan urutan: **cookie `token` dulu, lalu header `Authorization: Bearer`** — jadi browser cukup kirim request biasa (cookie otomatis ikut), sedangkan non-browser memakai header.

### Kenapa httpOnly Cookie, bukan localStorage?

| Aspek | httpOnly Cookie (dipakai) | localStorage |
|---|---|---|
| Dibaca JavaScript | ❌ Tidak (mitigasi XSS) | ✅ Ya — rentan dicuri skrip |
| Kirim otomatis oleh browser | ✅ Ya | ❌ Manual (rawan lupa/luput) |
| Risiko CSRF | ⚠️ Ada (ditangani SameSite=Lax) | Tidak ada |
| Cocok untuk non-browser | ⚠️ Perlu header manual | Tidak relevan |

Karena itu sistem memberikan **keduanya**: browser memakai cookie yang aman dari XSS, non-browser memakai Bearer header.

### Detail Keamanan Lain

- **JWT HS256**, TTL **24 jam**, klaim: `sub` (user ID), `username`, `iat`, `exp`.
- Middleware memaksa algoritma `HS256` saja (`WithValidMethods`) — menutup serangan *algorithm confusion* (mis. downgrade ke `none`).
- **bcrypt** (`DefaultCost = 10`) untuk hash password; field hash tidak pernah ikut di-serialize JSON (`json:"-"`).
- Username/password salah → **401 sama persis** (tidak membocorkan keberadaan akun).
- `JWT_SECRET` kosong → server menolak start (panic).
- Pesan error internal database **tidak pernah bocor** ke client (selalu 500 generik).
- Container app berjalan sebagai **non-root user**; Postgres tidak dipublikasikan ke host.

### Keterbatasan (belum ada)

- **Stateless** — tidak ada denylist/refresh token, sehingga logout global/pencabutan token sebelum 24 jam belum bisa dilakukan (token tetap valid sampai `exp`).
- Belum ada rate limiting, CORS policy, dan endpoint `/api/logout` (terdokumentasi di API.md tetapi belum didefinisikan di kode).

## Penanganan Konkurensi

Sistem ini tidak memakai transaksi multi-statement atau row locking eksplisit. Keamanan dari **race condition** (khususnya pada stok) ditopang tiga mekanisme:

### 1. UPDATE Atomik Satu-Statement (Anti Lost-Update)

Intinya ada di `internal/repository/barang.go` — penyesuaian stok adalah operasi **delta** (`+n` / `-n`), bukan baca-lalu-tulis:

```sql
UPDATE barang
SET stok = stok + $2, updated_at = now()
WHERE id = $1 AND stok + $2 >= 0
```

Pengecekan **"stok tidak boleh negatif" menjadi bagian dari WHERE clause perintah UPDATE itu sendiri**. PostgreSQL mengevaluasi dan mengunci barisnya secara atomik per-statement, sehingga:

- Dua request penarikan stok bersamaan **tidak akan pernah** menghasilkan stok negatif — walau keduanya lolos cek di aplikasi, hanya kombinasi yang valid yang bisa diterapkan database.
- **Tidak ada lost update** — tidak ada pola `SELECT` stok → hitung di Go → `UPDATE` nilai absolut, yang bisa saling menimpa.

Jika `RowsAffected() == 0`, repository membedakan dua kemungkinan: barang tidak ada (→ 404) vs stok tidak mencukupi (→ 409 `stok tidak mencukupi`).

> **Catatan klien:** gunakan `PUT /api/barang/:id/stok` (delta) untuk operasi stok. `PUT /api/barang/:id` menulis `stok` secara absolut (last-write-wins) dan bukan jalur yang aman untuk stok bersamaan.

### 2. Constraint CHECK di Database

Sebagai jaring pengaman terakhir di level skema (`migrations/schema.sql`):

```sql
stok INTEGER NOT NULL DEFAULT 0 CHECK (stok >= 0),
harga BIGINT NOT NULL DEFAULT 0 CHECK (harga >= 0),
```

Bahkan bila ada jalur tulis lain yang lolos validasi aplikasi, database tetap menolak nilai negatif.

### 3. UNIQUE Constraint + Deteksi Duplikat

- `kode barang` dan `username` dideklarasikan `UNIQUE` di database.
- Saat insert, duplikat terdeteksi via error code PostgreSQL **23505** (`isDuplicate` di `internal/repository/barang.go`) dan dipetakan ke respons **409 Conflict** (`kode barang sudah digunakan` / `username sudah digunakan`).

Jadi otoritas keputusan akhir selalu dipegang database, bukan urutan arrival request dari aplikasi.

## Konfigurasi Environment Variable

| Var | Default | Keterangan |
|---|---|---|
| `PORT` | `8080` | Port internal server Go |
| `JWT_SECRET` | *(wajib)* | Kunci tanda tangan JWT; kosong = panic |
| `COOKIE_SECURE` | `false` | Flag `Secure` pada cookie; `true` saat production HTTPS |
| `POSTGRES_USER` / `POSTGRES_PASSWORD` / `POSTGRES_DB` | – | Kredensial & nama database Postgres |
| `DB_HOST` | – | Host database (di compose: `postgres`) |
| `DB_PORT` | – | Port database (`5432`) |
| `DB_SSLMODE` | – | `disable` untuk lokal |
| `POSTGRES_URI` | *(opsional)* | DSN penuh; bila diisi, menggantikan `DB_*` |
