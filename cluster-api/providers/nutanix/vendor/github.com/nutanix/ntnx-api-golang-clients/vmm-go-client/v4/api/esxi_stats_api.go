//Api classes for vmm's golang SDK
package api

import (
	"encoding/json"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/client"
	import1 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/common/v1/stats"
	import2 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/esxi/stats"
	"net/http"
	"net/url"
	"strings"
)

type EsxiStatsApi struct {
	ApiClient     *client.ApiClient
	headersToSkip map[string]bool
}

func NewEsxiStatsApi(apiClient *client.ApiClient) *EsxiStatsApi {
	if apiClient == nil {
		apiClient = client.NewApiClient()
	}

	a := &EsxiStatsApi{
		ApiClient: apiClient,
	}

	headers := []string{"authorization", "cookie", "host", "user-agent"}
	a.headersToSkip = make(map[string]bool)
	for _, header := range headers {
		a.headersToSkip[header] = true
	}

	return a
}

// Get VM stats for a given VM. Users can include the following in the request query: 1) '$select': comma-separated attributes with the prefix 'stats/', e.g. 'stats/controllerNumIo,stats/hypervisorNumIo'. 2) '$startTime': the start time of the period for which stats should be reported, where the value should be in extended ISO-8601 format, e.g. '2023-01-01T12:00:00.000-08:00'; 3) '$endTime': the end time of the period for which stats should be reported, where the value should be in extended ISO-8601 format; 4) '$samplingInterval': the sampling interval in seconds at which statistical data should be collected; 5) '$statType': the down-sampling operator to use while performing down-sampling on stats data, where the allowed values are SUM, MIN, MAX, AVG, COUNT and LAST.
func (api *EsxiStatsApi) GetVmStatsById(extId *string, startTime_ *string, endTime_ *string, samplingInterval_ *int, statType_ *import1.DownSamplingOperator, select_ *string, args ...map[string]interface{}) (*import2.GetVmStatsApiResponse, error) {
	argMap := make(map[string]interface{})
	if len(args) > 0 {
		argMap = args[0]
	}

	uri := "/api/vmm/v4.0.b1/esxi/stats/vms/{extId}"

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

	unmarshalledResp := new(import2.GetVmStatsApiResponse)
	json.Unmarshal(responseBody.([]byte), &unmarshalledResp)
	return unmarshalledResp, err
}

// List VM stats for all VMs. Users can include the following in the request query: 1) '$select': comma-separated attributes with the prefix 'stats/', e.g. 'stats/controllerNumIo,stats/hypervisorNumIo'. 2) '$startTime': the start time of the period for which stats should be reported, where the value should be in extended ISO-8601 format, e.g. '2023-01-01T12:00:00.000-08:00'; 3) '$endTime': the end time of the period for which stats should be reported, where the value should be in extended ISO-8601 format; 4) '$samplingInterval': the sampling interval in seconds at which statistical data should be collected; 5) '$statType': the down-sampling operator to use while performing down-sampling on stats data, where the allowed values are SUM, MIN, MAX, AVG, COUNT and LAST; 6) '$orderby'; 7) '$page'; 8) '$limit'; and 9) '$filter': the OData filter to use, e.g. 'stats/hypervisorCpuUsagePpm gt 100000 and stats/guestMemoryUsagePpm lt 2000000.'
func (api *EsxiStatsApi) ListVmStats(startTime_ *string, endTime_ *string, samplingInterval_ *int, statType_ *import1.DownSamplingOperator, page_ *int, limit_ *int, filter_ *string, orderby_ *string, select_ *string, args ...map[string]interface{}) (*import2.ListVmStatsApiResponse, error) {
	argMap := make(map[string]interface{})
	if len(args) > 0 {
		argMap = args[0]
	}

	uri := "/api/vmm/v4.0.b1/esxi/stats/vms"

	// verify the required parameter 'startTime_' is set
	if nil == startTime_ {
		return nil, client.ReportError("startTime_ is required and must be specified")
	}
	// verify the required parameter 'endTime_' is set
	if nil == endTime_ {
		return nil, client.ReportError("endTime_ is required and must be specified")
	}

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

	unmarshalledResp := new(import2.ListVmStatsApiResponse)
	json.Unmarshal(responseBody.([]byte), &unmarshalledResp)
	return unmarshalledResp, err
}
