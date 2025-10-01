package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Database   DatabaseConfig
	Migrations MigrationsConfig
}

// DatabaseConfig represents database connection configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

// MigrationsConfig represents migrations configuration
type MigrationsConfig struct {
	Dir string
}

// Load loads the configuration from viper
func Load() (*Config, error) {
	cfg := &Config{
		Database: DatabaseConfig{
			Host:     viper.GetString("db.host"),
			Port:     viper.GetInt("db.port"),
			User:     viper.GetString("db.user"),
			Password: viper.GetString("db.password"),
			Name:     viper.GetString("db.name"),
			SSLMode:  viper.GetString("db.sslmode"),
		},
		Migrations: MigrationsConfig{
			Dir: viper.GetString("migrations.dir"),
		},
	}

	// Validate required fields
	if cfg.Database.User == "" {
		return nil, fmt.Errorf("database user is required")
	}
	if cfg.Database.Name == "" {
		return nil, fmt.Errorf("database name is required")
	}

	// Set default SSL mode if not provided
	if cfg.Database.SSLMode == "" {
		cfg.Database.SSLMode = "disable"
	}

	return cfg, nil
}

// GetConnectionString returns the PostgreSQL connection string
func (c *Config) GetConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
	)
}
