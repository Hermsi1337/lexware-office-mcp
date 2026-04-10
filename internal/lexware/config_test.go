package lexware

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoadConfigFromEnv_RequiresToken(t *testing.T) {
	t.Setenv("LEXWARE_API_TOKEN", "")

	_, err := LoadConfigFromEnv()
	require.Error(t, err)
}

func TestLoadConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("LEXWARE_API_TOKEN", "test-token")

	// Clear optional vars so defaults apply.
	for _, key := range []string{"LEXWARE_BASE_URL", "LEXWARE_USER_AGENT", "LEXWARE_FINALIZE_INVOICES"} {
		t.Setenv(key, "")
	}

	cfg, err := LoadConfigFromEnv()
	require.NoError(t, err)

	require.Equal(t, "test-token", cfg.APIToken)
	require.Equal(t, defaultBaseURL, cfg.BaseURL)
	require.Equal(t, defaultUserAgent(), cfg.UserAgent)
	require.Equal(t, 30*time.Second, cfg.HTTPTimeout)
	require.False(t, cfg.FinalizeInvoices)
}

func TestLoadConfigFromEnv_CustomValues(t *testing.T) {
	t.Setenv("LEXWARE_API_TOKEN", "custom-token")
	t.Setenv("LEXWARE_BASE_URL", "https://custom.api.test/")
	t.Setenv("LEXWARE_USER_AGENT", "my-agent/2.0")
	t.Setenv("LEXWARE_FINALIZE_INVOICES", "true")

	cfg, err := LoadConfigFromEnv()
	require.NoError(t, err)

	require.Equal(t, "custom-token", cfg.APIToken)
	require.Equal(t, "https://custom.api.test", cfg.BaseURL, "trailing slash should be trimmed")
	require.Equal(t, "my-agent/2.0", cfg.UserAgent)
	require.True(t, cfg.FinalizeInvoices)
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
			require.Equal(t, tt.want, got)
		})
	}
}
