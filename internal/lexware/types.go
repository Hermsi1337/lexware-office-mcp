package lexware

const TaxRatePercentage = 19

type Contact struct {
	Version int    `json:"version"`
	Roles   Roles  `json:"roles"`
	Person  Person `json:"person"`
	Note    string `json:"note,omitempty"`
}

type Profile struct {
	OrganizationID string `json:"organizationId"`
}

type Person struct {
	Salutation string `json:"salutation,omitempty"`
	FirstName  string `json:"firstName,omitempty"`
	LastName   string `json:"lastName"`
}

type Roles struct {
	Customer map[string]any `json:"customer,omitempty"`
}

type Invoice struct {
	Archived      bool         `json:"archived,omitempty"`
	VoucherDate   string       `json:"voucherDate"`
	Address       Address      `json:"address"`
	LineItems     []LineItem   `json:"lineItems"`
	TotalPrice    TotalPrice   `json:"totalPrice"`
	TaxConditions TaxCondition `json:"taxConditions"`
	PaymentTerms  PaymentTerms `json:"paymentConditions"`
	ShippingTerms ShippingTerm `json:"shippingConditions"`
	Title         string       `json:"title,omitempty"`
	Introduction  string       `json:"introduction,omitempty"`
	Remark        string       `json:"remark,omitempty"`
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
	GrossAmount       float64 `json:"grossAmount"`
	TaxRatePercentage float64 `json:"taxRatePercentage"`
}

type TotalPrice struct {
	Currency              string   `json:"currency"`
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
