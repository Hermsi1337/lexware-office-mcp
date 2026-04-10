package lexware

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type InvoiceSuite struct{ baseSuite }

func TestInvoiceSuite(t *testing.T) { suite.Run(t, new(InvoiceSuite)) }

func (s *InvoiceSuite) TestCreateInvoice() {
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

func (s *InvoiceSuite) TestGetInvoice() {
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

func (s *InvoiceSuite) TestFinalizeUsesConfigDefault() {
	s.mux.HandleFunc("POST /v1/invoices", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(s.T(), "true", r.URL.Query().Get("finalize"), "should use config default")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(CreateInvoiceResult{ID: "inv-default"})
	})

	// Rebuild client with FinalizeInvoices=true to test the default.
	client := NewClient(Config{
		APIToken:         "test-token",
		BaseURL:          s.srv.URL,
		UserAgent:        "test-agent",
		HTTPTimeout:      5 * time.Second,
		FinalizeInvoices: true,
	})

	result, err := client.CreateInvoice(context.Background(), Invoice{
		VoucherDate:   "2026-01-01",
		TaxConditions: TaxConditionGross(),
		ShippingTerms: ShippingTermNone(),
		TotalPrice:    TotalPrice{Currency: "EUR"},
	}, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), "inv-default", result.ID)
}
