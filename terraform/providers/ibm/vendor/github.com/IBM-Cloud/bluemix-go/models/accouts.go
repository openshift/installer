package models

import "time"

type V2Account struct {
	Guid          string
	Name          string
	Type          string
	State         string
	OwnerIamId    string
	CountryCode   string
	Organizations []AccountOrganization
}

type AccountOrganization struct {
	GUID   string `json:"guid"`
	Region string `json:"region"`
}

type AccountUser struct {
	UserId      string `json:"userId"`
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	State       string `json:"state"`
	IbmUniqueId string `json:"ibmUniqueId"`
	Email       string `json:"email"`
	Phonenumber string `json:"phonenumber"`
	Id          string `json:"id"`
	UaaGuid     string `json:"uaaGuid"`
	AccountId   string `json:"accountId"`
	Role        string `json:"role"`
	AddedOn     string `json:"added_on"`
	InvitedOn   string `json:"invitedOn"`
	Photo       string `json:"photo"`
}

type V1Account struct {
	Metadata Metadata      `json:"metadata"`
	Entity   AccountEntity `json:"entity"`
}

type AccountEntity struct {
	BillingCountryCode   string                `json:"billing_country_code"`
	BluemixSubscriptions []BluemixSubscription `json:"bluemix_subscriptions"`
	ConfigurationID      *string               `json:"configuration_id"` // Using pointer to handle null
	CountryCode          string                `json:"country_code"`
	CurrencyCode         string                `json:"currency_code"`
	CurrentBillingSystem string                `json:"current_billing_system"`
	CustomerID           string                `json:"customer_id"`
	IsIBMer              bool                  `json:"isIBMer"`
	Linkages             []AccountLinkage      `json:"linkages"`
	Name                 string                `json:"name"`
	OfferTemplate        string                `json:"offer_template"`
	Onboarded            int                   `json:"onboarded"`
	OrganizationsRegion  []OrganizationsRegion `json:"organizations_region"`
	Origin               string                `json:"origin"`
	Owner                string                `json:"owner"`
	OwnerIAMID           string                `json:"owner_iam_id"`
	OwnerUniqueID        string                `json:"owner_unique_id"`
	OwnerUserID          string                `json:"owner_userid"`
	State                string                `json:"state"`
	SubscriptionID       string                `json:"subscription_id"`
	Tags                 []interface{}         `json:"tags"` // Using interface{} since type is not specified
	TeamDirectoryEnabled bool                  `json:"team_directory_enabled"`
	TermsAndConditions   TermsAndConditions    `json:"terms_and_conditions"`
	Type                 string                `json:"type"`
}

type BluemixSubscription struct {
	BillToContact         string        `json:"bill_to_contact"`
	BillingSystem         string        `json:"billing_system"`
	CatalogID             string        `json:"catalog_id"`
	CurrentStateTimestamp time.Time     `json:"current_state_timestamp"`
	DistributionChannel   string        `json:"distribution_channel"`
	History               []History     `json:"history"`
	OrderID               string        `json:"order_id"`
	PartNumber            string        `json:"part_number"`
	PaygPendingTimestamp  time.Time     `json:"payg_pending_timestamp"`
	PaymentMethod         PaymentMethod `json:"payment_method"`
	SoftlayerAccountID    string        `json:"softlayer_account_id"`
	SoldToContact         string        `json:"sold_to_contact"`
	State                 string        `json:"state"`
	SubscriptionTags      []interface{} `json:"subscriptionTags"`
	SubscriptionID        string        `json:"subscription_id"`
	Type                  string        `json:"type"`
}

type History struct {
	BillingCountryCode string    `json:"billingCountryCode"`
	BillingSystem      string    `json:"billingSystem"`
	CountryCode        string    `json:"countryCode"`
	CurrencyCode       string    `json:"currencyCode"`
	EndTime            time.Time `json:"endTime"`
	StartTime          time.Time `json:"startTime"`
	State              string    `json:"state"`
	Type               string    `json:"type"`
	BillToContact      *string   `json:"billToContact,omitempty"`
	OrderId            *string   `json:"orderId,omitempty"`
	PaymentMethodType  *string   `json:"paymentMethodType,omitempty"`
	SoldToContact      *string   `json:"soldToContact,omitempty"`
	StateEndComments   *string   `json:"stateEndComments,omitempty"`
	StateEndedBy       *string   `json:"stateEndedBy,omitempty"`
	WalletId           *string   `json:"walletId,omitempty"`
}

type PaymentMethod struct {
	Ended    *string   `json:"ended"` // Using pointer to handle null
	Started  time.Time `json:"started"`
	Type     string    `json:"type"`
	WalletID string    `json:"wallet_id"`
}

type AccountLinkage struct {
	Origin string `json:"origin"`
	State  string `json:"state"`
}

type OrganizationsRegion struct {
	GUID   string `json:"guid"`
	Region string `json:"region"`
}

type TermsAndConditions struct {
	Accepted  bool      `json:"accepted"`
	Required  bool      `json:"required"`
	Timestamp time.Time `json:"timestamp"`
}

type Metadata struct {
	CreatedAt      time.Time `json:"created_at"`
	GUID           string    `json:"guid"`
	UpdateComments string    `json:"update_comments"`
	UpdatedAt      time.Time `json:"updated_at"`
	UpdatedBy      string    `json:"updated_by"`
	URL            string    `json:"url"`
}
