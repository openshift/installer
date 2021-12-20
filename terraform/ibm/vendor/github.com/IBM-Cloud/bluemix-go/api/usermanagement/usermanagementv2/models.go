package usermanagementv2

import "github.com/IBM-Cloud/bluemix-go/api/iampap/iampapv1"

// User ...
type User struct {
	Email       string `json:"email"`
	AccountRole string `json:"account_role"`
}

// UserInfo contains user info
type UserInfo struct {
	ID             string `json:"id"`
	IamID          string `json:"iam_id"`
	Realm          string `json:"realm"`
	UserID         string `json:"user_id"`
	Firstname      string `json:"firstname"`
	Lastname       string `json:"lastname"`
	State          string `json:"state"`
	Email          string `json:"email"`
	Phonenumber    string `json:"phonenumber"`
	Altphonenumber string `json:"altphonenumber"`
	Photo          string `json:"photo"`
	AccountID      string `json:"account_id"`
}

// UserSettings ...
type UserSettingOptions struct {

	//The console UI language
	Language string `json:"language"`

	//The language for email and phone notifications.
	NotificationLanguage string `json:"notification_language"`

	//The IP addresses listed here are the only ones from which this user can log in to IBM Cloud.
	AllowedIPAddresses string `json:"allowed_ip_addresses"`

	//Whether user-managed login is enabled.
	SelfManage bool `json:"self_manage"`
}

// UserInvite ...
type UserInvite struct {
	Users               []User            `json:"users"`
	IAMPolicy           []UserPolicy      `json:"iam_policy,omitempty"`
	AccessGroup         []string          `json:"access_groups,omitempty"`
	InfrastructureRoles *InfraPermissions `json:"infrastructure_roles,omitempty"`
	OrganizationRoles   []OrgRole         `json:"organization_roles,omitempty"`
}

// UsersList to get list of users
type UsersList struct {
	TotalUsers int        `json:"total_results"`
	Limit      int        `json:"limit"`
	FistURL    string     `json:"fist_url"`
	Resources  []UserInfo `json:"resources"`
}

// UserPolicy ...
type UserPolicy struct {
	Type      string              `json:"type"`
	Roles     []iampapv1.Role     `json:"roles"`
	Resources []iampapv1.Resource `json:"resources"`
}

//InfraPermissions ...
type InfraPermissions struct {
	Permissions []string `json:"permissions"`
}

//OrgRole ...
type OrgRole struct {
	Users           []string `json:"users"`
	Region          string   `json:"region"`
	Auditors        []string `json:"auditors,omitempty"`
	Managers        []string `json:"managers,omitempty"`
	BillingManagers []string `json:"billing_managers,omitempty"`
	ID              string   `json:"id"`
	Spaces          []Space  `json:"spaces"`
}

//Space ...
type Space struct {
	ID         string   `json:"id"`
	Managers   []string `json:"managers,omitempty"`
	Developers []string `json:"developers,omitempty"`
	Auditors   []string `json:"auditors,omitempty" `
}
