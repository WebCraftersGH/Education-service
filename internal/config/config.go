package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	AppEnv   string
	HTTPPort string
	LogLevel string

	DBHost         string
	DBPort         string
	DBUser         string
	DBPass         string
	DBName         string
	DBSSLMode      string
	DBTimeZone     string
	TokenCookie    string
	AuthServiceURL string
}

func Load(envPath string) (Config, error) {
	if err := loadEnvFile(envPath); err != nil {
		return Config{}, err
	}

	cfg := Config{
		AppEnv:         getEnv("APP_ENV", "local"),
		HTTPPort:       getEnv("HTTP_PORT", "8080"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "eduuser"),
		DBPass:         getEnv("DB_PASSWORD", "edupass"),
		DBName:         getEnv("DB_NAME", "education_service"),
		DBSSLMode:      getEnv("DB_SSLMODE", "disable"),
		DBTimeZone:     getEnv("DB_TIMEZONE", "UTC"),
		TokenCookie:    getEnv("TOKEN_COOKIE_NAME", "token"),
		AuthServiceURL: getEnv("AUTH_SERVICE_URL", "http://localhost:8080/api/v1"),
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	fmt.Println(cfg)

	return cfg, nil
}

func (c Config) Validate() error {
	missing := make([]string, 0)

	if strings.TrimSpace(c.DBHost) == "" {
		missing = append(missing, "DB_HOST")
	}
	if strings.TrimSpace(c.DBPort) == "" {
		missing = append(missing, "DB_PORT")
	}
	if strings.TrimSpace(c.DBUser) == "" {
		missing = append(missing, "DB_USER")
	}
	if strings.TrimSpace(c.DBName) == "" {
		missing = append(missing, "DB_NAME")
	}

	if len(missing) > 0 {
		return fmt.Errorf("required env variables are missing: %s", strings.Join(missing, ", "))
	}

	return nil
}

func (c Config) HTTPAddress() string {
	port := strings.TrimSpace(c.HTTPPort)
	if port == "" {
		port = "8080"
	}

	if strings.HasPrefix(port, ":") {
		return port
	}

	return ":" + port
}

func (c Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPass,
		c.DBName,
		c.DBSSLMode,
		c.DBTimeZone,
	)
}

func getEnv(key, def string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return def
	}
	return value
}

func loadEnvFile(envPath string) error {
	if strings.TrimSpace(envPath) == "" {
		envPath = ".env"
	}

	file, err := os.Open(envPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, `"'`)

		if key == "" {
			continue
		}

		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}

	return scanner.Err()
}
