package lexware

type VoucherListItem struct {
	VoucherID     string  `json:"voucherId"`
	VoucherType   string  `json:"voucherType"`
	VoucherStatus string  `json:"voucherStatus"`
	VoucherNumber string  `json:"voucherNumber"`
	VoucherDate   string  `json:"voucherDate"`
	DueDate       string  `json:"dueDate,omitempty"`
	ContactID     string  `json:"contactId,omitempty"`
	ContactName   string  `json:"contactName,omitempty"`
	TotalAmount   float64 `json:"totalAmount"`
	OpenAmount    float64 `json:"openAmount,omitempty"`
	Currency      string  `json:"currency"`
	Archived      bool    `json:"archived,omitempty"`
	CreatedDate   string  `json:"createdDate,omitempty"`
	UpdatedDate   string  `json:"updatedDate,omitempty"`
}
