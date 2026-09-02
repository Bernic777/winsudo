package executor

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modshell32          = syscall.NewLazyDLL("shell32.dll")
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

type ExecuteResult struct {
	Command  string
	ExitCode int
	Error    error
	PID      int
}

func Execute(command string, args []string, elevated bool) ExecuteResult {
	if elevated {
		return executeElevated(command, args)
	}
	return executeNormal(command, args)
}

func executeNormal(command string, args []string) ExecuteResult {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			return ExecuteResult{
				Command:  command,
				ExitCode: -1,
				Error:    err,
			}
		}
	}

	return ExecuteResult{
		Command:  command,
		ExitCode: exitCode,
		PID:      cmd.Process.Pid,
	}
}

func executeElevated(command string, args []string) ExecuteResult {
	exe, err := exec.LookPath(command)
	if err != nil {
		exe = command
	}

	verb := windows.StringToUTF16Ptr("runas")
	exePtr, _ := syscall.UTF16PtrFromString(exe)

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
		LpFile:       exePtr,
		LpParameters: argPtr,
		LpDirectory:  cwdPtr,
		NShow:        SW_SHOW,
	}

	ret, _, _ := procShellExecuteExW.Call(uintptr(unsafe.Pointer(si)))
	if ret == 0 {
		return ExecuteResult{
			Command:  command,
			ExitCode: -1,
			Error:    fmt.Errorf("ShellExecuteEx failed"),
		}
	}

	syscall.WaitForSingleObject(syscall.Handle(si.HInstApp), syscall.INFINITE)

	var exitCode uint32
	windows.GetExitCodeProcess(windows.Handle(si.HInstApp), &exitCode)
	syscall.CloseHandle(syscall.Handle(si.HInstApp))

	return ExecuteResult{
		Command:  command,
		ExitCode: int(exitCode),
		PID:      int(si.HInstApp),
	}
}

func ExecuteShell(command string) error {
	cmd := exec.Command("cmd.exe", "/c", command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
