package lexware

type PaymentConditionItem struct {
	ID                   string  `json:"id"`
	OrganizationID       string  `json:"organizationId"`
	PaymentConditionName string  `json:"paymentConditionName"`
	PaymentConditionType string  `json:"paymentConditionType"`
	DueDays              int     `json:"dueDays"`
	DiscountPercentage   float64 `json:"discountPercentage,omitempty"`
	DiscountRange        int     `json:"discountRange,omitempty"`
}
