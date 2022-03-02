package containerv2

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
)

type SatelliteLocationInfo struct {
	ID                string      `json:"id"`
	Name              string      `json:"name"`
	Region            string      `json:"region"`
	ResourceGroup     string      `json:"resourceGroup"`
	ResourceGroupName string      `json:"resourceGroupName"`
	PodSubnet         string      `json:"podSubnet"`
	ServiceSubnet     string      `json:"serviceSubnet"`
	CreatedDate       string      `json:"createdDate"`
	MasterKubeVersion string      `json:"masterKubeVersion"`
	TargetVersion     string      `json:"targetVersion"`
	WorkerCount       int         `json:"workerCount"`
	Location          string      `json:"location"`
	Datacenter        string      `json:"datacenter"`
	MultiAzCapable    bool        `json:"multiAzCapable"`
	Provider          string      `json:"provider"`
	State             string      `json:"state"`
	Status            string      `json:"status"`
	VersionEOS        string      `json:"versionEOS"`
	IsPaid            bool        `json:"isPaid"`
	Entitlement       string      `json:"entitlement"`
	Type              string      `json:"type"`
	Addons            interface{} `json:"addons"`
	EtcdPort          string      `json:"etcdPort"`
	MasterURL         string      `json:"masterURL"`
	Ingress           struct {
		Hostname   string `json:"hostname"`
		SecretName string `json:"secretName"`
		Status     string `json:"status"`
		Message    string `json:"message"`
	} `json:"ingress"`
	CaCertRotationStatus struct {
		Status              string `json:"status"`
		ActionTriggerDate   string `json:"actionTriggerDate"`
		ActionCompletedDate string `json:"actionCompletedDate"`
	} `json:"caCertRotationStatus"`
	ImageSecurityEnabled bool     `json:"imageSecurityEnabled"`
	DisableAutoUpdate    bool     `json:"disableAutoUpdate"`
	Crn                  string   `json:"crn"`
	WorkerZones          []string `json:"workerZones"`
	Lifecycle            struct {
		MasterStatus             string `json:"masterStatus"`
		MasterStatusModifiedDate string `json:"masterStatusModifiedDate"`
		MasterHealth             string `json:"masterHealth"`
		MasterState              string `json:"masterState"`
		ModifiedDate             string `json:"modifiedDate"`
	} `json:"lifecycle"`
	ServiceEndpoints struct {
		PrivateServiceEndpointEnabled bool   `json:"privateServiceEndpointEnabled"`
		PrivateServiceEndpointURL     string `json:"privateServiceEndpointURL"`
		PublicServiceEndpointEnabled  bool   `json:"publicServiceEndpointEnabled"`
		PublicServiceEndpointURL      string `json:"publicServiceEndpointURL"`
	} `json:"serviceEndpoints"`
	Features struct {
		KeyProtectEnabled bool `json:"keyProtectEnabled"`
		PullSecretApplied bool `json:"pullSecretApplied"`
	} `json:"features"`
	Vpcs      interface{} `json:"vpcs"`
	CosConfig struct {
		Region          string `json:"region"`
		Bucket          string `json:"bucket"`
		Endpoint        string `json:"endpoint"`
		ServiceInstance struct {
			Crn string `json:"crn"`
		} `json:"serviceInstance"`
	} `json:"cos_config"`
	Description string `json:"description"`
	Deployments struct {
		Enabled bool   `json:"enabled"`
		Message string `json:"message"`
	} `json:"deployments"`
	Hosts struct {
		Total     int `json:"total"`
		Available int `json:"available"`
	} `json:"hosts"`
	Iaas struct {
		Provider string `json:"provider"`
		Region   string `json:"region"`
	} `json:"iaas"`
	OpenVpnServerPort int `json:"open_vpn_server_port"`
}

type Satellite interface {
	GetLocationInfo(name string, target ClusterTargetHeader) (*SatelliteLocationInfo, error)
}

type satellite struct {
	client     *client.Client
	pathPrefix string
}

func newSatelliteAPI(c *client.Client) Satellite {
	return &satellite{
		client: c,
	}
}

func (s *satellite) GetLocationInfo(name string, target ClusterTargetHeader) (*SatelliteLocationInfo, error) {
	SatLocationInfo := &SatelliteLocationInfo{}
	rawURL := fmt.Sprintf("/v2/satellite/getController?controller=%s", name)
	_, err := s.client.Get(rawURL, &SatLocationInfo, target.ToMap())
	if err != nil {
		return nil, err
	}
	return SatLocationInfo, err
}
