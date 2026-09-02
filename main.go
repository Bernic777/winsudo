package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"

	"winsudo/internal/audit"
	"winsudo/internal/auth"
	"winsudo/internal/config"
	"winsudo/internal/executor"
	"winsudo/internal/platform"

	"golang.org/x/sys/windows"
)

const VERSION = "1.0.0"

var usage = `WinSudo v1.0.0 - sudo for Windows

Usage:
  winsudo <command> [args...]    Run command as administrator
  winsudo -u <user> <command>    Run as specific user
  winsudo --clear-cache          Clear cached credentials
  winsudo --list-cache           List cached users
  winsudo --version              Show version
  winsudo --help                 Show this help

Examples:
  winsudo cmd                    Open elevated command prompt
  winsudo powershell             Open elevated PowerShell
  winsudo notepad.exe            Run Notepad as admin
  winsudo net user admin123 /add Create new user
  winsudo "dir C:\Windows"       Run directory listing elevated

Config: config/winsudo.json
Logs:   logs/audit.log`

func printBanner() {
	banner := `
╔══════════════════════════════════════════════╗
║           WinSudo v%s - sudo for Windows           ║
║           Run commands with elevation             ║
╚══════════════════════════════════════════════╝`
	fmt.Printf(banner+"\n", VERSION)
}

func promptPassword() string {
	fmt.Print("[*] Password: ")
	
	handle := windows.Handle(syscall.Stdin)
	var mode uint32
	windows.GetConsoleMode(handle, &mode)
	windows.SetConsoleMode(handle, mode&^windows.ENABLE_ECHO_INPUT)

	var buf [256]uint16
	var nRead uint32
	bufBytes := make([]byte, 512)
	syscall.ReadFile(syscall.Handle(handle), bufBytes, &nRead, nil)
	buf = *(*[256]uint16)(unsafe.Pointer(&bufBytes[0]))
	
	windows.SetConsoleMode(handle, mode)
	fmt.Println()
	
	return strings.TrimRight(windows.UTF16ToString(buf[:nRead/2]), "\r\n")
}

func main() {
	printBanner()

	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(1)
	}

	cfg, err := config.LoadConfig("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] Config error: %v\n", err)
		fmt.Println("[*] Using default settings")
		cfg = &config.Config{
			Auth: config.AuthConfig{
				TimeoutSeconds:  300,
				MaxAttempts:     3,
				RequirePassword: true,
			},
			AllowedCommands: []string{},
			BlockedCommands: []string{},
			Audit: config.AuditConfig{
				Enabled: true,
				LogFile: "logs/audit.log",
			},
			Elevation: config.ElevationConfig{
				UseUAC: true,
			},
		}
	}

	if cfg.Audit.Enabled {
		audit.Init(cfg.Audit.LogFile)
		defer audit.Close()
	}

	args := os.Args[1:]
	
	switch args[0] {
	case "--version", "-v":
		fmt.Printf("WinSudo v%s\n", VERSION)
		os.Exit(0)
		
	case "--help", "-h":
		fmt.Println(usage)
		os.Exit(0)
		
	case "--clear-cache":
		auth.ClearCache()
		fmt.Println("[+] Cache cleared")
		os.Exit(0)
		
	case "--list-cache":
		users := auth.GetCachedUsers()
		if len(users) == 0 {
			fmt.Println("[*] No cached users")
		} else {
			fmt.Println("[*] Cached users:")
			for _, u := range users {
				fmt.Printf("  - %s\n", u)
			}
		}
		os.Exit(0)
		
	case "--admin":
		if !platform.IsAdmin() {
			fmt.Println("[!] Not running as Administrator")
			os.Exit(1)
		}
		fmt.Println("[+] Running as Administrator")
		os.Exit(0)
	}

	cmdArgs := args
	
	if args[0] == "-u" && len(args) > 2 {
		cmdArgs = args[2:]
	}

	command := cmdArgs[0]
	commandArgs := cmdArgs[1:]

	if !cfg.Elevation.UseUAC && !platform.IsAdmin() {
		fmt.Fprintln(os.Stderr, "[!] UAC elevation required. Run from elevated prompt or enable UAC in config.")
		os.Exit(1)
	}

	if !cfg.IsCommandAllowed(command) {
		fmt.Fprintf(os.Stderr, "[!] Command '%s' is not allowed by policy\n", command)
		audit.LogEvent("system", fmt.Sprintf("BLOCKED: %s", command))
		os.Exit(1)
	}

	var authResult auth.AuthResult
	
	if cfg.Auth.RequirePassword {
		for attempts := 0; attempts < cfg.Auth.MaxAttempts; attempts++ {
			password := promptPassword()
			authResult = auth.Authenticate(password, cfg.Auth.TimeoutSeconds)
			
			if authResult.Success {
				break
			}
			
			remaining := cfg.Auth.MaxAttempts - attempts - 1
			if remaining > 0 {
				fmt.Fprintf(os.Stderr, "[!] Authentication failed. %d attempts remaining\n", remaining)
			}
		}
		
		if !authResult.Success {
			fmt.Fprintln(os.Stderr, "[!] Authentication failed. Access denied.")
			audit.LogEvent(authResult.Username, "AUTH_FAILED")
			os.Exit(1)
		}
	} else {
		username, _ := auth.GetCurrentUser()
		authResult = auth.AuthResult{Success: true, Username: username}
	}

	fmt.Printf("[*] Executing: %s %s\n", command, strings.Join(commandArgs, " "))
	
	result := executor.Execute(command, commandArgs, cfg.Elevation.UseUAC)
	
	audit.LogCommand(
		authResult.Username,
		command,
		commandArgs,
		result.ExitCode == 0,
		result.PID,
	)

	if result.Error != nil {
		fmt.Fprintf(os.Stderr, "[!] Execution error: %v\n", result.Error)
		os.Exit(1)
	}

	os.Exit(result.ExitCode)
}
