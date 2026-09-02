<div align="center">

<img src="https://img.shields.io/badge/Windows-0078D4?style=for-the-badge&logo=windows&logoColor=white" alt="Windows"/>
<img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go"/>

# 🛡️ WinSudo

### sudo for Windows

**Jalankan perintah dengan hak akses administrator secara aman dan terdokumentasi.**

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg?style=flat)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/Bernic777/winsudo?color=green&logo=github)](https://github.com/Bernic777/winsudo/releases)
[![Downloads](https://img.shields.io/github/downloads/Bernic777/winsudo/total?color=purple)](https://github.com/Bernic777/winsudo/releases)

---

[features](#-fitur) • [instalasi](#-instalasi) • [penggunaan](#-penggunaan) • [konfigurasi](#%EF%B8%8F-konfigurasi) • [kontribusi](#-kontribusi)

</div>

---

## 📖 Tentang

WinSudo adalah utilitas command-line untuk Windows yang terinspirasi dari `sudo` di Unix/Linux. Tool ini memungkinkan Anda menjalankan perintah dengan hak akses administrator menggunakan autentikasi password, pencatatan log, dan konfigurasi keamanan.

---

## ✨ Fitur

| Fitur | Deskripsi |
|-------|-----------|
| 🔐 **Autentikasi Password** | Verifikasi identitas sebelum elevasi hak akses |
| 🛡️ **Integrasi UAC** | Menggunakan mekanisme keamanan bawaan Windows |
| 📝 **Pencatatan Log** | Mencatat setiap perintah yang dijalankan |
| ⚡ **Cache Kredensial** | Menyimpan password sementara (5 menit) |
| 🔒 **Kebijakan Perintah** | Whitelist/blacklist perintah tertentu |
| 👤 **Kontrol Pengguna** | Batasi siapa yang boleh menggunakan sudo |

---

## 📦 Instalasi

### Opsi 1: Download Binary (Disarankan)

1. Kunjungi halaman [Releases](https://github.com/Bernic777/winsudo/releases)
2. Download `winsudo.zip`
3. Ekstrak ke folder yang diinginkan
4. Jalankan `sudo.exe` dari command prompt

### Opsi 2: Build dari Source

```powershell
# Clone repository
git clone https://github.com/Bernic777/winsudo.git
cd winsudo

# Build
go build -o sudo.exe .

# Tambahkan ke PATH (opsional)
$env:PATH += ";$(Get-Location)"
```

### Opsi 3: Install dengan Go

```powershell
go install github.com/Bernic777/winsudo@latest
```

### Persyaratan

| Komponen | Versi Minimum |
|----------|---------------|
| Go | 1.21 atau lebih baru |
| Windows | 10 / 11 |
| Arsitektur | amd64 |

---

## 🚀 Penggunaan

### Perintah Dasar

```powershell
# Buka command prompt dengan hak admin
sudo cmd

# Buka PowerShell dengan hak admin
sudo powershell

# Jalankan program sebagai admin
sudo notepad.exe
sudo explorer.exe
```

### Jalankan Perintah Tertentu

```powershell
# Buat user baru
sudo net user admin123 /add

# Lihat isi direktori
sudo dir C:\Windows

# Jalankan dengan argument
sudo tasklist /svc
```

### Perintah Utilitas

```powershell
# Tampilkan versi
sudo --version

# Tampilkan bantuan
sudo --help

# Cek status admin
sudo --admin

# Kelola cache kredensial
sudo --clear-cache
sudo --list-cache
```

---

## ⚙️ Konfigurasi

Edit file `config/winsudo.json` untuk menyesuaikan perilaku:

```json
{
  "version": "1.0.0",
  "auth": {
    "timeout_seconds": 300,
    "max_attempts": 3,
    "require_password": true
  },
  "allowed_users": ["Administrator", "admin"],
  "allowed_commands": ["cmd", "powershell", "notepad"],
  "blocked_commands": ["format", "rd /s"],
  "audit": {
    "enabled": true,
    "log_file": "logs/audit.log"
  },
  "elevation": {
    "use_uac": true
  }
}
```

### Referensi Konfigurasi

| Seksi | Kunci | Deskripsi | Default |
|-------|-------|-----------|---------|
| `auth` | `timeout_seconds` | Durasi cache kredensial | `300` |
| | `max_attempts` | Maksimal percobaan password | `3` |
| | `require_password` | Aktifkan autentikasi | `true` |
| `allowed_users` | - | Pengguna yang boleh menggunakan sudo | `[]` (semua) |
| `allowed_commands` | - | Perintah yang diizinkan | `[]` (semua) |
| `blocked_commands` | - | Perintah yang diblokir | `[]` |
| `audit` | `enabled` | Aktifkan pencatatan log | `true` |
| | `log_file` | Lokasi file log | `logs/audit.log` |
| `elevation` | `use_uac` | Gunakan UAC untuk elevasi | `true` |

---

## 🏗️ Arsitektur

```
winsudo/
├── main.go                    # Titik masuk & parsing CLI
├── internal/
│   ├── auth/
│   │   └── auth.go            # API LogonUser Windows
│   ├── audit/
│   │   └── audit.go           # Pencatatan berbasis file
│   ├── config/
│   │   └── config.go          # Pemuat konfigurasi JSON
│   ├── executor/
│   │   └── executor.go        # Pembuatan proses
│   └── platform/
│       └── windows.go         # Panggilan API Win32
└── config/
    └── winsudo.json           # Konfigurasi default
```

---

## 🔒 Keamanan

- **Eksekusi tanpa admin**: Jalankan sudo dari terminal biasa untuk keamanan optimal
- **Cache kredensial**: Password disimpan selama 5 menit (dapat dikonfigurasi)
- **Jejak audit**: Semua perintah dicatat di `logs/audit.log`
- **Filter perintah**: Blokir perintah berbahaya melalui kebijakan
- **Integrasi UAC**: Menggunakan prompt elevasi bawaan Windows

---

## 🤝 Kontribusi

Kontribusi sangat diterima! Silakan baca [Panduan Kontribusi](CONTRIBUTING.md) terlebih dahulu.

```powershell
# Fork dan clone
git clone https://github.com/YOUR_USERNAME/winsudo.git
cd winsudo

# Buat branch fitur
git checkout -b feature/fitur-menarik

# Buat perubahan dan uji
go build -o sudo.exe .
.\sudo.exe --version

# Commit dan push
git add .
git commit -m "Tambah fitur menarik"
git push origin feature/fitur-menarik
```

---

## 📋 Changelog

### v1.0.0 (2026-09-02)
- 🎉 Rilis pertama
- 🔐 Autentikasi password
- 🛡️ Elevasi UAC
- 📝 Pencatatan log
- ⚡ Cache kredensial
- 🔒 Kebijakan perintah

---

## 📄 Lisensi

Proyek ini dilisensikan di bawah Lisensi MIT - lihat file [LICENSE](LICENSE) untuk detailnya.

---

<div align="center">

Dibuat dengan ❤️ untuk komunitas Windows

[⬆ Kembali ke Atas](#-winsudo)

</div>
