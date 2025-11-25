package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	PingURLs         []string
	PingInterval     time.Duration
	RequestTimeout   time.Duration
	MaxRetries       int
	LogLevel         string
	LogResponseBody  bool
}

func Load() (*Config, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	// Parse PING_URLS
	urlsStr := os.Getenv("PING_URLS")
	if urlsStr == "" {
		return nil, fmt.Errorf("PING_URLS is required")
	}
	urls := strings.Split(urlsStr, ",")
	var cleanURLs []string
	for _, url := range urls {
		trimmed := strings.TrimSpace(url)
		if trimmed != "" {
			cleanURLs = append(cleanURLs, trimmed)
		}
	}
	if len(cleanURLs) == 0 {
		return nil, fmt.Errorf("no valid URLs found in PING_URLS")
	}

	// Parse PING_INTERVAL
	intervalStr := getEnvOrDefault("PING_INTERVAL", "60")
	intervalSec, err := strconv.Atoi(intervalStr)
	if err != nil {
		return nil, fmt.Errorf("invalid PING_INTERVAL: %w", err)
	}

	// Parse REQUEST_TIMEOUT
	timeoutStr := getEnvOrDefault("REQUEST_TIMEOUT", "10")
	timeoutSec, err := strconv.Atoi(timeoutStr)
	if err != nil {
		return nil, fmt.Errorf("invalid REQUEST_TIMEOUT: %w", err)
	}

	// Parse MAX_RETRIES
	retriesStr := getEnvOrDefault("MAX_RETRIES", "3")
	maxRetries, err := strconv.Atoi(retriesStr)
	if err != nil {
		return nil, fmt.Errorf("invalid MAX_RETRIES: %w", err)
	}

	// Parse LOG_LEVEL
	logLevel := getEnvOrDefault("LOG_LEVEL", "info")

	// Parse LOG_RESPONSE_BODY
	logBodyStr := getEnvOrDefault("LOG_RESPONSE_BODY", "false")
	logBody, err := strconv.ParseBool(logBodyStr)
	if err != nil {
		logBody = false
	}

	return &Config{
		PingURLs:        cleanURLs,
		PingInterval:    time.Duration(intervalSec) * time.Second,
		RequestTimeout:  time.Duration(timeoutSec) * time.Second,
		MaxRetries:      maxRetries,
		LogLevel:        logLevel,
		LogResponseBody: logBody,
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}