package auth

import (
	"fmt"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	credCache     = make(map[string]time.Time)
	procGetUserNameW = syscall.NewLazyDLL("advapi32.dll").NewProc("GetUserNameW")
	procLogonUserW   = syscall.NewLazyDLL("advapi32.dll").NewProc("LogonUserW")
)

type AuthResult struct {
	Success  bool
	Username string
	Error    error
}

func GetCurrentUser() (string, error) {
	var bufSize uint32 = 256
	buf := make([]uint16, bufSize)

	ret, _, err := procGetUserNameW.Call(
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&bufSize)),
	)
	if ret == 0 {
		return "", fmt.Errorf("GetUserName failed: %v", err)
	}

	return windows.UTF16ToString(buf[:bufSize-1]), nil
}

func VerifyPassword(username, password string) bool {
	var token uintptr

	ret, _, _ := procLogonUserW.Call(
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(username))),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr("."))),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(password))),
		3, // LOGON32_LOGON_INTERACTIVE
		0, // LOGON32_PROVIDER_DEFAULT
		uintptr(unsafe.Pointer(&token)),
	)
	if ret == 0 {
		return false
	}
	syscall.CloseHandle(syscall.Handle(token))
	return true
}

func Authenticate(password string, timeoutSeconds int) AuthResult {
	username, err := GetCurrentUser()
	if err != nil {
		return AuthResult{Success: false, Error: err}
	}

	if cached, ok := credCache[username]; ok {
		if time.Since(cached).Seconds() < float64(timeoutSeconds) {
			return AuthResult{Success: true, Username: username}
		}
		delete(credCache, username)
	}

	if !VerifyPassword(username, password) {
		return AuthResult{
			Success:  false,
			Username: username,
			Error:    fmt.Errorf("authentication failed"),
		}
	}

	credCache[username] = time.Now()
	return AuthResult{Success: true, Username: username}
}

func readPassword() string {
	var nRead uint32

	handle := syscall.Handle(syscall.Stdin)

	var mode uint32
	windows.GetConsoleMode(windows.Handle(handle), &mode)
	windows.SetConsoleMode(windows.Handle(handle), mode&^windows.ENABLE_ECHO_INPUT)
	defer windows.SetConsoleMode(windows.Handle(handle), mode)

	bufBytes := make([]byte, 512)
	syscall.ReadFile(handle, bufBytes, &nRead, nil)

	buf := (*[256]uint16)(unsafe.Pointer(&bufBytes[0]))
	return strings.TrimRight(windows.UTF16ToString(buf[:nRead/2]), "\r\n")
}

func AuthenticateWithPrompt(timeoutSeconds int) AuthResult {
	username, err := GetCurrentUser()
	if err != nil {
		return AuthResult{Success: false, Error: err}
	}

	if cached, ok := credCache[username]; ok {
		if time.Since(cached).Seconds() < float64(timeoutSeconds) {
			fmt.Println("[*] Using cached credentials")
			return AuthResult{Success: true, Username: username}
		}
		delete(credCache, username)
	}

	fmt.Printf("[*] Password for %s: ", username)
	password := readPassword()
	fmt.Println()

	if !VerifyPassword(username, password) {
		return AuthResult{
			Success:  false,
			Username: username,
			Error:    fmt.Errorf("authentication failed"),
		}
	}

	credCache[username] = time.Now()
	return AuthResult{Success: true, Username: username}
}

func ClearCache() {
	credCache = make(map[string]time.Time)
}

func GetCachedUsers() []string {
	users := make([]string, 0, len(credCache))
	for u, t := range credCache {
		if time.Since(t).Seconds() < 300 {
			users = append(users, u)
		}
	}
	return users
}
