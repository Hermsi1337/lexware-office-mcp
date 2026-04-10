package lexware

type Quotation struct {
	VoucherDate    string       `json:"voucherDate"`
	ExpirationDate string       `json:"expirationDate,omitempty"`
	Address        Address      `json:"address"`
	LineItems      []LineItem   `json:"lineItems"`
	TotalPrice     TotalPrice   `json:"totalPrice"`
	TaxConditions  TaxCondition `json:"taxConditions"`
	PaymentTerms   PaymentTerms `json:"paymentConditions,omitempty"`
	ShippingTerms  ShippingTerm `json:"shippingConditions"`
	Title          string       `json:"title,omitempty"`
	Introduction   string       `json:"introduction,omitempty"`
	Remark         string       `json:"remark,omitempty"`
}

type QuotationDetail struct {
	ID             string       `json:"id"`
	OrganizationID string       `json:"organizationId"`
	VoucherStatus  string       `json:"voucherStatus"`
	VoucherNumber  string       `json:"voucherNumber"`
	VoucherDate    string       `json:"voucherDate"`
	ExpirationDate string       `json:"expirationDate,omitempty"`
	Address        Address      `json:"address"`
	LineItems      []LineItem   `json:"lineItems"`
	TotalPrice     TotalPrice   `json:"totalPrice"`
	TaxConditions  TaxCondition `json:"taxConditions"`
	ShippingTerms  ShippingTerm `json:"shippingConditions"`
	Title          string       `json:"title,omitempty"`
	Introduction   string       `json:"introduction,omitempty"`
	Remark         string       `json:"remark,omitempty"`
	Version        int          `json:"version"`
	CreatedDate    string       `json:"createdDate,omitempty"`
	UpdatedDate    string       `json:"updatedDate,omitempty"`
}
