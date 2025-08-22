package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Address      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	Secret       string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Address:      GetEnv("SERVER_ADDRESS", ":443"),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  1 * time.Minute,
			Secret:       GetEnv("JWT_SECRET", ""),
		},
		Database: DatabaseConfig{
			Host:     GetEnv("DB_HOST", "localhost"),
			Port:     GetEnv("DB_PORT", "5432"),
			User:     GetEnv("DB_USER", "grintheone"),
			Password: GetEnv("DB_PASSWORD", ""),
			Name:     GetEnv("DB_NAME", "foxygendb"),
			SSLMode:  GetEnv("DB_SSLMODE", "disable"),
		},
	}
}

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func (dc *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dc.Host, dc.Port, dc.User, dc.Password, dc.Name, dc.SSLMode)
}
