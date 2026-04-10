package lexware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("Authorization = %q, want Bearer test-token", got)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Profile{OrganizationID: "org-123"})
	})

	client, _ := testClient(t, mux)
	profile, err := client.GetProfile(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if profile.OrganizationID != "org-123" {
		t.Errorf("OrganizationID = %q, want %q", profile.OrganizationID, "org-123")
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if contact.ID != "abc-123" {
		t.Errorf("ID = %q, want %q", contact.ID, "abc-123")
	}
}

func TestListContacts(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/contacts", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("name"); got != "Muster%" {
			t.Errorf("name param = %q, want %q", got, "Muster%")
		}
		if got := r.URL.Query().Get("customer"); got != "true" {
			t.Errorf("customer param = %q, want %q", got, "true")
		}

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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Content) != 1 {
		t.Fatalf("got %d contacts, want 1", len(result.Content))
	}
	if result.Content[0].ID != "c1" {
		t.Errorf("contact ID = %q, want %q", result.Content[0].ID, "c1")
	}
}

func TestCreateInvoice(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/invoices", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("finalize"); got != "true" {
			t.Errorf("finalize param = %q, want %q", got, "true")
		}

		var invoice Invoice
		if err := json.NewDecoder(r.Body).Decode(&invoice); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if invoice.VoucherDate != "2026-01-01" {
			t.Errorf("VoucherDate = %q, want %q", invoice.VoucherDate, "2026-01-01")
		}

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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "inv-1" {
		t.Errorf("ID = %q, want %q", result.ID, "inv-1")
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inv.VoucherStatus != "open" {
		t.Errorf("VoucherStatus = %q, want %q", inv.VoucherStatus, "open")
	}
}

func TestListArticles(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/articles", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("type"); got != "SERVICE" {
			t.Errorf("type param = %q, want %q", got, "SERVICE")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Page[ArticleDetail]{
			Content:       []ArticleDetail{{ID: "a1", Title: "Consulting"}},
			TotalElements: 1,
		})
	})

	client, _ := testClient(t, mux)
	result, err := client.ListArticles(context.Background(), ArticleFilter{Type: "SERVICE"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Content) != 1 || result.Content[0].Title != "Consulting" {
		t.Errorf("unexpected articles: %+v", result.Content)
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(countries) != 2 {
		t.Fatalf("got %d countries, want 2", len(countries))
	}
	if countries[0].CountryCode != "DE" {
		t.Errorf("first country = %q, want %q", countries[0].CountryCode, "DE")
	}
}

func TestAPIError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/profile", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"invalid token"}`))
	})

	client, _ := testClient(t, mux)
	_, err := client.GetProfile(context.Background())
	if err == nil {
		t.Fatal("expected error for 401 response")
	}
	if got := err.Error(); got == "" {
		t.Error("error message should not be empty")
	}
}

func TestListVouchers(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/voucherlist", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("voucherType"); got != "salesinvoice" {
			t.Errorf("voucherType = %q, want %q", got, "salesinvoice")
		}

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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Content) != 1 || result.Content[0].VoucherID != "v1" {
		t.Errorf("unexpected vouchers: %+v", result.Content)
	}
}

func TestFinalizeDefault(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/invoices", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("finalize"); got != "true" {
			t.Errorf("finalize = %q, want %q (from config default)", got, "true")
		}

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
		FinalizeInvoices: true, // default finalize = true
	}
	client := NewClient(cfg)

	// Pass nil to use the config default.
	result, err := client.CreateInvoice(context.Background(), Invoice{
		VoucherDate:   "2026-01-01",
		TaxConditions: TaxConditionGross(),
		ShippingTerms: ShippingTermNone(),
		TotalPrice:    TotalPrice{Currency: "EUR"},
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "inv-default" {
		t.Errorf("ID = %q, want %q", result.ID, "inv-default")
	}
}
