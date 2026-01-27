# Wallet API - Bruno Collection

Bruno API collection untuk testing Wallet API endpoints.

## 🚀 Cara Import

1. **Install Bruno** dari [usebruno.com](https://www.usebruno.com/downloads)
2. **Buka Bruno** aplikasinya
3. **Klik Import** → **"Import Folder"**
4. **Pilih folder**: `docs/api/`
5. Selesai! Semua request akan terimport otomatis

## ⚙️ Setup Environment

Setelah import:

1. **Lihat sidebar kiri** (bagian bawah)
2. **Klik "Environments"** atau ikon globe 🌐
3. **Pilih environment**: "Local", "Development", atau "Production"
4. Environment sudah aktif ✓

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

**2. Login** → Dapatkan cookies untuk autentikasi
```json
{
  "username": "testuser",
  "password": "password123"
}
```

**3. Set Cookies Manual** (Wajib setelah Login)

Setelah Login berhasil, extract token dari response dan set ke variables:

**Cara 1: Dari Response Headers**
1. Jalankan request **Login**
2. Lihat **Response Headers** → cari `Set-Cookie`
3. Copy values dari `access_token=...` dan `refresh_token=...`
4. Di Bruno, klik panel **Variables** (atau Environment → Local)
5. Tambahkan:
   ```
   access_token = "paste-value-disini"
   refresh_token = "paste-value-disini"
   ```

✅ Sekarang semua authenticated request akan pakai cookies ini!

### Langkah 2: Buat Wallet

**4. Create Account** → Buat wallet pertama
```json
{
  "wallet_name": "Dompet Utama",
  "currency": "IDR"
}
```

**5. Copy ID** dari response dan simpan ke variable `wallet_id`

### Langkah 3: Isi Saldo

**6. Deposit** → Tambah saldo ke wallet (amount sebagai string dengan presisi desimal)
```json
{
  "amount": "100000.50",
  "description": "Deposit awal"
}
```

### Langkah 4: Transaksi

**7. Withdraw** → Tarik dana (amount sebagai string dengan presisi desimal)
```json
{
  "amount": "50000.25",
  "description": "Tarik tunai"
}
```

**8. Transfer** → Kirim ke wallet lain
```json
{
  "to_wallet_id": "uuid-wallet-tujuan",
  "amount": "25000.75",
  "description": "Transfer"
}
```

### Langkah 5: Lihat Riwayat

**9. Get Transactions** → Lihat semua transaksi

## 🔧 Setup Variables

### Authentication Variables (Set after Login)

Setelah menjalankan **Login**, extract token dari response headers:

```
access_token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
refresh_token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Lokasi set variable:**
- Klik **Environments** di sidebar kiri
- Pilih **Local** (atau environment aktif)
- Edit dan paste token values

### Wallet Variables

Untuk request yang butuh `wallet_id`:

1. Jalankan **Create Account**
2. Copy `id` dari response body
3. Set di Environment Variables

## 📚 API Endpoints

### Authentication (Tanpa Login)
- `POST /v1/auth/register` - Daftar user baru (dengan email)
- `POST /v1/auth/login` - Login dan dapatkan cookies
- `POST /v1/auth/refresh` - Refresh token

### Authentication (Butuh Login)
- `POST /v1/auth/logout` - Logout

### User Profile (Butuh Login)
- `GET /v1/users/profile` - Lihat profil user (dengan email)
- `PUT /v1/users/profile` - Update profil user (username & email)

### Wallet Management (Butuh Login)
- `POST /v1/wallets` - Buat wallet baru
- `GET /v1/wallets` - Lihat semua wallet user
- `GET /v1/wallets/:id` - Lihat detail wallet
- `POST /v1/wallets/:id/deposit` - Deposit saldo (amount sebagai string)
- `POST /v1/wallets/:id/withdraw` - Tarik saldo (amount sebagai string)
- `POST /v1/wallets/:id/transfer` - Transfer ke wallet lain (amount sebagai string)
- `GET /v1/wallets/:id/transactions` - Lihat riwayat transaksi

## 💡 Tips

- **Urutan testing**: Register → Login → Set Cookies Manual → Create Account → Deposit → [Transfer/Withdraw]
- **Set Cookies Selalu**: Setiap kali restart Bruno atau ganti user, jalankan Login dan set cookies manual
- **Simpan ID**: Selalu copy ID dari response untuk dipakai di request lain
- **Cek Cookies**: Setelah set cookies, jalankan **Get Profile** untuk verifikasi
- **Amount Format**: Selalu gunakan string untuk amount dengan presisi desimal (contoh: "100.50", "1000.00")

## ❓ Troubleshooting

### "401 Unauthorized" pada authenticated requests?

**Solusi:**
1. Jalankan **Login** request dulu
2. Copy cookies dari response headers
3. Set manual di Environment Variables (access_token & refresh_token)
4. Coba lagi request yang tadi gagal

### "Wallet not found"?

**Solusi:**
1. Pastikan sudah Create Account (wallet)
2. Copy `id` dari response body
3. Set variable `wallet_id` dengan ID tersebut

### "Invalid amount format"?

**Solusi:**
1. Pastikan amount dikirim sebagai string (bukan number/integer)
2. Gunakan format desimal yang valid: "100.50", "1000.00", "50"

---

**Selamat Testing! 🚀**
