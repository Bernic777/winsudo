package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type AuthConfig struct {
	TimeoutSeconds  int `json:"timeout_seconds"`
	MaxAttempts     int `json:"max_attempts"`
	RequirePassword bool `json:"require_password"`
}

type AuditConfig struct {
	Enabled         bool   `json:"enabled"`
	LogFile         string `json:"log_file"`
	LogCommandArgs  bool   `json:"log_command_args"`
}

type ElevationConfig struct {
	UseUAC       bool `json:"use_uac"`
	PersistToken bool `json:"persist_token"`
}

type Config struct {
	Version        string           `json:"version"`
	Auth           AuthConfig       `json:"auth"`
	AllowedUsers   []string         `json:"allowed_users"`
	AllowedCommands []string        `json:"allowed_commands"`
	BlockedCommands []string        `json:"blocked_commands"`
	Audit          AuditConfig      `json:"audit"`
	Elevation      ElevationConfig  `json:"elevation"`
}

var GlobalConfig *Config

func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		exePath, _ := os.Executable()
		configPath = filepath.Join(filepath.Dir(exePath), "config", "winsudo.json")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	GlobalConfig = &cfg
	return &cfg, nil
}

func (c *Config) IsUserAllowed(username string) bool {
	for _, u := range c.AllowedUsers {
		if strings.EqualFold(u, username) {
			return true
		}
	}
	return false
}

func (c *Config) IsCommandAllowed(command string) bool {
	cmdBase := strings.ToLower(strings.Fields(command)[0])

	for _, blocked := range c.BlockedCommands {
		if strings.Contains(strings.ToLower(command), strings.ToLower(blocked)) {
			return false
		}
	}

	if len(c.AllowedCommands) == 0 {
		return true
	}

	for _, allowed := range c.AllowedCommands {
		if strings.EqualFold(allowed, cmdBase) {
			return true
		}
	}
	return false
}
