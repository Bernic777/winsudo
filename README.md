<div align="center">

# 🛡️ WinSudo

### **sudo for Windows**

**Run any command with administrator privileges — securely, auditable, effortlessly.**

[![Go Report Card](https://goreportcard.com/badge/github.com/Bernic777/winsudo)](https://goreportcard.com/report/github.com/Bernic777/winsudo)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/Bernic777/winsudo?color=green&logo=github)](https://github.com/Bernic777/winsudo/releases)
[![Platform](https://img.shields.io/badge/Platform-Windows-0078D4?logo=windows&logoColor=white)](https://github.com/Bernic777/winsudo)
[![Downloads](https://img.shields.io/github/downloads/Bernic777/winsudo/total?color=purple&logo=docusign)](https://github.com/Bernic777/winsudo/releases)

---

**WinSudo** brings the power of Unix `sudo` to Windows. Execute commands with elevated privileges using a simple, secure, and auditable interface.

[Getting Started](#-quick-start) • [Features](#-features) • [Installation](#-installation) • [Usage](#-usage) • [Configuration](#%EF%B8%8F-configuration) • [Contributing](#-contributing)

</div>

---

## ⚡ Quick Start

```powershell
# Download and run
winsudo cmd
```

That's it. You now have an elevated command prompt.

---

## 🎯 Features

<table>
<tr>
<td>

🔐 **Password Authentication**
Verify identity before elevation

</td>
<td>

🛡️ **UAC Integration**
Windows-native security

</td>
<td>

📝 **Audit Logging**
Track every command executed

</td>
</tr>
<tr>
<td>

⚡ **Credential Caching**
5-minute timeout (configurable)

</td>
<td>

🔒 **Command Policy**
Whitelist/blacklist commands

</td>
<td>

👤 **User Control**
Authorize specific users only

</td>
</tr>
</table>

---

## 📦 Installation

### Option 1: Download Binary (Recommended)

Download the latest `winsudo.exe` from [Releases](https://github.com/Bernic777/winsudo/releases) and place it in your `PATH`.

### Option 2: Build from Source

```powershell
# Clone repository
git clone https://github.com/Bernic777/winsudo.git
cd winsudo

# Build
go build -o winsudo.exe .

# Optional: Add to PATH
$env:PATH += ";$(Get-Location)"
```

### Option 3: Install with Go

```powershell
go install github.com/Bernic777/winsudo@latest
```

### Requirements

| Requirement | Version |
|-------------|---------|
| Go | 1.21+ |
| Windows | 10/11 |
| Architecture | amd64 |

---

## 🚀 Usage

### Basic Commands

```powershell
# Open elevated command prompt
winsudo cmd

# Open elevated PowerShell
winsudo powershell

# Run any executable as admin
winsudo notepad.exe
winsudo explorer.exe
```

### Run Specific Commands

```powershell
# Create a new user
winsudo net user admin123 /add

# List directory contents
winsudo dir C:\Windows

# Run with arguments
winsudo tasklist /svc
```

### Utility Commands

```powershell
# Show version
winsudo --version

# Show help
winsudo --help

# Check admin status
winsudo --admin

# Manage credential cache
winsudo --clear-cache
winsudo --list-cache
```

---

## ⚙️ Configuration

Edit `config/winsudo.json` to customize behavior:

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

### Configuration Reference

| Section | Key | Description | Default |
|---------|-----|-------------|---------|
| `auth` | `timeout_seconds` | Credential cache duration | `300` |
| | `max_attempts` | Max password tries | `3` |
| | `require_password` | Enable authentication | `true` |
| `allowed_users` | - | Users who can use winsudo | `[]` (all) |
| `allowed_commands` | - | Commands allowed to run | `[]` (all) |
| `blocked_commands` | - | Commands blocked from running | `[]` |
| `audit` | `enabled` | Enable logging | `true` |
| | `log_file` | Audit log path | `logs/audit.log` |
| `elevation` | `use_uac` | Use UAC elevation | `true` |

---

## 🏗️ Architecture

```
winsudo/
├── main.go                    # Entry point & CLI parsing
├── internal/
│   ├── auth/
│   │   └── auth.go            # Windows LogonUser API
│   ├── audit/
│   │   └── audit.go           # File-based logging
│   ├── config/
│   │   └── config.go          # JSON config loader
│   ├── executor/
│   │   └── executor.go        # Process creation
│   └── platform/
│       └── windows.go         # Win32 API calls
└── config/
    └── winsudo.json           # Default configuration
```

### How It Works

```
┌─────────────┐    ┌──────────────┐    ┌─────────────┐
│   User      │───▶│  Authenticate │───▶│   Execute   │
│   Command   │    │  (Password)  │    │  (Elevated) │
└─────────────┘    └──────────────┘    └─────────────┘
                          │                     │
                          ▼                     ▼
                   ┌──────────────┐    ┌─────────────┐
                   │  Log Event   │    │  Audit Log  │
                   └──────────────┘    └─────────────┘
```

---

## 🔒 Security

- **Non-admin execution**: Run winsudo from a regular terminal for proper security
- **Credential caching**: Passwords cached for 5 minutes (configurable)
- **Audit trail**: All commands logged to `logs/audit.log`
- **Command filtering**: Block dangerous commands via policy
- **UAC integration**: Uses Windows native elevation prompts

---

## 🤝 Contributing

Contributions are welcome! Please read our [Contributing Guidelines](CONTRIBUTING.md) first.

```powershell
# Fork and clone
git clone https://github.com/YOUR_USERNAME/winsudo.git
cd winsudo

# Create feature branch
git checkout -b feature/amazing-feature

# Make changes and test
go build -o winsudo.exe .
.\winsudo.exe --version

# Commit and push
git add .
git commit -m "Add amazing feature"
git push origin feature/amazing-feature
```

---

## 📋 Changelog

### v1.0.0 (2026-09-02)
- 🎉 Initial release
- 🔐 Password authentication
- 🛡️ UAC elevation
- 📝 Audit logging
- ⚡ Credential caching
- 🔒 Command policy

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<div align="center">

**Made with ❤️ for the Windows community**

[⬆ Back to Top](#-winsudo)

</div>
