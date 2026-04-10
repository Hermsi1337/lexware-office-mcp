package lexware

import (
	"context"
	"strconv"
)

type CreateArticleResult struct {
	ID string `json:"id"`
}

// ArticleFilter holds optional query parameters for listing articles.
type ArticleFilter struct {
	ArticleNumber string
	Gtin          string
	Type          string
	Page          int
}

func (c *Client) CreateArticle(ctx context.Context, article Article) (*CreateArticleResult, error) {
	result := &CreateArticleResult{}
	resp, err := c.newRequest(ctx).
		SetBody(article).
		SetResult(result).
		Post("/v1/articles")
	if apiErr := wrapAPIError("create article", resp, err); apiErr != nil {
		return result, apiErr
	}

	return result, nil
}

func (c *Client) GetArticle(ctx context.Context, id string) (*ArticleDetail, error) {
	result := &ArticleDetail{}
	resp, err := c.newRequest(ctx).
		SetResult(result).
		Get("/v1/articles/" + id)
	if apiErr := wrapAPIError("get article", resp, err); apiErr != nil {
		return nil, apiErr
	}

	return result, nil
}

func (c *Client) ListArticles(ctx context.Context, filter ArticleFilter) (*Page[ArticleDetail], error) {
	req := c.newRequest(ctx)
	if filter.ArticleNumber != "" {
		req.SetQueryParam("articleNumber", filter.ArticleNumber)
	}
	if filter.Gtin != "" {
		req.SetQueryParam("gtin", filter.Gtin)
	}
	if filter.Type != "" {
		req.SetQueryParam("type", filter.Type)
	}
	if filter.Page > 0 {
		req.SetQueryParam("page", strconv.Itoa(filter.Page))
	}

	result := &Page[ArticleDetail]{}
	resp, err := req.SetResult(result).Get("/v1/articles")
	if apiErr := wrapAPIError("list articles", resp, err); apiErr != nil {
		return nil, apiErr
	}

	return result, nil
}
