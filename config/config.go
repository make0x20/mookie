package config

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/viper"
	"os"
	"strings"
)

/*
	Package config provides configuration management using Viper and TOML.
	Supports environment variables, default values, and file-based config.

	How to use:
	1. Create a config file or use defaults
	2. Load config with NewWithPath
	3. Access values through Config struct

	Example usage:
		// Load config
		cfg, err := config.NewWithPath(*path)
		if cfg == nil {
			log.Fatalf("error loading config: %v", err)
		}

	Environment variables:
	- Prefix: MOOKIE_ (customize as needed)
	- Format: MOOKIE_BINDADDRESS, MOOKIE_PORT, etc.
	- Overrides file config when present

	Config precedence:
	1. Environment variables
	2. Config file values
	3. Default values

	Default values:
	- BindAddress: "0.0.0.0"
	- Port: 8080
	- DatabasePath: "app.db"
	- LogFile: "" (stdout)
	- LogLevel: "normal"
*/

// Config defines the application configuration
type Config struct {
	BindAddress  string `mapstructure:"BindAddress"`
	Port         int    `mapstructure:"Port"`
	DatabasePath string `mapstructure:"DatabasePath"`
	LogFile      string `mapstructure:"LogFile"`
	LogLevel     string `mapstructure:"LogLevel"`
}

// NewWithPath creates a new config from the given path.
func NewWithPath(configPath string) (*Config, error) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		cfg := getDefaultConfig()
		data, err := toml.Marshal(cfg)
		if err != nil {
			return nil, fmt.Errorf("error creating default config: %w", err)
		}
		if err := os.WriteFile(configPath, data, 0644); err != nil {
			return nil, fmt.Errorf("error writing default config: %w", err)
		}
	}
	return loadConfig(configPath)
}

// loadConfig loads the config from the given path.
// If the file does not exist, it creates a default config file.
func loadConfig(configPath string) (*Config, error) {
	v := viper.New()

	// Set some defaults
	v.SetDefault("BindAddress", "0.0.0.0")
	v.SetDefault("Port", 8080)
	v.SetDefault("DatabasePath", "app.db")
	v.SetDefault("LogFile", "")
	v.SetDefault("LogLevel", "normal")

	v.SetConfigFile(configPath)
	v.SetConfigType("toml")
	v.AutomaticEnv()
	v.SetEnvPrefix("MOOKIE") // Change this to your app's name
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// getDefaultConfig returns the default config.
func getDefaultConfig() *Config {
	return &Config{
		BindAddress:  "0.0.0.0",
		Port:         8080,
		DatabasePath: "app.db",
		LogFile:      "",
		LogLevel:     "normal",
	}
}
