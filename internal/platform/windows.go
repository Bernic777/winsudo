package platform

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modshell32        = syscall.NewLazyDLL("shell32.dll")
	procShellExecuteExW = modshell32.NewProc("ShellExecuteExW")
)

type shellExecuteInfo struct {
	CbSize       uint32
	FMask        uint32
	Hwnd         uintptr
	LpVerb       *uint16
	LpFile       *uint16
	LpParameters *uint16
	LpDirectory  *uint16
	NShow        int32
	HInstApp     uintptr
	HInstItem    uintptr
	RID          uintptr
	HMonitor     uintptr
	Rest         [8]byte
}

const (
	SEE_MASK_NOCLOSEPROCESS = 0x00000040
	SW_SHOW                 = 5
)

func IsAdmin() bool {
	var sid *windows.SID
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		1,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		0, 0, 0, 0, 0, 0, 0,
		&sid,
	)
	if err != nil {
		return false
	}
	defer windows.FreeSid(sid)

	token, err := windows.OpenCurrentProcessToken()
	if err != nil {
		return false
	}
	defer token.Close()

	member, err := token.IsMember(sid)
	if err != nil {
		return false
	}

	return member
}

func RunElevated(command string, args []string) error {
	verb := windows.StringToUTF16Ptr("runas")
	exe, _ := syscall.UTF16PtrFromString(command)
	
	argStr := ""
	for _, a := range args {
		argStr += a + " "
	}
	argPtr, _ := syscall.UTF16PtrFromString(argStr)
	cwd, _ := os.Getwd()
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)

	si := &shellExecuteInfo{
		CbSize:       uint32(unsafe.Sizeof(shellExecuteInfo{})),
		FMask:        SEE_MASK_NOCLOSEPROCESS,
		LpVerb:       verb,
		LpFile:       exe,
		LpParameters: argPtr,
		LpDirectory:  cwdPtr,
		NShow:        SW_SHOW,
	}

	ret, _, _ := procShellExecuteExW.Call(uintptr(unsafe.Pointer(si)))
	if ret == 0 {
		return fmt.Errorf("ShellExecuteEx failed")
	}

	syscall.WaitForSingleObject(syscall.Handle(si.HInstApp), syscall.INFINITE)
	syscall.CloseHandle(syscall.Handle(si.HInstApp))

	return nil
}
