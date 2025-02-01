//Api classes for volumes's golang SDK
package api

import (
	"encoding/json"
	"github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/client"
	import2 "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/common/v1/stats"
	import1 "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/volumes/v4/config"
	import3 "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/volumes/v4/stats"
	"net/http"
	"net/url"
	"strings"
)

type VolumeGroupsApi struct {
	ApiClient     *client.ApiClient
	headersToSkip map[string]bool
}

func NewVolumeGroupsApi(apiClient *client.ApiClient) *VolumeGroupsApi {
	if apiClient == nil {
		apiClient = client.NewApiClient()
	}

	a := &VolumeGroupsApi{
		ApiClient: apiClient,
	}

	headers := []string{"authorization", "cookie", "host", "user-agent"}
	a.headersToSkip = make(map[string]bool)
	for _, header := range headers {
		a.headersToSkip[header] = true
	}

	return a
}

// Attaches iSCSI initiator to a Volume Group identified by {extId}.
func (api *VolumeGroupsApi) AttachIscsiClient(extId *string, body *import1.IscsiClient, args ...map[string]interface{}) (*import1.AttachIscsiClientApiResponse, error) {
	argMap := make(map[string]interface{})
	if len(args) > 0 {
		argMap = args[0]
	}

	uri := "/api/volumes/v4.0.b1/config/volume-groups/{extId}/$actions/attach-iscsi-client"

	// verify the required parameter 'extId' is set
	if nil == extId {
		return nil, client.ReportError("extId is required and must be specified")
	}
	// verify the required parameter 'body' is set
	if nil == body {
		return nil, client.ReportError("body is required and must be specified")
	}

	// Path Params

	uri = strings.Replace(uri, "{"+"extId"+"}", url.PathEscape(client.ParameterToString(*extId, "")), -1)
	headerParams := make(map[string]string)
	queryParams := url.Values{}
	formParams := url.Values{}

	// to determine the Content-Type header
	contentTypes := []string{"application/json"}

	// to determine the Accept header
	accepts := []string{"application/json"}

	// Headers provided explicitly on operation takes precedence
	for headerKey, value := range argMap {
		// Skip platform generated headers
		if !api.headersToSkip[strings.ToLower(headerKey)] {
			if value != nil {
				if headerValue, headerValueOk := value.(string); headerValueOk {
					headerParams[headerKey] = headerValue
				}
			}
		}
	}

	authNames := []string{"basicAuthScheme"}

	responseBody, err := api.ApiClient.CallApi(&uri, http.MethodPost, body, queryParams, headerParams, formParams, accepts, contentTypes, authNames)
	if nil != err || nil == responseBody {
		return nil, err
	}

	unmarshalledResp := new(import1.AttachIscsiClientApiResponse)
	json.Unmarshal(responseBody.([]byte), &unmarshalledResp)
	return unmarshalledResp, err
}

// Attaches VM to a Volume Group identified by {extId}.
func (api *VolumeGroupsApi) AttachVm(extId *string, body *import1.VmAttachment, args ...map[string]interface{}) (*import1.AttachVmApiResponse, error) {
	argMap := make(map[string]interface{})
	if len(args) > 0 {
		argMap = args[0]
	}

	uri := "/api/volumes/v4.0.b1/config/volume-groups/{extId}/$actions/attach-vm"

	// verify the required parameter 'extId' is set
	if nil == extId {
		return nil, client.ReportError("extId is required and must be specified")
	}
	// verify the required parameter 'body' is set
	if nil == body {
		return nil, client.ReportError("body is required and must be specified")
	}

	// Path Params

	uri = strings.Replace(uri, "{"+"extId"+"}", url.PathEscape(client.ParameterToString(*extId, "")), -1)
	headerParams := make(map[string]string)
	queryParams := url.Values{}
	formParams := url.Values{}

	// to determine the Content-Type header
	contentTypes := []string{"application/json"}

	// to determine the Accept header
	accepts := []string{"application/json"}

	// Headers provided explicitly on operation takes precedence
	for headerKey, value := range argMap {
		// Skip platform generated headers
		if !api.headersToSkip[strings.ToLower(headerKey)] {
			if value != nil {
				if headerValue, headerValueOk := value.(string); headerValueOk {
					headerParams[headerKey] = headerValue
				}
			}
		}
	}

	authNames := []string{"basicAuthScheme"}

	responseBody, err := api.ApiClient.CallApi(&uri, http.MethodPost, body, queryParams, headerParams, formParams, accepts, contentTypes, authNames)
	if nil != err || nil == responseBody {
		return nil, err
	}

	unmarshalledResp := new(import1.AttachVmApiResponse)
	json.Unmarshal(responseBody.([]byte), &unmarshalledResp)
	return unmarshalledResp, err
}

// Creates a new Volume Disk.
func (api *VolumeGroupsApi) CreateVolumeDisk(volumeGroupExtId *string, body *import1.VolumeDisk, args ...map[string]interface{}) (*import1.CreateVolumeDiskApiResponse, error) {
	argMap := make(map[string]interface{})
	if len(args) > 0 {
		argMap = args[0]
	}

	uri := "/api/volumes/v4.0.b1/config/volume-groups/{volumeGroupExtId}/disks"

	// verify the required parameter 'volumeGroupExtId' is set
	if nil == volumeGroupExtId {
		return nil, client.ReportError("volumeGroupExtId is required and must be specified")
	}
	// verify the required parameter 'body' is set
	if nil == body {
		return nil, client.ReportError("body is required and must be specified")
	}

	// Path Params

	uri = strings.Replace(uri, "{"+"volumeGroupExtId"+"}", url.PathEscape(client.ParameterToString(*volumeGroupExtId, "")), -1)
	headerParams := make(map[string]string)
	queryParams := url.Values{}
	formParams := url.Values{}

	// to determine the Content-Type header
	contentTypes := []string{"application/json"}

	// to determine the Accept header
	accepts := []string{"application/json"}

	// Headers provided explicitly on operation takes precedence
	for headerKey, value := range argMap {
		// Skip platform generated headers
		if !api.headersToSkip[strings.ToLower(headerKey)] {
			if value != nil {
				if headerValue, headerValueOk := value.(string); headerValueOk {
					headerParams[headerKey] = headerValue
				}
			}
		}
	}

	authNames := []string{"basicAuthScheme"}

	responseBody, err := api.ApiClient.CallApi(&uri, http.MethodPost, body, queryParams, headerParams, formParams, accepts, contentTypes, authNames)
	if nil != err || nil == responseBody {
		return nil, err
	}

	unmarshalledResp := new(import1.CreateVolumeDiskApiResponse)
	json.Unmarshal(responseBody.([]byte), &unmarshalledResp)
	return unmarshalledResp, err
}

// Creates a new Volume Group.
func (api *VolumeGroupsApi) CreateVolumeGroup(body *import1.VolumeGroup, args ...map[string]interface{}) (*import1.CreateVolumeGroupApiResponse, error) {
	argMap := make(map[string]interface{})
	if len(args) > 0 {
		argMap = args[0]
	}

	uri := "/api/volumes/v4.0.b1/config/volume-groups"

	// verify the required parameter 'body' is set
	if nil == body {
		return nil, client.ReportError("body is required and must be specified")
	}

	headerParams := make(map[string]string)
	queryParams := url.Values{}
	formParams := url.Values{}

	// to determine the Content-Type header
	contentTypes := []string{"application/json"}

	// to determine the Accept header
	accepts := []string{"application/json"}

	// Headers provided explicitly on operation takes precedence
	for headerKey, value := range argMap {
		// Skip platform generated headers
		if !api.headersToSkip[strings.ToLower(headerKey)] {
			if value != nil {
				if headerValue, headerValueOk := value.(string); headerValueOk {
					headerParams[headerKey] = headerValue
				}
			}
		}
	}

	authNames := []string{"basicAuthScheme"}

	responseBody, err := api.ApiClient.CallApi(&uri, http.MethodPost, body, queryParams, headerParams, formParams, accepts, contentTypes, authNames)
	if nil != err || nil == responseBody {
		return nil, err
	}

	unmarshalledResp := new(import1.CreateVolumeGroupApiResponse)
	json.Unmarshal(responseBody.([]byte), &unmarshalledResp)
	return unmarshalledResp, err
}

// Deletes a Volume Disk identified by {extId} in the Volume Group identified by {volumeGroupExtId}.
func (api *VolumeGroupsApi) DeleteVolumeDiskById(volumeGroupExtId *string, extId *string, args ...map[string]interface{}) (*import1.DeleteVolumeDiskApiResponse, error) {
	argMap := make(map[string]interface{})
	if len(args) > 0 {
		argMap = args[0]
	}

	uri := "/api/volumes/v4.0.b1/config/volume-groups/{volumeGroupExtId}/disks/{extId}"

	// verify the required parameter 'volumeGroupExtId' is set
	if nil == volumeGroupExtId {
		return nil, client.ReportError("volumeGroupExtId is required and must be specified")
	}
	// verify the required parameter 'extId' is set
	if nil == extId {
		return nil, client.ReportError("extId is required and must be specified")
	}

	// Path Params

	uri = strings.Replace(uri, "{"+"volumeGroupExtId"+"}", url.PathEscape(client.ParameterToString(*volumeGroupExtId, "")), -1)

	uri = strings.Replace(uri, "{"+"extId"+"}", url.PathEscape(client.ParameterToString(*extId, "")), -1)
	headerParams := make(map[string]string)
	queryParams := url.Values{}
	formParams := url.Values{}

	// to determine the Content-Type header
	contentTypes := []string{}

	// to determine the Accept header
	accepts := []string{"application/json"}

	// Headers provided explicitly on operation takes precedence
	for headerKey, value := range argMap {
		// Skip platform generated headers
		if !api.headersToSkip[strings.ToLower(headerKey)] {
			if value != nil {
				if headerValue, headerValueOk := value.(string); headerValueOk {
					headerParams[headerKey] = headerValue
				}
			}
		}
	}

	authNames := []string{"basicAuthScheme"}

	responseBody, err := api.ApiClient.CallApi(&uri, http.MethodDelete, nil, queryParams, headerParams, formParams, accepts, contentTypes, authNames)
	if nil != err || nil == responseBody {
		return nil, err
	}

	unmarshalledResp := new(import1.DeleteVolumeDiskApiResponse)
	json.Unmarshal(responseBody.([]byte), &unmarshalledResp)
	return unmarshalledResp, err
}

// Delete the Volume Group identified by {extId}.
func (api *VolumeGroupsApi) DeleteVolumeGroupById(extId *string, args ...map[string]interface{}) (*import1.DeleteVolumeGroupApiResponse, error) {
	argMap := make(map[string]interface{})
	if len(args) > 0 {
		argMap = args[0]
	}

	uri := "/api/volumes/v4.0.b1/config/volume-groups/{extId}"

	// verify the required parameter 'extId' is set
	if nil == extId {
		return nil, client.ReportError("extId is required and must be specified")
	}

	// Path Params

	uri = strings.Replace(uri, "{"+"extId"+"}", url.PathEscape(client.ParameterToString(*extId, "")), -1)
	headerParams := make(map[string]string)
	queryParams := url.Values{}
	formParams := url.Values{}

	// to determine the Content-Type header
	contentTypes := []string{}

	// to determine the Accept header
	accepts := []string{"application/json"}

	// Headers provided explicitly on operation takes precedence
	for headerKey, value := range argMap {
		// Skip platform generated headers
		if !api.headersToSkip[strings.ToLower(headerKey)] {
			if value != nil {
				if headerValue, headerValueOk := value.(string); headerValueOk {
					headerParams[headerKey] = headerValue
				}
			}
		}
	}

	authNames := []string{"basicAuthScheme"}

	responseBody, err := api.ApiClient.CallApi(&uri, http.MethodDelete, nil, queryParams, headerParams, formParams, accepts, contentTypes, authNames)
	if nil != err || nil == responseBody {
		return nil, err
	}

	unmarshalledResp := new(import1.DeleteVolumeGroupApiResponse)
	json.Unmarshal(responseBody.([]byte), &unmarshalledResp)
	return unmarshalledResp, err
}

// Detaches iSCSI initiator identified by {extId} from a Volume Group identified by {volumeGroupExtId}.
func (api *VolumeGroupsApi) DetachIscsiClient(extId *string, body *import1.IscsiClientAttachment, args ...map[string]interface{}) (*import1.DetachIscsiClientApiResponse, error) {
	argMap := make(map[string]interface{})
	if len(args) > 0 {
		argMap = args[0]
	}

	uri := "/api/volumes/v4.0.b1/config/volume-groups/{extId}/$actions/detach-iscsi-client"

	// verify the required parameter 'extId' is set
	if nil == extId {
		return nil, client.ReportError("extId is required and must be specified")
	}
	// verify the required parameter 'body' is set
	if nil == body {
		return nil, client.ReportError("body is required and must be specified")
	}

	// Path Params

	uri = strings.Replace(uri, "{"+"extId"+"}", url.PathEscape(client.ParameterToString(*extId, "")), -1)
	headerParams := make(map[string]string)
	queryParams := url.Values{}
	formParams := url.Values{}

	// to determine the Content-Type header
	contentTypes := []string{"application/json"}

	// to determine the Accept header
	accepts := []string{"application/json"}

	// Headers provided explicitly on operation takes precedence
	for headerKey, value := range argMap {
		// Skip platform generated headers
		if !api.headersToSkip[strings.ToLower(headerKey)] {
			if value != nil {
				if headerValue, headerValueOk := value.(string); headerValueOk {
					headerParams[headerKey] = headerValue
				}
			}
		}
	}

	authNames := []string{"basicAuthScheme"}

	responseBody, err := api.ApiClient.CallApi(&uri, http.MethodPost, body, queryParams, headerParams, formParams, accepts, contentTypes, authNames)
	if nil != err || nil == responseBody {
		return nil, err
	}

	unmarshalledResp := new(import1.DetachIscsiClientApiResponse)
	json.Unmarshal(responseBody.([]byte), &unmarshalledResp)
	return unmarshalledResp, err
}

// Detaches VM identified by {extId} from a Volume Group identified by {volumeGroupExtId}.
func (api *VolumeGroupsApi) DetachVm(extId *string, body *import1.VmAttachment, args ...map[string]interface{}) (*import1.DetachVmApiResponse, error) {
	argMap := make(map[string]interface{})
	if len(args) > 0 {
		argMap = args[0]
	}

	uri := "/api/volumes/v4.0.b1/config/volume-groups/{extId}/$actions/detach-vm"

	// verify the required parameter 'extId' is set
	if nil == extId {
		return nil, client.ReportError("extId is required and must be specified")
	}
	// verify the required parameter 'body' is set
	if nil == body {
		return nil, client.ReportError("body is required and must be specified")
	}

	// Path Params

	uri = strings.Replace(uri, "{"+"extId"+"}", url.PathEscape(client.ParameterToString(*extId, "")), -1)
	headerParams := make(map[string]string)
	queryParams := url.Values{}
	formParams := url.Values{}

	// to determine the Content-Type header
	contentTypes := []string{"application/json"}

	// to determine the Accept header
	accepts := []string{"application/json"}

	// Headers provided explicitly on operation takes precedence
	for headerKey, value := range argMap {
		// Skip platform generated headers
		if !api.headersToSkip[strings.ToLower(headerKey)] {
			if value != nil {
				if headerValue, headerValueOk := value.(string); headerValueOk {
					headerParams[headerKey] = headerValue
				}
			}
		}
	}

	authNames := []string{"basicAuthScheme"}

	responseBody, err := api.ApiClient.CallApi(&uri, http.MethodPost, body, queryParams, headerParams, formParams, accepts, contentTypes, authNames)
	if nil != err || nil == responseBody {
		return nil, err
	}

	unmarshalledResp := new(import1.DetachVmApiResponse)
	json.Unmarshal(responseBody.([]byte), &unmarshalledResp)
	return unmarshalledResp, err
}

// Query the Volume Disk identified by {extId} in the Volume Group identified by {volumeGroupExtId}.
func (api *VolumeGroupsApi) GetVolumeDiskById(volumeGroupExtId *string, extId *string, args ...map[string]interface{}) (*import1.GetVolumeDiskApiResponse, error) {
	argMap := make(map[string]interface{})
	if len(args) > 0 {
		argMap = args[0]
	}

	uri := "/api/volumes/v4.0.b1/config/volume-groups/{volumeGroupExtId}/disks/{extId}"

	// verify the required parameter 'volumeGroupExtId' is set
	if nil == volumeGroupExtId {
		return nil, client.ReportError("volumeGroupExtId is required and must be specified")
	}
	// verify the required parameter 'extId' is set
	if nil == extId {
		return nil, client.ReportError("extId is required and must be specified")
	}

	// Path Params

	uri = strings.Replace(uri, "{"+"volumeGroupExtId"+"}", url.PathEscape(client.ParameterToString(*volumeGroupExtId, "")), -1)

	uri = strings.Replace(uri, "{"+"extId"+"}", url.PathEscape(client.ParameterToString(*extId, "")), -1)
	headerParams := make(map[string]string)
	queryParams := url.Values{}
	formParams := url.Values{}

	// to determine the Content-Type header
	contentTypes := []string{}

	// to determine the Accept header
	accepts := []string{"application/json"}

	// Headers provided explicitly on operation takes precedence
	for headerKey, value := range argMap {
		// Skip platform generated headers
		if !api.headersToSkip[strings.ToLower(headerKey)] {
			if value != nil {
				if headerValue, headerValueOk := value.(string); headerValueOk {
					headerParams[headerKey] = headerValue
				}
			}
		}
	}

	authNames := []string{"basicAuthScheme"}

	responseBody, err := api.ApiClient.CallApi(&uri, http.MethodGet, nil, queryParams, headerParams, formParams, accepts, contentTypes, authNames)
	if nil != err || nil == responseBody {
		return nil, err
	}

	unmarshalledResp := new(import1.GetVolumeDiskApiResponse)
	json.Unmarshal(responseBody.([]byte), &unmarshalledResp)
	return unmarshalledResp, err
}

// Query the Volume Disk stats identified by {diskExtId}.
func (api *VolumeGroupsApi) GetVolumeDiskStats(volumeGroupExtId *string, extId *string, startTime_ *string, endTime_ *string, samplingInterval_ *int, statType_ *import2.DownSamplingOperator, args ...map[string]interface{}) (*import3.GetVolumeDiskStatsApiResponse, error) {
	argMap := make(map[string]interface{})
	if len(args) > 0 {
		argMap = args[0]
	}

	uri := "/api/volumes/v4.0.b1/stats/volume-groups/{volumeGroupExtId}/disks/{extId}"

	// verify the required parameter 'volumeGroupExtId' is set
	if nil == volumeGroupExtId {
		return nil, client.ReportError("volumeGroupExtId is required and must be specified")
	}
	// verify the required parameter 'extId' is set
	if nil == extId {
		return nil, client.ReportError("extId is required and must be specified")
	}
	// verify the required parameter 'startTime_' is set
	if nil == startTime_ {
		return nil, client.ReportError("startTime_ is required and must be specified")
	}
	// verify the required parameter 'endTime_' is set
	if nil == endTime_ {
		return nil, client.ReportError("endTime_ is required and must be specified")
	}

	// Path Params

	uri = strings.Replace(uri, "{"+"volumeGroupExtId"+"}", url.PathEscape(client.ParameterToString(*volumeGroupExtId, "")), -1)

	uri = strings.Replace(uri, "{"+"extId"+"}", url.PathEscape(client.ParameterToString(*extId, "")), -1)
	headerParams := make(map[string]string)
	queryParams := url.Values{}
	formParams := url.Values{}

	// to determine the Content-Type header
	contentTypes := []string{}

	// to determine the Accept header
	accepts := []string{"application/json"}

	// Query Params

	queryParams.Add("$startTime", client.ParameterToString(*startTime_, ""))

	queryParams.Add("$endTime", client.ParameterToString(*endTime_, ""))
	if samplingInterval_ != nil {

		queryParams.Add("$samplingInterval", client.ParameterToString(*samplingInterval_, ""))
	}
	if statType_ != nil {
		statType_QueryParamEnumVal := statType_.GetName()
		queryParams.Add("$statType", client.ParameterToString(statType_QueryParamEnumVal, ""))
	}
	// Headers provided explicitly on operation takes precedence
	for headerKey, value := range argMap {
		// Skip platform generated headers
		if !api.headersToSkip[strings.ToLower(headerKey)] {
			if value != nil {
				if headerValue, headerValueOk := value.(string); headerValueOk {
					headerParams[headerKey] = headerValue
				}
			}
		}
	}

	authNames := []string{"basicAuthScheme"}

	responseBody, err := api.ApiClient.CallApi(&uri, http.MethodGet, nil, queryParams, headerParams, formParams, accepts, contentTypes, authNames)
	if nil != err || nil == responseBody {
		return nil, err
	}

	unmarshalledResp := new(import3.GetVolumeDiskStatsApiResponse)
	json.Unmarshal(responseBody.([]byte), &unmarshalledResp)
	return unmarshalledResp, err
}

// Query the Volume Group identified by {extId}.
func (api *VolumeGroupsApi) GetVolumeGroupById(extId *string, args ...map[string]interface{}) (*import1.GetVolumeGroupApiResponse, error) {
	argMap := make(map[string]interface{})
	if len(args) > 0 {
		argMap = args[0]
	}

	uri := "/api/volumes/v4.0.b1/config/volume-groups/{extId}"

	// verify the required parameter 'extId' is set
	if nil == extId {
		return nil, client.ReportError("extId is required and must be specified")
	}

	// Path Params

	uri = strings.Replace(uri, "{"+"extId"+"}", url.PathEscape(client.ParameterToString(*extId, "")), -1)
	headerParams := make(map[string]string)
	queryParams := url.Values{}
	formParams := url.Values{}

	// to determine the Content-Type header
	contentTypes := []string{}

	// to determine the Accept header
	accepts := []string{"application/json"}

	// Headers provided explicitly on operation takes precedence
	for headerKey, value := range argMap {
		// Skip platform generated headers
		if !api.headersToSkip[strings.ToLower(headerKey)] {
			if value != nil {
				if headerValue, headerValueOk := value.(string); headerValueOk {
					headerParams[headerKey] = headerValue
				}
			}
		}
	}

	authNames := []string{"basicAuthScheme"}

	responseBody, err := api.ApiClient.CallApi(&uri, http.MethodGet, nil, queryParams, headerParams, formParams, accepts, contentTypes, authNames)
	if nil != err || nil == responseBody {
		return nil, err
	}

	unmarshalledResp := new(import1.GetVolumeGroupApiResponse)
	json.Unmarshal(responseBody.([]byte), &unmarshalledResp)
	return unmarshalledResp, err
}

// Query for metadata information which is associated with the Volume Group identified by {extId}.
func (api *VolumeGroupsApi) GetVolumeGroupMetadataById(extId *string, args ...map[string]interface{}) (*import1.GetVolumeGroupMetadataApiResponse, error) {
	argMap := make(map[string]interface{})
	if len(args) > 0 {
		argMap = args[0]
	}

	uri := "/api/volumes/v4.0.b1/config/volume-groups/{extId}/metadata"

	// verify the required parameter 'extId' is set
	if nil == extId {
		return nil, client.ReportError("extId is required and must be specified")
	}

	// Path Params

	uri = strings.Replace(uri, "{"+"extId"+"}", url.PathEscape(client.ParameterToString(*extId, "")), -1)
	headerParams := make(map[string]string)
	queryParams := url.Values{}
	formParams := url.Values{}

	// to determine the Content-Type header
	contentTypes := []string{}

	// to determine the Accept header
	accepts := []string{"application/json"}

	// Headers provided explicitly on operation takes precedence
	for headerKey, value := range argMap {
		// Skip platform generated headers
		if !api.headersToSkip[strings.ToLower(headerKey)] {
			if value != nil {
				if headerValue, headerValueOk := value.(string); headerValueOk {
					headerParams[headerKey] = headerValue
				}
			}
		}
	}

	authNames := []string{"basicAuthScheme"}

	responseBody, err := api.ApiClient.CallApi(&uri, http.MethodGet, nil, queryParams, headerParams, formParams, accepts, contentTypes, authNames)
	if nil != err || nil == responseBody {
		return nil, err
	}

	unmarshalledResp := new(import1.GetVolumeGroupMetadataApiResponse)
	json.Unmarshal(responseBody.([]byte), &unmarshalledResp)
	return unmarshalledResp, err
}

// Query the Volume Group stats identified by {extId}.
func (api *VolumeGroupsApi) GetVolumeGroupStats(extId *string, startTime_ *string, endTime_ *string, samplingInterval_ *int, statType_ *import2.DownSamplingOperator, args ...map[string]interface{}) (*import3.GetVolumeGroupStatsApiResponse, error) {
	argMap := make(map[string]interface{})
	if len(args) > 0 {
		argMap = args[0]
	}

	uri := "/api/volumes/v4.0.b1/stats/volume-groups/{extId}"

	// verify the required parameter 'extId' is set
	if nil == extId {
		return nil, client.ReportError("extId is required and must be specified")
	}
	// verify the required parameter 'startTime_' is set
	if nil == startTime_ {
		return nil, client.ReportError("startTime_ is required and must be specified")
	}
	// verify the required parameter 'endTime_' is set
	if nil == endTime_ {
		return nil, client.ReportError("endTime_ is required and must be specified")
	}

	// Path Params

	uri = strings.Replace(uri, "{"+"extId"+"}", url.PathEscape(client.ParameterToString(*extId, "")), -1)
	headerParams := make(map[string]string)
	queryParams := url.Values{}
	formParams := url.Values{}

	// to determine the Content-Type header
	contentTypes := []string{}

	// to determine the Accept header
	accepts := []string{"application/json"}

	// Query Params

	queryParams.Add("$startTime", client.ParameterToString(*startTime_, ""))

	queryParams.Add("$endTime", client.ParameterToString(*endTime_, ""))
	if samplingInterval_ != nil {

		queryParams.Add("$samplingInterval", client.ParameterToString(*samplingInterval_, ""))
	}
	if statType_ != nil {
		statType_QueryParamEnumVal := statType_.GetName()
		queryParams.Add("$statType", client.ParameterToString(statType_QueryParamEnumVal, ""))
	}
	// Headers provided explicitly on operation takes precedence
	for headerKey, value := range argMap {
		// Skip platform generated headers
		if !api.headersToSkip[strings.ToLower(headerKey)] {
			if value != nil {
				if headerValue, headerValueOk := value.(string); headerValueOk {
					headerParams[headerKey] = headerValue
				}
			}
		}
	}

	authNames := []string{"basicAuthScheme"}

	responseBody, err := api.ApiClient.CallApi(&uri, http.MethodGet, nil, queryParams, headerParams, formParams, accepts, contentTypes, authNames)
	if nil != err || nil == responseBody {
		return nil, err
	}

	unmarshalledResp := new(import3.GetVolumeGroupStatsApiResponse)
	json.Unmarshal(responseBody.([]byte), &unmarshalledResp)
	return unmarshalledResp, err
}

// Query the category details that are associated with the Volume Group identified by {volumeGroupExtId}.
func (api *VolumeGroupsApi) ListCategoryAssociationsByVolumeGroupId(extId *string, page_ *int, limit_ *int, args ...map[string]interface{}) (*import1.ListCategoryAssociationsApiResponse, error) {
	argMap := make(map[string]interface{})
	if len(args) > 0 {
		argMap = args[0]
	}

	uri := "/api/volumes/v4.0.b1/config/volume-groups/{extId}/category-associations"

	// verify the required parameter 'extId' is set
	if nil == extId {
		return nil, client.ReportError("extId is required and must be specified")
	}

	// Path Params

	uri = strings.Replace(uri, "{"+"extId"+"}", url.PathEscape(client.ParameterToString(*extId, "")), -1)
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
	// Headers provided explicitly on operation takes precedence
	for headerKey, value := range argMap {
		// Skip platform generated headers
		if !api.headersToSkip[strings.ToLower(headerKey)] {
			if value != nil {
				if headerValue, headerValueOk := value.(string); headerValueOk {
					headerParams[headerKey] = headerValue
				}
			}
		}
	}

	authNames := []string{"basicAuthScheme"}

	responseBody, err := api.ApiClient.CallApi(&uri, http.MethodGet, nil, queryParams, headerParams, formParams, accepts, contentTypes, authNames)
	if nil != err || nil == responseBody {
		return nil, err
	}

	unmarshalledResp := new(import1.ListCategoryAssociationsApiResponse)
	json.Unmarshal(responseBody.([]byte), &unmarshalledResp)
	return unmarshalledResp, err
}

// Query the list of external iSCSI attachments for a Volume Group identified by {extId}.
func (api *VolumeGroupsApi) ListExternalIscsiAttachmentsByVolumeGroupId(extId *string, page_ *int, limit_ *int, filter_ *string, orderby_ *string, expand_ *string, select_ *string, args ...map[string]interface{}) (*import1.ListExternalIscsiAttachmentsApiResponse, error) {
	argMap := make(map[string]interface{})
	if len(args) > 0 {
		argMap = args[0]
	}

	uri := "/api/volumes/v4.0.b1/config/volume-groups/{extId}/external-iscsi-attachments"

	// verify the required parameter 'extId' is set
	if nil == extId {
		return nil, client.ReportError("extId is required and must be specified")
	}

	// Path Params

	uri = strings.Replace(uri, "{"+"extId"+"}", url.PathEscape(client.ParameterToString(*extId, "")), -1)
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
	if expand_ != nil {

		queryParams.Add("$expand", client.ParameterToString(*expand_, ""))
	}
	if select_ != nil {

		queryParams.Add("$select", client.ParameterToString(*select_, ""))
	}
	// Headers provided explicitly on operation takes precedence
	for headerKey, value := range argMap {
		// Skip platform generated headers
		if !api.headersToSkip[strings.ToLower(headerKey)] {
			if value != nil {
				if headerValue, headerValueOk := value.(string); headerValueOk {
					headerParams[headerKey] = headerValue
				}
			}
		}
	}

	authNames := []string{"basicAuthScheme"}

	responseBody, err := api.ApiClient.CallApi(&uri, http.MethodGet, nil, queryParams, headerParams, formParams, accepts, contentTypes, authNames)
	if nil != err || nil == responseBody {
		return nil, err
	}

	unmarshalledResp := new(import1.ListExternalIscsiAttachmentsApiResponse)
	json.Unmarshal(responseBody.([]byte), &unmarshalledResp)
	return unmarshalledResp, err
}

// Query the list of VM attachments for a Volume Group identified by {extId}.
func (api *VolumeGroupsApi) ListVmAttachmentsByVolumeGroupId(extId *string, page_ *int, limit_ *int, filter_ *string, orderby_ *string, args ...map[string]interface{}) (*import1.ListVmAttachmentsApiResponse, error) {
	argMap := make(map[string]interface{})
	if len(args) > 0 {
		argMap = args[0]
	}

	uri := "/api/volumes/v4.0.b1/config/volume-groups/{extId}/vm-attachments"

	// verify the required parameter 'extId' is set
	if nil == extId {
		return nil, client.ReportError("extId is required and must be specified")
	}

	// Path Params

	uri = strings.Replace(uri, "{"+"extId"+"}", url.PathEscape(client.ParameterToString(*extId, "")), -1)
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
	// Headers provided explicitly on operation takes precedence
	for headerKey, value := range argMap {
		// Skip platform generated headers
		if !api.headersToSkip[strings.ToLower(headerKey)] {
			if value != nil {
				if headerValue, headerValueOk := value.(string); headerValueOk {
					headerParams[headerKey] = headerValue
				}
			}
		}
	}

	authNames := []string{"basicAuthScheme"}

	responseBody, err := api.ApiClient.CallApi(&uri, http.MethodGet, nil, queryParams, headerParams, formParams, accepts, contentTypes, authNames)
	if nil != err || nil == responseBody {
		return nil, err
	}

	unmarshalledResp := new(import1.ListVmAttachmentsApiResponse)
	json.Unmarshal(responseBody.([]byte), &unmarshalledResp)
	return unmarshalledResp, err
}

// Query the list of disks corresponding to a Volume Group identified by {volumeGroupExtId}.
func (api *VolumeGroupsApi) ListVolumeDisksByVolumeGroupId(volumeGroupExtId *string, page_ *int, limit_ *int, filter_ *string, orderby_ *string, select_ *string, args ...map[string]interface{}) (*import1.ListVolumeDisksApiResponse, error) {
	argMap := make(map[string]interface{})
	if len(args) > 0 {
		argMap = args[0]
	}

	uri := "/api/volumes/v4.0.b1/config/volume-groups/{volumeGroupExtId}/disks"

	// verify the required parameter 'volumeGroupExtId' is set
	if nil == volumeGroupExtId {
		return nil, client.ReportError("volumeGroupExtId is required and must be specified")
	}

	// Path Params

	uri = strings.Replace(uri, "{"+"volumeGroupExtId"+"}", url.PathEscape(client.ParameterToString(*volumeGroupExtId, "")), -1)
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
	// Headers provided explicitly on operation takes precedence
	for headerKey, value := range argMap {
		// Skip platform generated headers
		if !api.headersToSkip[strings.ToLower(headerKey)] {
			if value != nil {
				if headerValue, headerValueOk := value.(string); headerValueOk {
					headerParams[headerKey] = headerValue
				}
			}
		}
	}

	authNames := []string{"basicAuthScheme"}

	responseBody, err := api.ApiClient.CallApi(&uri, http.MethodGet, nil, queryParams, headerParams, formParams, accepts, contentTypes, authNames)
	if nil != err || nil == responseBody {
		return nil, err
	}

	unmarshalledResp := new(import1.ListVolumeDisksApiResponse)
	json.Unmarshal(responseBody.([]byte), &unmarshalledResp)
	return unmarshalledResp, err
}

// Query the list of Volume Groups.
func (api *VolumeGroupsApi) ListVolumeGroups(page_ *int, limit_ *int, filter_ *string, orderby_ *string, expand_ *string, select_ *string, args ...map[string]interface{}) (*import1.ListVolumeGroupsApiResponse, error) {
	argMap := make(map[string]interface{})
	if len(args) > 0 {
		argMap = args[0]
	}

	uri := "/api/volumes/v4.0.b1/config/volume-groups"

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
	if expand_ != nil {

		queryParams.Add("$expand", client.ParameterToString(*expand_, ""))
	}
	if select_ != nil {

		queryParams.Add("$select", client.ParameterToString(*select_, ""))
	}
	// Headers provided explicitly on operation takes precedence
	for headerKey, value := range argMap {
		// Skip platform generated headers
		if !api.headersToSkip[strings.ToLower(headerKey)] {
			if value != nil {
				if headerValue, headerValueOk := value.(string); headerValueOk {
					headerParams[headerKey] = headerValue
				}
			}
		}
	}

	authNames := []string{"basicAuthScheme"}

	responseBody, err := api.ApiClient.CallApi(&uri, http.MethodGet, nil, queryParams, headerParams, formParams, accepts, contentTypes, authNames)
	if nil != err || nil == responseBody {
		return nil, err
	}

	unmarshalledResp := new(import1.ListVolumeGroupsApiResponse)
	json.Unmarshal(responseBody.([]byte), &unmarshalledResp)
	return unmarshalledResp, err
}

// Reverts a Volume Group identified by {volumeGroupExtId}. This API will perform an in-place restore from a specified Volume Group recovery point.
func (api *VolumeGroupsApi) RevertVolumeGroup(extId *string, body *import1.RevertSpec, args ...map[string]interface{}) (*import1.RevertVolumeGroupApiResponse, error) {
	argMap := make(map[string]interface{})
	if len(args) > 0 {
		argMap = args[0]
	}

	uri := "/api/volumes/v4.0.b1/config/volume-groups/{extId}/$actions/revert"

	// verify the required parameter 'extId' is set
	if nil == extId {
		return nil, client.ReportError("extId is required and must be specified")
	}
	// verify the required parameter 'body' is set
	if nil == body {
		return nil, client.ReportError("body is required and must be specified")
	}

	// Path Params

	uri = strings.Replace(uri, "{"+"extId"+"}", url.PathEscape(client.ParameterToString(*extId, "")), -1)
	headerParams := make(map[string]string)
	queryParams := url.Values{}
	formParams := url.Values{}

	// to determine the Content-Type header
	contentTypes := []string{"application/json"}

	// to determine the Accept header
	accepts := []string{"application/json"}

	// Headers provided explicitly on operation takes precedence
	for headerKey, value := range argMap {
		// Skip platform generated headers
		if !api.headersToSkip[strings.ToLower(headerKey)] {
			if value != nil {
				if headerValue, headerValueOk := value.(string); headerValueOk {
					headerParams[headerKey] = headerValue
				}
			}
		}
	}

	authNames := []string{"basicAuthScheme"}

	responseBody, err := api.ApiClient.CallApi(&uri, http.MethodPost, body, queryParams, headerParams, formParams, accepts, contentTypes, authNames)
	if nil != err || nil == responseBody {
		return nil, err
	}

	unmarshalledResp := new(import1.RevertVolumeGroupApiResponse)
	json.Unmarshal(responseBody.([]byte), &unmarshalledResp)
	return unmarshalledResp, err
}
