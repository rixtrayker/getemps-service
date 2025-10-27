package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Cache    CacheConfig
	Security SecurityConfig
	Logging  LoggingConfig
}

type AppConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	DBName       string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type CacheConfig struct {
	Enabled bool
	TTL     int // seconds
}

type SecurityConfig struct {
	APISecretKey  string
	TokenRequired bool
}

type LoggingConfig struct {
	Level string
	ToDB  bool
}

func Load() (*Config, error) {
	// Try to load .env file (ignore error if file doesn't exist)
	_ = godotenv.Load()

	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres123")
	viper.SetDefault("DB_NAME", "getemps_db")
	viper.SetDefault("DB_SSL_MODE", "disable")
	viper.SetDefault("DB_MAX_OPEN_CONNS", 25)
	viper.SetDefault("DB_MAX_IDLE_CONNS", 5)
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", 6379)
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 0)
	viper.SetDefault("CACHE_ENABLED", true)
	viper.SetDefault("CACHE_TTL", 300)
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_TO_DB", true)
	viper.SetDefault("API_SECRET_KEY", "your-super-secret-key-change-in-production")
	viper.SetDefault("TOKEN_REQUIRED", true)

	viper.AutomaticEnv()

	config := &Config{
		App: AppConfig{
			Port: getEnvStr("APP_PORT", "8080"),
			Env:  getEnvStr("APP_ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:         getEnvStr("DB_HOST", "localhost"),
			Port:         getEnvInt("DB_PORT", 5432),
			User:         getEnvStr("DB_USER", "postgres"),
			Password:     getEnvStr("DB_PASSWORD", "postgres123"),
			DBName:       getEnvStr("DB_NAME", "getemps_db"),
			SSLMode:      getEnvStr("DB_SSL_MODE", "disable"),
			MaxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", 5),
		},
		Redis: RedisConfig{
			Host:     getEnvStr("REDIS_HOST", "localhost"),
			Port:     getEnvInt("REDIS_PORT", 6379),
			Password: getEnvStr("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		Cache: CacheConfig{
			Enabled: getEnvBool("CACHE_ENABLED", true),
			TTL:     getEnvInt("CACHE_TTL", 300),
		},
		Security: SecurityConfig{
			APISecretKey:  getEnvStr("API_SECRET_KEY", "your-super-secret-key-change-in-production"),
			TokenRequired: getEnvBool("TOKEN_REQUIRED", true),
		},
		Logging: LoggingConfig{
			Level: getEnvStr("LOG_LEVEL", "info"),
			ToDB:  getEnvBool("LOG_TO_DB", true),
		},
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func getEnvStr(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func validateConfig(config *Config) error {
	if config.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if config.Database.User == "" {
		return fmt.Errorf("database user is required")
	}
	if config.Database.DBName == "" {
		return fmt.Errorf("database name is required")
	}
	if config.Security.APISecretKey == "" || config.Security.APISecretKey == "your-super-secret-key-change-in-production" {
		if config.App.Env == "production" {
			return fmt.Errorf("API secret key must be set in production")
		}
	}
	return nil
}
