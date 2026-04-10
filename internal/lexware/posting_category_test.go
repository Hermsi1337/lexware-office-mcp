package lexware

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type PostingCategorySuite struct{ baseSuite }

func TestPostingCategorySuite(t *testing.T) { suite.Run(t, new(PostingCategorySuite)) }

func (s *PostingCategorySuite) TestListPostingCategories() {
	s.mux.HandleFunc("GET /v1/posting-categories", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]PostingCategory{
			{ID: "cat-1", Name: "Revenue", Type: "revenue"},
			{ID: "cat-2", Name: "Expense", Type: "expense"},
		})
	})

	categories, err := s.client.ListPostingCategories(context.Background())
	require.NoError(s.T(), err)
	require.Len(s.T(), categories, 2)
	require.Equal(s.T(), "Revenue", categories[0].Name)
}
