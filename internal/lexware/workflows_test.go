package lexware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// ---------- WorkflowSuite ----------

// WorkflowSuite tests the Lexware API client workflows against a local
// httptest server. SetupTest creates a fresh ServeMux before every test so
// handlers never leak between cases.
type WorkflowSuite struct {
	suite.Suite
	mux    *http.ServeMux
	srv    *httptest.Server
	client *Client
}

func TestWorkflowSuite(t *testing.T) {
	suite.Run(t, new(WorkflowSuite))
}

func (s *WorkflowSuite) SetupTest() {
	s.mux = http.NewServeMux()
	s.srv = httptest.NewServer(s.mux)

	cfg := Config{
		APIToken:    "test-token",
		BaseURL:     s.srv.URL,
		UserAgent:   "test-agent",
		HTTPTimeout: 5 * time.Second,
	}
	s.client = NewClient(cfg)
}

func (s *WorkflowSuite) TearDownTest() {
	s.srv.Close()
}

// ---------- Profile ----------

func (s *WorkflowSuite) TestGetProfile() {
	s.mux.HandleFunc("GET /v1/profile", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(s.T(), "Bearer test-token", r.Header.Get("Authorization"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Profile{OrganizationID: "org-123"})
	})

	profile, err := s.client.GetProfile(context.Background())
	require.NoError(s.T(), err)
	require.Equal(s.T(), "org-123", profile.OrganizationID)
}

// ---------- Contacts ----------

func (s *WorkflowSuite) TestGetContact() {
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

func (s *WorkflowSuite) TestListContacts() {
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

// ---------- Invoices ----------

func (s *WorkflowSuite) TestCreateInvoice() {
	s.mux.HandleFunc("POST /v1/invoices", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(s.T(), "true", r.URL.Query().Get("finalize"))

		var invoice Invoice
		err := json.NewDecoder(r.Body).Decode(&invoice)
		require.NoError(s.T(), err)
		require.Equal(s.T(), "2026-01-01", invoice.VoucherDate)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(CreateInvoiceResult{ID: "inv-1"})
	})

	finalize := true
	result, err := s.client.CreateInvoice(context.Background(), Invoice{
		VoucherDate:   "2026-01-01",
		TaxConditions: TaxConditionGross(),
		ShippingTerms: ShippingTermNone(),
		TotalPrice:    TotalPrice{Currency: "EUR"},
	}, &finalize)
	require.NoError(s.T(), err)
	require.Equal(s.T(), "inv-1", result.ID)
}

func (s *WorkflowSuite) TestGetInvoice() {
	s.mux.HandleFunc("GET /v1/invoices/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(InvoiceDetail{
			ID:            r.PathValue("id"),
			VoucherStatus: "open",
			VoucherNumber: "RE-001",
		})
	})

	inv, err := s.client.GetInvoice(context.Background(), "inv-99")
	require.NoError(s.T(), err)
	require.Equal(s.T(), "open", inv.VoucherStatus)
}

func (s *WorkflowSuite) TestFinalizeUsesConfigDefault() {
	s.mux.HandleFunc("POST /v1/invoices", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(s.T(), "true", r.URL.Query().Get("finalize"), "should use config default")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(CreateInvoiceResult{ID: "inv-default"})
	})

	// Rebuild client with FinalizeInvoices=true to test the default.
	cfg := Config{
		APIToken:         "test-token",
		BaseURL:          s.srv.URL,
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
	require.NoError(s.T(), err)
	require.Equal(s.T(), "inv-default", result.ID)
}

// ---------- Articles ----------

func (s *WorkflowSuite) TestListArticles() {
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

// ---------- Countries ----------

func (s *WorkflowSuite) TestListCountries() {
	s.mux.HandleFunc("GET /v1/countries", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Country{
			{CountryCode: "DE", CountryNameEN: "Germany", TaxClassification: "de"},
			{CountryCode: "AT", CountryNameEN: "Austria", TaxClassification: "intraCommunity"},
		})
	})

	countries, err := s.client.ListCountries(context.Background())
	require.NoError(s.T(), err)
	require.Len(s.T(), countries, 2)
	require.Equal(s.T(), "DE", countries[0].CountryCode)
}

// ---------- Voucherlist ----------

func (s *WorkflowSuite) TestListVouchers() {
	s.mux.HandleFunc("GET /v1/voucherlist", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(s.T(), "salesinvoice", r.URL.Query().Get("voucherType"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Page[VoucherListItem]{
			Content: []VoucherListItem{
				{VoucherID: "v1", VoucherType: "salesinvoice", TotalAmount: 119.0},
			},
			TotalElements: 1,
		})
	})

	result, err := s.client.ListVouchers(context.Background(), VoucherlistFilter{
		VoucherType: "salesinvoice",
	})
	require.NoError(s.T(), err)
	require.Len(s.T(), result.Content, 1)
	require.Equal(s.T(), "v1", result.Content[0].VoucherID)
}

// ---------- Error handling ----------

func (s *WorkflowSuite) TestAPIErrorReturnsStatusCode() {
	s.mux.HandleFunc("GET /v1/profile", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"invalid token"}`))
	})

	_, err := s.client.GetProfile(context.Background())
	require.Error(s.T(), err)
	require.Contains(s.T(), err.Error(), "401")
}
