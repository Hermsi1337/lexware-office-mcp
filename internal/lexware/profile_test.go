package lexware

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ProfileSuite struct{ baseSuite }

func TestProfileSuite(t *testing.T) { suite.Run(t, new(ProfileSuite)) }

func (s *ProfileSuite) TestGetProfile() {
	s.mux.HandleFunc("GET /v1/profile", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(s.T(), "Bearer test-token", r.Header.Get("Authorization"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Profile{OrganizationID: "org-123"})
	})

	profile, err := s.client.GetProfile(context.Background())
	require.NoError(s.T(), err)
	require.Equal(s.T(), "org-123", profile.OrganizationID)
}
