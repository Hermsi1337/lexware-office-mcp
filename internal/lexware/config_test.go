package lexware

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// ---------- ConfigSuite ----------

type ConfigSuite struct {
	suite.Suite
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}

func (s *ConfigSuite) TestRequiresToken() {
	s.T().Setenv("LEXWARE_API_TOKEN", "")

	_, err := LoadConfigFromEnv()
	require.Error(s.T(), err)
}

func (s *ConfigSuite) TestDefaults() {
	s.T().Setenv("LEXWARE_API_TOKEN", "test-token")
	for _, key := range []string{"LEXWARE_BASE_URL", "LEXWARE_USER_AGENT", "LEXWARE_FINALIZE_INVOICES"} {
		s.T().Setenv(key, "")
	}

	cfg, err := LoadConfigFromEnv()
	require.NoError(s.T(), err)

	require.Equal(s.T(), "test-token", cfg.APIToken)
	require.Equal(s.T(), defaultBaseURL, cfg.BaseURL)
	require.Equal(s.T(), defaultUserAgent(), cfg.UserAgent)
	require.Equal(s.T(), 30*time.Second, cfg.HTTPTimeout)
	require.False(s.T(), cfg.FinalizeInvoices)
}

func (s *ConfigSuite) TestCustomValues() {
	s.T().Setenv("LEXWARE_API_TOKEN", "custom-token")
	s.T().Setenv("LEXWARE_BASE_URL", "https://custom.api.test/")
	s.T().Setenv("LEXWARE_USER_AGENT", "my-agent/2.0")
	s.T().Setenv("LEXWARE_FINALIZE_INVOICES", "true")

	cfg, err := LoadConfigFromEnv()
	require.NoError(s.T(), err)

	require.Equal(s.T(), "custom-token", cfg.APIToken)
	require.Equal(s.T(), "https://custom.api.test", cfg.BaseURL, "trailing slash should be trimmed")
	require.Equal(s.T(), "my-agent/2.0", cfg.UserAgent)
	require.True(s.T(), cfg.FinalizeInvoices)
}

// ---------- ParseBoolEnvSuite ----------

type ParseBoolEnvSuite struct {
	suite.Suite
}

func TestParseBoolEnvSuite(t *testing.T) {
	suite.Run(t, new(ParseBoolEnvSuite))
}

func (s *ParseBoolEnvSuite) TestEmptyFallbackFalse() {
	os.Unsetenv("TEST_BOOL_EMPTY_FALSE")
	require.Equal(s.T(), false, parseBoolEnv("TEST_BOOL_EMPTY_FALSE", false))
}

func (s *ParseBoolEnvSuite) TestEmptyFallbackTrue() {
	os.Unsetenv("TEST_BOOL_EMPTY_TRUE")
	require.Equal(s.T(), true, parseBoolEnv("TEST_BOOL_EMPTY_TRUE", true))
}

func (s *ParseBoolEnvSuite) TestTrueString() {
	s.T().Setenv("TEST_BOOL_TRUE", "true")
	require.Equal(s.T(), true, parseBoolEnv("TEST_BOOL_TRUE", false))
}

func (s *ParseBoolEnvSuite) TestFalseString() {
	s.T().Setenv("TEST_BOOL_FALSE", "false")
	require.Equal(s.T(), false, parseBoolEnv("TEST_BOOL_FALSE", true))
}

func (s *ParseBoolEnvSuite) TestOneIsTrue() {
	s.T().Setenv("TEST_BOOL_ONE", "1")
	require.Equal(s.T(), true, parseBoolEnv("TEST_BOOL_ONE", false))
}

func (s *ParseBoolEnvSuite) TestZeroIsFalse() {
	s.T().Setenv("TEST_BOOL_ZERO", "0")
	require.Equal(s.T(), false, parseBoolEnv("TEST_BOOL_ZERO", true))
}

func (s *ParseBoolEnvSuite) TestInvalidUsesFallback() {
	s.T().Setenv("TEST_BOOL_INVALID", "maybe")
	require.Equal(s.T(), false, parseBoolEnv("TEST_BOOL_INVALID", false))
}
