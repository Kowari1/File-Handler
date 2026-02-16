package config

import (
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	InputDir          string
	OutputDir         string
	ScanInterval      time.Duration
	WorkerCount       int
	MaxDBConns        int
	MinDBConns        int
	MaxDBConnLifetime time.Duration
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		DBHost:            getEnvStr("DB_HOST", "localhost"),
		DBPort:            getEnvStr("DB_PORT", "5432"),
		DBUser:            getEnvStr("DB_USER", "postgres"),
		DBPassword:        getEnvStr("DB_PASSWORD", "123"),
		DBName:            getEnvStr("DB_NAME", "devices"),
		InputDir:          getEnvStr("INPUT_DIR", "./input"),
		OutputDir:         getEnvStr("OUTPUT_DIR", "./output"),
		ScanInterval:      getEnvDuration("SCAN_INTERVAL", 10*time.Second),
		WorkerCount:       getEnvInt("WORKER_COUNT", 5),
		MaxDBConns:        getEnvInt("MAX_CONNS", 10),
		MinDBConns:        getEnvInt("MIN_CONNS", 2),
		MaxDBConnLifetime: getEnvDuration("MAX_CONFIG_LIFE", 1*time.Hour),
	}, nil
}

func (c *Config) DatabaseURL() string {
	if full := os.Getenv("DATABASE_URL"); full != "" {
		return full
	}

	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.DBUser, c.DBPassword),
		Host:   c.DBHost + ":" + c.DBPort,
		Path:   c.DBName,
	}

	q := u.Query()
	q.Set("sslmode", "disable")
	u.RawQuery = q.Encode()

	return u.String()
}

func getEnvStr(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func getEnvInt(key string, fallback int) int {
	if s := os.Getenv(key); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			return int(v)
		}
	}

	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if s := os.Getenv(key); s != "" {
		if d, err := time.ParseDuration(s); err == nil {
			return d
		}
	}

	return fallback
}
