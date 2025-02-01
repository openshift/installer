/*
 * Generated file models/vmm/v4/images/config/config_model.go.
 *
 * Product version: 4.0.1-beta-1
 *
 * Part of the Nutanix VMM APIs
 *
 * (c) 2024 Nutanix Inc.  All rights reserved
 *
 */

/*
  Module vmm.v4.images.config of Nutanix VMM APIs
*/
package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	import3 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/common/v1/response"
	import1 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/prism/v4/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/error"
	"time"
)

/*
REST response for all response codes in API path /vmm/v4.0.b1/images/config/placement-policies Post operation
*/
type CreatePlacementPolicyApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfCreatePlacementPolicyApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewCreatePlacementPolicyApiResponse() *CreatePlacementPolicyApiResponse {
	p := new(CreatePlacementPolicyApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "vmm.v4.images.config.CreatePlacementPolicyApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *CreatePlacementPolicyApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *CreatePlacementPolicyApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfCreatePlacementPolicyApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /vmm/v4.0.b1/images/config/rate-limit-policies Post operation
*/
type CreateRateLimitPolicyApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfCreateRateLimitPolicyApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewCreateRateLimitPolicyApiResponse() *CreateRateLimitPolicyApiResponse {
	p := new(CreateRateLimitPolicyApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "vmm.v4.images.config.CreateRateLimitPolicyApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *CreateRateLimitPolicyApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *CreateRateLimitPolicyApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfCreateRateLimitPolicyApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /vmm/v4.0.b1/images/config/placement-policies/{extId} Delete operation
*/
type DeletePlacementPolicyApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfDeletePlacementPolicyApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewDeletePlacementPolicyApiResponse() *DeletePlacementPolicyApiResponse {
	p := new(DeletePlacementPolicyApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "vmm.v4.images.config.DeletePlacementPolicyApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *DeletePlacementPolicyApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *DeletePlacementPolicyApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfDeletePlacementPolicyApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /vmm/v4.0.b1/images/config/rate-limit-policies/{extId} Delete operation
*/
type DeleteRateLimitPolicyApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfDeleteRateLimitPolicyApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewDeleteRateLimitPolicyApiResponse() *DeleteRateLimitPolicyApiResponse {
	p := new(DeleteRateLimitPolicyApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "vmm.v4.images.config.DeleteRateLimitPolicyApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *DeleteRateLimitPolicyApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *DeleteRateLimitPolicyApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfDeleteRateLimitPolicyApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
The enforced rate limit on cluster when multiple image rate limit policies apply.
*/
type EffectiveRateLimitPolicy struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Cluster external identifier.
	*/
	ClusterExtId *string `json:"clusterExtId"`
	/*
	  Image rate limit policy external identifier.
	*/
	RateLimitExtId *string `json:"rateLimitExtId"`
}

func (p *EffectiveRateLimitPolicy) MarshalJSON() ([]byte, error) {
	type EffectiveRateLimitPolicyProxy EffectiveRateLimitPolicy
	return json.Marshal(struct {
		*EffectiveRateLimitPolicyProxy
		ClusterExtId   *string `json:"clusterExtId,omitempty"`
		RateLimitExtId *string `json:"rateLimitExtId,omitempty"`
	}{
		EffectiveRateLimitPolicyProxy: (*EffectiveRateLimitPolicyProxy)(p),
		ClusterExtId:                  p.ClusterExtId,
		RateLimitExtId:                p.RateLimitExtId,
	})
}

func NewEffectiveRateLimitPolicy() *EffectiveRateLimitPolicy {
	p := new(EffectiveRateLimitPolicy)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "vmm.v4.images.config.EffectiveRateLimitPolicy"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Enforcement status of the image placement policy.
*/
type EnforcementState int

const (
	ENFORCEMENTSTATE_UNKNOWN   EnforcementState = 0
	ENFORCEMENTSTATE_REDACTED  EnforcementState = 1
	ENFORCEMENTSTATE_ACTIVE    EnforcementState = 2
	ENFORCEMENTSTATE_SUSPENDED EnforcementState = 3
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *EnforcementState) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"ACTIVE",
		"SUSPENDED",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e EnforcementState) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"ACTIVE",
		"SUSPENDED",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *EnforcementState) index(name string) EnforcementState {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"ACTIVE",
		"SUSPENDED",
	}
	for idx := range names {
		if names[idx] == name {
			return EnforcementState(idx)
		}
	}
	return ENFORCEMENTSTATE_UNKNOWN
}

func (e *EnforcementState) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for EnforcementState:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *EnforcementState) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e EnforcementState) Ref() *EnforcementState {
	return &e
}

/*
Category-based entity filter.
*/
type Filter struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	CategoryExtIds []string `json:"categoryExtIds,omitempty"`

	Type *FilterMatchType `json:"type"`
}

func (p *Filter) MarshalJSON() ([]byte, error) {
	type FilterProxy Filter
	return json.Marshal(struct {
		*FilterProxy
		Type *FilterMatchType `json:"type,omitempty"`
	}{
		FilterProxy: (*FilterProxy)(p),
		Type:        p.Type,
	})
}

func NewFilter() *Filter {
	p := new(Filter)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "vmm.v4.images.config.Filter"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Filter matching type.
*/
type FilterMatchType int

const (
	FILTERMATCHTYPE_UNKNOWN              FilterMatchType = 0
	FILTERMATCHTYPE_REDACTED             FilterMatchType = 1
	FILTERMATCHTYPE_CATEGORIES_MATCH_ALL FilterMatchType = 2
	FILTERMATCHTYPE_CATEGORIES_MATCH_ANY FilterMatchType = 3
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *FilterMatchType) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"CATEGORIES_MATCH_ALL",
		"CATEGORIES_MATCH_ANY",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e FilterMatchType) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"CATEGORIES_MATCH_ALL",
		"CATEGORIES_MATCH_ANY",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *FilterMatchType) index(name string) FilterMatchType {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"CATEGORIES_MATCH_ALL",
		"CATEGORIES_MATCH_ANY",
	}
	for idx := range names {
		if names[idx] == name {
			return FilterMatchType(idx)
		}
	}
	return FILTERMATCHTYPE_UNKNOWN
}

func (e *FilterMatchType) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for FilterMatchType:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *FilterMatchType) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e FilterMatchType) Ref() *FilterMatchType {
	return &e
}

/*
REST response for all response codes in API path /vmm/v4.0.b1/images/config/placement-policies/{extId} Get operation
*/
type GetPlacementPolicyApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfGetPlacementPolicyApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewGetPlacementPolicyApiResponse() *GetPlacementPolicyApiResponse {
	p := new(GetPlacementPolicyApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "vmm.v4.images.config.GetPlacementPolicyApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *GetPlacementPolicyApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *GetPlacementPolicyApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfGetPlacementPolicyApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /vmm/v4.0.b1/images/config/rate-limit-policies/{extId} Get operation
*/
type GetRateLimitPolicyApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfGetRateLimitPolicyApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewGetRateLimitPolicyApiResponse() *GetRateLimitPolicyApiResponse {
	p := new(GetRateLimitPolicyApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "vmm.v4.images.config.GetRateLimitPolicyApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *GetRateLimitPolicyApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *GetRateLimitPolicyApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfGetRateLimitPolicyApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /vmm/v4.0.b1/images/config/effective-rate-limit-policies Get operation
*/
type ListEffectiveRateLimitPoliciesApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfListEffectiveRateLimitPoliciesApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewListEffectiveRateLimitPoliciesApiResponse() *ListEffectiveRateLimitPoliciesApiResponse {
	p := new(ListEffectiveRateLimitPoliciesApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "vmm.v4.images.config.ListEffectiveRateLimitPoliciesApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *ListEffectiveRateLimitPoliciesApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *ListEffectiveRateLimitPoliciesApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfListEffectiveRateLimitPoliciesApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /vmm/v4.0.b1/images/config/placement-policies Get operation
*/
type ListPlacementPoliciesApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfListPlacementPoliciesApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewListPlacementPoliciesApiResponse() *ListPlacementPoliciesApiResponse {
	p := new(ListPlacementPoliciesApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "vmm.v4.images.config.ListPlacementPoliciesApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *ListPlacementPoliciesApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *ListPlacementPoliciesApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfListPlacementPoliciesApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /vmm/v4.0.b1/images/config/rate-limit-policies Get operation
*/
type ListRateLimitPoliciesApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfListRateLimitPoliciesApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewListRateLimitPoliciesApiResponse() *ListRateLimitPoliciesApiResponse {
	p := new(ListRateLimitPoliciesApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "vmm.v4.images.config.ListRateLimitPoliciesApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *ListRateLimitPoliciesApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *ListRateLimitPoliciesApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfListRateLimitPoliciesApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

type PlacementPolicy struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	ClusterEntityFilter *Filter `json:"clusterEntityFilter"`
	/*
	  Create time of an image placement policy.
	*/
	CreateTime *time.Time `json:"createTime,omitempty"`
	/*
	  Description of the image placement policy.
	*/
	Description *string `json:"description,omitempty"`

	EnforcementState *EnforcementState `json:"enforcementState,omitempty"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`

	ImageEntityFilter *Filter `json:"imageEntityFilter"`
	/*
	  Last update time of an image placement policy.
	*/
	LastUpdateTime *time.Time `json:"lastUpdateTime,omitempty"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import3.ApiLink `json:"links,omitempty"`
	/*
	  Name of the image placement policy.
	*/
	Name *string `json:"name"`
	/*
	  External identifier of the owner of the image placement policy.
	*/
	OwnerExtId *string `json:"ownerExtId,omitempty"`

	PlacementType *PlacementType `json:"placementType"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`
}

func (p *PlacementPolicy) MarshalJSON() ([]byte, error) {
	type PlacementPolicyProxy PlacementPolicy
	return json.Marshal(struct {
		*PlacementPolicyProxy
		ClusterEntityFilter *Filter        `json:"clusterEntityFilter,omitempty"`
		ImageEntityFilter   *Filter        `json:"imageEntityFilter,omitempty"`
		Name                *string        `json:"name,omitempty"`
		PlacementType       *PlacementType `json:"placementType,omitempty"`
	}{
		PlacementPolicyProxy: (*PlacementPolicyProxy)(p),
		ClusterEntityFilter:  p.ClusterEntityFilter,
		ImageEntityFilter:    p.ImageEntityFilter,
		Name:                 p.Name,
		PlacementType:        p.PlacementType,
	})
}

func NewPlacementPolicy() *PlacementPolicy {
	p := new(PlacementPolicy)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "vmm.v4.images.config.PlacementPolicy"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Type of the image placement policy.
*/
type PlacementType int

const (
	PLACEMENTTYPE_UNKNOWN  PlacementType = 0
	PLACEMENTTYPE_REDACTED PlacementType = 1
	PLACEMENTTYPE_SOFT     PlacementType = 2
	PLACEMENTTYPE_HARD     PlacementType = 3
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *PlacementType) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"SOFT",
		"HARD",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e PlacementType) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"SOFT",
		"HARD",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *PlacementType) index(name string) PlacementType {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"SOFT",
		"HARD",
	}
	for idx := range names {
		if names[idx] == name {
			return PlacementType(idx)
		}
	}
	return PLACEMENTTYPE_UNKNOWN
}

func (e *PlacementType) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for PlacementType:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *PlacementType) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e PlacementType) Ref() *PlacementType {
	return &e
}

type RateLimitPolicy struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	ClusterEntityFilter *Filter `json:"clusterEntityFilter"`
	/*
	  Create time of a rate limit policy.
	*/
	CreateTime *time.Time `json:"createTime,omitempty"`
	/*
	  Image rate limit policy specification.
	*/
	Description *string `json:"description,omitempty"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  Last update time of a rate limit policy.
	*/
	LastUpdateTime *time.Time `json:"lastUpdateTime,omitempty"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import3.ApiLink `json:"links,omitempty"`
	/*
	  External identifier of the Prism Elements where a rate limit is the effective rate limit policy.
	*/
	MatchingClusterExtIds []string `json:"matchingClusterExtIds,omitempty"`
	/*
	  Name of the image rate limit policy.
	*/
	Name *string `json:"name"`
	/*
	  External identifier of the owner of the rate limit policy.
	*/
	OwnerExtId *string `json:"ownerExtId,omitempty"`
	/*
	  Network bandwidth in KBps that the rate limited image operation can utilize.
	*/
	RateLimitKbps *int64 `json:"rateLimitKbps"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`
}

func (p *RateLimitPolicy) MarshalJSON() ([]byte, error) {
	type RateLimitPolicyProxy RateLimitPolicy
	return json.Marshal(struct {
		*RateLimitPolicyProxy
		ClusterEntityFilter *Filter `json:"clusterEntityFilter,omitempty"`
		Name                *string `json:"name,omitempty"`
		RateLimitKbps       *int64  `json:"rateLimitKbps,omitempty"`
	}{
		RateLimitPolicyProxy: (*RateLimitPolicyProxy)(p),
		ClusterEntityFilter:  p.ClusterEntityFilter,
		Name:                 p.Name,
		RateLimitKbps:        p.RateLimitKbps,
	})
}

func NewRateLimitPolicy() *RateLimitPolicy {
	p := new(RateLimitPolicy)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "vmm.v4.images.config.RateLimitPolicy"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
REST response for all response codes in API path /vmm/v4.0.b1/images/config/placement-policies/{extId}/$actions/resume Post operation
*/
type ResumePlacementPolicyApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfResumePlacementPolicyApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewResumePlacementPolicyApiResponse() *ResumePlacementPolicyApiResponse {
	p := new(ResumePlacementPolicyApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "vmm.v4.images.config.ResumePlacementPolicyApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *ResumePlacementPolicyApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *ResumePlacementPolicyApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfResumePlacementPolicyApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /vmm/v4.0.b1/images/config/placement-policies/{extId}/$actions/suspend Post operation
*/
type SuspendPlacementPolicyApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfSuspendPlacementPolicyApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewSuspendPlacementPolicyApiResponse() *SuspendPlacementPolicyApiResponse {
	p := new(SuspendPlacementPolicyApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "vmm.v4.images.config.SuspendPlacementPolicyApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *SuspendPlacementPolicyApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *SuspendPlacementPolicyApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfSuspendPlacementPolicyApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
Indicates whether the tasks running to enforce policy should be canceled.
These would be checkout and uncheckout tasks to placement remove images from the cluster(s).
*/
type SuspendPlacementPolicyConfig struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Boolean to indicate if tasks running to enforce image placement policy should be suspended.
	*/
	ShouldCancelRunningTasks *bool `json:"shouldCancelRunningTasks,omitempty"`
}

func NewSuspendPlacementPolicyConfig() *SuspendPlacementPolicyConfig {
	p := new(SuspendPlacementPolicyConfig)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "vmm.v4.images.config.SuspendPlacementPolicyConfig"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	p.ShouldCancelRunningTasks = new(bool)
	*p.ShouldCancelRunningTasks = false

	return p
}

/*
REST response for all response codes in API path /vmm/v4.0.b1/images/config/placement-policies/{extId} Put operation
*/
type UpdatePlacementPolicyApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfUpdatePlacementPolicyApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewUpdatePlacementPolicyApiResponse() *UpdatePlacementPolicyApiResponse {
	p := new(UpdatePlacementPolicyApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "vmm.v4.images.config.UpdatePlacementPolicyApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *UpdatePlacementPolicyApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *UpdatePlacementPolicyApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfUpdatePlacementPolicyApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /vmm/v4.0.b1/images/config/rate-limit-policies/{extId} Put operation
*/
type UpdateRateLimitPolicyApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfUpdateRateLimitPolicyApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewUpdateRateLimitPolicyApiResponse() *UpdateRateLimitPolicyApiResponse {
	p := new(UpdateRateLimitPolicyApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "vmm.v4.images.config.UpdateRateLimitPolicyApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *UpdateRateLimitPolicyApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *UpdateRateLimitPolicyApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfUpdateRateLimitPolicyApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

type OneOfCreateRateLimitPolicyApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
	oneOfType0    *import1.TaskReference `json:"-"`
}

func NewOneOfCreateRateLimitPolicyApiResponseData() *OneOfCreateRateLimitPolicyApiResponseData {
	p := new(OneOfCreateRateLimitPolicyApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfCreateRateLimitPolicyApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfCreateRateLimitPolicyApiResponseData is nil"))
	}
	switch v.(type) {
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	case import1.TaskReference:
		if nil == p.oneOfType0 {
			p.oneOfType0 = new(import1.TaskReference)
		}
		*p.oneOfType0 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType0.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType0.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfCreateRateLimitPolicyApiResponseData) GetValue() interface{} {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return *p.oneOfType0
	}
	return nil
}

func (p *OneOfCreateRateLimitPolicyApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "vmm.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	vOneOfType0 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType0); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType0.ObjectType_ {
			if nil == p.oneOfType0 {
				p.oneOfType0 = new(import1.TaskReference)
			}
			*p.oneOfType0 = *vOneOfType0
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType0.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType0.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfCreateRateLimitPolicyApiResponseData"))
}

func (p *OneOfCreateRateLimitPolicyApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType0)
	}
	return nil, errors.New("No value to marshal for OneOfCreateRateLimitPolicyApiResponseData")
}

type OneOfListPlacementPoliciesApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
	oneOfType0    []PlacementPolicy      `json:"-"`
}

func NewOneOfListPlacementPoliciesApiResponseData() *OneOfListPlacementPoliciesApiResponseData {
	p := new(OneOfListPlacementPoliciesApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfListPlacementPoliciesApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfListPlacementPoliciesApiResponseData is nil"))
	}
	switch v.(type) {
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	case []PlacementPolicy:
		p.oneOfType0 = v.([]PlacementPolicy)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = "List<vmm.v4.images.config.PlacementPolicy>"
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = "List<vmm.v4.images.config.PlacementPolicy>"
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfListPlacementPoliciesApiResponseData) GetValue() interface{} {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	if "List<vmm.v4.images.config.PlacementPolicy>" == *p.Discriminator {
		return p.oneOfType0
	}
	return nil
}

func (p *OneOfListPlacementPoliciesApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "vmm.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	vOneOfType0 := new([]PlacementPolicy)
	if err := json.Unmarshal(b, vOneOfType0); err == nil {

		if len(*vOneOfType0) == 0 || "vmm.v4.images.config.PlacementPolicy" == *((*vOneOfType0)[0].ObjectType_) {
			p.oneOfType0 = *vOneOfType0
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = "List<vmm.v4.images.config.PlacementPolicy>"
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = "List<vmm.v4.images.config.PlacementPolicy>"
			return nil

		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfListPlacementPoliciesApiResponseData"))
}

func (p *OneOfListPlacementPoliciesApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	if "List<vmm.v4.images.config.PlacementPolicy>" == *p.Discriminator {
		return json.Marshal(p.oneOfType0)
	}
	return nil, errors.New("No value to marshal for OneOfListPlacementPoliciesApiResponseData")
}

type OneOfGetRateLimitPolicyApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType0    *RateLimitPolicy       `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfGetRateLimitPolicyApiResponseData() *OneOfGetRateLimitPolicyApiResponseData {
	p := new(OneOfGetRateLimitPolicyApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfGetRateLimitPolicyApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfGetRateLimitPolicyApiResponseData is nil"))
	}
	switch v.(type) {
	case RateLimitPolicy:
		if nil == p.oneOfType0 {
			p.oneOfType0 = new(RateLimitPolicy)
		}
		*p.oneOfType0 = v.(RateLimitPolicy)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType0.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType0.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfGetRateLimitPolicyApiResponseData) GetValue() interface{} {
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return *p.oneOfType0
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfGetRateLimitPolicyApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType0 := new(RateLimitPolicy)
	if err := json.Unmarshal(b, vOneOfType0); err == nil {
		if "vmm.v4.images.config.RateLimitPolicy" == *vOneOfType0.ObjectType_ {
			if nil == p.oneOfType0 {
				p.oneOfType0 = new(RateLimitPolicy)
			}
			*p.oneOfType0 = *vOneOfType0
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType0.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType0.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "vmm.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfGetRateLimitPolicyApiResponseData"))
}

func (p *OneOfGetRateLimitPolicyApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType0)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfGetRateLimitPolicyApiResponseData")
}

type OneOfSuspendPlacementPolicyApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
	oneOfType0    *import1.TaskReference `json:"-"`
}

func NewOneOfSuspendPlacementPolicyApiResponseData() *OneOfSuspendPlacementPolicyApiResponseData {
	p := new(OneOfSuspendPlacementPolicyApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfSuspendPlacementPolicyApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfSuspendPlacementPolicyApiResponseData is nil"))
	}
	switch v.(type) {
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	case import1.TaskReference:
		if nil == p.oneOfType0 {
			p.oneOfType0 = new(import1.TaskReference)
		}
		*p.oneOfType0 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType0.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType0.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfSuspendPlacementPolicyApiResponseData) GetValue() interface{} {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return *p.oneOfType0
	}
	return nil
}

func (p *OneOfSuspendPlacementPolicyApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "vmm.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	vOneOfType0 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType0); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType0.ObjectType_ {
			if nil == p.oneOfType0 {
				p.oneOfType0 = new(import1.TaskReference)
			}
			*p.oneOfType0 = *vOneOfType0
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType0.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType0.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfSuspendPlacementPolicyApiResponseData"))
}

func (p *OneOfSuspendPlacementPolicyApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType0)
	}
	return nil, errors.New("No value to marshal for OneOfSuspendPlacementPolicyApiResponseData")
}

type OneOfUpdateRateLimitPolicyApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
	oneOfType0    *import1.TaskReference `json:"-"`
}

func NewOneOfUpdateRateLimitPolicyApiResponseData() *OneOfUpdateRateLimitPolicyApiResponseData {
	p := new(OneOfUpdateRateLimitPolicyApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfUpdateRateLimitPolicyApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfUpdateRateLimitPolicyApiResponseData is nil"))
	}
	switch v.(type) {
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	case import1.TaskReference:
		if nil == p.oneOfType0 {
			p.oneOfType0 = new(import1.TaskReference)
		}
		*p.oneOfType0 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType0.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType0.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfUpdateRateLimitPolicyApiResponseData) GetValue() interface{} {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return *p.oneOfType0
	}
	return nil
}

func (p *OneOfUpdateRateLimitPolicyApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "vmm.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	vOneOfType0 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType0); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType0.ObjectType_ {
			if nil == p.oneOfType0 {
				p.oneOfType0 = new(import1.TaskReference)
			}
			*p.oneOfType0 = *vOneOfType0
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType0.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType0.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfUpdateRateLimitPolicyApiResponseData"))
}

func (p *OneOfUpdateRateLimitPolicyApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType0)
	}
	return nil, errors.New("No value to marshal for OneOfUpdateRateLimitPolicyApiResponseData")
}

type OneOfResumePlacementPolicyApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
	oneOfType0    *import1.TaskReference `json:"-"`
}

func NewOneOfResumePlacementPolicyApiResponseData() *OneOfResumePlacementPolicyApiResponseData {
	p := new(OneOfResumePlacementPolicyApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfResumePlacementPolicyApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfResumePlacementPolicyApiResponseData is nil"))
	}
	switch v.(type) {
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	case import1.TaskReference:
		if nil == p.oneOfType0 {
			p.oneOfType0 = new(import1.TaskReference)
		}
		*p.oneOfType0 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType0.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType0.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfResumePlacementPolicyApiResponseData) GetValue() interface{} {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return *p.oneOfType0
	}
	return nil
}

func (p *OneOfResumePlacementPolicyApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "vmm.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	vOneOfType0 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType0); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType0.ObjectType_ {
			if nil == p.oneOfType0 {
				p.oneOfType0 = new(import1.TaskReference)
			}
			*p.oneOfType0 = *vOneOfType0
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType0.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType0.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfResumePlacementPolicyApiResponseData"))
}

func (p *OneOfResumePlacementPolicyApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType0)
	}
	return nil, errors.New("No value to marshal for OneOfResumePlacementPolicyApiResponseData")
}

type OneOfGetPlacementPolicyApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
	oneOfType0    *PlacementPolicy       `json:"-"`
}

func NewOneOfGetPlacementPolicyApiResponseData() *OneOfGetPlacementPolicyApiResponseData {
	p := new(OneOfGetPlacementPolicyApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfGetPlacementPolicyApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfGetPlacementPolicyApiResponseData is nil"))
	}
	switch v.(type) {
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	case PlacementPolicy:
		if nil == p.oneOfType0 {
			p.oneOfType0 = new(PlacementPolicy)
		}
		*p.oneOfType0 = v.(PlacementPolicy)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType0.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType0.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfGetPlacementPolicyApiResponseData) GetValue() interface{} {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return *p.oneOfType0
	}
	return nil
}

func (p *OneOfGetPlacementPolicyApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "vmm.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	vOneOfType0 := new(PlacementPolicy)
	if err := json.Unmarshal(b, vOneOfType0); err == nil {
		if "vmm.v4.images.config.PlacementPolicy" == *vOneOfType0.ObjectType_ {
			if nil == p.oneOfType0 {
				p.oneOfType0 = new(PlacementPolicy)
			}
			*p.oneOfType0 = *vOneOfType0
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType0.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType0.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfGetPlacementPolicyApiResponseData"))
}

func (p *OneOfGetPlacementPolicyApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType0)
	}
	return nil, errors.New("No value to marshal for OneOfGetPlacementPolicyApiResponseData")
}

type OneOfDeletePlacementPolicyApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
	oneOfType0    *import1.TaskReference `json:"-"`
}

func NewOneOfDeletePlacementPolicyApiResponseData() *OneOfDeletePlacementPolicyApiResponseData {
	p := new(OneOfDeletePlacementPolicyApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfDeletePlacementPolicyApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfDeletePlacementPolicyApiResponseData is nil"))
	}
	switch v.(type) {
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	case import1.TaskReference:
		if nil == p.oneOfType0 {
			p.oneOfType0 = new(import1.TaskReference)
		}
		*p.oneOfType0 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType0.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType0.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfDeletePlacementPolicyApiResponseData) GetValue() interface{} {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return *p.oneOfType0
	}
	return nil
}

func (p *OneOfDeletePlacementPolicyApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "vmm.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	vOneOfType0 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType0); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType0.ObjectType_ {
			if nil == p.oneOfType0 {
				p.oneOfType0 = new(import1.TaskReference)
			}
			*p.oneOfType0 = *vOneOfType0
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType0.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType0.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfDeletePlacementPolicyApiResponseData"))
}

func (p *OneOfDeletePlacementPolicyApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType0)
	}
	return nil, errors.New("No value to marshal for OneOfDeletePlacementPolicyApiResponseData")
}

type OneOfUpdatePlacementPolicyApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
	oneOfType0    *import1.TaskReference `json:"-"`
}

func NewOneOfUpdatePlacementPolicyApiResponseData() *OneOfUpdatePlacementPolicyApiResponseData {
	p := new(OneOfUpdatePlacementPolicyApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfUpdatePlacementPolicyApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfUpdatePlacementPolicyApiResponseData is nil"))
	}
	switch v.(type) {
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	case import1.TaskReference:
		if nil == p.oneOfType0 {
			p.oneOfType0 = new(import1.TaskReference)
		}
		*p.oneOfType0 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType0.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType0.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfUpdatePlacementPolicyApiResponseData) GetValue() interface{} {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return *p.oneOfType0
	}
	return nil
}

func (p *OneOfUpdatePlacementPolicyApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "vmm.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	vOneOfType0 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType0); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType0.ObjectType_ {
			if nil == p.oneOfType0 {
				p.oneOfType0 = new(import1.TaskReference)
			}
			*p.oneOfType0 = *vOneOfType0
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType0.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType0.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfUpdatePlacementPolicyApiResponseData"))
}

func (p *OneOfUpdatePlacementPolicyApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType0)
	}
	return nil, errors.New("No value to marshal for OneOfUpdatePlacementPolicyApiResponseData")
}

type OneOfCreatePlacementPolicyApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
	oneOfType0    *import1.TaskReference `json:"-"`
}

func NewOneOfCreatePlacementPolicyApiResponseData() *OneOfCreatePlacementPolicyApiResponseData {
	p := new(OneOfCreatePlacementPolicyApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfCreatePlacementPolicyApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfCreatePlacementPolicyApiResponseData is nil"))
	}
	switch v.(type) {
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	case import1.TaskReference:
		if nil == p.oneOfType0 {
			p.oneOfType0 = new(import1.TaskReference)
		}
		*p.oneOfType0 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType0.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType0.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfCreatePlacementPolicyApiResponseData) GetValue() interface{} {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return *p.oneOfType0
	}
	return nil
}

func (p *OneOfCreatePlacementPolicyApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "vmm.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	vOneOfType0 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType0); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType0.ObjectType_ {
			if nil == p.oneOfType0 {
				p.oneOfType0 = new(import1.TaskReference)
			}
			*p.oneOfType0 = *vOneOfType0
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType0.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType0.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfCreatePlacementPolicyApiResponseData"))
}

func (p *OneOfCreatePlacementPolicyApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType0)
	}
	return nil, errors.New("No value to marshal for OneOfCreatePlacementPolicyApiResponseData")
}

type OneOfDeleteRateLimitPolicyApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
	oneOfType0    *import1.TaskReference `json:"-"`
}

func NewOneOfDeleteRateLimitPolicyApiResponseData() *OneOfDeleteRateLimitPolicyApiResponseData {
	p := new(OneOfDeleteRateLimitPolicyApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfDeleteRateLimitPolicyApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfDeleteRateLimitPolicyApiResponseData is nil"))
	}
	switch v.(type) {
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	case import1.TaskReference:
		if nil == p.oneOfType0 {
			p.oneOfType0 = new(import1.TaskReference)
		}
		*p.oneOfType0 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType0.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType0.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfDeleteRateLimitPolicyApiResponseData) GetValue() interface{} {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return *p.oneOfType0
	}
	return nil
}

func (p *OneOfDeleteRateLimitPolicyApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "vmm.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	vOneOfType0 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType0); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType0.ObjectType_ {
			if nil == p.oneOfType0 {
				p.oneOfType0 = new(import1.TaskReference)
			}
			*p.oneOfType0 = *vOneOfType0
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType0.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType0.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfDeleteRateLimitPolicyApiResponseData"))
}

func (p *OneOfDeleteRateLimitPolicyApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType0)
	}
	return nil, errors.New("No value to marshal for OneOfDeleteRateLimitPolicyApiResponseData")
}

type OneOfListEffectiveRateLimitPoliciesApiResponseData struct {
	Discriminator *string                    `json:"-"`
	ObjectType_   *string                    `json:"-"`
	oneOfType400  *import2.ErrorResponse     `json:"-"`
	oneOfType0    []EffectiveRateLimitPolicy `json:"-"`
}

func NewOneOfListEffectiveRateLimitPoliciesApiResponseData() *OneOfListEffectiveRateLimitPoliciesApiResponseData {
	p := new(OneOfListEffectiveRateLimitPoliciesApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfListEffectiveRateLimitPoliciesApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfListEffectiveRateLimitPoliciesApiResponseData is nil"))
	}
	switch v.(type) {
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	case []EffectiveRateLimitPolicy:
		p.oneOfType0 = v.([]EffectiveRateLimitPolicy)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = "List<vmm.v4.images.config.EffectiveRateLimitPolicy>"
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = "List<vmm.v4.images.config.EffectiveRateLimitPolicy>"
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfListEffectiveRateLimitPoliciesApiResponseData) GetValue() interface{} {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	if "List<vmm.v4.images.config.EffectiveRateLimitPolicy>" == *p.Discriminator {
		return p.oneOfType0
	}
	return nil
}

func (p *OneOfListEffectiveRateLimitPoliciesApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "vmm.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	vOneOfType0 := new([]EffectiveRateLimitPolicy)
	if err := json.Unmarshal(b, vOneOfType0); err == nil {

		if len(*vOneOfType0) == 0 || "vmm.v4.images.config.EffectiveRateLimitPolicy" == *((*vOneOfType0)[0].ObjectType_) {
			p.oneOfType0 = *vOneOfType0
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = "List<vmm.v4.images.config.EffectiveRateLimitPolicy>"
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = "List<vmm.v4.images.config.EffectiveRateLimitPolicy>"
			return nil

		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfListEffectiveRateLimitPoliciesApiResponseData"))
}

func (p *OneOfListEffectiveRateLimitPoliciesApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	if "List<vmm.v4.images.config.EffectiveRateLimitPolicy>" == *p.Discriminator {
		return json.Marshal(p.oneOfType0)
	}
	return nil, errors.New("No value to marshal for OneOfListEffectiveRateLimitPoliciesApiResponseData")
}

type OneOfListRateLimitPoliciesApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType0    []RateLimitPolicy      `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfListRateLimitPoliciesApiResponseData() *OneOfListRateLimitPoliciesApiResponseData {
	p := new(OneOfListRateLimitPoliciesApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfListRateLimitPoliciesApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfListRateLimitPoliciesApiResponseData is nil"))
	}
	switch v.(type) {
	case []RateLimitPolicy:
		p.oneOfType0 = v.([]RateLimitPolicy)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = "List<vmm.v4.images.config.RateLimitPolicy>"
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = "List<vmm.v4.images.config.RateLimitPolicy>"
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfListRateLimitPoliciesApiResponseData) GetValue() interface{} {
	if "List<vmm.v4.images.config.RateLimitPolicy>" == *p.Discriminator {
		return p.oneOfType0
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfListRateLimitPoliciesApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType0 := new([]RateLimitPolicy)
	if err := json.Unmarshal(b, vOneOfType0); err == nil {

		if len(*vOneOfType0) == 0 || "vmm.v4.images.config.RateLimitPolicy" == *((*vOneOfType0)[0].ObjectType_) {
			p.oneOfType0 = *vOneOfType0
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = "List<vmm.v4.images.config.RateLimitPolicy>"
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = "List<vmm.v4.images.config.RateLimitPolicy>"
			return nil

		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "vmm.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfListRateLimitPoliciesApiResponseData"))
}

func (p *OneOfListRateLimitPoliciesApiResponseData) MarshalJSON() ([]byte, error) {
	if "List<vmm.v4.images.config.RateLimitPolicy>" == *p.Discriminator {
		return json.Marshal(p.oneOfType0)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfListRateLimitPoliciesApiResponseData")
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
