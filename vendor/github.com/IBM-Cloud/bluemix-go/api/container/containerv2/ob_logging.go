package containerv2

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
)

const (
	authResourceAccount = "X-Auth-Resource-Account"
)

//LoggingTargetHeader ...
type LoggingTargetHeader struct {
	AccountID string
}

//ToMap ...
func (c LoggingTargetHeader) ToMap() map[string]string {
	m := make(map[string]string, 2)
	m[authResourceAccount] = c.AccountID
	return m
}

//LoggingCreateRequest ...
type LoggingCreateRequest struct {
	Cluster         string `json:"cluster"`
	IngestionKey    string `json:"ingestionKey,omitempty"`
	LoggingInstance string `json:"instance"`
	PrivateEndpoint bool   `json:"privateEndpoint,omitempty"`
}

//LoggingCreateResponse ...
type LoggingCreateResponse struct {
	DaemonsetName   string `json:"daemonsetName"`
	IngestionKey    string `json:"ingestionKey"`
	InstanceID      string `json:"instanceId"`
	InstanceName    string `json:"instanceName"`
	PrivateEndpoint bool   `json:"privateEndpoint"`
}

//LoggingUpdateRequest ...
type LoggingUpdateRequest struct {
	Cluster         string `json:"cluster"`
	IngestionKey    string `json:"ingestionKey"`
	Instance        string `json:"instance"`
	NewInstance     string `json:"newInstance"`
	PrivateEndpoint bool   `json:"privateEndpoint"`
}

//LoggingUpdateResponse ...
type LoggingUpdateResponse struct {
	AgentKey        string `json:"agentKey"`
	AgentNamespace  string `json:"agentNamespace"`
	CRN             string `json:"crn"`
	DaemonsetName   string `json:"daemonsetName"`
	DiscoveredAgent bool   `json:"discoveredAgent"`
	InstanceID      string `json:"instanceId"`
	InstanceName    string `json:"instanceName"`
	Namespace       string `json:"namespace"`
	PrivateEndpoint bool   `json:"privateEndpoint"`
}

//LoggingInfo ...
type LoggingInfo struct {
	AgentKey        string `json:"agentKey"`
	AgentNamespace  string `json:"agentNamespace"`
	CRN             string `json:"crn"`
	DaemonsetName   string `json:"daemonsetName"`
	DiscoveredAgent bool   `json:"discoveredAgent"`
	InstanceID      string `json:"instanceId"`
	InstanceName    string `json:"instanceName"`
	Namespace       string `json:"namespace"`
	PrivateEndpoint bool   `json:"privateEndpoint"`
}

//LoggingDeleteRequest ...
type LoggingDeleteRequest struct {
	Cluster  string `json:"cluster"`
	Instance string `json:"instance"`
}

//Logging interface
type Logging interface {
	CreateLoggingConfig(params LoggingCreateRequest, target LoggingTargetHeader) (LoggingCreateResponse, error)
	GetLoggingConfig(clusterName string, LoggingInstance string, target LoggingTargetHeader) (*LoggingInfo, error)
	ListLoggingInstances(clusterName string, target LoggingTargetHeader) ([]LoggingInfo, error)
	UpdateLoggingConfig(params LoggingUpdateRequest, target LoggingTargetHeader) (LoggingUpdateResponse, error)
	DeleteLoggingConfig(params LoggingDeleteRequest, target LoggingTargetHeader) (interface{}, error)
}
type logging struct {
	client *client.Client
}

func newLoggingAPI(c *client.Client) Logging {
	return &logging{
		client: c,
	}
}

//CreateLoggingConfig ...
//Create a Logging configuration for a cluster.
func (r *logging) CreateLoggingConfig(params LoggingCreateRequest, target LoggingTargetHeader) (LoggingCreateResponse, error) {
	var resp LoggingCreateResponse
	_, err := r.client.Post("/v2/observe/logging/createConfig", params, &resp, target.ToMap())
	return resp, err
}

//GetLoggingConfig ...
//Show the details of an existing Logging configuration.
func (r *logging) GetLoggingConfig(clusterName, loggingInstance string, target LoggingTargetHeader) (*LoggingInfo, error) {
	loggingInfo := &LoggingInfo{}
	rawURL := fmt.Sprintf("/v2/observe/logging/getConfig?cluster=%s&instance=%s", clusterName, loggingInstance)
	_, err := r.client.Get(rawURL, &loggingInfo, target.ToMap())
	if err != nil {
		return nil, err
	}
	return loggingInfo, err
}

//ListLoggingInstances...
//List all logging configurations for a cluster.
func (r *logging) ListLoggingInstances(clusterName string, target LoggingTargetHeader) ([]LoggingInfo, error) {
	logging := []LoggingInfo{}
	rawURL := fmt.Sprintf("/v2/observe/logging/getConfigs?cluster=%s", clusterName)
	_, err := r.client.Get(rawURL, &logging, target.ToMap())
	if err != nil {
		return nil, err
	}
	return logging, nil
}

//UpdateLoggingConfig ...
//Update a Logging configuration in the cluster.
func (r *logging) UpdateLoggingConfig(params LoggingUpdateRequest, target LoggingTargetHeader) (LoggingUpdateResponse, error) {
	var logging LoggingUpdateResponse
	_, err := r.client.Post("/v2/observe/logging/modifyConfig", params, &logging, target.ToMap())
	return logging, err
}

//DeleteLoggingConfig ...
//Remove a Logging configuration from a cluster.
func (r *logging) DeleteLoggingConfig(params LoggingDeleteRequest, target LoggingTargetHeader) (interface{}, error) {
	var response interface{}
	_, err := r.client.Post("/v2/observe/logging/removeConfig", params, &response, target.ToMap())
	if err != nil {
		return response, err
	}
	return response, err
}
