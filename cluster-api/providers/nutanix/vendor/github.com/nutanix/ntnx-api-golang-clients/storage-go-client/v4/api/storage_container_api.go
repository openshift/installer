//Api classes for storage's golang SDK
package api

import (
	"encoding/json"
	"github.com/nutanix/ntnx-api-golang-clients/storage-go-client/v4/client"
	import1 "github.com/nutanix/ntnx-api-golang-clients/storage-go-client/v4/models/storage/v4/config"
	"net/http"
	"net/url"
	"strings"
)

type StorageContainerApi struct {
	ApiClient *client.ApiClient
}

func NewStorageContainerApi(apiClient *client.ApiClient) *StorageContainerApi {
	if apiClient == nil {
		apiClient = client.NewApiClient()
	}

	a := &StorageContainerApi{
		ApiClient: apiClient,
	}
	return a
}

/**
  Add a Data Store instance to the cluster.
  Add a Data Store instance to the cluster.

  parameters:-
  -> body (storage.v4.config.DataStoreMount) (required)
  -> containerExtId (string) (required) : Uuid of the storage container.
  -> args (map[string]interface{}) (optional) : Additional Arguments

  returns: (*storage.v4.config.StorageContainerTaskResponse, error)
*/
func (api *StorageContainerApi) AddDataStoreForCluster(body *import1.DataStoreMount, containerExtId *string, args ...map[string]interface{}) (*import1.StorageContainerTaskResponse, error) {
	argMap := make(map[string]interface{})
	if len(args) > 0 {
		argMap = args[0]
	}

	uri := "/api/storage/v4.0.a3/config/storage-containers/{containerExtId}/$actions/mount"

	// verify the required parameter 'body' is set
	if nil == body {
		return nil, client.ReportError("body is required and must be specified")
	}
	// verify the required parameter 'containerExtId' is set
	if nil == containerExtId {
		return nil, client.ReportError("containerExtId is required and must be specified")
	}

	// Path Params
	uri = strings.Replace(uri, "{"+"containerExtId"+"}", url.PathEscape(client.ParameterToString(*containerExtId, "")), -1)
	headerParams := make(map[string]string)
	queryParams := url.Values{}
	formParams := url.Values{}

	// to determine the Content-Type header
	contentTypes := []string{"application/json"}

	// to determine the Accept header
	accepts := []string{"application/json"}

	// Header Params
	if ifMatch, ifMatchOk := argMap["If-Match"].(string); ifMatchOk {
		headerParams["If-Match"] = ifMatch
	}
	if ifNoneMatch, ifNoneMatchOk := argMap["If-None-Match"].(string); ifNoneMatchOk {
		headerParams["If-None-Match"] = ifNoneMatch
	}
	authNames := []string{"basicAuthScheme"}

	responseBody, err := api.ApiClient.CallApi(&uri, http.MethodPost, body, queryParams, headerParams, formParams, accepts, contentTypes, authNames)
	if nil != err || nil == responseBody {
		return nil, err
	}
	unmarshalledResp := new(import1.StorageContainerTaskResponse)
	json.Unmarshal(responseBody, &unmarshalledResp)
	return unmarshalledResp, err
}

/**
  Add a Storage Container instance to the cluster.
  Add a Storage Container instance to the cluster.

  parameters:-
  -> body (storage.v4.config.StorageContainer) (required)
  -> xClusterId (string) (required) : Remote cluster uuid to forward the request.
  -> args (map[string]interface{}) (optional) : Additional Arguments

  returns: (*storage.v4.config.StorageContainerTaskResponse, error)
*/
func (api *StorageContainerApi) AddStorageContainerForCluster(body *import1.StorageContainer, xClusterId *string, args ...map[string]interface{}) (*import1.StorageContainerTaskResponse, error) {
	argMap := make(map[string]interface{})
	if len(args) > 0 {
		argMap = args[0]
	}

	uri := "/api/storage/v4.0.a3/config/storage-containers"

	// verify the required parameter 'body' is set
	if nil == body {
		return nil, client.ReportError("body is required and must be specified")
	}
	// verify the required parameter 'xClusterId' is set
	if nil == xClusterId {
		return nil, client.ReportError("xClusterId is required and must be specified")
	}

	headerParams := make(map[string]string)
	queryParams := url.Values{}
	formParams := url.Values{}

	// to determine the Content-Type header
	contentTypes := []string{"application/json"}

	// to determine the Accept header
	accepts := []string{"application/json"}

	// Header Params
	if ifMatch, ifMatchOk := argMap["If-Match"].(string); ifMatchOk {
		headerParams["If-Match"] = ifMatch
	}
	if ifNoneMatch, ifNoneMatchOk := argMap["If-None-Match"].(string); ifNoneMatchOk {
		headerParams["If-None-Match"] = ifNoneMatch
	}
	headerParams["X-Cluster-Id"] = client.ParameterToString(*xClusterId, "")

	authNames := []string{"basicAuthScheme"}

	responseBody, err := api.ApiClient.CallApi(&uri, http.MethodPost, body, queryParams, headerParams, formParams, accepts, contentTypes, authNames)
	if nil != err || nil == responseBody {
		return nil, err
	}
	unmarshalledResp := new(import1.StorageContainerTaskResponse)
	json.Unmarshal(responseBody, &unmarshalledResp)
	return unmarshalledResp, err
}

/**
  delete a Data Store instance in the cluster.
  delete a Data Store instance in the cluster.

  parameters:-
  -> body (storage.v4.config.DataStoreUnmount) (required)
  -> containerExtId (string) (required) : Uuid of the storage container.
  -> args (map[string]interface{}) (optional) : Additional Arguments

  returns: (*storage.v4.config.DataStoreTaskResponse, error)
*/
func (api *StorageContainerApi) DeleteDataStoreForCluster(body *import1.DataStoreUnmount, containerExtId *string, args ...map[string]interface{}) (*import1.DataStoreTaskResponse, error) {
	argMap := make(map[string]interface{})
	if len(args) > 0 {
		argMap = args[0]
	}

	uri := "/api/storage/v4.0.a3/config/storage-containers/{containerExtId}/$actions/unmount"

	// verify the required parameter 'body' is set
	if nil == body {
		return nil, client.ReportError("body is required and must be specified")
	}
	// verify the required parameter 'containerExtId' is set
	if nil == containerExtId {
		return nil, client.ReportError("containerExtId is required and must be specified")
	}

	// Path Params
	uri = strings.Replace(uri, "{"+"containerExtId"+"}", url.PathEscape(client.ParameterToString(*containerExtId, "")), -1)
	headerParams := make(map[string]string)
	queryParams := url.Values{}
	formParams := url.Values{}

	// to determine the Content-Type header
	contentTypes := []string{"application/json"}

	// to determine the Accept header
	accepts := []string{"application/json"}

	// Header Params
	if ifMatch, ifMatchOk := argMap["If-Match"].(string); ifMatchOk {
		headerParams["If-Match"] = ifMatch
	}
	if ifNoneMatch, ifNoneMatchOk := argMap["If-None-Match"].(string); ifNoneMatchOk {
		headerParams["If-None-Match"] = ifNoneMatch
	}
	authNames := []string{"basicAuthScheme"}

	responseBody, err := api.ApiClient.CallApi(&uri, http.MethodPost, body, queryParams, headerParams, formParams, accepts, contentTypes, authNames)
	if nil != err || nil == responseBody {
		return nil, err
	}
	unmarshalledResp := new(import1.DataStoreTaskResponse)
	json.Unmarshal(responseBody, &unmarshalledResp)
	return unmarshalledResp, err
}

/**
  Delete an existing storage container.
  Delete an existing storage container.

  parameters:-
  -> containerExtId (string) (required) : Uuid of the storage container.
  -> ignoreSmallFiles (bool) (optional)
  -> args (map[string]interface{}) (optional) : Additional Arguments

  returns: (*storage.v4.config.StorageContainerTaskResponse, error)
*/
func (api *StorageContainerApi) DeleteStorageContainerByExtId(containerExtId *string, ignoreSmallFiles *bool, args ...map[string]interface{}) (*import1.StorageContainerTaskResponse, error) {
	argMap := make(map[string]interface{})
	if len(args) > 0 {
		argMap = args[0]
	}

	uri := "/api/storage/v4.0.a3/config/storage-containers/{containerExtId}"

	// verify the required parameter 'containerExtId' is set
	if nil == containerExtId {
		return nil, client.ReportError("containerExtId is required and must be specified")
	}

	// Path Params
	uri = strings.Replace(uri, "{"+"containerExtId"+"}", url.PathEscape(client.ParameterToString(*containerExtId, "")), -1)
	headerParams := make(map[string]string)
	queryParams := url.Values{}
	formParams := url.Values{}

	// to determine the Content-Type header
	contentTypes := []string{}

	// to determine the Accept header
	accepts := []string{"application/json"}

	// Query Params
	if ignoreSmallFiles != nil {
		queryParams.Add("ignoreSmallFiles", client.ParameterToString(*ignoreSmallFiles, ""))
	}

	// Header Params
	if ifMatch, ifMatchOk := argMap["If-Match"].(string); ifMatchOk {
		headerParams["If-Match"] = ifMatch
	}
	if ifNoneMatch, ifNoneMatchOk := argMap["If-None-Match"].(string); ifNoneMatchOk {
		headerParams["If-None-Match"] = ifNoneMatch
	}
	authNames := []string{"basicAuthScheme"}

	responseBody, err := api.ApiClient.CallApi(&uri, http.MethodDelete, nil, queryParams, headerParams, formParams, accepts, contentTypes, authNames)
	if nil != err || nil == responseBody {
		return nil, err
	}
	unmarshalledResp := new(import1.StorageContainerTaskResponse)
	json.Unmarshal(responseBody, &unmarshalledResp)
	return unmarshalledResp, err
}

/**
  Get the list of all storage containers configured in the cluster.
  Get the list of all storage containers configured in the cluster.

  parameters:-
  -> page_ (int) (optional) : A URL query parameter that specifies the page number of the result set.  Must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range will lead to no results being returned.
  -> limit_ (int) (optional) : A URL query parameter that specifies the total number of records returned in the result set.  Must be a positive integer between 0 and 100. Any number out of this range will lead to a validation error. If the limit is not provided a default value of 50 records will be returned in the result set.
  -> filter_ (string) (optional) : A URL query parameter that allows clients to filter a collection of resources. The expression specified with $filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the $filter must conform to the OData V4.01 URL conventions. The filter can be applied on the following fields: - affinityHostUuid - clusterExtId - containerExtId - containerId - name - ownerUuid - replicationFactor - storagePoolUuid
  -> orderby_ (string) (optional) : A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified the resources will be sorted in ascending order by default. For example, 'orderby=templateName desc' would get all templates sorted by templateName in desc order. The orderby can be applied to the following fields: - name
  -> select_ (string) (optional) : A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the $select must conform to the OData V4.01 URL conventions. If a $select expression consists of a single select item that is an asterisk (i.e. *), then all properties on the matching resource will be returned. - advertisedCapacity - affinityHostUuid - clusterExtId - compressionDelayInSecs - compressionEnabled - containerExtId - containerId - downMigrateTimesInSecs - enableSoftwareEncryption - encrypted - erasureCode - erasureCodeDelaySecs - extId - fingerPrintOnWrite - inlineEcEnabled - isNfsWhitelistInherited - isNutanixManaged - links - mappedRemoteContainers - markedForRemoval - maxCapacity - name - nfsWhitelistAddress - onDiskDedup - oplogReplicationFactor - ownerUuid - preferHigherECFaultDomain - randomIoPreference - replicationFactor - seqIoPreference - storagePoolUuid - tenantId - totalExplicitReservedCapacity - totalImplicitReservedCapacity - vstoreNameList
  -> args (map[string]interface{}) (optional) : Additional Arguments

  returns: (*storage.v4.config.StorageContainersResponse, error)
*/
func (api *StorageContainerApi) GetAllStorageContainers(page_ *int, limit_ *int, filter_ *string, orderby_ *string, select_ *string, args ...map[string]interface{}) (*import1.StorageContainersResponse, error) {
	argMap := make(map[string]interface{})
	if len(args) > 0 {
		argMap = args[0]
	}

	uri := "/api/storage/v4.0.a3/config/storage-containers"

	headerParams := make(map[string]string)
	queryParams := url.Values{}
	formParams := url.Values{}

	// to determine the Content-Type header
	contentTypes := []string{}

	// to determine the Accept header
	accepts := []string{"application/json"}

	// Query Params
	if page_ != nil {
		queryParams.Add("$page", client.ParameterToString(*page_, ""))
	}
	if limit_ != nil {
		queryParams.Add("$limit", client.ParameterToString(*limit_, ""))
	}
	if filter_ != nil {
		queryParams.Add("$filter", client.ParameterToString(*filter_, ""))
	}
	if orderby_ != nil {
		queryParams.Add("$orderby", client.ParameterToString(*orderby_, ""))
	}
	if select_ != nil {
		queryParams.Add("$select", client.ParameterToString(*select_, ""))
	}

	// Header Params
	if ifMatch, ifMatchOk := argMap["If-Match"].(string); ifMatchOk {
		headerParams["If-Match"] = ifMatch
	}
	if ifNoneMatch, ifNoneMatchOk := argMap["If-None-Match"].(string); ifNoneMatchOk {
		headerParams["If-None-Match"] = ifNoneMatch
	}
	authNames := []string{"basicAuthScheme"}

	responseBody, err := api.ApiClient.CallApi(&uri, http.MethodGet, nil, queryParams, headerParams, formParams, accepts, contentTypes, authNames)
	if nil != err || nil == responseBody {
		return nil, err
	}
	unmarshalledResp := new(import1.StorageContainersResponse)
	json.Unmarshal(responseBody, &unmarshalledResp)
	return unmarshalledResp, err
}

/**
  Get the datastore.
  Get the datastore.

  parameters:-
  -> filter_ (string) (optional) : A URL query parameter that allows clients to filter a collection of resources. The expression specified with $filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the $filter must conform to the OData V4.01 URL conventions. The filter can be applied on the following fields: - containerExtId
  -> args (map[string]interface{}) (optional) : Additional Arguments

  returns: (*storage.v4.config.DataStoreResponse, error)
*/
func (api *StorageContainerApi) GetDataStores(filter_ *string, args ...map[string]interface{}) (*import1.DataStoreResponse, error) {
	argMap := make(map[string]interface{})
	if len(args) > 0 {
		argMap = args[0]
	}

	uri := "/api/storage/v4.0.a3/config/storage-containers/datastores"

	headerParams := make(map[string]string)
	queryParams := url.Values{}
	formParams := url.Values{}

	// to determine the Content-Type header
	contentTypes := []string{}

	// to determine the Accept header
	accepts := []string{"application/json"}

	// Query Params
	if filter_ != nil {
		queryParams.Add("$filter", client.ParameterToString(*filter_, ""))
	}

	// Header Params
	if ifMatch, ifMatchOk := argMap["If-Match"].(string); ifMatchOk {
		headerParams["If-Match"] = ifMatch
	}
	if ifNoneMatch, ifNoneMatchOk := argMap["If-None-Match"].(string); ifNoneMatchOk {
		headerParams["If-None-Match"] = ifNoneMatch
	}
	authNames := []string{"basicAuthScheme"}

	responseBody, err := api.ApiClient.CallApi(&uri, http.MethodGet, nil, queryParams, headerParams, formParams, accepts, contentTypes, authNames)
	if nil != err || nil == responseBody {
		return nil, err
	}
	unmarshalledResp := new(import1.DataStoreResponse)
	json.Unmarshal(responseBody, &unmarshalledResp)
	return unmarshalledResp, err
}

/**
  Get the attributes of an existing storage container.
  Get the attributes of an existing storage container.

  parameters:-
  -> containerExtId (string) (required) : Uuid of the storage container.
  -> args (map[string]interface{}) (optional) : Additional Arguments

  returns: (*storage.v4.config.StorageContainerResponse, error)
*/
func (api *StorageContainerApi) GetStorageContainerByExtId(containerExtId *string, args ...map[string]interface{}) (*import1.StorageContainerResponse, error) {
	argMap := make(map[string]interface{})
	if len(args) > 0 {
		argMap = args[0]
	}

	uri := "/api/storage/v4.0.a3/config/storage-containers/{containerExtId}"

	// verify the required parameter 'containerExtId' is set
	if nil == containerExtId {
		return nil, client.ReportError("containerExtId is required and must be specified")
	}

	// Path Params
	uri = strings.Replace(uri, "{"+"containerExtId"+"}", url.PathEscape(client.ParameterToString(*containerExtId, "")), -1)
	headerParams := make(map[string]string)
	queryParams := url.Values{}
	formParams := url.Values{}

	// to determine the Content-Type header
	contentTypes := []string{}

	// to determine the Accept header
	accepts := []string{"application/json"}

	// Header Params
	if ifMatch, ifMatchOk := argMap["If-Match"].(string); ifMatchOk {
		headerParams["If-Match"] = ifMatch
	}
	if ifNoneMatch, ifNoneMatchOk := argMap["If-None-Match"].(string); ifNoneMatchOk {
		headerParams["If-None-Match"] = ifNoneMatch
	}
	authNames := []string{"basicAuthScheme"}

	responseBody, err := api.ApiClient.CallApi(&uri, http.MethodGet, nil, queryParams, headerParams, formParams, accepts, contentTypes, authNames)
	if nil != err || nil == responseBody {
		return nil, err
	}
	unmarshalledResp := new(import1.StorageContainerResponse)
	json.Unmarshal(responseBody, &unmarshalledResp)
	return unmarshalledResp, err
}

/**
  Updates the attributes of an existing storage container.
  Updates the attributes of an existing storage container.

  parameters:-
  -> body (storage.v4.config.StorageContainer) (required)
  -> containerExtId (string) (required) : Uuid of the storage container.
  -> args (map[string]interface{}) (optional) : Additional Arguments

  returns: (*storage.v4.config.StorageContainerTaskResponse, error)
*/
func (api *StorageContainerApi) UpdateStorageContainer(body *import1.StorageContainer, containerExtId *string, args ...map[string]interface{}) (*import1.StorageContainerTaskResponse, error) {
	argMap := make(map[string]interface{})
	if len(args) > 0 {
		argMap = args[0]
	}

	uri := "/api/storage/v4.0.a3/config/storage-containers/{containerExtId}"

	// verify the required parameter 'body' is set
	if nil == body {
		return nil, client.ReportError("body is required and must be specified")
	}
	// verify the required parameter 'containerExtId' is set
	if nil == containerExtId {
		return nil, client.ReportError("containerExtId is required and must be specified")
	}

	// Path Params
	uri = strings.Replace(uri, "{"+"containerExtId"+"}", url.PathEscape(client.ParameterToString(*containerExtId, "")), -1)
	headerParams := make(map[string]string)
	queryParams := url.Values{}
	formParams := url.Values{}

	// to determine the Content-Type header
	contentTypes := []string{"application/json"}

	// to determine the Accept header
	accepts := []string{"application/json"}

	// Header Params
	if ifMatch, ifMatchOk := argMap["If-Match"].(string); ifMatchOk {
		headerParams["If-Match"] = ifMatch
	}
	if ifNoneMatch, ifNoneMatchOk := argMap["If-None-Match"].(string); ifNoneMatchOk {
		headerParams["If-None-Match"] = ifNoneMatch
	}
	authNames := []string{"basicAuthScheme"}

	responseBody, err := api.ApiClient.CallApi(&uri, http.MethodPut, body, queryParams, headerParams, formParams, accepts, contentTypes, authNames)
	if nil != err || nil == responseBody {
		return nil, err
	}
	unmarshalledResp := new(import1.StorageContainerTaskResponse)
	json.Unmarshal(responseBody, &unmarshalledResp)
	return unmarshalledResp, err
}
