/*
 * Generated file models/common/v1/response/response_model.go.
 *
 * Product version: 4.0.2-alpha-3
 *
 * Part of the Nutanix Storage Versioned APIs
 *
 * (c) 2023 Nutanix Inc.  All rights reserved
 *
 */

/*
  Nutanix Standard Response Format
*/
package response

import (
	import1 "github.com/nutanix/ntnx-api-golang-clients/storage-go-client/v4/models/common/v1/config"
)

/**
A HATEOAS style link for the response.  Each link contains a user friendly name identifying the link and an address for retrieving the particular resource.
*/
type ApiLink struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/**
	  The URL at which the entity described by this link can be accessed.
	*/
	Href *string `json:"href,omitempty"`
	/**
	  A name that identifies the relationship of this link to the object that is returned by the URL.  The special value of "self" identifies the URL for the object.
	*/
	Rel *string `json:"rel,omitempty"`
}

func NewApiLink() *ApiLink {
	p := new(ApiLink)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "common.v1.response.ApiLink"
	p.Reserved_ = map[string]interface{}{"$fqObjectType": "common.v1.r0.a3.response.ApiLink"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/**
The metadata associated with an API response. Always present and minimally contains the self-link for the API request that produced this response.  Also contains pagination data for paginated requests.
*/
type ApiResponseMetadata struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/**
	  An array of entity-specific metadata
	*/
	ExtraInfo []import1.KVPair `json:"extraInfo,omitempty"`

	Flags []import1.Flag `json:"flags,omitempty"`

	Links []ApiLink `json:"links,omitempty"`

	Messages []import1.Message `json:"messages,omitempty"`

	TotalAvailableResults *int `json:"totalAvailableResults,omitempty"`
}

func NewApiResponseMetadata() *ApiResponseMetadata {
	p := new(ApiResponseMetadata)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "common.v1.response.ApiResponseMetadata"
	p.Reserved_ = map[string]interface{}{"$fqObjectType": "common.v1.r0.a3.response.ApiResponseMetadata"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/**
A model that represents an object instance that is accessible through an API endpoint.  Instances of this type get an extId field that contains the globally unique identifier for that instance.  Externally accessible instances are always tenant aware and therefore extend the TenantAwareModel
*/
type ExternalizableAbstractModel struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/**
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/**
	  A HATEOAS style link for the response.  Each link contains a user friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []ApiLink `json:"links,omitempty"`
	/**
	  A globally unique identifier that represents the tenant that owns this entity.  It is automatically assigned by the system and is immutable from an API consumer perspective (some use cases may cause this Id to change - for instance a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`
}

func NewExternalizableAbstractModel() *ExternalizableAbstractModel {
	p := new(ExternalizableAbstractModel)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "common.v1.response.ExternalizableAbstractModel"
	p.Reserved_ = map[string]interface{}{"$fqObjectType": "common.v1.r0.a3.response.ExternalizableAbstractModel"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}
