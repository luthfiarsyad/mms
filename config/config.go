package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Paseto   PasetoConfig   `mapstructure:"paseto"`
	Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
	Mode    string `mapstructure:"mode"`
	Address string `mapstructure:"address"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	DSN      string `mapstructure:"dsn"`
}

type PasetoConfig struct {
	SymmetricKey  string `mapstructure:"symmetric_key"`
	ExpireMinutes int    `mapstructure:"expire_minutes"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
}

// Cfg holds the loaded configuration for the whole application.
// After calling Load, other packages can read config via config.Get()
var Cfg *Config

// Load reads configuration from config/config.yaml (or a custom path if provided).
// It does NOT set any defaults — missing required fields will return an error.
func Load(configPath string) error {
	v := viper.New()

	if configPath == "" {
		// look for config/config.yaml relative to project root or current working dir
		v.AddConfigPath("./config")
		v.AddConfigPath(".")
		v.SetConfigName("config")
		v.SetConfigType("yaml")
	} else {
		v.SetConfigFile(configPath)
	}

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}

	// Basic required-field validation (fail early if values are missing)
	if cfg.Server.Address == "" || cfg.Server.Mode == "" {
		return fmt.Errorf("config: server.address and server.mode must be set")
	}
	// Database: either full DSN or host/user/name must be provided
	if cfg.Database.DSN == "" {
		if cfg.Database.Host == "" || cfg.Database.User == "" || cfg.Database.Name == "" {
			return fmt.Errorf("config: either database.dsn or (database.host, database.user, database.name) must be set")
		}
	}
	if cfg.Paseto.SymmetricKey == "" {
		return fmt.Errorf("config: paseto.symmetric_key must be set")
	}

	Cfg = &cfg
	return nil
}

// Get returns the loaded Config pointer (may be nil if Load wasn't called or failed)
func Get() *Config { return Cfg }
