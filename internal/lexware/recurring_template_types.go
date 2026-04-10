package lexware

type RecurringTemplateDetail struct {
	ID              string       `json:"id"`
	OrganizationID  string       `json:"organizationId"`
	Title           string       `json:"title,omitempty"`
	Address         Address      `json:"address"`
	LineItems       []LineItem   `json:"lineItems"`
	TotalPrice      TotalPrice   `json:"totalPrice"`
	TaxConditions   TaxCondition `json:"taxConditions"`
	PaymentTerms    PaymentTerms `json:"paymentConditions,omitempty"`
	ShippingTerms   ShippingTerm `json:"shippingConditions"`
	Introduction    string       `json:"introduction,omitempty"`
	Remark          string       `json:"remark,omitempty"`
	NextVoucherDate string       `json:"nextVoucherDate,omitempty"`
	RecurringCycle  string       `json:"recurringCycle,omitempty"`
	CreatedDate     string       `json:"createdDate,omitempty"`
	UpdatedDate     string       `json:"updatedDate,omitempty"`
}
