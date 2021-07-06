package containerv2

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
)

const (
	resourceAccount = "X-Auth-Resource-Account"
)

//MonitoringTargetHeader ...
type MonitoringTargetHeader struct {
	AccountID string
}

//ToMap ...
func (c MonitoringTargetHeader) ToMap() map[string]string {
	m := make(map[string]string, 2)
	m[resourceAccount] = c.AccountID
	return m
}

//MonitoringCreateRequest ...
type MonitoringCreateRequest struct {
	Cluster         string `json:"cluster"`
	IngestionKey    string `json:"ingestionKey,omitempty"`
	SysidigInstance string `json:"instance"`
	PrivateEndpoint bool   `json:"privateEndpoint,omitempty"`
}

//MonitoringCreateResponse ...
type MonitoringCreateResponse struct {
	DaemonsetName   string `json:"daemonsetName"`
	IngestionKey    string `json:"ingestionKey"`
	InstanceID      string `json:"instanceId"`
	InstanceName    string `json:"instanceName"`
	PrivateEndpoint bool   `json:"privateEndpoint"`
}

//MonitoringUpdateRequest ...
type MonitoringUpdateRequest struct {
	Cluster         string `json:"cluster"`
	IngestionKey    string `json:"ingestionKey"`
	Instance        string `json:"instance"`
	NewInstance     string `json:"newInstance"`
	PrivateEndpoint bool   `json:"privateEndpoint"`
}

//MonitoringUpdateResponse ...
type MonitoringUpdateResponse struct {
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

//MonitoringInfo ...
type MonitoringInfo struct {
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

//MonitoringDeleteRequest ...
type MonitoringDeleteRequest struct {
	Cluster  string `json:"cluster"`
	Instance string `json:"instance"`
}

//Monitoring interface
type Monitoring interface {
	CreateMonitoringConfig(params MonitoringCreateRequest, target MonitoringTargetHeader) (MonitoringCreateResponse, error)
	GetMonitoringConfig(clusterName string, monitoringInstance string, target MonitoringTargetHeader) (*MonitoringInfo, error)
	ListAllMonitors(clusterName string, target MonitoringTargetHeader) ([]MonitoringInfo, error)
	UpdateMonitoringConfig(params MonitoringUpdateRequest, target MonitoringTargetHeader) (MonitoringUpdateResponse, error)
	DeleteMonitoringConfig(params MonitoringDeleteRequest, target MonitoringTargetHeader) (interface{}, error)
}
type monitoring struct {
	client *client.Client
}

func newMonitoringAPI(c *client.Client) Monitoring {
	return &monitoring{
		client: c,
	}
}

//CreateMonitoringConfig ...
//Create a Sysdig monitoring configuration for a cluster.
func (r *monitoring) CreateMonitoringConfig(params MonitoringCreateRequest, target MonitoringTargetHeader) (MonitoringCreateResponse, error) {
	var monitoring MonitoringCreateResponse
	_, err := r.client.Post("/v2/observe/monitoring/createConfig", params, &monitoring, target.ToMap())
	return monitoring, err
}

//GetMonitoringConfig ...
//Show the details of an existing Sysdig monitoring configuration.
func (r *monitoring) GetMonitoringConfig(clusterName, monitoringInstance string, target MonitoringTargetHeader) (*MonitoringInfo, error) {
	monitoringInfo := &MonitoringInfo{}
	rawURL := fmt.Sprintf("/v2/observe/monitoring/getConfig?cluster=%s&instance=%s", clusterName, monitoringInstance)
	_, err := r.client.Get(rawURL, &monitoringInfo, target.ToMap())
	if err != nil {
		return nil, err
	}
	return monitoringInfo, err
}

//ListAllMonitors ...
//List all Sysdig monitoring configurations for a cluster.
func (r *monitoring) ListAllMonitors(clusterName string, target MonitoringTargetHeader) ([]MonitoringInfo, error) {
	monitors := []MonitoringInfo{}
	rawURL := fmt.Sprintf("v2/observe/monitoring/getConfigs?cluster=%s", clusterName)
	_, err := r.client.Get(rawURL, &monitors, target.ToMap())
	if err != nil {
		return nil, err
	}
	return monitors, nil
}

//UpdateMonitoringConfig ...
//Update a Sysdig monitoring configuration in the cluster.
func (r *monitoring) UpdateMonitoringConfig(params MonitoringUpdateRequest, target MonitoringTargetHeader) (MonitoringUpdateResponse, error) {
	var monitoring MonitoringUpdateResponse
	_, err := r.client.Post("/v2/observe/monitoring/modifyConfig", params, &monitoring, target.ToMap())
	return monitoring, err
}

//DeleteMonitoringConfig ...
//Remove a Sysdig monitoring configuration from a cluster.
func (r *monitoring) DeleteMonitoringConfig(params MonitoringDeleteRequest, target MonitoringTargetHeader) (interface{}, error) {
	var response interface{}
	_, err := r.client.Post("/v2/observe/monitoring/removeConfig", params, &response, target.ToMap())
	if err != nil {
		return response, err
	}
	return response, err

}
