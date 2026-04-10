package lexware

const TaxRatePercentage = 19

// ---------- Pagination ----------

// Page represents a paginated Lexware API response.
type Page[T any] struct {
	Content          []T `json:"content"`
	First            bool `json:"first"`
	Last             bool `json:"last"`
	TotalPages       int  `json:"totalPages"`
	TotalElements    int  `json:"totalElements"`
	NumberOfElements int  `json:"numberOfElements"`
	Size             int  `json:"size"`
	Number           int  `json:"number"`
}

// ---------- Profile ----------

type Profile struct {
	OrganizationID string `json:"organizationId"`
}

// ---------- Contacts ----------

type Contact struct {
	Version int    `json:"version"`
	Roles   Roles  `json:"roles"`
	Person  Person `json:"person"`
	Note    string `json:"note,omitempty"`
}

type ContactDetail struct {
	ID             string          `json:"id"`
	OrganizationID string          `json:"organizationId"`
	Version        int             `json:"version"`
	Roles          Roles           `json:"roles"`
	Company        *Company        `json:"company,omitempty"`
	Person         *Person         `json:"person,omitempty"`
	Addresses      *ContactLists   `json:"addresses,omitempty"`
	EmailAddresses *ContactLists   `json:"emailAddresses,omitempty"`
	PhoneNumbers   *ContactLists   `json:"phoneNumbers,omitempty"`
	Note           string          `json:"note,omitempty"`
	Archived       bool            `json:"archived,omitempty"`
}

type Company struct {
	Name               string `json:"name"`
	TaxNumber          string `json:"taxNumber,omitempty"`
	VatRegistrationID  string `json:"vatRegistrationId,omitempty"`
	ContactPersons     []Person `json:"contactPersons,omitempty"`
}

type Person struct {
	Salutation string `json:"salutation,omitempty"`
	FirstName  string `json:"firstName,omitempty"`
	LastName   string `json:"lastName"`
}

type Roles struct {
	Customer map[string]any `json:"customer,omitempty"`
	Vendor   map[string]any `json:"vendor,omitempty"`
}

type ContactLists struct {
	Billing  []ContactAddress `json:"billing,omitempty"`
	Shipping []ContactAddress `json:"shipping,omitempty"`
	Business []ContactEmail   `json:"business,omitempty"`
	Office   []ContactPhone   `json:"office,omitempty"`
	Private  []ContactEmail   `json:"private,omitempty"`
	Other    []ContactEmail   `json:"other,omitempty"`
	Fax      []ContactPhone   `json:"fax,omitempty"`
	Mobile   []ContactPhone   `json:"mobile,omitempty"`
}

type ContactAddress struct {
	Supplement  string `json:"supplement,omitempty"`
	Street      string `json:"street,omitempty"`
	Zip         string `json:"zip,omitempty"`
	City        string `json:"city,omitempty"`
	CountryCode string `json:"countryCode,omitempty"`
}

type ContactEmail struct {
	Email string `json:"email,omitempty"`
}

type ContactPhone struct {
	PhoneNumber string `json:"phoneNumber,omitempty"`
}

// ---------- Invoices ----------

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
	ID                string       `json:"id"`
	OrganizationID    string       `json:"organizationId"`
	VoucherStatus     string       `json:"voucherStatus"`
	VoucherNumber     string       `json:"voucherNumber"`
	VoucherDate       string       `json:"voucherDate"`
	DueDate           string       `json:"dueDate,omitempty"`
	Address           Address      `json:"address"`
	LineItems         []LineItem   `json:"lineItems"`
	TotalPrice        TotalPrice   `json:"totalPrice"`
	TaxConditions     TaxCondition `json:"taxConditions"`
	PaymentTerms      PaymentTerms `json:"paymentConditions"`
	ShippingTerms     ShippingTerm `json:"shippingConditions"`
	Title             string       `json:"title,omitempty"`
	Introduction      string       `json:"introduction,omitempty"`
	Remark            string       `json:"remark,omitempty"`
	Archived          bool         `json:"archived,omitempty"`
	Version           int          `json:"version"`
	CreatedDate       string       `json:"createdDate,omitempty"`
	UpdatedDate       string       `json:"updatedDate,omitempty"`
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

// ---------- Articles ----------

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

// ---------- Quotations ----------

type Quotation struct {
	VoucherDate   string       `json:"voucherDate"`
	ExpirationDate string      `json:"expirationDate,omitempty"`
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

// ---------- Credit Notes ----------

type CreditNote struct {
	VoucherDate   string       `json:"voucherDate"`
	Address       Address      `json:"address"`
	LineItems     []LineItem   `json:"lineItems"`
	TotalPrice    TotalPrice   `json:"totalPrice"`
	TaxConditions TaxCondition `json:"taxConditions"`
	ShippingTerms ShippingTerm `json:"shippingConditions"`
	Title         string       `json:"title,omitempty"`
	Introduction  string       `json:"introduction,omitempty"`
	Remark        string       `json:"remark,omitempty"`
}

type CreditNoteDetail struct {
	ID             string       `json:"id"`
	OrganizationID string       `json:"organizationId"`
	VoucherStatus  string       `json:"voucherStatus"`
	VoucherNumber  string       `json:"voucherNumber"`
	VoucherDate    string       `json:"voucherDate"`
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

// ---------- Countries ----------

type Country struct {
	CountryCode       string `json:"countryCode"`
	CountryNameDE     string `json:"countryNameDE"`
	CountryNameEN     string `json:"countryNameEN"`
	TaxClassification string `json:"taxClassification"`
}

// ---------- Helper constructors ----------

func TaxConditionGross() TaxCondition {
	return TaxCondition{TaxType: "gross"}
}

func ShippingTermNone() ShippingTerm {
	return ShippingTerm{ShippingType: "none"}
}
