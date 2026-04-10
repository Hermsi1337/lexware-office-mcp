package lexware

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfigFromEnv_RequiresToken(t *testing.T) {
	t.Setenv("LEXWARE_API_TOKEN", "")

	_, err := LoadConfigFromEnv()
	if err == nil {
		t.Fatal("expected error when LEXWARE_API_TOKEN is empty")
	}
}

func TestLoadConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("LEXWARE_API_TOKEN", "test-token")

	// Clear optional vars so defaults apply.
	for _, key := range []string{"LEXWARE_BASE_URL", "LEXWARE_USER_AGENT", "LEXWARE_FINALIZE_INVOICES"} {
		t.Setenv(key, "")
	}

	cfg, err := LoadConfigFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.APIToken != "test-token" {
		t.Errorf("APIToken = %q, want %q", cfg.APIToken, "test-token")
	}
	if cfg.BaseURL != defaultBaseURL {
		t.Errorf("BaseURL = %q, want %q", cfg.BaseURL, defaultBaseURL)
	}
	if cfg.UserAgent != defaultUserAgent() {
		t.Errorf("UserAgent = %q, want %q", cfg.UserAgent, defaultUserAgent())
	}
	if cfg.HTTPTimeout != 30*time.Second {
		t.Errorf("HTTPTimeout = %v, want %v", cfg.HTTPTimeout, 30*time.Second)
	}
	if cfg.FinalizeInvoices {
		t.Error("FinalizeInvoices should default to false")
	}
}

func TestLoadConfigFromEnv_CustomValues(t *testing.T) {
	t.Setenv("LEXWARE_API_TOKEN", "custom-token")
	t.Setenv("LEXWARE_BASE_URL", "https://custom.api.test/")
	t.Setenv("LEXWARE_USER_AGENT", "my-agent/2.0")
	t.Setenv("LEXWARE_FINALIZE_INVOICES", "true")

	cfg, err := LoadConfigFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.APIToken != "custom-token" {
		t.Errorf("APIToken = %q, want %q", cfg.APIToken, "custom-token")
	}
	// Trailing slash should be trimmed.
	if cfg.BaseURL != "https://custom.api.test" {
		t.Errorf("BaseURL = %q, want %q", cfg.BaseURL, "https://custom.api.test")
	}
	if cfg.UserAgent != "my-agent/2.0" {
		t.Errorf("UserAgent = %q, want %q", cfg.UserAgent, "my-agent/2.0")
	}
	if !cfg.FinalizeInvoices {
		t.Error("FinalizeInvoices should be true when env is 'true'")
	}
}

func TestParseBoolEnv(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		fallback bool
		want     bool
	}{
		{"empty uses fallback false", "", false, false},
		{"empty uses fallback true", "", true, true},
		{"true string", "true", false, true},
		{"false string", "false", true, false},
		{"1 is true", "1", false, true},
		{"0 is false", "0", true, false},
		{"invalid uses fallback", "maybe", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := "TEST_BOOL_" + tt.name
			if tt.envValue != "" {
				os.Setenv(key, tt.envValue)
				defer os.Unsetenv(key)
			} else {
				os.Unsetenv(key)
			}

			got := parseBoolEnv(key, tt.fallback)
			if got != tt.want {
				t.Errorf("parseBoolEnv(%q, %v) = %v, want %v", key, tt.fallback, got, tt.want)
			}
		})
	}
}
