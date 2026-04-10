package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/dennis/lexware-office-mcp/internal/lexware"
)

type createArticleInput struct {
	Article lexware.Article `json:"article" jsonschema:"Article payload with title, type (PRODUCT or SERVICE), unitName, and price"`
}

type getArticleInput struct {
	ID string `json:"id" jsonschema:"UUID of the article to retrieve"`
}

type listArticlesInput struct {
	ArticleNumber string `json:"articleNumber,omitempty" jsonschema:"Filter by exact article number"`
	Gtin          string `json:"gtin,omitempty" jsonschema:"Filter by GTIN"`
	Type          string `json:"type,omitempty" jsonschema:"Filter by type: PRODUCT or SERVICE"`
	Page          int    `json:"page,omitempty" jsonschema:"Page number for pagination (0-based)"`
}

func (s *Server) registerArticleTools() {
	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_create_article",
		Description: "Create an article (product or service) with title, unit, and pricing.",
	}, s.createArticle)

	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_get_article",
		Description: "Retrieve a single article by its UUID.",
	}, s.getArticle)

	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_list_articles",
		Description: "List articles with optional filters and pagination. Supports exact match on article number and GTIN, and filtering by type (PRODUCT or SERVICE). No full-text search available.",
	}, s.listArticles)
}

func (s *Server) createArticle(ctx context.Context, _ *mcp.CallToolRequest, input createArticleInput) (*mcp.CallToolResult, map[string]any, error) {
	result, err := s.client.CreateArticle(ctx, input.Article)
	return s.workflowResult("create article", result, err)
}

func (s *Server) getArticle(ctx context.Context, _ *mcp.CallToolRequest, input getArticleInput) (*mcp.CallToolResult, map[string]any, error) {
	if strings.TrimSpace(input.ID) == "" {
		return nil, nil, fmt.Errorf("id is required")
	}

	result, err := s.client.GetArticle(ctx, input.ID)
	return s.workflowResult("get article", result, err)
}

func (s *Server) listArticles(ctx context.Context, _ *mcp.CallToolRequest, input listArticlesInput) (*mcp.CallToolResult, map[string]any, error) {
	filter := lexware.ArticleFilter{
		ArticleNumber: input.ArticleNumber,
		Gtin:          input.Gtin,
		Type:          input.Type,
		Page:          input.Page,
	}

	result, err := s.client.ListArticles(ctx, filter)
	return s.workflowResult("list articles", result, err)
}
