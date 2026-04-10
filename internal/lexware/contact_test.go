package lexware

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ContactSuite struct{ baseSuite }

func TestContactSuite(t *testing.T) { suite.Run(t, new(ContactSuite)) }

func (s *ContactSuite) TestGetContact() {
	s.mux.HandleFunc("GET /v1/contacts/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ContactDetail{
			ID:      r.PathValue("id"),
			Version: 1,
		})
	})

	contact, err := s.client.GetContact(context.Background(), "abc-123")
	require.NoError(s.T(), err)
	require.Equal(s.T(), "abc-123", contact.ID)
}

func (s *ContactSuite) TestListContacts() {
	s.mux.HandleFunc("GET /v1/contacts", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(s.T(), "Muster%", r.URL.Query().Get("name"))
		require.Equal(s.T(), "true", r.URL.Query().Get("customer"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Page[ContactDetail]{
			Content:       []ContactDetail{{ID: "c1"}},
			TotalElements: 1,
			TotalPages:    1,
			First:         true,
			Last:          true,
		})
	})

	boolTrue := true
	result, err := s.client.ListContacts(context.Background(), ContactFilter{
		Name:     "Muster%",
		Customer: &boolTrue,
	})
	require.NoError(s.T(), err)
	require.Len(s.T(), result.Content, 1)
	require.Equal(s.T(), "c1", result.Content[0].ID)
}
