package lexware

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultBaseURL     = "https://api.lexware.io"
	defaultUserAgent   = "lexware-office-mcp/0.1.0"
	defaultMinInterval = 550 * time.Millisecond
)

type Config struct {
	APIToken         string
	BaseURL          string
	UserAgent        string
	MinInterval      time.Duration
	HTTPTimeout      time.Duration
	FinalizeInvoices bool
}

func LoadConfigFromEnv() (Config, error) {
	token := strings.TrimSpace(os.Getenv("LEXWARE_API_TOKEN"))
	if token == "" {
		return Config{}, fmt.Errorf("LEXWARE_API_TOKEN is required")
	}

	baseURL := strings.TrimSpace(os.Getenv("LEXWARE_BASE_URL"))
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	userAgent := strings.TrimSpace(os.Getenv("LEXWARE_USER_AGENT"))
	if userAgent == "" {
		userAgent = defaultUserAgent
	}

	minInterval := defaultMinInterval
	if raw := strings.TrimSpace(os.Getenv("LEXWARE_MIN_INTERVAL_MS")); raw != "" {
		ms, err := strconv.Atoi(raw)
		if err != nil || ms < 0 {
			return Config{}, fmt.Errorf("LEXWARE_MIN_INTERVAL_MS must be a non-negative integer")
		}
		minInterval = time.Duration(ms) * time.Millisecond
	}

	return Config{
		APIToken:         token,
		BaseURL:          strings.TrimRight(baseURL, "/"),
		UserAgent:        userAgent,
		MinInterval:      minInterval,
		HTTPTimeout:      30 * time.Second,
		FinalizeInvoices: parseBoolEnv("LEXWARE_FINALIZE_INVOICES", false),
	}, nil
}

func parseBoolEnv(key string, fallback bool) bool {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}

	value, err := strconv.ParseBool(raw)
	if err != nil {
		return fallback
	}

	return value
}
