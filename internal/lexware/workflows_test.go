package lexware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// testClient returns a Client pointing at a test HTTP server.
func testClient(t *testing.T, handler http.Handler) (*Client, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	cfg := Config{
		APIToken:    "test-token",
		BaseURL:     srv.URL,
		UserAgent:   "test-agent",
		HTTPTimeout: 5 * time.Second,
	}
	return NewClient(cfg), srv
}

func TestGetProfile(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/profile", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Profile{OrganizationID: "org-123"})
	})

	client, _ := testClient(t, mux)
	profile, err := client.GetProfile(context.Background())
	require.NoError(t, err)
	require.Equal(t, "org-123", profile.OrganizationID)
}

func TestGetContact(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/contacts/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ContactDetail{
			ID:      id,
			Version: 1,
		})
	})

	client, _ := testClient(t, mux)
	contact, err := client.GetContact(context.Background(), "abc-123")
	require.NoError(t, err)
	require.Equal(t, "abc-123", contact.ID)
}

func TestListContacts(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/contacts", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "Muster%", r.URL.Query().Get("name"))
		require.Equal(t, "true", r.URL.Query().Get("customer"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Page[ContactDetail]{
			Content:       []ContactDetail{{ID: "c1"}},
			TotalElements: 1,
			TotalPages:    1,
			First:         true,
			Last:          true,
		})
	})

	client, _ := testClient(t, mux)
	boolTrue := true
	result, err := client.ListContacts(context.Background(), ContactFilter{
		Name:     "Muster%",
		Customer: &boolTrue,
	})
	require.NoError(t, err)
	require.Len(t, result.Content, 1)
	require.Equal(t, "c1", result.Content[0].ID)
}

func TestCreateInvoice(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/invoices", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "true", r.URL.Query().Get("finalize"))

		var invoice Invoice
		err := json.NewDecoder(r.Body).Decode(&invoice)
		require.NoError(t, err)
		require.Equal(t, "2026-01-01", invoice.VoucherDate)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(CreateInvoiceResult{ID: "inv-1"})
	})

	client, _ := testClient(t, mux)
	finalize := true
	result, err := client.CreateInvoice(context.Background(), Invoice{
		VoucherDate:   "2026-01-01",
		TaxConditions: TaxConditionGross(),
		ShippingTerms: ShippingTermNone(),
		TotalPrice:    TotalPrice{Currency: "EUR"},
	}, &finalize)
	require.NoError(t, err)
	require.Equal(t, "inv-1", result.ID)
}

func TestGetInvoice(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/invoices/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(InvoiceDetail{
			ID:            r.PathValue("id"),
			VoucherStatus: "open",
			VoucherNumber: "RE-001",
		})
	})

	client, _ := testClient(t, mux)
	inv, err := client.GetInvoice(context.Background(), "inv-99")
	require.NoError(t, err)
	require.Equal(t, "open", inv.VoucherStatus)
}

func TestListArticles(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/articles", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "SERVICE", r.URL.Query().Get("type"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Page[ArticleDetail]{
			Content:       []ArticleDetail{{ID: "a1", Title: "Consulting"}},
			TotalElements: 1,
		})
	})

	client, _ := testClient(t, mux)
	result, err := client.ListArticles(context.Background(), ArticleFilter{Type: "SERVICE"})
	require.NoError(t, err)
	require.Len(t, result.Content, 1)
	require.Equal(t, "Consulting", result.Content[0].Title)
}

func TestListCountries(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/countries", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Country{
			{CountryCode: "DE", CountryNameEN: "Germany", TaxClassification: "de"},
			{CountryCode: "AT", CountryNameEN: "Austria", TaxClassification: "intraCommunity"},
		})
	})

	client, _ := testClient(t, mux)
	countries, err := client.ListCountries(context.Background())
	require.NoError(t, err)
	require.Len(t, countries, 2)
	require.Equal(t, "DE", countries[0].CountryCode)
}

func TestAPIError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/profile", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"invalid token"}`))
	})

	client, _ := testClient(t, mux)
	_, err := client.GetProfile(context.Background())
	require.Error(t, err)
	require.Contains(t, err.Error(), "401")
}

func TestListVouchers(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/voucherlist", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "salesinvoice", r.URL.Query().Get("voucherType"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Page[VoucherListItem]{
			Content: []VoucherListItem{
				{VoucherID: "v1", VoucherType: "salesinvoice", TotalAmount: 119.0},
			},
			TotalElements: 1,
		})
	})

	client, _ := testClient(t, mux)
	result, err := client.ListVouchers(context.Background(), VoucherlistFilter{
		VoucherType: "salesinvoice",
	})
	require.NoError(t, err)
	require.Len(t, result.Content, 1)
	require.Equal(t, "v1", result.Content[0].VoucherID)
}

func TestFinalizeDefault(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/invoices", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "true", r.URL.Query().Get("finalize"), "should use config default")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(CreateInvoiceResult{ID: "inv-default"})
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	cfg := Config{
		APIToken:         "test-token",
		BaseURL:          srv.URL,
		UserAgent:        "test-agent",
		HTTPTimeout:      5 * time.Second,
		FinalizeInvoices: true,
	}
	client := NewClient(cfg)

	result, err := client.CreateInvoice(context.Background(), Invoice{
		VoucherDate:   "2026-01-01",
		TaxConditions: TaxConditionGross(),
		ShippingTerms: ShippingTermNone(),
		TotalPrice:    TotalPrice{Currency: "EUR"},
	}, nil)
	require.NoError(t, err)
	require.Equal(t, "inv-default", result.ID)
}
