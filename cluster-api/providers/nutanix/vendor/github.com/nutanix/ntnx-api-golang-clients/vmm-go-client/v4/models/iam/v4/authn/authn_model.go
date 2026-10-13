/*
 * Generated file models/iam/v4/authn/authn_model.go.
 *
 * Product version: 4.0.1-beta-1
 *
 * Part of the Nutanix VMM APIs
 *
 * (c) 2024 Nutanix Inc.  All rights reserved
 *
 */

/*
  Module iam.v4.authn of Nutanix VMM APIs
*/
package authn

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	import2 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/common/v1/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/common/v1/response"
	"time"
)

/*
Information of Bucket Access Key.
*/
type BucketsAccessKey struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Name of the Bucket Access Key.
	*/
	AccessKeyName *string `json:"accessKeyName"`
	/*
	  Creation time for the Bucket Access Key.
	*/
	CreatedTime *time.Time `json:"createdTime,omitempty"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import1.ApiLink `json:"links,omitempty"`
	/*
	  Secret Access Key, it will be returned only during Bucket Access Key creation.
	*/
	SecretAccessKey *string `json:"secretAccessKey,omitempty"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`
	/*
	  User Identifier who owns the Bucket Access Key.
	*/
	UserId *string `json:"userId"`
}

func (p *BucketsAccessKey) MarshalJSON() ([]byte, error) {
	type BucketsAccessKeyProxy BucketsAccessKey
	return json.Marshal(struct {
		*BucketsAccessKeyProxy
		AccessKeyName *string `json:"accessKeyName,omitempty"`
		UserId        *string `json:"userId,omitempty"`
	}{
		BucketsAccessKeyProxy: (*BucketsAccessKeyProxy)(p),
		AccessKeyName:         p.AccessKeyName,
		UserId:                p.UserId,
	})
}

func NewBucketsAccessKey() *BucketsAccessKey {
	p := new(BucketsAccessKey)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "iam.v4.authn.BucketsAccessKey"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Information of the User.
*/
type User struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Any additional attribute for the User.
	*/
	AdditionalAttributes []import2.KVPair `json:"additionalAttributes,omitempty"`
	/*
	  Bucket Access Keys for the User.
	*/
	BucketsAccessKeys []BucketsAccessKey `json:"bucketsAccessKeys,omitempty"`
	/*
	  User or Service who created the User.
	*/
	CreatedBy *string `json:"createdBy,omitempty"`
	/*
	  Creation time of the User.
	*/
	CreatedTime *time.Time `json:"createdTime,omitempty"`
	/*
	  Display name for the User.
	*/
	DisplayName *string `json:"displayName,omitempty"`
	/*
	  Email Id for the User.
	*/
	EmailId *string `json:"emailId,omitempty"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  First name for the User.
	*/
	FirstName *string `json:"firstName,omitempty"`
	/*
	  Identifier of the IDP for the User.
	*/
	IdpId *string `json:"idpId,omitempty"`
	/*
	  Flag to force the User to reset password.
	*/
	IsForceResetPasswordEnabled *bool `json:"isForceResetPasswordEnabled,omitempty"`
	/*
	  Last successful logged in time for the User.
	*/
	LastLoginTime *time.Time `json:"lastLoginTime,omitempty"`
	/*
	  Last name for the User.
	*/
	LastName *string `json:"lastName,omitempty"`
	/*
	  Last updated by this User ID.
	*/
	LastUpdatedBy *string `json:"lastUpdatedBy,omitempty"`
	/*
	  Last updated time of the User.
	*/
	LastUpdatedTime *time.Time `json:"lastUpdatedTime,omitempty"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import1.ApiLink `json:"links,omitempty"`
	/*
	  Default locale for the User.
	*/
	Locale *string `json:"locale,omitempty"`
	/*
	  Middle name for the User.
	*/
	MiddleInitial *string `json:"middleInitial,omitempty"`
	/*
	  Password for the User.
	*/
	Password *string `json:"password,omitempty"`
	/*
	  Default Region for the User.
	*/
	Region *string `json:"region,omitempty"`

	Status *UserStatusType `json:"status,omitempty"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`

	UserType *UserType `json:"userType,omitempty"`
	/*
	  Identifier for the User in the form an email address.
	*/
	Username *string `json:"username,omitempty"`
}

func NewUser() *User {
	p := new(User)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "iam.v4.authn.User"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Status of the User.
*/
type UserStatusType int

const (
	USERSTATUSTYPE_UNKNOWN  UserStatusType = 0
	USERSTATUSTYPE_REDACTED UserStatusType = 1
	USERSTATUSTYPE_ACTIVE   UserStatusType = 2
	USERSTATUSTYPE_INACTIVE UserStatusType = 3
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *UserStatusType) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"ACTIVE",
		"INACTIVE",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e UserStatusType) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"ACTIVE",
		"INACTIVE",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *UserStatusType) index(name string) UserStatusType {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"ACTIVE",
		"INACTIVE",
	}
	for idx := range names {
		if names[idx] == name {
			return UserStatusType(idx)
		}
	}
	return USERSTATUSTYPE_UNKNOWN
}

func (e *UserStatusType) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for UserStatusType:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *UserStatusType) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e UserStatusType) Ref() *UserStatusType {
	return &e
}

/*
Type of the User.
*/
type UserType int

const (
	USERTYPE_UNKNOWN  UserType = 0
	USERTYPE_REDACTED UserType = 1
	USERTYPE_LOCAL    UserType = 2
	USERTYPE_SAML     UserType = 3
	USERTYPE_LDAP     UserType = 4
	USERTYPE_EXTERNAL UserType = 5
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *UserType) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"LOCAL",
		"SAML",
		"LDAP",
		"EXTERNAL",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e UserType) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"LOCAL",
		"SAML",
		"LDAP",
		"EXTERNAL",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *UserType) index(name string) UserType {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"LOCAL",
		"SAML",
		"LDAP",
		"EXTERNAL",
	}
	for idx := range names {
		if names[idx] == name {
			return UserType(idx)
		}
	}
	return USERTYPE_UNKNOWN
}

func (e *UserType) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for UserType:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *UserType) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e UserType) Ref() *UserType {
	return &e
}

type FileDetail struct {
	Path        *string `json:"-"`
	ObjectType_ *string `json:"-"`
}

func NewFileDetail() *FileDetail {
	p := new(FileDetail)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "FileDetail"

	return p
}
