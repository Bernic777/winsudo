package audit

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	logFile *os.File
	mu      sync.Mutex
)

type AuditEntry struct {
	Timestamp time.Time
	Username  string
	Command   string
	Args      string
	Success   bool
	PID       int
}

func Init(logPath string) error {
	dir := filepath.Dir(logPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	var err error
	logFile, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return err
}

func Close() {
	if logFile != nil {
		logFile.Close()
	}
}

func Log(entry AuditEntry) {
	mu.Lock()
	defer mu.Unlock()

	if logFile == nil {
		return
	}

	status := "SUCCESS"
	if !entry.Success {
		status = "FAILED"
	}

	line := fmt.Sprintf("[%s] %s | User: %s | Command: %s %s | Status: %s | PID: %d\n",
		entry.Timestamp.Format("2006-01-02 15:04:05"),
		entry.Timestamp.Format("-0700"),
		entry.Username,
		entry.Command,
		entry.Args,
		status,
		entry.PID,
	)

	logFile.WriteString(line)
	logFile.Sync()
}

func LogCommand(username, command string, args []string, success bool, pid int) {
	fullArgs := ""
	for _, a := range args {
		fullArgs += a + " "
	}

	Log(AuditEntry{
		Timestamp: time.Now(),
		Username:  username,
		Command:   command,
		Args:      fullArgs,
		Success:   success,
		PID:       pid,
	})
}

func LogEvent(username, event string) {
	Log(AuditEntry{
		Timestamp: time.Now(),
		Username:  username,
		Command:   "EVENT",
		Args:      event,
		Success:   true,
		PID:       os.Getpid(),
	})
}
