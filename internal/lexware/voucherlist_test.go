package lexware

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type VoucherlistSuite struct{ baseSuite }

func TestVoucherlistSuite(t *testing.T) { suite.Run(t, new(VoucherlistSuite)) }

func (s *VoucherlistSuite) TestListVouchers() {
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
