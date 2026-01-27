# Wallet API - Postman Collection

Postman collection untuk testing Wallet API endpoints dengan lengkap.

## 🚀 Cara Import

1. **Install Postman** dari [postman.com](https://www.postman.com/downloads/)
2. **Buka Postman** aplikasinya
3. **Klik Import** → **"Upload Files"**
4. **Pilih file**: `docs/postman/Wallet_API.postman_collection.json`
5. Selesai! Collection akan terimport otomatis

## ⚙️ Setup Variables

| Variable | Value | Auto-Set? |
|----------|-------|-----------|
| `base_url` | `http://localhost:8080` | Manual (set 1x) |
| `access_token` | - | ✨ Auto dari Login |
| `refresh_token` | - | ✨ Auto dari Login |
| `wallet_id` | - | ✨ Auto dari Create Wallet |
| `to_wallet_id` | - | Manual (untuk transfer) |

**Setup:**
1. Import collection ke Postman
2. Klik collection → tab **Variables**
3. Set `base_url` = `http://localhost:8080`
4. Jalankan **Login** → tokens otomatis tersimpan
5. Jalankan **Create Wallet** → wallet_id otomatis tersimpan

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

## ❓ Troubleshooting

**401 Unauthorized?** → Jalankan **Login** dulu (tokens auto-set)

**Wallet not found?** → Jalankan **Create Wallet** dulu (wallet_id auto-set)

**Invalid amount format?** → Pastikan amount pakai **string**: `"100.50"`

**Transfer gagal?** → Pastikan punya 2 wallet, set `to_wallet_id` manual

---

**Selamat Testing! 🚀**