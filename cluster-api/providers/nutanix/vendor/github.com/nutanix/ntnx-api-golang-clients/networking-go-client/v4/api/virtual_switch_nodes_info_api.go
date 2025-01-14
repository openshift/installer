//Api classes for networking's golang SDK
package api

import (
	"encoding/json"
	"github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/client"
	import1 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
	"net/http"
	"net/url"
	"strings"
)

type VirtualSwitchNodesInfoApi struct {
	ApiClient     *client.ApiClient
	headersToSkip map[string]bool
}

func NewVirtualSwitchNodesInfoApi(apiClient *client.ApiClient) *VirtualSwitchNodesInfoApi {
	if apiClient == nil {
		apiClient = client.NewApiClient()
	}

	a := &VirtualSwitchNodesInfoApi{
		ApiClient: apiClient,
	}

	headers := []string{"authorization", "cookie", "host", "user-agent"}
	a.headersToSkip = make(map[string]bool)
	for _, header := range headers {
		a.headersToSkip[header] = true
	}

	return a
}

// Check to see whether a node in a cluster is a storage-only node or not
func (api *VirtualSwitchNodesInfoApi) ListNodeSchedulableStatus(xClusterId *string, args ...map[string]interface{}) (*import1.ListNodeSchedulableStatusesApiResponse, error) {
	argMap := make(map[string]interface{})
	if len(args) > 0 {
		argMap = args[0]
	}

	uri := "/api/networking/v4.0.b1/config/node-schedulable-status"

	headerParams := make(map[string]string)
	queryParams := url.Values{}
	formParams := url.Values{}

	// to determine the Content-Type header
	contentTypes := []string{}

	// to determine the Accept header
	accepts := []string{"application/json"}

	if xClusterId != nil {
		headerParams["X-Cluster-Id"] = client.ParameterToString(*xClusterId, "")
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

	unmarshalledResp := new(import1.ListNodeSchedulableStatusesApiResponse)
	json.Unmarshal(responseBody.([]byte), &unmarshalledResp)
	return unmarshalledResp, err
}
