package models

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
