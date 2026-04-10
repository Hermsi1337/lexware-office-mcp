package lexware

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RecurringTemplateSuite struct{ baseSuite }

func TestRecurringTemplateSuite(t *testing.T) { suite.Run(t, new(RecurringTemplateSuite)) }

func (s *RecurringTemplateSuite) TestGetRecurringTemplate() {
	s.mux.HandleFunc("GET /v1/recurring-templates/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(RecurringTemplateDetail{
			ID:             r.PathValue("id"),
			Title:          "Monthly Hosting",
			RecurringCycle: "MONTHLY",
		})
	})

	tpl, err := s.client.GetRecurringTemplate(context.Background(), "rt-42")
	require.NoError(s.T(), err)
	require.Equal(s.T(), "rt-42", tpl.ID)
	require.Equal(s.T(), "Monthly Hosting", tpl.Title)
}
