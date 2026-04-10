package lexware

type Contact struct {
	Version int    `json:"version"`
	Roles   Roles  `json:"roles"`
	Person  Person `json:"person"`
	Note    string `json:"note,omitempty"`
}

type ContactDetail struct {
	ID             string        `json:"id"`
	OrganizationID string        `json:"organizationId"`
	Version        int           `json:"version"`
	Roles          Roles         `json:"roles"`
	Company        *Company      `json:"company,omitempty"`
	Person         *Person       `json:"person,omitempty"`
	Addresses      *ContactLists `json:"addresses,omitempty"`
	EmailAddresses *ContactLists `json:"emailAddresses,omitempty"`
	PhoneNumbers   *ContactLists `json:"phoneNumbers,omitempty"`
	Note           string        `json:"note,omitempty"`
	Archived       bool          `json:"archived,omitempty"`
}

type Company struct {
	Name              string   `json:"name"`
	TaxNumber         string   `json:"taxNumber,omitempty"`
	VatRegistrationID string   `json:"vatRegistrationId,omitempty"`
	ContactPersons    []Person `json:"contactPersons,omitempty"`
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
