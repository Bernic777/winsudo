# WinSudo

**sudo for Windows** - Run commands with administrator privileges

[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/winsudo)](https://goreportcard.com/report/github.com/yourusername/winsudo)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/yourusername/winsudo)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/yourusername/winsudo)](https://github.com/yourusername/winsudo/releases)

## Features

- 🔐 Password authentication before elevation
- 🛡️ UAC integration for secure elevation
- 📝 Audit logging of all commands
- ⚡ Credential caching (sudo-like timeout)
- 🔒 Command whitelist/blacklist policy
- 👤 User authorization control

## Installation

### Download

Download the latest `winsudo.exe` from [Releases](https://github.com/yourusername/winsudo/releases).

### Build from Source

```bash
# Clone repository
git clone https://github.com/yourusername/winsudo.git
cd winsudo

# Build
go build -o winsudo.exe .
```

### Requirements

- Go 1.21 or later
- Windows 10/11

## Usage

```powershell
# Basic usage - run command as admin
winsudo cmd
winsudo powershell
winsudo notepad.exe

# Run specific command with args
winsudo net user admin123 /add
winsudo dir C:\Windows

# Run as specific user
winsudo -u Administrator cmd

# Utility commands
winsudo --version
winsudo --help
winsudo --admin
winsudo --clear-cache
winsudo --list-cache
```

## Configuration

Edit `config/winsudo.json`:

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

### Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `auth.timeout_seconds` | Cache duration for credentials | 300 |
| `auth.max_attempts` | Max password attempts | 3 |
| `auth.require_password` | Require password authentication | true |
| `allowed_users` | Users allowed to use winsudo | [] (all) |
| `allowed_commands` | Commands allowed to run | [] (all) |
| `blocked_commands` | Commands blocked from running | [] |
| `audit.enabled` | Enable audit logging | true |
| `audit.log_file` | Path to audit log | logs/audit.log |
| `elevation.use_uac` | Use UAC for elevation | true |

## Project Structure

```
winsudo/
├── main.go                    # Entry point
├── go.mod
├── go.sum
├── winsudo.exe                # Built binary
├── config/
│   └── winsudo.json           # Configuration
├── internal/
│   ├── auth/
│   │   └── auth.go            # Authentication
│   ├── audit/
│   │   └── audit.go           # Logging
│   ├── config/
│   │   └── config.go          # Config management
│   ├── executor/
│   │   └── executor.go        # Command execution
│   └── platform/
│       └── windows.go         # Windows API
├── logs/                      # Audit logs
├── README.md
├── LICENSE
├── CONTRIBUTING.md
└── .gitignore
```

## Security

- Always run from a non-admin terminal for proper security
- Credentials are cached for 5 minutes (configurable)
- All commands are logged to `logs/audit.log`
- Blocked commands cannot be executed via winsudo

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by Unix `sudo`
- Built with Go and Windows API
