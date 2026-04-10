package lexware

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ErrorSuite struct{ baseSuite }

func TestErrorSuite(t *testing.T) { suite.Run(t, new(ErrorSuite)) }

func (s *ErrorSuite) TestAPIErrorReturnsStatusCode() {
	s.mux.HandleFunc("GET /v1/profile", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"invalid token"}`))
	})

	_, err := s.client.GetProfile(context.Background())
	require.Error(s.T(), err)
	require.Contains(s.T(), err.Error(), "401")
}
