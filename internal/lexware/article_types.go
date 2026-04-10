package lexware

type Article struct {
	Title         string        `json:"title"`
	Type          string        `json:"type"`
	UnitName      string        `json:"unitName"`
	ArticleNumber string        `json:"articleNumber,omitempty"`
	Description   string        `json:"description,omitempty"`
	Gtin          string        `json:"gtin,omitempty"`
	Note          string        `json:"note,omitempty"`
	Price         *ArticlePrice `json:"price,omitempty"`
}

type ArticleDetail struct {
	ID             string        `json:"id"`
	OrganizationID string        `json:"organizationId"`
	Version        int           `json:"version"`
	Title          string        `json:"title"`
	Type           string        `json:"type"`
	UnitName       string        `json:"unitName"`
	ArticleNumber  string        `json:"articleNumber,omitempty"`
	Description    string        `json:"description,omitempty"`
	Gtin           string        `json:"gtin,omitempty"`
	Note           string        `json:"note,omitempty"`
	Price          *ArticlePrice `json:"price,omitempty"`
	CreatedDate    string        `json:"createdDate,omitempty"`
	UpdatedDate    string        `json:"updatedDate,omitempty"`
	Archived       bool          `json:"archived,omitempty"`
}

type ArticlePrice struct {
	NetPrice     float64 `json:"netPrice,omitempty"`
	GrossPrice   float64 `json:"grossPrice,omitempty"`
	LeadingPrice string  `json:"leadingPrice"`
	TaxRate      float64 `json:"taxRate"`
}
