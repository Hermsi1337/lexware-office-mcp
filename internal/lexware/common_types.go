package lexware

const TaxRatePercentage = 19

// Page represents a paginated Lexware API response.
type Page[T any] struct {
	Content          []T  `json:"content"`
	First            bool `json:"first"`
	Last             bool `json:"last"`
	TotalPages       int  `json:"totalPages"`
	TotalElements    int  `json:"totalElements"`
	NumberOfElements int  `json:"numberOfElements"`
	Size             int  `json:"size"`
	Number           int  `json:"number"`
}

type Address struct {
	ContactID   string `json:"contactId,omitempty"`
	Supplement  string `json:"supplement,omitempty"`
	Name        string `json:"name"`
	Street      string `json:"street"`
	City        string `json:"city"`
	Zip         string `json:"zip"`
	CountryCode string `json:"countryCode"`
}

type LineItem struct {
	Type      string    `json:"type"`
	Name      string    `json:"name"`
	Quantity  int       `json:"quantity"`
	UnitName  string    `json:"unitName"`
	UnitPrice UnitPrice `json:"unitPrice"`
	Discount  float64   `json:"discountPercentage,omitempty"`
}

type UnitPrice struct {
	Currency          string  `json:"currency"`
	NetAmount         float64 `json:"netAmount,omitempty"`
	GrossAmount       float64 `json:"grossAmount,omitempty"`
	TaxRatePercentage float64 `json:"taxRatePercentage"`
}

type TotalPrice struct {
	Currency              string   `json:"currency"`
	TotalNetAmount        float64  `json:"totalNetAmount,omitempty"`
	TotalGrossAmount      float64  `json:"totalGrossAmount,omitempty"`
	TotalTaxAmount        float64  `json:"totalTaxAmount,omitempty"`
	TotalDiscountAbsolute *float64 `json:"totalDiscountAbsolute,omitempty"`
}

type TaxCondition struct {
	TaxType string `json:"taxType"`
}

type PaymentTerms struct {
	PaymentTermLabel    string `json:"paymentTermLabel"`
	PaymentTermDuration int    `json:"paymentTermDuration"`
}

type ShippingTerm struct {
	ShippingType string `json:"shippingType"`
}

func TaxConditionGross() TaxCondition {
	return TaxCondition{TaxType: "gross"}
}

func ShippingTermNone() ShippingTerm {
	return ShippingTerm{ShippingType: "none"}
}
