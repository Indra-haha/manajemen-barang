# API Manajemen Barang

Base URL: `http://localhost:8080`

Semua request dan response menggunakan `Content-Type: application/json`.

Semua error dikembalikan dengan bentuk yang sama:

```json
{ "error": "pesan kesalahan" }
```

## Cara menjalankan contoh curl

- **Salin baris `curl` yang satu baris saja.** Blok JSON berisi payload adalah referensi — jangan di-paste ke CMD (CMD tidak bisa memproses JSON multi-baris).
- Endpoint bertanda **Auth: Ya** memerlukan header `Authorization: Bearer <TOKEN>`. Ambil token dari response `POST /api/login`. Di CMD Windows, simpan dulu di variabel agar tidak menulis token berulang kali:

```
set TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9....
```

Kemudian gunakan `-H "Authorization: Bearer %TOKEN%"` pada tiap contoh curl.


## 1. Register

Membuat akun baru.

| | |
|---|---|
| **Method** | `POST` |
| **Path** | `/api/register` |
| **Auth** | Tidak perlu |

### Payload

```json
{
  "username": "admin2",
  "password": "rahasia123"
}
```

### Contoh curl

```bash
curl -X POST http://localhost:8080/api/register -H "Content-Type: application/json" -d "{\"username\":\"admin2\",\"password\":\"rahasia123\"}"
```

### Response `201 Created`

```json
{
  "id": 1,
  "username": "admin2"
}
```

### Error

| Status | Pesan | Penyebab |
|---|---|---|
| `400` | `body tidak valid` | JSON rusak / bukan object |
| `400` | `username atau password tidak valid` | Panjang username/password di luar batas |
| `409` | `username sudah digunakan` | Username sudah ada |
| `500` | `terjadi kesalahan pada server` | Gagal hash / error tak terduga |

---

## 2. Login

Menukar kredensial dengan JWT token.

| | |
|---|---|
| **Method** | `POST` |
| **Path** | `/api/login` |
| **Auth** | Tidak perlu |

### Payload

```json
{
  "username": "admin2",
  "password": "rahasia123"
}
```

### Contoh curl

```bash
curl -X POST http://localhost:8080/api/login -H "Content-Type: application/json" -d "{\"username\":\"admin2\",\"password\":\"rahasia123\"}"
```

### Response `200 OK`

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9....",
}
```

Token berlaku **24 jam** (`expires_in` dalam detik). Kirim token ini pada request ke endpoint yang dilindungi:

```
Authorization: Bearer <token>
```

### Error

| Status | Pesan | Penyebab |
|---|---|---|
| `400` | `body tidak valid` | JSON rusak |
| `401` | `username atau password salah` | Username tidak ada / password salah |
| `500` | `terjadi kesalahan pada server` | Error tak terduga |

---

## 3. Logout

Menutup sesi di sisi server (menghapus cookie `token` bila dikirim browser).

| | |
|---|---|
| **Method** | `POST` |
| **Path** | `/api/logout` |
| **Auth** | Tidak perlu |

### Contoh curl

```bash
curl -X POST http://localhost:8080/api/logout
```

### Response `200 OK`

```json
{
  "message": "logout berhasil"
}
```

> Catatan: JWT bersifat stateless — server tidak menyimpan daftar token aktif. Logout menghapus cookie pada browser; token di `Authorization` header tetap sah sampai masa berlakunya habis. Untuk pencabutan token sebelum kedaluwarsa, diperlukan mekanisme denylist/refresh.

---

## 4. List Barang

Mengambil semua barang, mendukung filter.

| | |
|---|---|
| **Method** | `GET` |
| **Path** | `/api/barang` |
| **Auth** | Ya (JWT) |

### Query Parameter

| Nama | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `kategori` | string | tidak | Filter exact match kategori |
| `q` | string | tidak | Pencarian `nama` atau `kode` (LIKE %...%) |

Contoh: `GET /api/barang?kategori=elektronik&q=samsung`

### Contoh curl

```bash
curl http://localhost:8080/api/barang -H "Authorization: Bearer %TOKEN%"


curl "http://localhost:8080/api/barang?kategori=elektronik&q=samsung" -H "Authorization: Bearer %TOKEN%"
```

### Response `200 OK`

```json
[
  {
    "id": 1,
    "kode": "BRG001",
    "nama": "Samsung A54",
    "kategori": "elektronik",
    "stok": 10,
    "harga": 5000000,
    "created_at": "2026-10-20T10:00:00Z",
    "updated_at": "2026-10-20T10:00:00Z"
  }
]
```

Urutan: `id` menurun (terbaru di atas). Jika tidak ada data, response berisi array kosong `[]`.

### Error

| Status | Pesan | Penyebab |
|---|---|---|
| `401` | `username atau password salah` | Header token kosong / tidak valid / kedaluwarsa |
| `500` | `gagal mengambil data` | Error database |

---

## 5. Detail Barang

| | |
|---|---|
| **Method** | `GET` |
| **Path** | `/api/barang/:id` |
| **Auth** | Ya (JWT) |

### Path Parameter

| Nama | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `id` | integer | ya | ID barang |

### Contoh curl

```bash
curl http://localhost:8080/api/barang/1 -H "Authorization: Bearer %TOKEN%"
```

### Response `200 OK`

```json
{
  "id": 1,
  "kode": "BRG001",
  "nama": "Samsung A54",
  "kategori": "elektronik",
  "stok": 10,
  "harga": 5000000,
  "created_at": "2026-10-20T10:00:00Z",
  "updated_at": "2026-10-20T10:00:00Z"
}
```

### Error

| Status | Pesan | Penyebab |
|---|---|---|
| `400` | `id tidak valid` | `id` bukan angka |
| `401` | `username atau password salah` | Header token kosong / tidak valid / kedaluwarsa |
| `404` | `barang tidak ditemukan` | ID tidak ada |
| `500` | `terjadi kesalahan pada server` | Error tak terduga |

---

## 6. Tambah Barang

| | |
|---|---|
| **Method** | `POST` |
| **Path** | `/api/barang` |
| **Auth** | Ya (JWT) |

### Payload

```json
{
  "kode": "BRG001",
  "nama": "Samsung A54",
  "kategori": "elektronik",
  "stok": 10,
  "harga": 5000000
}
```

`kode` dan `nama` wajib diisi. `kategori` opsional (default `umum`).

### Contoh curl

```bash
curl -X POST http://localhost:8080/api/barang -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"kode\":\"BRG001\",\"nama\":\"Samsung A54\",\"kategori\":\"elektronik\",\"stok\":10,\"harga\":5000000}"
```

### Response `201 Created`

Mengembalikan objek barang lengkap dengan `id` dan timestamp:

```json
{
  "id": 1,
  "kode": "BRG001",
  "nama": "Samsung A54",
  "kategori": "elektronik",
  "stok": 10,
  "harga": 5000000,
  "created_at": "2026-10-20T10:00:00Z",
  "updated_at": "2026-10-20T10:00:00Z"
}
```

### Error

| Status | Pesan | Penyebab |
|---|---|---|
| `400` | `body tidak valid` | JSON rusak |
| `400` | `data barang tidak valid` | `kode`/`nama` kosong |
| `400` | `stok tidak boleh negatif` | `stok < 0` |
| `400` | `harga tidak boleh negatif` | `harga < 0` |
| `400` | `data melebihi panjang maksimal` | Panjang field melebihi batas |
| `401` | `username atau password salah` | Header token kosong / tidak valid / kedaluwarsa |
| `409` | `kode barang sudah digunakan` | Duplikat `kode` |
| `500` | `terjadi kesalahan pada server` | Error tak terduga |

---

## 7. Ubah Barang

| | |
|---|---|
| **Method** | `PUT` |
| **Path** | `/api/barang/:id` |
| **Auth** | Ya (JWT) |

### Payload

Sama seperti tambah barang. `kode` wajib dikirim (tetapi tidak diubah di database), `id` diambil dari URL.

```json
{
  "kode": "BRG001",
  "nama": "Samsung A54 5G",
  "kategori": "elektronik",
  "stok": 8,
  "harga": 4500000
}
```

### Contoh curl

```bash
curl -X PUT http://localhost:8080/api/barang/1 -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"kode\":\"BRG001\",\"nama\":\"Samsung A54 5G\",\"kategori\":\"elektronik\",\"stok\":8,\"harga\":4500000}"
```

### Response `200 OK`

Mengembalikan objek terbaru dari database (termasuk `updated_at` yang sudah diperbarui):

```json
{
  "id": 1,
  "kode": "BRG001",
  "nama": "Samsung A54 5G",
  "kategori": "elektronik",
  "stok": 8,
  "harga": 4500000,
  "created_at": "2026-10-20T10:00:00Z",
  "updated_at": "2026-10-20T11:30:00Z"
}
```

> Catatan: field yang diperbarui adalah `nama`, `kategori`, `stok`, `harga`, dan `updated_at` (otomatis). `kode` dan `created_at` tidak berubah — karena `kode` tidak diubah di database, error `409` duplikat tidak mungkin terjadi di endpoint ini.

### Error

| Status | Pesan | Penyebab |
|---|---|---|
| `400` | `id tidak valid` | `id` bukan angka |
| `400` | `body tidak valid` | JSON rusak |
| `400` | `data barang tidak valid` | `kode`/`nama` kosong |
| `400` | `stok tidak boleh negatif` | `stok < 0` |
| `400` | `harga tidak boleh negatif` | `harga < 0` |
| `400` | `data melebihi panjang maksimal` | Panjang field melebihi batas |
| `401` | `username atau password salah` | Header token kosong / tidak valid / kedaluwarsa |
| `404` | `barang tidak ditemukan` | ID tidak ada |
| `500` | `terjadi kesalahan pada server` | Error tak terduga |

---

## 8. Penyesuaian Stok

Mengubah stok secara atomik (aman dari race condition): delta positif menambah, negatif mengurangi.

| | |
|---|---|
| **Method** | `PUT` |
| **Path** | `/api/barang/:id/stok` |
| **Auth** | Ya (JWT) |

### Payload

```json
{
  "jumlah": -3
}
```

| Nama | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `jumlah` | integer | ya | Delta stok; positif = tambah, negatif = kurangi, `0` tidak valid |

### Contoh curl

```bash
curl -X PUT http://localhost:8080/api/barang/1/stok -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"jumlah\":-3}"
```

### Response `200 OK`

Mengembalikan objek barang terbaru dari database:

```json
{
  "id": 1,
  "kode": "BRG001",
  "nama": "Samsung A54",
  "kategori": "elektronik",
  "stok": 7,
  "harga": 5000000,
  "created_at": "2026-10-20T10:00:00Z",
  "updated_at": "2026-10-20T12:00:00Z"
}
```

### Error

| Status | Pesan | Penyebab |
|---|---|---|
| `400` | `id tidak valid` | `id` bukan angka |
| `400` | `body tidak valid` | JSON rusak |
| `400` | `data barang tidak valid` | `jumlah` = 0 |
| `401` | `username atau password salah` | Header token kosong / tidak valid / kedaluwarsa |
| `404` | `barang tidak ditemukan` | ID tidak ada |
| `409` | `stok tidak mencukupi` | Pengurangan membuat stok < 0 |
| `500` | `terjadi kesalahan pada server` | Error tak terduga |

---

## 9. Hapus Barang

| | |
|---|---|
| **Method** | `DELETE` |
| **Path** | `/api/barang/:id` |
| **Auth** | Ya (JWT) |

### Contoh curl

```bash
curl -X DELETE http://localhost:8080/api/barang/1 -H "Authorization: Bearer %TOKEN%"
```

### Response `204 No Content`

Tanpa body.

### Error

| Status | Pesan | Penyebab |
|---|---|---|
| `400` | `id tidak valid` | `id` bukan angka |
| `401` | `username atau password salah` | Header token kosong / tidak valid / kedaluwarsa |
| `404` | `barang tidak ditemukan` | ID tidak ada |
| `500` | `terjadi kesalahan pada server` | Error tak terduga |

---

## Ringkasan Endpoint

| Method | Path | Auth | Keterangan |
|---|---|---|---|
| `POST` | `/api/register` | Tidak | Buat akun |
| `POST` | `/api/login` | Tidak | Dapat token |
| `POST` | `/api/logout` | Tidak | Hapus cookie sesi |
| `GET` | `/api/barang` | Ya | List + filter `kategori`, `q` |
| `POST` | `/api/barang` | Ya | Tambah barang |
| `GET` | `/api/barang/:id` | Ya | Detail barang |
| `PUT` | `/api/barang/:id` | Ya | Ubah barang |
| `PUT` | `/api/barang/:id/stok` | Ya | Penyesuaian stok atomik |
| `DELETE` | `/api/barang/:id` | Ya | Hapus barang |
