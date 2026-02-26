package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Storage  StorageConfig
}

type ServerConfig struct {
	Address      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	Secret       string
	SSLKey       string
	SSLSert      string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type StorageConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
	Location  string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Address:      GetEnv("SERVER_ADDRESS", ":443"),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  1 * time.Minute,
			Secret:       GetEnv("JWT_SECRET", ""),
			SSLKey:       GetEnv("SSL_KEY_PATH", ""),
			SSLSert:      GetEnv("SSL_CERT_PATH", ""),
		},
		Database: DatabaseConfig{
			Host:     GetEnv("DB_HOST", "db"),
			Port:     GetEnv("DB_PORT", "5432"),
			User:     GetEnv("DB_USER", ""),
			Password: GetEnv("DB_PASSWORD", ""),
			Name:     GetEnv("DB_NAME", ""),
			SSLMode:  GetEnv("DB_SSLMODE", "disable"),
		},
		Storage: StorageConfig{
			Endpoint:  GetEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey: GetEnv("MINIO_ACCESS_KEY", ""),
			SecretKey: GetEnv("MINIO_SECRET_KEY", ""),
			Bucket:    GetEnv("MINIO_BUCKET", "attachments"),
			UseSSL:    GetEnvBool("MINIO_USE_SSL", false),
			Location:  GetEnv("MINIO_LOCATION", "us-east-1"),
		},
	}
}

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func GetEnvBool(key string, defaultValue bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return parsed
}

func (dc *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dc.Host, dc.Port, dc.User, dc.Password, dc.Name, dc.SSLMode)
}
