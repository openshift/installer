/*
 * Generated file models/prism/v4/config/config_model.go.
 *
 * Product version: 4.0.1-beta-1
 *
 * Part of the Nutanix Prism Versioned APIs
 *
 * (c) 2024 Nutanix Inc.  All rights reserved
 *
 */

/*
  Configure Tasks and Monitoring
*/
package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	import3 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/common/v1/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/common/v1/response"
	import1 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/error"
	"time"
)

/*
This attribute contains the list of entities and policies which have been assigned the given category.<br>
These entities are grouped by entity types (like VM or HOST) or policy types (like PROTECTION_POLICY or NGT_POLICY).<br>
Each associated object contains the total entities belonging to the given entity type, count, category extId, and
references (for example for VM it'd be VM uuid).
*/
type AssociationDetail struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  External identifier for the given category, used across all v4 apis/entities/resources where categories are referenced.<br>
	The field has UUID format.<br>
	A type 4 UUID is generated during category creation.
	*/
	CategoryId *string `json:"categoryId,omitempty"`

	ResourceGroup *ResourceGroup `json:"resourceGroup,omitempty"`
	/*
	  The UUID of the entity or policy associated with the particular category.
	*/
	ResourceId *string `json:"resourceId,omitempty"`

	ResourceType *ResourceType `json:"resourceType,omitempty"`
}

func NewAssociationDetail() *AssociationDetail {
	p := new(AssociationDetail)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.config.AssociationDetail"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
This attribute contains the list of entities which have been assigned the given category.<br>
These entities are grouped by entity types (like VM or HOST) or policy types (like PROTECTION_POLICY or
NGT_POLICY).<br>
Each associated object contains the total entities belonging to the given entity type, category extId, and
references (for example for VM it'd be VM UUID).
*/
type AssociationDetailOld struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  External identifier for the given category.
	*/
	CategoryId *string `json:"categoryId,omitempty"`
	/*
	  Denotes the type of a category.<br>
	There are three types of categories: SYSTEM, INTERNAL, and USER.<br>
	This field is immutable.
	*/
	Count *int `json:"count,omitempty"`

	ResourceGroup *ResourceGroup `json:"resourceGroup,omitempty"`

	ResourceReferences []Reference `json:"resourceReferences,omitempty"`

	ResourceType *ResourceType `json:"resourceType,omitempty"`
}

func NewAssociationDetailOld() *AssociationDetailOld {
	p := new(AssociationDetailOld)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.config.AssociationDetailOld"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

type AssociationDetailOldProjection struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  External identifier for the given category.
	*/
	CategoryId *string `json:"categoryId,omitempty"`
	/*
	  Denotes the type of a category.<br>
	There are three types of categories: SYSTEM, INTERNAL, and USER.<br>
	This field is immutable.
	*/
	Count *int `json:"count,omitempty"`

	ResourceGroup *ResourceGroup `json:"resourceGroup,omitempty"`

	ResourceReferences []Reference `json:"resourceReferences,omitempty"`

	ResourceType *ResourceType `json:"resourceType,omitempty"`
}

func NewAssociationDetailOldProjection() *AssociationDetailOldProjection {
	p := new(AssociationDetailOldProjection)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.config.AssociationDetailOldProjection"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

type AssociationDetailProjection struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  External identifier for the given category, used across all v4 apis/entities/resources where categories are referenced.<br>
	The field has UUID format.<br>
	A type 4 UUID is generated during category creation.
	*/
	CategoryId *string `json:"categoryId,omitempty"`

	ResourceGroup *ResourceGroup `json:"resourceGroup,omitempty"`
	/*
	  The UUID of the entity or policy associated with the particular category.
	*/
	ResourceId *string `json:"resourceId,omitempty"`

	ResourceType *ResourceType `json:"resourceType,omitempty"`
}

func NewAssociationDetailProjection() *AssociationDetailProjection {
	p := new(AssociationDetailProjection)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.config.AssociationDetailProjection"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
This attribute contains the list of entities and policies which have been assigned the given category.<br>
These entities are grouped by entity types (like VM or HOST) or policy types (like PROTECTION_POLICY or NGT_POLICY).<br>
Each associated object contains the total entities belonging to the given entity type, count, category extId, and
references (for example for VM it'd be VM uuid).
*/
type AssociationSummary struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  External identifier for the given category, used across all v4 apis/entities/resources where categories are referenced.<br>
	The field has UUID format.<br>
	A type 4 UUID is generated during category creation.
	*/
	CategoryId *string `json:"categoryId,omitempty"`
	/*
	  Count of associations of a particular type of entity or policy
	*/
	Count *int `json:"count,omitempty"`

	ResourceGroup *ResourceGroup `json:"resourceGroup,omitempty"`

	ResourceType *ResourceType `json:"resourceType,omitempty"`
}

func NewAssociationSummary() *AssociationSummary {
	p := new(AssociationSummary)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.config.AssociationSummary"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

type AssociationSummaryProjection struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  External identifier for the given category, used across all v4 apis/entities/resources where categories are referenced.<br>
	The field has UUID format.<br>
	A type 4 UUID is generated during category creation.
	*/
	CategoryId *string `json:"categoryId,omitempty"`
	/*
	  Count of associations of a particular type of entity or policy
	*/
	Count *int `json:"count,omitempty"`

	ResourceGroup *ResourceGroup `json:"resourceGroup,omitempty"`

	ResourceType *ResourceType `json:"resourceType,omitempty"`
}

func NewAssociationSummaryProjection() *AssociationSummaryProjection {
	p := new(AssociationSummaryProjection)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.config.AssociationSummaryProjection"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
REST response for all response codes in API path /prism/v4.0.b1/config/tasks/{taskExtId}/$actions/cancel Post operation
*/
type CancelTaskApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfCancelTaskApiResponseData `json:"data,omitempty"`

	Metadata *import2.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewCancelTaskApiResponse() *CancelTaskApiResponse {
	p := new(CancelTaskApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.config.CancelTaskApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *CancelTaskApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *CancelTaskApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfCancelTaskApiResponseData()
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

type Category struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  This field gives basic information about resources that are associated to the category.<br>
	The results present under this field summarize the counts of various kinds of resources associated to the category.<br>
	For more detailed information about the UUIDs of the resources, please look into the field `detailedAssociations`.<br>
	This field will be ignored, if given in the payload of `updateCategoryById` or `createCategory` APIs.<br>
	This field will not be present by default in `listCategories` API, unless the parameter $expand=associations is present in the url.
	*/
	Associations []AssociationSummary `json:"associations,omitempty"`
	/*
	  A string consisting of the description of the category as defined by the user.<br>
	Description can be optionally provided in the payload of `createCategory` and `updateCategoryById` APIs.<br>
	Description field can be updated through `updateCategoryById` API.<br>
	The server does not validate this value nor does it enforce the uniqueness or any other constraints.<br>
	It is the responsibility of the user to ensure that any semantic or syntactic constraints are retained when mutating
	this field.
	*/
	Description *string `json:"description,omitempty"`
	/*
	  This field gives detailed information about resources that are associated to the category.<br>
	The results present under this field contain the UUIDs of the entities and policies of various kinds associated to the category.<br>
	This field will be ignored, if given in the payload of `updateCategoryById` or `createCategory` APIs.<br>
	This field will not be present by default in `listCategories` or `getCategoryById` APIs, unless the parameter $expand=detailedAssociations is present in the url.
	*/
	DetailedAssociations []AssociationDetail `json:"detailedAssociations,omitempty"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  The key of a category when it is represented in `key:value` format.

	Constraints applicable when field is given in the payload during create and update:
	* A string of maxlength of 64
	* Character at the start cannot be `$`
	* Character `/` is not allowed anywhere

	It is a mandatory field in the payload of `createCategory` and `updateCategoryById` APIs.<br>
	This field can't be updated through `updateCategoryById` API.
	*/
	Key *string `json:"key"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import2.ApiLink `json:"links,omitempty"`
	/*
	  This field contains the UUID of a user who owns the category.<br>
	This field will be ignored, if given in the payload of `createCategory` API. Hence, when a category is created, the logged-in user automatically becomes the owner of the category.<br>
	This field can be updated through `updateCategoryById` API, in which case, should be provided, UUID of a valid user present in the system.<br>
	Validity of the user UUID can be checked by invoking the api: authn/users/{extId} in the 'Identity and Access Management' or 'IAM' namespace.<br>
	It is used for enabling RBAC access to self-owned categories.
	*/
	OwnerUuid *string `json:"ownerUuid,omitempty"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`

	Type *CategoryType `json:"type,omitempty"`
	/*
	  The value of a category when it is represented in `key:value` format.

	Constraints applicable when field is given in the payload during create and update:
	* A string of maxlength 64
	* Character at the start cannot be `$`
	* Character `/` is not allowed anywhere

	It is a mandatory input field in the payload of `createCategory` and `updateCategoryById` APIs.<br>
	This field can be updated through `updateCategoryById` API.<br>
	Updating the value will not change the extId of the category.
	*/
	Value *string `json:"value"`
}

func (p *Category) MarshalJSON() ([]byte, error) {
	type CategoryProxy Category
	return json.Marshal(struct {
		*CategoryProxy
		Key   *string `json:"key,omitempty"`
		Value *string `json:"value,omitempty"`
	}{
		CategoryProxy: (*CategoryProxy)(p),
		Key:           p.Key,
		Value:         p.Value,
	})
}

func NewCategory() *Category {
	p := new(Category)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.config.Category"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
This attribute contains the list of entities which have been assigned the given category.<br>
These entities are grouped by entity types (like VM or HOST).<br>
Each associated object contains the total entities belonging to the given entity type, and category extId.
*/
type CategoryAssociationSummaryOld struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  External identifier for the given category.
	*/
	CategoryId *string `json:"categoryId,omitempty"`
	/*
	  Denotes the type of a category.<br>
	There are three types of categories: SYSTEM, INTERNAL, and USER.<br>
	This field is immutable.
	*/
	Count *int `json:"count,omitempty"`

	ResourceGroup *ResourceGroup `json:"resourceGroup,omitempty"`

	ResourceType *ResourceType `json:"resourceType,omitempty"`
}

func NewCategoryAssociationSummaryOld() *CategoryAssociationSummaryOld {
	p := new(CategoryAssociationSummaryOld)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.config.CategoryAssociationSummaryOld"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

type CategoryAssociationSummaryOldProjection struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  External identifier for the given category.
	*/
	CategoryId *string `json:"categoryId,omitempty"`
	/*
	  Denotes the type of a category.<br>
	There are three types of categories: SYSTEM, INTERNAL, and USER.<br>
	This field is immutable.
	*/
	Count *int `json:"count,omitempty"`

	ResourceGroup *ResourceGroup `json:"resourceGroup,omitempty"`

	ResourceType *ResourceType `json:"resourceType,omitempty"`
}

func NewCategoryAssociationSummaryOldProjection() *CategoryAssociationSummaryOldProjection {
	p := new(CategoryAssociationSummaryOldProjection)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.config.CategoryAssociationSummaryOldProjection"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Denotes the type of a category.<br>
There are three types of categories: SYSTEM, INTERNAL, and USER.<br>
This field is immutable.
*/
type CategoryOld struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	Associations []CategoryAssociationSummaryOld `json:"associations,omitempty"`

	ChildCategories []CategorySummaryOld `json:"childCategories,omitempty"`
	/*
	  A string consisting of the description of the category as defined by the user.

	The server does not validate this value nor does it enforce the uniqueness or any other constraints.<br>
	It is the responsibility of the user to ensure that any semantic or syntactic constraints are retained when mutating
	this field.
	*/
	Description *string `json:"description,omitempty"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  The fully qualified name of this category. It is unique for each category.<br>
	It is a read-only field.
	The service constructs it from the name-parentExtId combination.
	An example of a fqName would be `Location/Bangalore`, where `Location` is the
	parent category's name and `Bangalore` is the category name.<br>
	This field is immutable.<br>
	*/
	FqName *string `json:"fqName,omitempty"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import2.ApiLink `json:"links,omitempty"`
	/*
	  Opaque metadata which can be associated to a category.<br>
	It is a list of key-value pairs.<br>
	For example, for a category 'California/SanJose' we can associate a geographical coordinate based metadata
	like: {'latitude': '37.3382° N' , 'longitude': '121.8863° W'}.

	The server does not validate this value nor does it enforce the uniqueness or any other constraints.
	It is the responsibility of the user to ensure that any semantic or syntactic constraints are retained when mutating
	this field.
	*/
	Metadata []import3.KVPair `json:"metadata,omitempty"`
	/*
	  The short name of this category. It may not be unique for each category.<br>
	It is a mandatory field which must be specified inside the post/put request body.<br>
	This field is immutable.
	*/
	Name *string `json:"name"`
	/*
	  It is a read-only field inserted into category entity at the time of category creation, and which contains the UUID of
	the user who created this category. It is used for enabling authorization of a particular kind where the user has no
	access to view/create/update/delete any categories other than the category created by oneself.
	*/
	OwnerUuid *string `json:"ownerUuid,omitempty"`
	/*
	  The parent category of this category (may be null if this category is not part of a hierarchy).<br>
	Each category can have at most one parent.<br>
	A parent cannot be deleted until all the children categories are deleted first.<br>
	Must be specified inside the post/put request body for child categories (if not specified, the service assumes
	the category to be a parent category).<br>
	This field is immutable.
	*/
	ParentExtId *string `json:"parentExtId,omitempty"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`

	Type *CategoryType `json:"type,omitempty"`
	/*
	  The user specified name is a string that the user can specify; with syntax and semantics controlled by the user.

	The server does not validate this value nor does it enforce the uniqueness or any other constraints.<br>
	It is the responsibility of the user to ensure that any semantic or syntactic constraints are retained when mutating
	this field. Unlike the name of the categories, which is immutable, the user name can be changed by the user to meet their needs.
	*/
	UserSpecifiedName *string `json:"userSpecifiedName,omitempty"`
}

func (p *CategoryOld) MarshalJSON() ([]byte, error) {
	type CategoryOldProxy CategoryOld
	return json.Marshal(struct {
		*CategoryOldProxy
		Name *string `json:"name,omitempty"`
	}{
		CategoryOldProxy: (*CategoryOldProxy)(p),
		Name:             p.Name,
	})
}

func NewCategoryOld() *CategoryOld {
	p := new(CategoryOld)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.config.CategoryOld"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

type CategoryOldProjection struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	Associations []CategoryAssociationSummaryOld `json:"associations,omitempty"`

	CategoryAssociationSummaryOldProjection *CategoryAssociationSummaryOldProjection `json:"categoryAssociationSummaryOldProjection,omitempty"`

	CategorySummaryOldProjection *CategorySummaryOldProjection `json:"categorySummaryOldProjection,omitempty"`

	ChildCategories []CategorySummaryOld `json:"childCategories,omitempty"`
	/*
	  A string consisting of the description of the category as defined by the user.

	The server does not validate this value nor does it enforce the uniqueness or any other constraints.<br>
	It is the responsibility of the user to ensure that any semantic or syntactic constraints are retained when mutating
	this field.
	*/
	Description *string `json:"description,omitempty"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  The fully qualified name of this category. It is unique for each category.<br>
	It is a read-only field.
	The service constructs it from the name-parentExtId combination.
	An example of a fqName would be `Location/Bangalore`, where `Location` is the
	parent category's name and `Bangalore` is the category name.<br>
	This field is immutable.<br>
	*/
	FqName *string `json:"fqName,omitempty"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import2.ApiLink `json:"links,omitempty"`
	/*
	  Opaque metadata which can be associated to a category.<br>
	It is a list of key-value pairs.<br>
	For example, for a category 'California/SanJose' we can associate a geographical coordinate based metadata
	like: {'latitude': '37.3382° N' , 'longitude': '121.8863° W'}.

	The server does not validate this value nor does it enforce the uniqueness or any other constraints.
	It is the responsibility of the user to ensure that any semantic or syntactic constraints are retained when mutating
	this field.
	*/
	Metadata []import3.KVPair `json:"metadata,omitempty"`
	/*
	  The short name of this category. It may not be unique for each category.<br>
	It is a mandatory field which must be specified inside the post/put request body.<br>
	This field is immutable.
	*/
	Name *string `json:"name"`
	/*
	  It is a read-only field inserted into category entity at the time of category creation, and which contains the UUID of
	the user who created this category. It is used for enabling authorization of a particular kind where the user has no
	access to view/create/update/delete any categories other than the category created by oneself.
	*/
	OwnerUuid *string `json:"ownerUuid,omitempty"`
	/*
	  The parent category of this category (may be null if this category is not part of a hierarchy).<br>
	Each category can have at most one parent.<br>
	A parent cannot be deleted until all the children categories are deleted first.<br>
	Must be specified inside the post/put request body for child categories (if not specified, the service assumes
	the category to be a parent category).<br>
	This field is immutable.
	*/
	ParentExtId *string `json:"parentExtId,omitempty"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`

	Type *CategoryType `json:"type,omitempty"`
	/*
	  The user specified name is a string that the user can specify; with syntax and semantics controlled by the user.

	The server does not validate this value nor does it enforce the uniqueness or any other constraints.<br>
	It is the responsibility of the user to ensure that any semantic or syntactic constraints are retained when mutating
	this field. Unlike the name of the categories, which is immutable, the user name can be changed by the user to meet their needs.
	*/
	UserSpecifiedName *string `json:"userSpecifiedName,omitempty"`
}

func NewCategoryOldProjection() *CategoryOldProjection {
	p := new(CategoryOldProjection)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.config.CategoryOldProjection"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

type CategoryProjection struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	AssociationDetailProjection *AssociationDetailProjection `json:"associationDetailProjection,omitempty"`

	AssociationSummaryProjection *AssociationSummaryProjection `json:"associationSummaryProjection,omitempty"`
	/*
	  This field gives basic information about resources that are associated to the category.<br>
	The results present under this field summarize the counts of various kinds of resources associated to the category.<br>
	For more detailed information about the UUIDs of the resources, please look into the field `detailedAssociations`.<br>
	This field will be ignored, if given in the payload of `updateCategoryById` or `createCategory` APIs.<br>
	This field will not be present by default in `listCategories` API, unless the parameter $expand=associations is present in the url.
	*/
	Associations []AssociationSummary `json:"associations,omitempty"`
	/*
	  A string consisting of the description of the category as defined by the user.<br>
	Description can be optionally provided in the payload of `createCategory` and `updateCategoryById` APIs.<br>
	Description field can be updated through `updateCategoryById` API.<br>
	The server does not validate this value nor does it enforce the uniqueness or any other constraints.<br>
	It is the responsibility of the user to ensure that any semantic or syntactic constraints are retained when mutating
	this field.
	*/
	Description *string `json:"description,omitempty"`
	/*
	  This field gives detailed information about resources that are associated to the category.<br>
	The results present under this field contain the UUIDs of the entities and policies of various kinds associated to the category.<br>
	This field will be ignored, if given in the payload of `updateCategoryById` or `createCategory` APIs.<br>
	This field will not be present by default in `listCategories` or `getCategoryById` APIs, unless the parameter $expand=detailedAssociations is present in the url.
	*/
	DetailedAssociations []AssociationDetail `json:"detailedAssociations,omitempty"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  The key of a category when it is represented in `key:value` format.

	Constraints applicable when field is given in the payload during create and update:
	* A string of maxlength of 64
	* Character at the start cannot be `$`
	* Character `/` is not allowed anywhere

	It is a mandatory field in the payload of `createCategory` and `updateCategoryById` APIs.<br>
	This field can't be updated through `updateCategoryById` API.
	*/
	Key *string `json:"key"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import2.ApiLink `json:"links,omitempty"`
	/*
	  This field contains the UUID of a user who owns the category.<br>
	This field will be ignored, if given in the payload of `createCategory` API. Hence, when a category is created, the logged-in user automatically becomes the owner of the category.<br>
	This field can be updated through `updateCategoryById` API, in which case, should be provided, UUID of a valid user present in the system.<br>
	Validity of the user UUID can be checked by invoking the api: authn/users/{extId} in the 'Identity and Access Management' or 'IAM' namespace.<br>
	It is used for enabling RBAC access to self-owned categories.
	*/
	OwnerUuid *string `json:"ownerUuid,omitempty"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`

	Type *CategoryType `json:"type,omitempty"`
	/*
	  The value of a category when it is represented in `key:value` format.

	Constraints applicable when field is given in the payload during create and update:
	* A string of maxlength 64
	* Character at the start cannot be `$`
	* Character `/` is not allowed anywhere

	It is a mandatory input field in the payload of `createCategory` and `updateCategoryById` APIs.<br>
	This field can be updated through `updateCategoryById` API.<br>
	Updating the value will not change the extId of the category.
	*/
	Value *string `json:"value"`
}

func (p *CategoryProjection) MarshalJSON() ([]byte, error) {
	type CategoryProjectionProxy CategoryProjection
	return json.Marshal(struct {
		*CategoryProjectionProxy
		Key   *string `json:"key,omitempty"`
		Value *string `json:"value,omitempty"`
	}{
		CategoryProjectionProxy: (*CategoryProjectionProxy)(p),
		Key:                     p.Key,
		Value:                   p.Value,
	})
}

func NewCategoryProjection() *CategoryProjection {
	p := new(CategoryProjection)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.config.CategoryProjection"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

type CategorySummaryOld struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	Associations []CategoryAssociationSummaryOld `json:"associations,omitempty"`

	ChildCategories []CategorySummaryOld `json:"childCategories,omitempty"`
	/*
	  A string consisting of the description of the category as defined by the user.

	The server does not validate this value nor does it enforce the uniqueness or any other constraints.<br>
	It is the responsibility of the user to ensure that any semantic or syntactic constraints are retained when mutating
	this field.
	*/
	Description *string `json:"description,omitempty"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  The fully qualified name of this category. It is unique for each category.<br>
	It is a read-only field.
	The service constructs it from the name-parentExtId combination.
	An example of a fqName would be `Location/Bangalore`, where `Location` is the
	parent category's name and `Bangalore` is the category name.<br>
	This field is immutable.<br>
	*/
	FqName *string `json:"fqName,omitempty"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import2.ApiLink `json:"links,omitempty"`
	/*
	  The short name of this category. It may not be unique for each category.<br>
	It is a mandatory field which must be specified inside the post/put request body.<br>
	This field is immutable.
	*/
	Name *string `json:"name"`
	/*
	  It is a read-only field inserted into category entity at the time of category creation, and which contains the UUID of
	the user who created this category. It is used for enabling authorization of a particular kind where the user has no
	access to view/create/update/delete any categories other than the category created by oneself.
	*/
	OwnerUuid *string `json:"ownerUuid,omitempty"`
	/*
	  The parent category of this category (may be null if this category is not part of a hierarchy).<br>
	Each category can have at most one parent.<br>
	A parent cannot be deleted until all the children categories are deleted first.<br>
	Must be specified inside the post/put request body for child categories (if not specified, the service assumes
	the category to be a parent category).<br>
	This field is immutable.
	*/
	ParentExtId *string `json:"parentExtId,omitempty"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`

	Type *CategoryType `json:"type,omitempty"`
	/*
	  The user specified name is a string that the user can specify; with syntax and semantics controlled by the user.

	The server does not validate this value nor does it enforce the uniqueness or any other constraints.<br>
	It is the responsibility of the user to ensure that any semantic or syntactic constraints are retained when mutating
	this field. Unlike the name of the categories, which is immutable, the user name can be changed by the user to meet their needs.
	*/
	UserSpecifiedName *string `json:"userSpecifiedName,omitempty"`
}

func (p *CategorySummaryOld) MarshalJSON() ([]byte, error) {
	type CategorySummaryOldProxy CategorySummaryOld
	return json.Marshal(struct {
		*CategorySummaryOldProxy
		Name *string `json:"name,omitempty"`
	}{
		CategorySummaryOldProxy: (*CategorySummaryOldProxy)(p),
		Name:                    p.Name,
	})
}

func NewCategorySummaryOld() *CategorySummaryOld {
	p := new(CategorySummaryOld)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.config.CategorySummaryOld"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

type CategorySummaryOldProjection struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	Associations []CategoryAssociationSummaryOld `json:"associations,omitempty"`

	CategoryAssociationSummaryOldProjection *CategoryAssociationSummaryOldProjection `json:"categoryAssociationSummaryOldProjection,omitempty"`

	CategorySummaryOldProjection *CategorySummaryOldProjection `json:"categorySummaryOldProjection,omitempty"`

	ChildCategories []CategorySummaryOld `json:"childCategories,omitempty"`
	/*
	  A string consisting of the description of the category as defined by the user.

	The server does not validate this value nor does it enforce the uniqueness or any other constraints.<br>
	It is the responsibility of the user to ensure that any semantic or syntactic constraints are retained when mutating
	this field.
	*/
	Description *string `json:"description,omitempty"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  The fully qualified name of this category. It is unique for each category.<br>
	It is a read-only field.
	The service constructs it from the name-parentExtId combination.
	An example of a fqName would be `Location/Bangalore`, where `Location` is the
	parent category's name and `Bangalore` is the category name.<br>
	This field is immutable.<br>
	*/
	FqName *string `json:"fqName,omitempty"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import2.ApiLink `json:"links,omitempty"`
	/*
	  The short name of this category. It may not be unique for each category.<br>
	It is a mandatory field which must be specified inside the post/put request body.<br>
	This field is immutable.
	*/
	Name *string `json:"name"`
	/*
	  It is a read-only field inserted into category entity at the time of category creation, and which contains the UUID of
	the user who created this category. It is used for enabling authorization of a particular kind where the user has no
	access to view/create/update/delete any categories other than the category created by oneself.
	*/
	OwnerUuid *string `json:"ownerUuid,omitempty"`
	/*
	  The parent category of this category (may be null if this category is not part of a hierarchy).<br>
	Each category can have at most one parent.<br>
	A parent cannot be deleted until all the children categories are deleted first.<br>
	Must be specified inside the post/put request body for child categories (if not specified, the service assumes
	the category to be a parent category).<br>
	This field is immutable.
	*/
	ParentExtId *string `json:"parentExtId,omitempty"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`

	Type *CategoryType `json:"type,omitempty"`
	/*
	  The user specified name is a string that the user can specify; with syntax and semantics controlled by the user.

	The server does not validate this value nor does it enforce the uniqueness or any other constraints.<br>
	It is the responsibility of the user to ensure that any semantic or syntactic constraints are retained when mutating
	this field. Unlike the name of the categories, which is immutable, the user name can be changed by the user to meet their needs.
	*/
	UserSpecifiedName *string `json:"userSpecifiedName,omitempty"`
}

func (p *CategorySummaryOldProjection) MarshalJSON() ([]byte, error) {
	type CategorySummaryOldProjectionProxy CategorySummaryOldProjection
	return json.Marshal(struct {
		*CategorySummaryOldProjectionProxy
		Name *string `json:"name,omitempty"`
	}{
		CategorySummaryOldProjectionProxy: (*CategorySummaryOldProjectionProxy)(p),
		Name:                              p.Name,
	})
}

func NewCategorySummaryOldProjection() *CategorySummaryOldProjection {
	p := new(CategorySummaryOldProjection)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.config.CategorySummaryOldProjection"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Denotes the type of a category.<br>
There are three types of categories: SYSTEM, INTERNAL, and USER.<br>
This field is immutable.
*/
type CategoryType int

const (
	CATEGORYTYPE_UNKNOWN  CategoryType = 0
	CATEGORYTYPE_REDACTED CategoryType = 1
	CATEGORYTYPE_USER     CategoryType = 2
	CATEGORYTYPE_SYSTEM   CategoryType = 3
	CATEGORYTYPE_INTERNAL CategoryType = 4
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *CategoryType) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"USER",
		"SYSTEM",
		"INTERNAL",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e CategoryType) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"USER",
		"SYSTEM",
		"INTERNAL",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *CategoryType) index(name string) CategoryType {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"USER",
		"SYSTEM",
		"INTERNAL",
	}
	for idx := range names {
		if names[idx] == name {
			return CategoryType(idx)
		}
	}
	return CATEGORYTYPE_UNKNOWN
}

func (e *CategoryType) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for CategoryType:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *CategoryType) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e CategoryType) Ref() *CategoryType {
	return &e
}

/*
REST response for all response codes in API path /prism/v4.0.b1/config/categories Post operation
*/
type CreateCategoryApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfCreateCategoryApiResponseData `json:"data,omitempty"`

	Metadata *import2.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewCreateCategoryApiResponse() *CreateCategoryApiResponse {
	p := new(CreateCategoryApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.config.CreateCategoryApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *CreateCategoryApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *CreateCategoryApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfCreateCategoryApiResponseData()
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
REST response for all response codes in API path /prism/v4.0.b1/config/categories/{extId} Delete operation
*/
type DeleteCategoryApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfDeleteCategoryApiResponseData `json:"data,omitempty"`

	Metadata *import2.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewDeleteCategoryApiResponse() *DeleteCategoryApiResponse {
	p := new(DeleteCategoryApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.config.DeleteCategoryApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *DeleteCategoryApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *DeleteCategoryApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfDeleteCategoryApiResponseData()
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
Details of the entity.
*/
type EntityReference struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  A globally unique identifier of the entity.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  Entity type identified as 'namespace:module[:submodule]:entityType'. For example - vmm:ahv:vm, where vmm is the namepsace, ahv is the module and vm is the entitytype.
	*/
	Rel *string `json:"rel,omitempty"`
}

func NewEntityReference() *EntityReference {
	p := new(EntityReference)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.config.EntityReference"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
REST response for all response codes in API path /prism/v4.0.b1/config/categories/{extId} Get operation
*/
type GetCategoryApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfGetCategoryApiResponseData `json:"data,omitempty"`

	Metadata *import2.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewGetCategoryApiResponse() *GetCategoryApiResponse {
	p := new(GetCategoryApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.config.GetCategoryApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *GetCategoryApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *GetCategoryApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfGetCategoryApiResponseData()
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
REST response for all response codes in API path /prism/v4.0.b1/config/tasks/{extId} Get operation
*/
type GetTaskApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfGetTaskApiResponseData `json:"data,omitempty"`

	Metadata *import2.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewGetTaskApiResponse() *GetTaskApiResponse {
	p := new(GetTaskApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.config.GetTaskApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *GetTaskApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *GetTaskApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfGetTaskApiResponseData()
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
REST response for all response codes in API path /prism/v4.0.b1/config/categories Get operation
*/
type ListCategoriesApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfListCategoriesApiResponseData `json:"data,omitempty"`

	Metadata *import2.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewListCategoriesApiResponse() *ListCategoriesApiResponse {
	p := new(ListCategoriesApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.config.ListCategoriesApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *ListCategoriesApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *ListCategoriesApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfListCategoriesApiResponseData()
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
REST response for all response codes in API path /prism/v4.0.b1/config/tasks Get operation
*/
type ListTasksApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfListTasksApiResponseData `json:"data,omitempty"`

	Metadata *import2.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewListTasksApiResponse() *ListTasksApiResponse {
	p := new(ListTasksApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.config.ListTasksApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *ListTasksApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *ListTasksApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfListTasksApiResponseData()
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
Reference to the owner of the task.
*/
type OwnerReference struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  A globally unique identifier of the task owner.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  Username of the task owner.
	*/
	Name *string `json:"name,omitempty"`
}

func NewOwnerReference() *OwnerReference {
	p := new(OwnerReference)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.config.OwnerReference"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Contains references for entities given in EntityAssociation.
This contains the entity ID and a list of links to fetch the associated entities.
*/
type Reference struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	Links []import2.ApiLink `json:"links,omitempty"`
	/*
	  The external identifier of the resource which uniquely identifies it.
	*/
	ResourceId *string `json:"resourceId,omitempty"`
}

func NewReference() *Reference {
	p := new(Reference)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.config.Reference"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
An enum denoting the resource group.<br>
Resources can be organised into either an entity or a policy, hence it supports two possible values:
  * ENTITY
  * POLICY
*/
type ResourceGroup int

const (
	RESOURCEGROUP_UNKNOWN  ResourceGroup = 0
	RESOURCEGROUP_REDACTED ResourceGroup = 1
	RESOURCEGROUP_ENTITY   ResourceGroup = 2
	RESOURCEGROUP_POLICY   ResourceGroup = 3
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *ResourceGroup) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"ENTITY",
		"POLICY",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e ResourceGroup) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"ENTITY",
		"POLICY",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *ResourceGroup) index(name string) ResourceGroup {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"ENTITY",
		"POLICY",
	}
	for idx := range names {
		if names[idx] == name {
			return ResourceGroup(idx)
		}
	}
	return RESOURCEGROUP_UNKNOWN
}

func (e *ResourceGroup) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for ResourceGroup:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *ResourceGroup) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e ResourceGroup) Ref() *ResourceGroup {
	return &e
}

/*
An enum denoting the associated resource types.<br>
Resource types are further grouped into 2 types - entity or a policy.
*/
type ResourceType int

const (
	RESOURCETYPE_UNKNOWN                 ResourceType = 0
	RESOURCETYPE_REDACTED                ResourceType = 1
	RESOURCETYPE_VM                      ResourceType = 2
	RESOURCETYPE_MH_VM                   ResourceType = 3
	RESOURCETYPE_IMAGE                   ResourceType = 4
	RESOURCETYPE_SUBNET                  ResourceType = 5
	RESOURCETYPE_CLUSTER                 ResourceType = 6
	RESOURCETYPE_HOST                    ResourceType = 7
	RESOURCETYPE_REPORT                  ResourceType = 8
	RESOURCETYPE_MARKETPLACE_ITEM        ResourceType = 9
	RESOURCETYPE_BLUEPRINT               ResourceType = 10
	RESOURCETYPE_APP                     ResourceType = 11
	RESOURCETYPE_VOLUMEGROUP             ResourceType = 12
	RESOURCETYPE_IMAGE_PLACEMENT_POLICY  ResourceType = 13
	RESOURCETYPE_NETWORK_SECURITY_POLICY ResourceType = 14
	RESOURCETYPE_NETWORK_SECURITY_RULE   ResourceType = 15
	RESOURCETYPE_VM_HOST_AFFINITY_POLICY ResourceType = 16
	RESOURCETYPE_QOS_POLICY              ResourceType = 17
	RESOURCETYPE_NGT_POLICY              ResourceType = 18
	RESOURCETYPE_PROTECTION_RULE         ResourceType = 19
	RESOURCETYPE_ACCESS_CONTROL_POLICY   ResourceType = 20
	RESOURCETYPE_STORAGE_POLICY          ResourceType = 21
	RESOURCETYPE_IMAGE_RATE_LIMIT        ResourceType = 22
	RESOURCETYPE_RECOVERY_PLAN           ResourceType = 23
	RESOURCETYPE_BUNDLE                  ResourceType = 24
	RESOURCETYPE_POLICY_SCHEMA           ResourceType = 25
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *ResourceType) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"VM",
		"MH_VM",
		"IMAGE",
		"SUBNET",
		"CLUSTER",
		"HOST",
		"REPORT",
		"MARKETPLACE_ITEM",
		"BLUEPRINT",
		"APP",
		"VOLUMEGROUP",
		"IMAGE_PLACEMENT_POLICY",
		"NETWORK_SECURITY_POLICY",
		"NETWORK_SECURITY_RULE",
		"VM_HOST_AFFINITY_POLICY",
		"QOS_POLICY",
		"NGT_POLICY",
		"PROTECTION_RULE",
		"ACCESS_CONTROL_POLICY",
		"STORAGE_POLICY",
		"IMAGE_RATE_LIMIT",
		"RECOVERY_PLAN",
		"BUNDLE",
		"POLICY_SCHEMA",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e ResourceType) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"VM",
		"MH_VM",
		"IMAGE",
		"SUBNET",
		"CLUSTER",
		"HOST",
		"REPORT",
		"MARKETPLACE_ITEM",
		"BLUEPRINT",
		"APP",
		"VOLUMEGROUP",
		"IMAGE_PLACEMENT_POLICY",
		"NETWORK_SECURITY_POLICY",
		"NETWORK_SECURITY_RULE",
		"VM_HOST_AFFINITY_POLICY",
		"QOS_POLICY",
		"NGT_POLICY",
		"PROTECTION_RULE",
		"ACCESS_CONTROL_POLICY",
		"STORAGE_POLICY",
		"IMAGE_RATE_LIMIT",
		"RECOVERY_PLAN",
		"BUNDLE",
		"POLICY_SCHEMA",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *ResourceType) index(name string) ResourceType {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"VM",
		"MH_VM",
		"IMAGE",
		"SUBNET",
		"CLUSTER",
		"HOST",
		"REPORT",
		"MARKETPLACE_ITEM",
		"BLUEPRINT",
		"APP",
		"VOLUMEGROUP",
		"IMAGE_PLACEMENT_POLICY",
		"NETWORK_SECURITY_POLICY",
		"NETWORK_SECURITY_RULE",
		"VM_HOST_AFFINITY_POLICY",
		"QOS_POLICY",
		"NGT_POLICY",
		"PROTECTION_RULE",
		"ACCESS_CONTROL_POLICY",
		"STORAGE_POLICY",
		"IMAGE_RATE_LIMIT",
		"RECOVERY_PLAN",
		"BUNDLE",
		"POLICY_SCHEMA",
	}
	for idx := range names {
		if names[idx] == name {
			return ResourceType(idx)
		}
	}
	return RESOURCETYPE_UNKNOWN
}

func (e *ResourceType) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for ResourceType:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *ResourceType) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e ResourceType) Ref() *ResourceType {
	return &e
}

/*
The task object tracking an asynchronous operation.
*/
type Task struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  List of globally unique identifiers for clusters associated with the task or any of its subtasks.
	*/
	ClusterExtIds []string `json:"clusterExtIds,omitempty"`
	/*
	  UTC date and time in RFC-3339 format when the task was completed.
	*/
	CompletedTime *time.Time `json:"completedTime,omitempty"`
	/*
	  Additional details on the task to aid the user with further actions post completion of the task.
	*/
	CompletionDetails []import3.KVPair `json:"completionDetails,omitempty"`
	/*
	  UTC date and time in RFC-3339 format when the task was created.
	*/
	CreatedTime *time.Time `json:"createdTime,omitempty"`
	/*
	  Reference to entities associated with the task.
	*/
	EntitiesAffected []EntityReference `json:"entitiesAffected,omitempty"`
	/*
	  Error details explaining a task failure. These would be populated only in the case of task failures.
	*/
	ErrorMessages []import1.AppMessage `json:"errorMessages,omitempty"`
	/*
	  A globally unique identifier of a task.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  Signifies if the task can be cancelled.
	*/
	IsCancelable *bool `json:"isCancelable,omitempty"`
	/*
	  UTC date and time in RFC-3339 format when the task was last updated.
	*/
	LastUpdatedTime *time.Time `json:"lastUpdatedTime,omitempty"`
	/*
	  Provides an error message in the absence of a well-defined error message for the tasks created through legacy APIs.
	*/
	LegacyErrorMessage *string `json:"legacyErrorMessage,omitempty"`
	/*
	  The operation name being tracked by the task.
	*/
	Operation *string `json:"operation,omitempty"`
	/*
	  Description of the operation being tracked by the task.
	*/
	OperationDescription *string `json:"operationDescription,omitempty"`

	OwnedBy *OwnerReference `json:"ownedBy,omitempty"`

	ParentTask *TaskReferenceInternal `json:"parentTask,omitempty"`
	/*
	  Task progress expressed as a percentage.
	*/
	ProgressPercentage *int `json:"progressPercentage,omitempty"`
	/*
	  UTC date and time in RFC-3339 format when the task was started.
	*/
	StartedTime *time.Time `json:"startedTime,omitempty"`

	Status *TaskStatus `json:"status,omitempty"`
	/*
	  List of steps completed as part of the task.
	*/
	SubSteps []TaskStep `json:"subSteps,omitempty"`
	/*
	  Reference to tasks spawned as children of the current task. The task get response would contain a limited number of subtask references. To get the entire list of subtasks for a task, use the parent task filter in the task list API.
	*/
	SubTasks []TaskReferenceInternal `json:"subTasks,omitempty"`
	/*
	  Warning messages to alert the user of issues which did not directly cause task failure. These can be populated for any task.
	*/
	Warnings []import1.AppMessage `json:"warnings,omitempty"`
}

func NewTask() *Task {
	p := new(Task)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.config.Task"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
A reference to a task tracking an asynchronous operation. The status of the task can be queried by making a GET request to the task URI provided in the metadata section of the API response.
*/
type TaskReference struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  A globally unique identifier of a task.
	*/
	ExtId *string `json:"extId,omitempty"`
}

func NewTaskReference() *TaskReference {
	p := new(TaskReference)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.config.TaskReference"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Reference to the parent task associated with the current task.
*/
type TaskReferenceInternal struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  A globally unique identifier of the task.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  The URL at which the entity described by the link can be accessed.
	*/
	Href *string `json:"href,omitempty"`
	/*
	  A name that identifies the relationship of the link to the object that is returned by the URL.  The unique value of "self" identifies the URL for the object.
	*/
	Rel *string `json:"rel,omitempty"`
}

func NewTaskReferenceInternal() *TaskReferenceInternal {
	p := new(TaskReferenceInternal)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.config.TaskReferenceInternal"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Status of the task.
*/
type TaskStatus int

const (
	TASKSTATUS_UNKNOWN   TaskStatus = 0
	TASKSTATUS_REDACTED  TaskStatus = 1
	TASKSTATUS_QUEUED    TaskStatus = 2
	TASKSTATUS_RUNNING   TaskStatus = 3
	TASKSTATUS_CANCELING TaskStatus = 4
	TASKSTATUS_SUCCEEDED TaskStatus = 5
	TASKSTATUS_FAILED    TaskStatus = 6
	TASKSTATUS_CANCELED  TaskStatus = 7
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *TaskStatus) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"QUEUED",
		"RUNNING",
		"CANCELING",
		"SUCCEEDED",
		"FAILED",
		"CANCELED",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e TaskStatus) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"QUEUED",
		"RUNNING",
		"CANCELING",
		"SUCCEEDED",
		"FAILED",
		"CANCELED",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *TaskStatus) index(name string) TaskStatus {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"QUEUED",
		"RUNNING",
		"CANCELING",
		"SUCCEEDED",
		"FAILED",
		"CANCELED",
	}
	for idx := range names {
		if names[idx] == name {
			return TaskStatus(idx)
		}
	}
	return TASKSTATUS_UNKNOWN
}

func (e *TaskStatus) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for TaskStatus:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *TaskStatus) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e TaskStatus) Ref() *TaskStatus {
	return &e
}

/*
A single step in the task.
*/
type TaskStep struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Message describing the completed steps for the task.
	*/
	Name *string `json:"name,omitempty"`
}

func NewTaskStep() *TaskStep {
	p := new(TaskStep)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.config.TaskStep"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
REST response for all response codes in API path /prism/v4.0.b1/config/categories/{extId} Put operation
*/
type UpdateCategoryApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfUpdateCategoryApiResponseData `json:"data,omitempty"`

	Metadata *import2.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewUpdateCategoryApiResponse() *UpdateCategoryApiResponse {
	p := new(UpdateCategoryApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.config.UpdateCategoryApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *UpdateCategoryApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *UpdateCategoryApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfUpdateCategoryApiResponseData()
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

type OneOfCreateCategoryApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType400  *import1.ErrorResponse `json:"-"`
	oneOfType0    *Category              `json:"-"`
}

func NewOneOfCreateCategoryApiResponseData() *OneOfCreateCategoryApiResponseData {
	p := new(OneOfCreateCategoryApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfCreateCategoryApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfCreateCategoryApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import1.ErrorResponse)
		}
		*p.oneOfType400 = v.(import1.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	case Category:
		if nil == p.oneOfType0 {
			p.oneOfType0 = new(Category)
		}
		*p.oneOfType0 = v.(Category)
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

func (p *OneOfCreateCategoryApiResponseData) GetValue() interface{} {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return *p.oneOfType0
	}
	return nil
}

func (p *OneOfCreateCategoryApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType400 := new(import1.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "prism.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import1.ErrorResponse)
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
	vOneOfType0 := new(Category)
	if err := json.Unmarshal(b, vOneOfType0); err == nil {
		if "prism.v4.config.Category" == *vOneOfType0.ObjectType_ {
			if nil == p.oneOfType0 {
				p.oneOfType0 = new(Category)
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
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfCreateCategoryApiResponseData"))
}

func (p *OneOfCreateCategoryApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType0)
	}
	return nil, errors.New("No value to marshal for OneOfCreateCategoryApiResponseData")
}

type OneOfListCategoriesApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType401  []CategoryProjection   `json:"-"`
	oneOfType400  *import1.ErrorResponse `json:"-"`
	oneOfType0    []Category             `json:"-"`
}

func NewOneOfListCategoriesApiResponseData() *OneOfListCategoriesApiResponseData {
	p := new(OneOfListCategoriesApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfListCategoriesApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfListCategoriesApiResponseData is nil"))
	}
	switch v.(type) {
	case []CategoryProjection:
		p.oneOfType401 = v.([]CategoryProjection)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = "List<prism.v4.config.CategoryProjection>"
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = "List<prism.v4.config.CategoryProjection>"
	case import1.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import1.ErrorResponse)
		}
		*p.oneOfType400 = v.(import1.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	case []Category:
		p.oneOfType0 = v.([]Category)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = "List<prism.v4.config.Category>"
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = "List<prism.v4.config.Category>"
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfListCategoriesApiResponseData) GetValue() interface{} {
	if "List<prism.v4.config.CategoryProjection>" == *p.Discriminator {
		return p.oneOfType401
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	if "List<prism.v4.config.Category>" == *p.Discriminator {
		return p.oneOfType0
	}
	return nil
}

func (p *OneOfListCategoriesApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType401 := new([]CategoryProjection)
	if err := json.Unmarshal(b, vOneOfType401); err == nil {

		if len(*vOneOfType401) == 0 || "prism.v4.config.CategoryProjection" == *((*vOneOfType401)[0].ObjectType_) {
			p.oneOfType401 = *vOneOfType401
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = "List<prism.v4.config.CategoryProjection>"
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = "List<prism.v4.config.CategoryProjection>"
			return nil

		}
	}
	vOneOfType400 := new(import1.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "prism.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import1.ErrorResponse)
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
	vOneOfType0 := new([]Category)
	if err := json.Unmarshal(b, vOneOfType0); err == nil {

		if len(*vOneOfType0) == 0 || "prism.v4.config.Category" == *((*vOneOfType0)[0].ObjectType_) {
			p.oneOfType0 = *vOneOfType0
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = "List<prism.v4.config.Category>"
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = "List<prism.v4.config.Category>"
			return nil

		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfListCategoriesApiResponseData"))
}

func (p *OneOfListCategoriesApiResponseData) MarshalJSON() ([]byte, error) {
	if "List<prism.v4.config.CategoryProjection>" == *p.Discriminator {
		return json.Marshal(p.oneOfType401)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	if "List<prism.v4.config.Category>" == *p.Discriminator {
		return json.Marshal(p.oneOfType0)
	}
	return nil, errors.New("No value to marshal for OneOfListCategoriesApiResponseData")
}

type OneOfUpdateCategoryApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType400  *import1.ErrorResponse `json:"-"`
	oneOfType0    *Category              `json:"-"`
}

func NewOneOfUpdateCategoryApiResponseData() *OneOfUpdateCategoryApiResponseData {
	p := new(OneOfUpdateCategoryApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfUpdateCategoryApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfUpdateCategoryApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import1.ErrorResponse)
		}
		*p.oneOfType400 = v.(import1.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	case Category:
		if nil == p.oneOfType0 {
			p.oneOfType0 = new(Category)
		}
		*p.oneOfType0 = v.(Category)
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

func (p *OneOfUpdateCategoryApiResponseData) GetValue() interface{} {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return *p.oneOfType0
	}
	return nil
}

func (p *OneOfUpdateCategoryApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType400 := new(import1.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "prism.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import1.ErrorResponse)
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
	vOneOfType0 := new(Category)
	if err := json.Unmarshal(b, vOneOfType0); err == nil {
		if "prism.v4.config.Category" == *vOneOfType0.ObjectType_ {
			if nil == p.oneOfType0 {
				p.oneOfType0 = new(Category)
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
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfUpdateCategoryApiResponseData"))
}

func (p *OneOfUpdateCategoryApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType0)
	}
	return nil, errors.New("No value to marshal for OneOfUpdateCategoryApiResponseData")
}

type OneOfCancelTaskApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *import1.AppMessage    `json:"-"`
	oneOfType400  *import1.ErrorResponse `json:"-"`
}

func NewOneOfCancelTaskApiResponseData() *OneOfCancelTaskApiResponseData {
	p := new(OneOfCancelTaskApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfCancelTaskApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfCancelTaskApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.AppMessage:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(import1.AppMessage)
		}
		*p.oneOfType2001 = v.(import1.AppMessage)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import1.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import1.ErrorResponse)
		}
		*p.oneOfType400 = v.(import1.ErrorResponse)
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

func (p *OneOfCancelTaskApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfCancelTaskApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(import1.AppMessage)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "prism.v4.error.AppMessage" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(import1.AppMessage)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import1.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "prism.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import1.ErrorResponse)
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
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfCancelTaskApiResponseData"))
}

func (p *OneOfCancelTaskApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfCancelTaskApiResponseData")
}

type OneOfGetTaskApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *Task                  `json:"-"`
	oneOfType400  *import1.ErrorResponse `json:"-"`
}

func NewOneOfGetTaskApiResponseData() *OneOfGetTaskApiResponseData {
	p := new(OneOfGetTaskApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfGetTaskApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfGetTaskApiResponseData is nil"))
	}
	switch v.(type) {
	case Task:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(Task)
		}
		*p.oneOfType2001 = v.(Task)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import1.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import1.ErrorResponse)
		}
		*p.oneOfType400 = v.(import1.ErrorResponse)
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

func (p *OneOfGetTaskApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfGetTaskApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(Task)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "prism.v4.config.Task" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(Task)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import1.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "prism.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import1.ErrorResponse)
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
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfGetTaskApiResponseData"))
}

func (p *OneOfGetTaskApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfGetTaskApiResponseData")
}

type OneOfListTasksApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 []Task                 `json:"-"`
	oneOfType400  *import1.ErrorResponse `json:"-"`
}

func NewOneOfListTasksApiResponseData() *OneOfListTasksApiResponseData {
	p := new(OneOfListTasksApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfListTasksApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfListTasksApiResponseData is nil"))
	}
	switch v.(type) {
	case []Task:
		p.oneOfType2001 = v.([]Task)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = "List<prism.v4.config.Task>"
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = "List<prism.v4.config.Task>"
	case import1.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import1.ErrorResponse)
		}
		*p.oneOfType400 = v.(import1.ErrorResponse)
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

func (p *OneOfListTasksApiResponseData) GetValue() interface{} {
	if "List<prism.v4.config.Task>" == *p.Discriminator {
		return p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfListTasksApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new([]Task)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {

		if len(*vOneOfType2001) == 0 || "prism.v4.config.Task" == *((*vOneOfType2001)[0].ObjectType_) {
			p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = "List<prism.v4.config.Task>"
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = "List<prism.v4.config.Task>"
			return nil

		}
	}
	vOneOfType400 := new(import1.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "prism.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import1.ErrorResponse)
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
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfListTasksApiResponseData"))
}

func (p *OneOfListTasksApiResponseData) MarshalJSON() ([]byte, error) {
	if "List<prism.v4.config.Task>" == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfListTasksApiResponseData")
}

type OneOfGetCategoryApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType400  *import1.ErrorResponse `json:"-"`
	oneOfType0    *Category              `json:"-"`
}

func NewOneOfGetCategoryApiResponseData() *OneOfGetCategoryApiResponseData {
	p := new(OneOfGetCategoryApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfGetCategoryApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfGetCategoryApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import1.ErrorResponse)
		}
		*p.oneOfType400 = v.(import1.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	case Category:
		if nil == p.oneOfType0 {
			p.oneOfType0 = new(Category)
		}
		*p.oneOfType0 = v.(Category)
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

func (p *OneOfGetCategoryApiResponseData) GetValue() interface{} {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return *p.oneOfType0
	}
	return nil
}

func (p *OneOfGetCategoryApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType400 := new(import1.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "prism.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import1.ErrorResponse)
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
	vOneOfType0 := new(Category)
	if err := json.Unmarshal(b, vOneOfType0); err == nil {
		if "prism.v4.config.Category" == *vOneOfType0.ObjectType_ {
			if nil == p.oneOfType0 {
				p.oneOfType0 = new(Category)
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
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfGetCategoryApiResponseData"))
}

func (p *OneOfGetCategoryApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType0)
	}
	return nil, errors.New("No value to marshal for OneOfGetCategoryApiResponseData")
}

type OneOfDeleteCategoryApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType1    *interface{}           `json:"-"`
	oneOfType400  *import1.ErrorResponse `json:"-"`
}

func NewOneOfDeleteCategoryApiResponseData() *OneOfDeleteCategoryApiResponseData {
	p := new(OneOfDeleteCategoryApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfDeleteCategoryApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfDeleteCategoryApiResponseData is nil"))
	}
	if nil == v {
		if nil == p.oneOfType1 {
			p.oneOfType1 = new(interface{})
		}
		*p.oneOfType1 = nil
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = "EMPTY"
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = "EMPTY"
		return nil
	}
	switch v.(type) {
	case import1.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import1.ErrorResponse)
		}
		*p.oneOfType400 = v.(import1.ErrorResponse)
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

func (p *OneOfDeleteCategoryApiResponseData) GetValue() interface{} {
	if "EMPTY" == *p.Discriminator {
		return *p.oneOfType1
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfDeleteCategoryApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType1 := new(interface{})
	if err := json.Unmarshal(b, vOneOfType1); err == nil {
		if nil == *vOneOfType1 {
			if nil == p.oneOfType1 {
				p.oneOfType1 = new(interface{})
			}
			*p.oneOfType1 = nil
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = "EMPTY"
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = "EMPTY"
			return nil
		}
	}
	vOneOfType400 := new(import1.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "prism.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import1.ErrorResponse)
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
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfDeleteCategoryApiResponseData"))
}

func (p *OneOfDeleteCategoryApiResponseData) MarshalJSON() ([]byte, error) {
	if "EMPTY" == *p.Discriminator {
		return json.Marshal(p.oneOfType1)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfDeleteCategoryApiResponseData")
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
