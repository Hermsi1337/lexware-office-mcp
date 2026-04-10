package lexware

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ArticleSuite struct{ baseSuite }

func TestArticleSuite(t *testing.T) { suite.Run(t, new(ArticleSuite)) }

func (s *ArticleSuite) TestListArticles() {
	s.mux.HandleFunc("GET /v1/articles", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(s.T(), "SERVICE", r.URL.Query().Get("type"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Page[ArticleDetail]{
			Content:       []ArticleDetail{{ID: "a1", Title: "Consulting"}},
			TotalElements: 1,
		})
	})

	result, err := s.client.ListArticles(context.Background(), ArticleFilter{Type: "SERVICE"})
	require.NoError(s.T(), err)
	require.Len(s.T(), result.Content, 1)
	require.Equal(s.T(), "Consulting", result.Content[0].Title)
}
