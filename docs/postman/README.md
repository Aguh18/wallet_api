# Wallet API - Postman Collection

Postman collection untuk testing Wallet API endpoints dengan lengkap.

## 🚀 Cara Import

1. **Install Postman** dari [postman.com](https://www.postman.com/downloads/)
2. **Buka Postman** aplikasinya
3. **Klik Import** → **"Upload Files"**
4. **Pilih file**: `docs/postman/Wallet_API.postman_collection.json`
5. Selesai! Collection akan terimport otomatis

## ⚙️ Setup Environment & Variables

Postman collection ini menggunakan variables untuk memudahkan testing:

### Global Variables (Otomatis)
- `base_url` - Base URL API (default: http://localhost:8080)
- `access_token` - JWT access token (set otomatis setelah login)
- `refresh_token` - JWT refresh token (set otomatis setelah login)
- `wallet_id` - Wallet ID (set otomatis setelah create wallet)
- `to_wallet_id` - Target wallet ID untuk transfer (manual)

### Setup Manual

Setelah import collection:

1. Klik collection **"Wallet API"**
2. Pilih tab **Variables**
3. Pastikan `base_url` sesuai dengan environment Anda:
   - Local: `http://localhost:8080`
   - Development: sesuaikan
   - Production: sesuaikan

## 🧪 Cara Testing

### Langkah 1: Authentication

**1. Register** → Buat akun baru
```json
{
  "username": "testuser",
  "email": "testuser@example.com",
  "password": "password123"
}
```

**2. Login** → Dapatkan auth cookies
```json
{
  "username": "testuser",
  "password": "password123"
}
```

**Catatan**: Setelah Login berhasil, cookies akan otomatis di-set ke variables.

### Langkah 2: Buat Wallet

**3. Create Wallet** → Buat wallet pertama
```json
{
  "wallet_name": "Dompet Utama",
  "currency": "IDR"
}
```

**Catatan**: Wallet ID akan otomatis disimpan ke variable `wallet_id`.

### Langkah 3: Transaksi

**4. Deposit** → Tambah saldo (amount sebagai string dengan desimal)
```json
{
  "amount": "100000.50",
  "description": "Deposit awal"
}
```

**5. Get All Wallets** → Lihat semua wallet Anda

**6. Get Wallet by ID** → Lihat detail wallet

**7. Withdraw** → Tarik dana (amount sebagai string)
```json
{
  "amount": "50000.25",
  "description": "Tarik tunai"
}
```

**8. Transfer** → Kirim ke wallet lain (perlu 2 wallet)

Untuk test transfer:
1. Buat wallet kedua (Register → Login → Create Wallet baru)
2. Copy wallet ID kedua
3. Set ke variable `to_wallet_id`
4. Jalankan request Transfer

```json
{
  "to_wallet_id": "{{to_wallet_id}}",
  "amount": "25000.75",
  "description": "Transfer untuk jajan"
}
```

**9. Get Transactions** → Lihat riwayat transaksi
- Query params: `limit=10`, `offset=0`

## 💡 Perbedaan dengan Bruno

| Fitur | Bruno | Postman |
|-------|-------|---------|
| Auto-set cookies | Manual dari headers | Manual otomatis via script |
| Variables | Manual | Semi-otomatis (test scripts) |
| Response validation | Manual | Otomatis (test scripts) |
| Environment | Built-in environments | Collection variables |
| Interface | Lebih simpel | Lebih feature-rich |

## 📝 Format Amount (PENTING)

Semua request dengan **amount** harus menggunakan **string** dengan format desimal:

✅ **BENAR:**
```json
{
  "amount": "100000.50"
}
```

❌ **SALAH:**
```json
{
  "amount": 100000.50
}
```

Ini menggunakan `shopspring/decimal` untuk presisi monetik exact!

## 🔧 Test Scripts

Collection ini dilengkapi dengan test scripts untuk:

- **Register**: Validasi status 201 dan response structure
- **Login**: Validasi status 200 dan cookies
- **Create Wallet**: Validasi response dan auto-set `wallet_id`

Jalankan tests dengan:
1. Klik collection atau request
2. Klik tab **Tests**
3. Klik **Send**
4. Lihat hasil test di tab **Test Results**

## ❓ Troubleshooting

### "401 Unauthorized"?
**Solusi:**
1. Pastikan sudah **Login** request
2. Cookies otomatis di-set ke variables
3. Coba lagi request yang gagal

### "Wallet not found"?
**Solusi:**
1. Pastikan sudah **Create Wallet**
2. `wallet_id` variable otomatis di-set
3. Bisa juga manual set di tab Variables

### "Invalid amount format"?
**Solusi:**
1. Pastikan amount dikirim sebagai **string**
2. Gunakan format desimal: "100.50", "1000.00", "50"

### Transfer gagal?
**Solusi:**
1. Pastikan punya 2 wallet berbeda
2. Set `to_wallet_id` dengan ID wallet kedua
3. Pastikan wallet pertama cukup saldo

---

**Selamat Testing dengan Postman! 🚀**