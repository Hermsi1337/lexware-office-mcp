package lexware

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

type InvoiceDetail struct {
	ID             string       `json:"id"`
	OrganizationID string       `json:"organizationId"`
	VoucherStatus  string       `json:"voucherStatus"`
	VoucherNumber  string       `json:"voucherNumber"`
	VoucherDate    string       `json:"voucherDate"`
	DueDate        string       `json:"dueDate,omitempty"`
	Address        Address      `json:"address"`
	LineItems      []LineItem   `json:"lineItems"`
	TotalPrice     TotalPrice   `json:"totalPrice"`
	TaxConditions  TaxCondition `json:"taxConditions"`
	PaymentTerms   PaymentTerms `json:"paymentConditions"`
	ShippingTerms  ShippingTerm `json:"shippingConditions"`
	Title          string       `json:"title,omitempty"`
	Introduction   string       `json:"introduction,omitempty"`
	Remark         string       `json:"remark,omitempty"`
	Archived       bool         `json:"archived,omitempty"`
	Version        int          `json:"version"`
	CreatedDate    string       `json:"createdDate,omitempty"`
	UpdatedDate    string       `json:"updatedDate,omitempty"`
}
