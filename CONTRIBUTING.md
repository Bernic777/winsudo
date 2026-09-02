# Contributing to WinSudo

Thank you for your interest in contributing to WinSudo! This document provides guidelines and information for contributors.

## Getting Started

1. Fork the repository
2. Clone your fork
3. Create a feature branch
4. Make your changes
5. Submit a pull request

## Development Setup

### Prerequisites

- Go 1.21 or later
- Windows 10/11
- Git

### Building

```bash
# Clone repository
git clone https://github.com/Bernic777/winsudo.git
cd winsudo

# Build
go build -o sudo.exe .

# Run tests
go test ./...
```

### Project Structure

```
winsudo/
├── main.go                    # Entry point
├── internal/
│   ├── auth/                  # Authentication
│   ├── audit/                 # Logging
│   ├── config/                # Config management
│   ├── executor/              # Command execution
│   └── platform/              # Windows API
└── config/
    └── winsudo.json           # Configuration
```

## Code Style

- Follow standard Go conventions
- Use `gofmt` to format code
- Run `go vet` before committing
- Add comments for complex logic
- Keep functions focused and small

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test
go test -v ./internal/auth/...
```

## Pull Request Process

1. Update documentation if needed
2. Add tests for new features
3. Ensure all tests pass
4. Update CHANGELOG.md (if exists)
5. Request review from maintainers

## Reporting Issues

- Use GitHub Issues
- Include Windows version
- Include Go version
- Provide minimal reproduction steps
- Include error messages

## Code of Conduct

- Be respectful
- Constructive feedback
- Help others learn
- Focus on the project

## Questions?

Open an issue for questions about contributing.
