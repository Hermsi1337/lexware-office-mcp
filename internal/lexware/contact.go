package lexware

import (
	"context"
	"strconv"
	"strings"
)

type CreateContactResult struct {
	ID string `json:"id"`
}

// ContactFilter holds optional query parameters for listing contacts.
type ContactFilter struct {
	Email    string
	Name     string
	Number   *int
	Customer *bool
	Vendor   *bool
	Page     int
}

func (c *Client) CreateSimpleContact(ctx context.Context, name, sourceReference string) (*CreateContactResult, error) {
	contact := Contact{
		Version: 0,
		Roles: Roles{
			Customer: map[string]any{},
		},
		Person: Person{
			LastName: name,
		},
	}
	if strings.TrimSpace(sourceReference) != "" {
		contact.Note = sourceReference
	}

	result := &CreateContactResult{}
	resp, err := c.newRequest(ctx).
		SetBody(contact).
		SetResult(result).
		Post("/v1/contacts")
	if apiErr := wrapAPIError("create contact", resp, err); apiErr != nil {
		return result, apiErr
	}

	return result, nil
}

func (c *Client) GetContact(ctx context.Context, id string) (*ContactDetail, error) {
	result := &ContactDetail{}
	resp, err := c.newRequest(ctx).
		SetResult(result).
		Get("/v1/contacts/" + id)
	if apiErr := wrapAPIError("get contact", resp, err); apiErr != nil {
		return nil, apiErr
	}

	return result, nil
}

func (c *Client) ListContacts(ctx context.Context, filter ContactFilter) (*Page[ContactDetail], error) {
	req := c.newRequest(ctx)
	if filter.Email != "" {
		req.SetQueryParam("email", filter.Email)
	}
	if filter.Name != "" {
		req.SetQueryParam("name", filter.Name)
	}
	if filter.Number != nil {
		req.SetQueryParam("number", strconv.Itoa(*filter.Number))
	}
	if filter.Customer != nil {
		req.SetQueryParam("customer", strconv.FormatBool(*filter.Customer))
	}
	if filter.Vendor != nil {
		req.SetQueryParam("vendor", strconv.FormatBool(*filter.Vendor))
	}
	if filter.Page > 0 {
		req.SetQueryParam("page", strconv.Itoa(filter.Page))
	}

	result := &Page[ContactDetail]{}
	resp, err := req.SetResult(result).Get("/v1/contacts")
	if apiErr := wrapAPIError("list contacts", resp, err); apiErr != nil {
		return nil, apiErr
	}

	return result, nil
}
