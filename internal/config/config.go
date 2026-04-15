package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type Config struct {
	Discord  DiscordConfig  `json:"discord"`
	Model    ModelConfig    `json:"model"`
	Database DatabaseConfig `json:"database"`
	Exec     ExecConfig     `json:"exec"`
	MCP      MCPConfig      `json:"mcp"`
	Health   HealthConfig   `json:"health"`
	Logging  LoggingConfig  `json:"logging"`
}

type DiscordConfig struct {
	Token           string   `json:"token"`
	GuildID         string   `json:"guild_id"`
	AllowedChannels []string `json:"allowed_channels"`
}

type ModelConfig struct {
	Provider    string  `json:"provider"`
	APIKey      string  `json:"api_key"`
	Endpoint    string  `json:"endpoint"`
	Model       string  `json:"model"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}

type DatabaseConfig struct {
	Path string `json:"path"`
}

type ExecConfig struct {
	Shell       string `json:"shell"`
	CWD         string `json:"cwd"`
	TimeoutSecs int    `json:"timeout_seconds"`
	Backend     string `json:"backend"` // "host" or "docker"
}

type MCPConfig struct {
	Servers []MCPServer `json:"servers"`
}

type MCPServer struct {
	Name string `json:"name"`
	Type string `json:"type"` // "stdio" or "sse"
	URL  string `json:"url"`
}

type HealthConfig struct {
	Port    int  `json:"port"`
	Enabled bool `json:"enabled"`
}

type LoggingConfig struct {
	Level string `json:"level"`
	Path  string `json:"path"`
}

func Load(configPath string) (*Config, error) {
	if configPath == "" {
		home, err := homedir.Dir()
		if err != nil {
			return nil, fmt.Errorf("cannot find home directory: %w", err)
		}
		configPath = filepath.Join(home, ".saia", "config.json")
	}

	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("json")

	// Environment variable overrides
	v.AutomaticEnv()
	v.BindEnv("discord.token", "SAIA_DISCORD_TOKEN")
	v.BindEnv("model.api_key", "SAIA_MODEL_API_KEY")
	v.BindEnv("database.path", "SAIA_DB_PATH")
	v.BindEnv("logging.level", "SAIA_LOG_LEVEL")
	v.BindEnv("health.port", "SAIA_PORT")

	if err := v.ReadInConfig(); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file not found: %s (see config.json.example)", configPath)
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	// Expand ~ in paths
	cfg.Database.Path = expandPath(cfg.Database.Path)
	cfg.Exec.CWD = expandPath(cfg.Exec.CWD)
	cfg.Logging.Path = expandPath(cfg.Logging.Path)

	return &cfg, nil
}

func expandPath(p string) string {
	if len(p) == 0 || p[0] != '~' {
		return p
	}
	home, err := homedir.Dir()
	if err != nil {
		return p
	}
	return filepath.Join(home, p[1:])
}
