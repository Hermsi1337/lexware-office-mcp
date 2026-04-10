package lexware

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type PaymentConditionSuite struct{ baseSuite }

func TestPaymentConditionSuite(t *testing.T) { suite.Run(t, new(PaymentConditionSuite)) }

func (s *PaymentConditionSuite) TestListPaymentConditions() {
	s.mux.HandleFunc("GET /v1/payment-conditions", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]PaymentConditionItem{
			{ID: "pc-1", PaymentConditionName: "Net 30", DueDays: 30},
			{ID: "pc-2", PaymentConditionName: "Immediate", DueDays: 0},
		})
	})

	conditions, err := s.client.ListPaymentConditions(context.Background())
	require.NoError(s.T(), err)
	require.Len(s.T(), conditions, 2)
	require.Equal(s.T(), "Net 30", conditions[0].PaymentConditionName)
}
