package lexware

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/stretchr/testify/suite"
)

// baseSuite provides a fresh httptest server and Lexware Client for every
// test. Embed it in resource-specific suites to avoid boilerplate.
type baseSuite struct {
	suite.Suite
	mux    *http.ServeMux
	srv    *httptest.Server
	client *Client
}

func (s *baseSuite) SetupTest() {
	s.mux = http.NewServeMux()
	s.srv = httptest.NewServer(s.mux)
	s.client = NewClient(Config{
		APIToken:    "test-token",
		BaseURL:     s.srv.URL,
		UserAgent:   "test-agent",
		HTTPTimeout: 5 * time.Second,
	})
}

func (s *baseSuite) TearDownTest() {
	s.srv.Close()
}
