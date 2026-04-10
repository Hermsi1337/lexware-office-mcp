package lexware

type OrderConfirmation struct {
	VoucherDate   string       `json:"voucherDate"`
	Address       Address      `json:"address"`
	LineItems     []LineItem   `json:"lineItems"`
	TotalPrice    TotalPrice   `json:"totalPrice"`
	TaxConditions TaxCondition `json:"taxConditions"`
	PaymentTerms  PaymentTerms `json:"paymentConditions,omitempty"`
	ShippingTerms ShippingTerm `json:"shippingConditions"`
	Title         string       `json:"title,omitempty"`
	Introduction  string       `json:"introduction,omitempty"`
	Remark        string       `json:"remark,omitempty"`
}

type OrderConfirmationDetail struct {
	ID             string       `json:"id"`
	OrganizationID string       `json:"organizationId"`
	VoucherStatus  string       `json:"voucherStatus"`
	VoucherNumber  string       `json:"voucherNumber"`
	VoucherDate    string       `json:"voucherDate"`
	Address        Address      `json:"address"`
	LineItems      []LineItem   `json:"lineItems"`
	TotalPrice     TotalPrice   `json:"totalPrice"`
	TaxConditions  TaxCondition `json:"taxConditions"`
	PaymentTerms   PaymentTerms `json:"paymentConditions,omitempty"`
	ShippingTerms  ShippingTerm `json:"shippingConditions"`
	Title          string       `json:"title,omitempty"`
	Introduction   string       `json:"introduction,omitempty"`
	Remark         string       `json:"remark,omitempty"`
	Version        int          `json:"version"`
	CreatedDate    string       `json:"createdDate,omitempty"`
	UpdatedDate    string       `json:"updatedDate,omitempty"`
}
