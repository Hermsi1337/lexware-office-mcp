package lexware

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CountrySuite struct{ baseSuite }

func TestCountrySuite(t *testing.T) { suite.Run(t, new(CountrySuite)) }

func (s *CountrySuite) TestListCountries() {
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
