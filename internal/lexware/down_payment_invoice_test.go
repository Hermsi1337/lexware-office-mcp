package lexware

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DownPaymentInvoiceSuite struct{ baseSuite }

func TestDownPaymentInvoiceSuite(t *testing.T) { suite.Run(t, new(DownPaymentInvoiceSuite)) }

func (s *DownPaymentInvoiceSuite) TestGetDownPaymentInvoice() {
	s.mux.HandleFunc("GET /v1/down-payment-invoices/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DownPaymentInvoiceDetail{
			ID:            r.PathValue("id"),
			VoucherStatus: "open",
			VoucherNumber: "AR-001",
		})
	})

	detail, err := s.client.GetDownPaymentInvoice(context.Background(), "dpi-99")
	require.NoError(s.T(), err)
	require.Equal(s.T(), "dpi-99", detail.ID)
	require.Equal(s.T(), "open", detail.VoucherStatus)
}
